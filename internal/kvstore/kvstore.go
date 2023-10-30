package kvstore

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/gKits/PavoSQL/internal/freelist"
	"github.com/gKits/PavoSQL/pkg/vcache"
)

const PAGE_SIZE = 4096

type KVStore struct {
	kvOpts                                // options type embeded
	f       *os.File                      // the database file
	root    uint64                        // pointer to the root of the btree
	free    freelist.Freelist             // freelist managing free pages
	cache   vcache.VCache[uint64, []byte] // cache of freed pages still to be read
	version uint64                        // latest version of the kv store
	wLock   sync.Mutex                    // write lock allowing only a single concurrent writer a time
	fLock   sync.RWMutex                  // file lock making sure file is not read and written at the same time
}

func New(opts ...kvOptFunc) *KVStore {
	opt := defaultOpts()
	for _, fn := range opts {
		fn(&opt)
	}

	return &KVStore{
		kvOpts: opt,
		cache:  vcache.New[uint64, []byte](0),
	}
}

func (kv *KVStore) Open() error {
	f, err := os.OpenFile(kv.path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		kv.Close()
		return err
	}
	kv.f = f

	if err := kv.loadMaster(); err != nil {
		kv.Close()
		return err
	}

	return nil
}

func (kv *KVStore) Close() {
	kv.f.Close()
}

func (kv *KVStore) Read() (*KVReader, error) {
	r := newReader(kv.version, kv.root, kv.mmap.chunks, kv.endRead)

	return r, nil
}

func (kv *KVStore) Write() (*KVWriter, error) {
	kv.wLock.Lock()

	w := newWriter(kv.root, kv.commitWrite, kv.abortWrite)

	w.version = kv.version

	return w, nil
}

func (kv *KVStore) endRead(r *KVReader) {
}

func (kv *KVStore) commitWrite(w *KVWriter) error {
	kv.root = w.tree.Root
	return nil
}

func (kv *KVStore) abortWrite(w *KVWriter) {
	kv.wLock.Unlock()
}

func (kv *KVStore) loadMaster() error {
	data := kv.mmap.chunks[0]
	root := binary.BigEndian.Uint64(data[16:24])
	free := binary.BigEndian.Uint64(data[24:32])
	npages := binary.BigEndian.Uint64(data[32:40])
	version := binary.BigEndian.Uint64(data[40:48])

	if !bytes.Equal([]byte(kv.sig), data[:16]) {

	}

	if npages < 1 || npages > uint64(kv.mmap.fileSize/PAGE_SIZE) {
	}

	kv.root = root
	kv.version = version

	return nil
}

func (kv *KVStore) flushMaster() error {
	var data [40]byte

	copy(data[:16], []byte(kv.sig))
	binary.BigEndian.PutUint64(data[16:], kv.root)
	binary.BigEndian.PutUint64(data[24:], kv.root)
	binary.BigEndian.PutUint64(data[32:], kv.page.flushed)

	_, err := kv.f.WriteAt(data[:], 0)
	return err
}

func (kv *KVStore) writePages(changes map[uint64][]byte) (err error) {
	path, name := filepath.Split(kv.path)

	tmp, err := os.CreateTemp(path, name)
	if err != nil {
		return err
	}
	defer func() {
		tmp.Close()
		if err != nil {
			os.Remove(tmp.Name())
		}
	}()

	tmpName := tmp.Name()

	if _, err := io.Copy(tmp, kv.f); err != nil {
		return err
	}

	for ptr, page := range changes {
		if _, err := tmp.WriteAt(page, int64(ptr)); err != nil {
			return err
		}
	}

	if err := tmp.Sync(); err != nil {
		return err
	}

	destInfo, err := kv.f.Stat()
	if err != nil {
		return err
	}

	if err := tmp.Chmod(destInfo.Mode()); err != nil {
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	kv.fLock.Lock()
	defer kv.fLock.Unlock()

	if err := kv.f.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpName, kv.path); err != nil {
		return err
	}

	kv.f, err = os.Open(kv.path)
	if err != nil {
		return err
	}

	return nil
}

func (kv *KVStore) getFilePage(ptr uint64) ([]byte, error) {
	kv.fLock.RLock()
	defer kv.fLock.RUnlock()

	page := make([]byte, PAGE_SIZE)
	if _, err := kv.f.ReadAt(page, int64(ptr)); err != nil {
		return nil, err
	}
	return page, nil
}

func (kv *KVStore) pullPage(ptr uint64) ([]byte, error) {
	return nil, nil
}

func (kv *KVStore) allocPage(page []byte) (uint64, error) {
	return 0, nil
}

func (kv *KVStore) freeFilePage(ptr uint64) error {
	page, err := kv.getFilePage(ptr)
	if err != nil {
		return err
	}

	kv.cache.Cache(ptr, page) // store freed page in versioned cache
	return nil
}

/*
kvOpts struct is embeded into KVStore to implement the functional options
pattern for KVStore
*/
type kvOptFunc func(*kvOpts)

type kvOpts struct {
	path string // path to the database file
	sig  string // the file signature
}

func defaultOpts() kvOpts {
	return kvOpts{
		path: "/var/lib/pavosql",
		sig:  "PavoSQL_DB_File:",
	}
}

func WithPath(path string) kvOptFunc {
	return func(opts *kvOpts) {
		opts.path = path
	}
}

func WithSignature(sig string) kvOptFunc {
	s := make([]byte, 16)
	copy(s, sig)

	return func(opts *kvOpts) {
		opts.sig = string(s)
	}
}
