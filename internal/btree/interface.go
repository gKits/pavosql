package btree

import (
	"io"

	"github.com/gKits/PavoSQL/internal/btree/node"
)

type noder interface {
	Type() node.Type
	Len() int
	Size() int
	Key(i int) ([]byte, error)
	Search(key []byte) (int, bool)
	Bytes() ([]byte, error)
}

type Backend interface {
	GetNode(off uint64) (noder, error)
	NewReader() io.ReaderAt
	NewWriter() io.WriterAt
}
