package btree

import (
	"fmt"
)

type nodeType uint16

const (
	btreePointer nodeType = iota
	btreeLeaf
)

type node interface {
	Type() nodeType
	Total() int
	Size() int
	Key(int) ([]byte, error)
	Search([]byte) (int, bool)
	Encode() []byte
}

func decodeNode(d []byte) (node, error) {
	var n node
	var err error
	errMsg := "node: cannot decode node: %v"

	switch nodeType(d[0]) {
	case btreePointer:
		n, err = DecodePointer(d)
		break
	case btreeLeaf:
		n, err = DecodeLeaf(d)
		break
	default:
		return nil, fmt.Errorf(errMsg, "invalid node type")
	}

	return n, err
}
