package kvstore

import (
	"github.com/gKits/PavoSQL/internal/btree"
	"github.com/gKits/PavoSQL/internal/freelist"
)

type KVWriter struct {
	KVReader
	kv      *KVStore
	free    freelist.Freelist
	nappend int
	changes map[uint64][]byte
	commit  func(*KVWriter) error
	abort   func(*KVWriter)
}

func newWriter(root uint64, commit func(*KVWriter) error, abort func(*KVWriter)) *KVWriter {
	w := &KVWriter{}
	w.tree = btree.New(root, PAGE_SIZE, w.getWriterPage, w.pullPage, w.allocPage, w.freePage)
	return w
}

func (w *KVWriter) Set(k, v []byte) error {
	return w.tree.Set(k, v)
}

func (w *KVWriter) Del(k []byte) (bool, error) {
	return w.tree.Delete(k)
}

func (w *KVWriter) Abort() {
	w.abort(w)
}

func (w *KVWriter) Commit() error {
	return w.commit(w)
}

func (w *KVWriter) getWriterPage(ptr uint64) ([]byte, error) {
	page, ok := w.changes[ptr]
	if !ok {
		return w.getPage(ptr)
	}
	return page, nil
}

func (w *KVWriter) pullPage(ptr uint64) ([]byte, error) {
	page, err := w.getWriterPage(ptr)
	if err != nil {
		return nil, err
	}

	if err := w.freePage(ptr); err != nil {
		return nil, err
	}

	return page, nil
}

func (w *KVWriter) allocPage(d []byte) (uint64, error) {
	return 0, nil
}

func (w *KVWriter) freePage(ptr uint64) error {
	w.changes[ptr] = nil
	return nil
}
