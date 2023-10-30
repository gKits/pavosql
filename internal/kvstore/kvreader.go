package kvstore

import (
	"github.com/gKits/PavoSQL/internal/btree"
)

type KVReader struct {
	version uint64
	tree    btree.BTree
	chunks  [][]byte
	get     func(uint64) ([]byte, error)
	close   func(*KVReader)
	idx     int
}

func newReader(v, root uint64, chunks [][]byte, close func(*KVReader)) *KVReader {
	r := &KVReader{version: v, chunks: chunks, close: close}
	r.tree = btree.NewReadOnly(root, PAGE_SIZE, r.getPage)
	return r
}

func (r *KVReader) Get(k []byte) ([]byte, error) {
	return r.tree.Get(k)
}

func (r *KVReader) Seek(k []byte) (*btree.Iterator, error) {
	return nil, nil
}

func (r *KVReader) Close() {
	r.close(r)
}

func (r *KVReader) getPage(ptr uint64) ([]byte, error) {
	start := uint64(0)
	for _, chunk := range r.chunks {
		end := start + uint64(len(chunk)/PAGE_SIZE)
		if ptr < end {
			offset := PAGE_SIZE * (ptr - start)
			return chunk[offset : offset+PAGE_SIZE], nil
		}
		start = end
	}
	return r.get(ptr)
}
