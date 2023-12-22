package tree

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
)

type PointerNode struct {
	keys [][]byte
	ptrs []uint64
}

func (n *PointerNode) Decode(b []byte) error {
	off := uint16(0)
	if nodeType(binary.BigEndian.Uint16(b[off:])) != nodePointer {
		return errInvalNodeType
	}
	off += 2

	nKeys := binary.BigEndian.Uint32(b[off:])
	n.keys = make([][]byte, nKeys)
	n.ptrs = make([]uint64, nKeys)
	off += 4

	for i := uint32(0); i < nKeys; i++ {
		kSize := binary.BigEndian.Uint16(b[off:])
		off += 2

		n.keys[i] = b[off : off+kSize]
		off += kSize
		n.ptrs[i] = binary.BigEndian.Uint64(b[off : off+8])
		off += 8
	}

	return nil
}

func (n *PointerNode) Encode() []byte {
	b := make([]byte, n.Size())
	off := 0

	binary.BigEndian.PutUint16(b, uint16(nodePointer))
	off += 2

	binary.BigEndian.PutUint32(b[off:], uint32(len(n.keys)))
	off += 4

	for i := 0; i < len(n.keys); i++ {
		binary.BigEndian.PutUint16(b[off:], uint16(len(n.keys[i])))
		off += 2

		copy(b[off:], n.keys[i])
		off += len(n.keys[i])

		binary.BigEndian.PutUint64(b[off:], n.ptrs[i])
		off += 8
	}

	return b
}

func (n *PointerNode) NKeys() int {
	return len(n.keys)
}

func (n *PointerNode) Size() int {
	size := NodeHeader
	for i := 0; i < len(n.keys); i++ {
		size += 2 + len(n.keys[i]) + 8
	}
	return size
}

func (n *PointerNode) Type() nodeType {
	return nodePointer
}

func (n *PointerNode) Key(idx int) ([]byte, error) {
	if idx >= len(n.keys) {
		return nil, errIdxOutOfRange
	}
	return n.keys[idx], nil
}

func (n *PointerNode) PtrAt(idx int) (uint64, error) {
	if idx >= len(n.ptrs) {
		return 0, errIdxOutOfRange
	}
	return n.ptrs[idx], nil
}

func (n *PointerNode) Ptr(k []byte) (uint64, error) {
	idx, exists := n.Find(k)
	if !exists {
		return 0, errKeyNotExists
	}
	return n.ptrs[idx], nil
}

func (n *PointerNode) Insert(k []byte, ptr uint64) error {
	idx, exists := n.Find(k)
	if exists {
		return errKeyExists
	}

	n.keys = slices.Insert(n.keys, idx, k)
	n.ptrs = slices.Insert(n.ptrs, idx, ptr)

	return nil
}

func (n *PointerNode) Update(k []byte, ptr uint64) error {
	idx, exists := n.Find(k)
	if !exists {
		return errKeyNotExists
	}

	n.ptrs[idx] = ptr
	return nil
}

func (n *PointerNode) Delete(k []byte) error {
	idx, exists := n.Find(k)
	if !exists {
		return errKeyNotExists
	}

	n.keys = slices.Delete(n.keys, idx, idx+1)
	n.ptrs = slices.Delete(n.ptrs, idx, idx+1)

	return nil
}

func (n *PointerNode) Find(k []byte) (idx int, exists bool) {
	return slices.BinarySearchFunc(n.keys, k, bytes.Compare)
}

func (n *PointerNode) Split() (Node, Node) {
	l := &PointerNode{
		keys: slices.Clone(n.keys[:len(n.keys)/2]),
		ptrs: slices.Clone(n.ptrs[:len(n.ptrs)/2]),
	}
	r := &PointerNode{
		keys: slices.Clone(n.keys[len(n.keys)/2:]),
		ptrs: slices.Clone(n.ptrs[len(n.ptrs)/2:]),
	}
	return l, r
}

func (n *PointerNode) Merge(m Node) error {
	mPtr, ok := m.(*PointerNode)
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

	n.keys = append(n.keys, mPtr.keys...)
	n.ptrs = append(n.ptrs, mPtr.ptrs...)
	return nil
}

func (n *PointerNode) String() string {
	return fmt.Sprintf("PointerNode{keys: %s ptrs: %v}", n.keys, n.ptrs)
}
