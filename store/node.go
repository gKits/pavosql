package store

import (
	"encoding/binary"
	"errors"
)

type nodeType uint16

const (
	ptrNode nodeType = iota
	lfNode
	flNode
)

var (
	errNodeIdx     = errors.New("node: index out of key range")
	errNodeType    = errors.New("node: unknown node type")
	errNodeDecode  = errors.New("node: cannot decode, type does not match")
	errNodeUpdate  = errors.New("node: cannot update k-v, keys must match")
	errNodeMerge   = errors.New("node: cannot merge, lefts last key must be less than rights first key")
	errNodeUseless = errors.New("node: useless method, used to implement interface for FreelistNode")
)

type node interface {
	Type() nodeType
	Total() int
	// decode([]byte) error
	Encode() []byte
	Size() int
	Key(int) ([]byte, error)
	Search([]byte) (int, bool)
	Merge(node) (node, error)
	Split() (node, node)
}

func decodeNode(d []byte) (node, error) {
	typ := nodeType(binary.BigEndian.Uint16(d[0:2]))

	switch typ {
	case ptrNode:
		n := pointerNode{}
		n.Decode(d)
		return &n, nil

	case lfNode:
		n := leafNode{}
		n.Decode(d)
		return &n, nil

	case flNode:
		n := freelistNode{}
		n.Decode(d)
		return &n, nil
	}

	return nil, errNodeType
}
