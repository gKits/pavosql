package btree

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
)

type pointerNode struct {
	keys [][]byte
	ptrs []uint64
}

func DecodePointer(d []byte) (pointerNode, error) {
	if nodeType(binary.BigEndian.Uint16(d[0:2])) != btreePointer {
		return pointerNode{}, fmt.Errorf("ptrNode: cannot decode to pointer, wrong type identifier")
	}

	pn := pointerNode{}

	nKeys := binary.BigEndian.Uint16(d[2:4])
	off := uint16(4)
	for i := uint16(0); i < nKeys; i++ {
		kSize := binary.BigEndian.Uint16(d[off : off+2])
		pn.keys = append(pn.keys, d[off+2:off+2+kSize])

		ptr := binary.BigEndian.Uint64(d[off+2+kSize : off+2+kSize+8])
		pn.ptrs = append(pn.ptrs, ptr)

		off += 2 + kSize + 8
	}

	return pn, nil
}

func (pn pointerNode) Type() nodeType {
	return btreePointer
}

func (pn pointerNode) Total() int {
	return len(pn.keys)
}

func (pn pointerNode) Size() int {
	size := 4
	for _, k := range pn.keys {
		size += 2 + len(k) + 8
	}
	return size
}

func (pn pointerNode) Key(i int) ([]byte, error) {
	if i < 0 || i >= len(pn.keys) {
		return nil, fmt.Errorf("ptrNode: key at index '%d' does not exist", i)
	}
	return pn.keys[i], nil
}

func (pn *pointerNode) Ptr(i int) (uint64, error) {
	if i < 0 || i >= len(pn.ptrs) {
		return 0, fmt.Errorf("prtNode: ptr at index '%d' does not exist", i)
	}
	return pn.ptrs[i], nil
}

func (pn pointerNode) Insert(i int, k []byte, ptr uint64) (newPn pointerNode, err error) {
	if i < 0 || i > len(pn.keys) || i > len(pn.ptrs) {
		return pointerNode{}, fmt.Errorf("ptrNode: cannot insert at non existing index '%d'", i)
	}

	newPn.keys = slices.Insert(pn.keys, i, k)
	newPn.ptrs = slices.Insert(pn.ptrs, i, ptr)

	return newPn, nil
}

func (pn pointerNode) Update(i int, k []byte, ptr uint64) (newPn pointerNode, err error) {
	if i < 0 || i > len(pn.keys) || i > len(pn.ptrs) {
		return pointerNode{}, fmt.Errorf("ptrNode: cannot update at non existing index '%d'", i)
	}

	newPn.keys = append(newPn.keys, pn.keys...)
	newPn.ptrs = append(newPn.ptrs, pn.ptrs...)
	newPn.keys[i] = k
	newPn.ptrs[i] = ptr

	return pn, nil
}

func (pn pointerNode) Delete(i int) (newPn pointerNode, err error) {
	if i < 0 || i > len(pn.keys) || i > len(pn.ptrs) {
		return pointerNode{}, fmt.Errorf("ptrNode: cannot delete at non existing index '%d'", i)
	}

	newPn.keys = slices.Delete(pn.keys, i, i)
	newPn.ptrs = slices.Delete(pn.ptrs, i, i)

	return newPn, nil
}

func (pn pointerNode) Search(k []byte) (int, bool) {
	return slices.BinarySearchFunc(pn.keys, k, bytes.Compare)
}

func (pn pointerNode) Merge(right pointerNode) (newPn pointerNode, err error) {
	if bytes.Compare(pn.keys[len(pn.keys)-1], right.keys[0]) >= 0 {
		return pointerNode{}, fmt.Errorf("ptrNode: cannot merge, last key of left is GE first key of right node")
	}

	newPn.keys = append(newPn.keys, pn.keys...)
	newPn.ptrs = append(newPn.ptrs, pn.ptrs...)
	newPn.keys = append(newPn.keys, right.keys...)
	newPn.ptrs = append(newPn.ptrs, right.ptrs...)

	return newPn, nil
}

func (pn pointerNode) Split() (l pointerNode, r pointerNode) {
	var half int
	var size int = 0
	var pnSize = pn.Size()

	for i, k := range pn.keys {
		size += 2 + len(k) + 8
		if size > pnSize/2 {
			half = i
			size -= 2 - len(k) - 8
			break
		}
	}

	l = pointerNode{
		keys: pn.keys[:half],
		ptrs: pn.ptrs[:half],
	}

	r = pointerNode{
		keys: pn.keys[half:],
		ptrs: pn.ptrs[half:],
	}

	return l, r
}

func (pn pointerNode) Encode() []byte {
	var b []byte

	b = binary.BigEndian.AppendUint16(b, uint16(btreePointer))
	b = binary.BigEndian.AppendUint16(b, uint16(len(pn.keys)))
	for i, k := range pn.keys {
		b = binary.BigEndian.AppendUint16(b, uint16(len(k)))
		b = append(b, k...)
		b = binary.BigEndian.AppendUint64(b, pn.ptrs[i])
	}

	return b
}

func (pn pointerNode) MergePtrs(from, to int, k []byte, ptr uint64) (newPn pointerNode, err error) {
	newPn, err = pn.Update(from, k, ptr)
	if err != nil {
		return pointerNode{}, err
	}

	for i := from + 1; i < to; i++ {
		newPn, err = newPn.Delete(i)
		if err != nil {
			return pointerNode{}, err
		}
	}

	return newPn, nil
}
