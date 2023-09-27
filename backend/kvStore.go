package backend

import (
	"encoding/binary"
	"errors"
	"os"
	"slices"
	"syscall"
)

var (
	errKVBadMaster    = errors.New("kv: malformed master page")
	errKVBadPtr       = errors.New("kv: bad pointer")
	errKVPageTooLarge = errors.New("kv: page too large")
)

type KVStore struct {
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
	temp map[uint64][pageSize]byte
	Sig  [16]byte
}

func (kv *KVStore) Open() error {
	f, err := os.OpenFile(kv.path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		goto err
	}
	kv.f = f

	if err = kv.mmap.init(kv.f); err != nil {
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

func (kv *KVStore) Close() {
	if err := kv.mmap.close(); err != nil {
		panic(err)
	}
	kv.f.Close()
}

func (kv *KVStore) Get(k []byte) ([]byte, error) {
	return kv.bt.Get(k)
}

func (kv *KVStore) Set(k, v []byte) error {
	if err := kv.bt.Insert(k, v); err != nil {
		return err
	}
	return kv.flush()
}

func (kv *KVStore) Del(k []byte) (bool, error) {
	del, err := kv.bt.Delete(k)
	if err != nil {
		return false, err
	}
	return del, kv.flush()
}

func (kv *KVStore) flush() error {
	if err := kv.writePages(); err != nil {
		return err
	}
	return kv.syncPages()
}

func (kv *KVStore) writePages() error {
	nPages := int(kv.page.flushed) + len(kv.page.temp)

	if err := kv.extendFile(nPages); err != nil {
		return err
	}

	if err := kv.mmap.extend(kv.f, nPages); err != nil {
		return err
	}

	for i, page := range kv.page.temp {
		n, err := kv.getPage(kv.page.flushed + uint64(i))
		if err != nil {
			return err
		}
		copy(n.encode(), page)
	}

	return nil
}

func (kv *KVStore) syncPages() error {
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

func (kv *KVStore) loadMaster() error {
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

	if 1 > nPages || nPages > (uint64(kv.mmap.fileSize/pageSize)) || 0 > btRoot || btRoot >= nPages {
		return errKVBadMaster
	}

	kv.bt.root = btRoot
	kv.fl.root = flRoot

	return nil
}

func (kv *KVStore) writeMaster() error {
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

func (kv *KVStore) extendFile(n int) error {
	filePages := kv.mmap.fileSize / pageSize
	if filePages >= n {
		return nil
	}

	for filePages < n {
		inc := filePages / 8
		if inc < 1 {
			inc = 1
		}
		filePages += inc
	}

	fileSize := filePages * pageSize
	if err := syscall.Fallocate(int(kv.f.Fd()), 0, 0, int64(fileSize)); err != nil {
		return err
	}

	kv.mmap.fileSize = fileSize
	return nil
}

func (kv *KVStore) getPage(ptr uint64) (node, error) {
	page, ok := kv.page.updates[ptr]
	if !ok {
		start := uint64(0)
		for _, chunk := range kv.mmap.chunks {
			end := start + uint64(len(chunk))/pageSize
			if ptr < end {
				off := pageSize * (ptr - start)
				page = chunk[off : off+pageSize]
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

func (kv *KVStore) pullPage(ptr uint64) (node, error) {
	n, err := kv.getPage(ptr)
	if err != nil {
		return nil, err
	}

	if err := kv.freePage(ptr); err != nil {
		return nil, err
	}

	return n, nil
}

func (kv *KVStore) allocPage(n node) (uint64, error) {
	if n.size() > pageSize {
		return 0, errKVPageTooLarge
	}

	ptr, err := kv.fl.pop()
	if errors.Is(err, errFLPopEmpty) {
		ptr = kv.page.flushed + uint64(kv.page.nappend)
		kv.page.nappend++
	} else if err != nil {
		return 0, err
	}

	kv.page.updates[ptr] = n.encode()
	return ptr, nil
}

func (kv *KVStore) freePage(ptr uint64) error {
	kv.page.updates[ptr] = nil
	return nil
}
