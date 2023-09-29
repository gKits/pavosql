package store

import (
	"bytes"
	"encoding/binary"
	"slices"
)

type pointerNode struct {
	keys [][]byte
	ptrs []uint64
}

func (pn *pointerNode) Decode(d []byte) error {
	if nodeType(binary.BigEndian.Uint16(d[0:2])) != ptrNode {
		return errNodeDecode
	}

	nKeys := binary.BigEndian.Uint16(d[2:4])
	off := uint16(4)
	for i := uint16(0); i < nKeys; i++ {
		kSize := binary.BigEndian.Uint16(d[off : off+2])
		pn.keys = append(pn.keys, d[off+2:off+2+kSize])

		ptr := binary.BigEndian.Uint64(d[off+2+kSize : off+2+kSize+8])
		pn.ptrs = append(pn.ptrs, ptr)

		off += 2 + kSize + 8
	}

	return nil
}

// node interface methods

func (pn pointerNode) Type() nodeType {
	return ptrNode
}

func (pn pointerNode) Total() int {
	return len(pn.keys)
}

func (pn pointerNode) Encode() []byte {
	var b []byte

	b = binary.BigEndian.AppendUint16(b, uint16(ptrNode))
	b = binary.BigEndian.AppendUint16(b, uint16(len(pn.keys)))
	for i, k := range pn.keys {
		b = binary.BigEndian.AppendUint16(b, uint16(len(k)))
		b = append(b, k...)
		b = binary.BigEndian.AppendUint64(b, pn.ptrs[i])
	}

	return b
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
		return nil, errNodeIdx
	}
	return pn.keys[i], nil
}

func (ln pointerNode) Search(k []byte) (int, bool) {
	l := 0
	r := len(ln.keys)

	i := r / 2
	for {
		if cmp := bytes.Compare(k, ln.keys[i]); cmp < 0 {
			r = i
		} else if cmp > 0 {
			l = i + 1
		} else {
			return i, true
		}

		i = (l + r) / 2
		if l >= r {
			break
		}
	}

	if i < len(ln.keys) && bytes.Equal(ln.keys[i], k) {
		return i, true
	}

	return i, false
}

func (pn pointerNode) Merge(toMerge node) (node, error) {
	right, ok := toMerge.(*pointerNode)
	if !ok {
		return nil, errNodeMerge
	}

	if bytes.Compare(pn.keys[len(pn.keys)-1], right.keys[0]) >= 0 {
		return nil, errNodeMerge
	}

	pn.keys = append(pn.keys, right.keys...)
	pn.ptrs = append(pn.ptrs, right.ptrs...)

	return &pn, nil
}

func (pn pointerNode) Split() (node, node) {
	var half int
	var size int = 0
	var pnSize = pn.size()

	for i, k := range pn.keys {
		size += 2 + len(k) + 8
		if size > pnSize/2 {
			half = i
			size -= 2 - len(k) - 8
			break
		}
	}

	split := pointerNode{
		keys: pn.keys[half:],
		ptrs: pn.ptrs[half:],
	}

	pn.keys = pn.keys[:half]
	pn.ptrs = pn.ptrs[:half]

	return &pn, &split
}

// pointerNode specific methods

func (pn *pointerNode) Ptr(i int) (uint64, error) {
	if i < 0 || i >= len(pn.ptrs) {
		return 0, errNodeIdx
	}
	return pn.ptrs[i], nil
}

func (pn *pointerNode) KeyPtr(i int) ([]byte, uint64, error) {
	if i < 0 || i >= len(pn.keys) || i >= len(pn.ptrs) {
		return nil, 0, errNodeIdx
	}
	return pn.keys[i], pn.ptrs[i], nil
}

func (pn pointerNode) Insert(i int, k []byte, ptr uint64) (pointerNode, error) {
	if i < 0 || i > len(pn.keys) || i > len(pn.ptrs) {
		return pointerNode{}, errNodeIdx
	}

	pn.keys = slices.Insert(pn.keys, i, k)
	pn.ptrs = slices.Insert(pn.ptrs, i, ptr)

	return pn, nil
}

func (pn pointerNode) Update(i int, k []byte, ptr uint64) (pointerNode, error) {
	if i < 0 || i > len(pn.keys) || i > len(pn.ptrs) {
		return pointerNode{}, errNodeIdx
	} else if slices.Equal(k, pn.keys[i]) {
		return pointerNode{}, errNodeUpdate
	}

	pn.ptrs[i] = ptr

	return pn, nil
}

func (pn pointerNode) Delete(i int) (pointerNode, error) {
	if i < 0 || i > len(pn.keys) || i > len(pn.ptrs) {
		return pointerNode{}, errNodeIdx
	}

	pn.keys = slices.Delete(pn.keys, i, i)
	pn.ptrs = slices.Delete(pn.ptrs, i, i)

	return pn, nil
}
