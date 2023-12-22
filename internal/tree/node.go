package tree

import (
	"encoding/binary"
	"errors"
)

const NodeHeader = 6

var (
	errKeyExists     = errors.New("key already exists")
	errKeyNotExists  = errors.New("key does not exist")
	errInvalNodeType = errors.New("invalid node type")
	errIdxOutOfRange = errors.New("index is out of node range")
	errLeafHasNoPtr  = errors.New("leaf node cannot contain ptr")
	errMergeType     = errors.New("merging nodes need to have equal type")
	errMergeOrder    = errors.New("last key of left node needs to greater than first key of right node")
	errNodeAssert    = errors.New("node assertion failed")
)

type nodeType uint16

const (
	nodePointer nodeType = 100
	nodeLeaf    nodeType = 101
)

type Node interface {
	Type() nodeType                       // Returns node type.
	Key(idx int) ([]byte, error)          // Returns key at given index idx.
	NKeys() int                           // Returns the number of keys stored in the node.
	Size() int                            // Returns the encoded size of the node in bytes.
	Delete(k []byte) error                // Deletes the key and its value from the node.
	Find(k []byte) (idx int, exists bool) // Returns the index and existing status of the given key k.
	Encode() []byte                       // Encodes the node and returns the encoded byte stream.
	Decode([]byte) error                  // Decodes the give bytes stream into the node.
	Split() (l Node, r Node)              // Returns two nodes containing each half of the original node.
	Merge(Node) error                     // Right merges all keys and values onto node.
}

func NewNode(b []byte) (Node, error) {
	var n Node
	switch nodeType(binary.BigEndian.Uint16(b)) {
	case nodePointer:
		n = &PointerNode{}
	case nodeLeaf:
		n = &LeafNode{}
	default:
		return nil, errInvalNodeType
	}

	if err := n.Decode(b); err != nil {
		return nil, err
	}
	return n, nil
}
