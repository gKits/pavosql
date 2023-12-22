package tree

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
)

type LeafNode struct {
	keys [][]byte
	vals [][]byte
}

func (n *LeafNode) Decode(b []byte) error {
	off := uint16(0)
	if nodeType(binary.BigEndian.Uint16(b[off:])) != nodeLeaf {
		return errInvalNodeType
	}
	off += 2

	nKeys := binary.BigEndian.Uint32(b[off:])
	n.keys = make([][]byte, nKeys)
	n.vals = make([][]byte, nKeys)
	off += 4

	for i := uint32(0); i < nKeys; i++ {
		kSize := binary.BigEndian.Uint16(b[off:])
		off += 2
		vSize := binary.BigEndian.Uint16(b[off:])
		off += 2

		n.keys[i] = b[off : off+kSize]
		off += kSize
		n.vals[i] = b[off : off+vSize]
		off += vSize
	}

	return nil
}

func (n *LeafNode) Encode() []byte {
	b := make([]byte, n.Size())
	off := 0

	binary.BigEndian.PutUint16(b[off:], uint16(nodeLeaf))
	off += 2

	binary.BigEndian.PutUint32(b[off:], uint32(len(n.keys)))
	off += 4

	for i := 0; i < len(n.keys); i++ {
		binary.BigEndian.PutUint16(b[off:], uint16(len(n.keys[i])))
		off += 2

		binary.BigEndian.PutUint16(b[off:], uint16(len(n.vals[i])))
		off += 2

		copy(b[off:], n.keys[i])
		off += len(n.keys[i])

		copy(b[off:], n.vals[i])
		off += len(n.vals[i])
	}

	return b
}

func (n *LeafNode) NKeys() int {
	return len(n.keys)
}

func (n *LeafNode) Size() int {
	size := NodeHeader
	for i := 0; i < len(n.keys); i++ {
		size += 4 + len(n.keys[i]) + len(n.vals[i])
	}
	return size
}

func (n *LeafNode) Type() nodeType {
	return nodeLeaf
}

func (n *LeafNode) Key(idx int) ([]byte, error) {
	if idx >= len(n.keys) {
		return nil, errIdxOutOfRange
	}
	return n.keys[idx], nil
}

func (n *LeafNode) ValAt(idx int) ([]byte, error) {
	if idx >= len(n.vals) {
		return nil, errIdxOutOfRange
	}
	return n.vals[idx], nil
}

func (n *LeafNode) Val(k []byte) ([]byte, error) {
	idx, exists := n.Find(k)
	if !exists {
		return nil, errKeyNotExists
	}
	return n.vals[idx], nil
}

func (n *LeafNode) Insert(k, v []byte) error {
	idx, exists := n.Find(k)
	if exists {
		return errKeyExists
	}

	n.keys = slices.Insert(n.keys, idx, k)
	n.vals = slices.Insert(n.vals, idx, v)

	return nil
}

func (n *LeafNode) Update(k, v []byte) error {
	idx, exists := n.Find(k)
	if !exists {
		return errKeyNotExists
	}

	n.vals[idx] = v
	return nil
}

func (n *LeafNode) Delete(k []byte) error {
	idx, exists := n.Find(k)
	if !exists {
		return errKeyNotExists
	}

	n.keys = slices.Delete(n.keys, idx, idx+1)
	n.vals = slices.Delete(n.vals, idx, idx+1)

	return nil
}

func (n *LeafNode) Find(k []byte) (int, bool) {
	return slices.BinarySearchFunc(n.keys, k, bytes.Compare)
}

func (n *LeafNode) Split() (Node, Node) {
	l := &LeafNode{
		keys: slices.Clone(n.keys[:len(n.keys)/2]),
		vals: slices.Clone(n.vals[:len(n.vals)/2]),
	}
	r := &LeafNode{
		keys: slices.Clone(n.keys[len(n.keys)/2:]),
		vals: slices.Clone(n.vals[len(n.vals)/2:]),
	}
	return l, r
}

func (n *LeafNode) SplitBy(size int) []Node {
	parts := size / n.Size()
	nodes := make([]Node, parts)

	var p, accum, prev int
	for i := 0; i < len(n.keys) && p < parts; i++ {
		sizeOf := NodeHeader + len(n.keys[i]) + len(n.vals[i])
		accum += sizeOf
		if accum > size {
			nodes[p] = &LeafNode{
				keys: slices.Clone(n.keys[prev:i]),
				vals: slices.Clone(n.vals[prev:i]),
			}
			prev = i
		}
	}
	return nodes
}

func (n *LeafNode) Merge(m Node) error {
	mLeaf, ok := m.(*LeafNode)
	if !ok {
		return errMergeType
	}

	last, err := n.Key(n.NKeys() - 1)
	if err != nil {
		return err
	}

	if next, err := m.Key(0); err != nil {
		return err
	} else if bytes.Compare(last, next) >= 0 {
		return errMergeOrder
	}

	n.keys = append(n.keys, mLeaf.keys...)
	n.vals = append(n.vals, mLeaf.vals...)
	return nil
}

func (n *LeafNode) String() string {
	return fmt.Sprintf("LeafNode{keys: %s vals: %s}", n.keys, n.vals)
}
