package backend

import (
	"encoding/binary"
	"errors"
	"os"
	"slices"
)

var (
	errKVBadMaster    = errors.New("kv: malformed master page")
	errKVBadPtr       = errors.New("kv: bad pointer")
	errKVPageTooLarge = errors.New("kv: page too large")
)

type Store struct {
	path string
	f    *os.File
	bt   bTree
	fl   freelist
	mmap mmap
	page struct {
		flushed uint64
		temp    [][]byte
		nappend int
		updates map[uint64][]byte
	}
	temp map[uint64][PageSize]byte
	Sig  [16]byte
}

func (kv *Store) Open() error {
	f, err := os.OpenFile(kv.path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		goto err
	}
	kv.f = f

	if err = kv.mmap.Init(kv.f); err != nil {
		goto err
	}

	kv.bt = bTree{0, kv.getPage, kv.pullPage, kv.allocPage, kv.freePage}
	kv.fl = freelist{
		0,
		func(ptr uint64) (freelistNode, error) {
			n, err := kv.getPage(ptr)
			if err != nil {
				return freelistNode{}, err
			}
			return n.(freelistNode), nil
		},
		func(ptr uint64) (freelistNode, error) {
			n, err := kv.pullPage(ptr)
			if err != nil {
				return freelistNode{}, err
			}
			return n.(freelistNode), nil
		},
		kv.allocPage,
		kv.freePage,
	}

	if err = kv.loadMaster(); err != nil {
		goto err
	}

	return nil

err:
	kv.Close()
	return err
}

func (kv *Store) Close() {
	if err := kv.mmap.Close(); err != nil {
		panic(err)
	}
	kv.f.Close()
}

func (kv *Store) Get(k []byte) ([]byte, error) {
	return kv.bt.Get(k)
}

func (kv *Store) Set(k, v []byte) error {
	if err := kv.bt.Insert(k, v); err != nil {
		return err
	}
	return kv.flush()
}

func (kv *Store) Del(k []byte) (bool, error) {
	del, err := kv.bt.Delete(k)
	if err != nil {
		return false, err
	}
	return del, kv.flush()
}

func (kv *Store) flush() error {
	if err := kv.writePages(); err != nil {
		return err
	}
	return kv.syncPages()
}

func (kv *Store) writePages() error {
	nPages := int(kv.page.flushed) + len(kv.page.temp)

	if err := kv.mmap.ExtendFile(kv.f, nPages); err != nil {
		return err
	}

	if err := kv.mmap.Extend(kv.f, nPages); err != nil {
		return err
	}

	for i, page := range kv.page.temp {
		n, err := kv.getPage(kv.page.flushed + uint64(i))
		if err != nil {
			return err
		}
		copy(n.Encode(), page)
	}

	return nil
}

func (kv *Store) syncPages() error {
	if err := kv.f.Sync(); err != nil {
		return err
	}

	kv.page.flushed += uint64(len(kv.page.temp))
	clear(kv.page.temp)

	if err := kv.writeMaster(); err != nil {
		return err
	}

	if err := kv.f.Sync(); err != nil {
		return err
	}

	return nil
}

func (kv *Store) loadMaster() error {
	if kv.mmap.fileSize == 0 {
		kv.page.flushed = 1
		return nil
	}

	b := kv.mmap.chunks[0]

	if slices.Equal(kv.Sig[:], b[:16]) {
		return errKVBadMaster
	}

	btRoot := binary.BigEndian.Uint64(b[16:24])
	flRoot := binary.BigEndian.Uint64(b[24:32])
	nPages := binary.BigEndian.Uint64(b[32:40])

	if 1 > nPages || nPages > (uint64(kv.mmap.fileSize/PageSize)) || 0 > btRoot || btRoot >= nPages {
		return errKVBadMaster
	}

	kv.bt.root = btRoot
	kv.fl.root = flRoot

	return nil
}

func (kv *Store) writeMaster() error {
	d := [40]byte{}
	copy(d[:16], kv.Sig[:])
	binary.BigEndian.PutUint64(d[16:], kv.bt.root)
	binary.BigEndian.PutUint64(d[24:], kv.fl.root)
	binary.BigEndian.PutUint64(d[32:], kv.page.flushed)

	if _, err := kv.f.WriteAt(nil, 0); err != nil {
		return err
	}

	return nil
}

func (kv *Store) getPage(ptr uint64) (node, error) {
	page, ok := kv.page.updates[ptr]
	if !ok {
		start := uint64(0)
		for _, chunk := range kv.mmap.chunks {
			end := start + uint64(len(chunk))/PageSize
			if ptr < end {
				off := PageSize * (ptr - start)
				page = chunk[off : off+PageSize]
				break
			}
			start = end
		}
	}
	if page == nil {
		return nil, errKVBadPtr
	}
	return decodeNode(page)
}

func (kv *Store) pullPage(ptr uint64) (node, error) {
	n, err := kv.getPage(ptr)
	if err != nil {
		return nil, err
	}

	if err := kv.freePage(ptr); err != nil {
		return nil, err
	}

	return n, nil
}

func (kv *Store) allocPage(n node) (uint64, error) {
	if n.Size() > PageSize {
		return 0, errKVPageTooLarge
	}

	ptr, err := kv.fl.Pop()
	if errors.Is(err, errFLPopEmpty) {
		ptr = kv.page.flushed + uint64(kv.page.nappend)
		kv.page.nappend++
	} else if err != nil {
		return 0, err
	}

	kv.page.updates[ptr] = n.Encode()
	return ptr, nil
}

func (kv *Store) freePage(ptr uint64) error {
	kv.page.updates[ptr] = nil
	if err := kv.fl.Push(ptr); err != nil {
		return err
	}
	return nil
}
