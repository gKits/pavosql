package node

import (
	"bytes"
	"encoding/binary"
	"slices"
)

type PointerNode struct {
	keys [][]byte
	ptrs []uint64
	size int
}

// Node interface functions

func (pn *PointerNode) Type() NodeType {
	return PNTR_NODE
}

func (pn PointerNode) NKeys() int {
	return len(pn.keys)
}

func (pn *PointerNode) Encode() []byte {
	var b []byte

	binary.BigEndian.AppendUint16(b, uint16(PNTR_NODE))
	binary.BigEndian.AppendUint16(b, uint16(len(pn.keys)))
	for i, k := range pn.keys {
		ptr := pn.ptrs[i]
		binary.BigEndian.AppendUint16(b, uint16(len(k)))
		b = append(b, k...)
		binary.BigEndian.AppendUint64(b, ptr)
	}

	return b
}

func (pn *PointerNode) Decode(d []byte) error {
	if NodeType(binary.BigEndian.Uint16(d[0:2])) != PNTR_NODE {
		return errNodeDecode
	}

	pn.size = len(d)

	nKeys := binary.BigEndian.Uint16(d[2:4])
	off := uint16(4)
	for i := uint16(0); i < nKeys; i++ {
		kSize := binary.BigEndian.Uint16(d[off : off+2])
		pn.keys = append(pn.keys, d[off+4:off+4+kSize])

		ptr := binary.BigEndian.Uint64(d[off+4+kSize : off+4+kSize+8])
		pn.ptrs = append(pn.ptrs, ptr)

		off += 4 + kSize + 8
	}

	return nil
}

func (pn *PointerNode) Size() int {
	return pn.size
}

func (pn *PointerNode) Key(i int) ([]byte, error) {
	if i >= len(pn.keys) {
		return nil, errNodeIdx
	}
	return pn.keys[i], nil
}

func (ln *PointerNode) Search(k []byte) (int, bool) {
	l := 0
	r := len(ln.keys)

	var i int
	var cmp int

	for i = r / 2; l < r; i = (l + r) / 2 {
		cmp = bytes.Compare(k, ln.keys[i])
		if cmp < 0 {
			r = i - 1
		} else if cmp > 0 {
			l = i + 1
		} else {
			return i, true
		}
	}

	if cmp > 1 {
		i++
	}

	return i, false
}

func (pn PointerNode) Merge(toMerge Node) (Node, error) {
	right, ok := toMerge.(*PointerNode)
	if !ok {
		return nil, errNodeMerge
	}

	if bytes.Compare(pn.keys[len(pn.keys)-1], right.keys[0]) >= 0 {
		return nil, errNodeMerge
	}

	pn.keys = append(pn.keys, right.keys...)
	pn.size += right.size
	pn.ptrs = append(pn.ptrs, right.ptrs...)

	return &pn, nil
}

func (pn PointerNode) Split() (Node, Node) {
	var half int
	var size int = 0

	for i, k := range pn.keys {
		size += 2 + len(k) + 8
		if size > pn.size/2 {
			half = i
			size -= 2 - len(k) - 8
			break
		}
	}

	split := PointerNode{
		keys: pn.keys[half:],
		ptrs: pn.ptrs[half:],
		size: pn.size - size,
	}

	pn.keys = pn.keys[:half]
	pn.ptrs = pn.ptrs[:half]
	pn.size = size

	return &pn, &split
}

// PointerNode specific functions

func (pn *PointerNode) Ptr(i int) (uint64, error) {
	if i >= len(pn.ptrs) {
		return 0, errNodeIdx
	}
	return pn.ptrs[i], nil
}

func (pn *PointerNode) KeyPtr(i int) ([]byte, uint64, error) {
	if i >= len(pn.keys) || i >= len(pn.ptrs) {
		return nil, 0, errNodeIdx
	}
	return pn.keys[i], pn.ptrs[i], nil
}

func (pn PointerNode) Insert(i int, k []byte, ptr uint64) (*PointerNode, error) {
	if i > len(pn.keys) || i > len(pn.ptrs) {
		return nil, errNodeIdx
	}

	pn.keys = slices.Insert(pn.keys, i, k)
	pn.ptrs = slices.Insert(pn.ptrs, i, ptr)
	pn.size += len(k) + 8 + 2

	return &pn, nil
}

func (pn PointerNode) Update(i int, k []byte, ptr uint64) (*PointerNode, error) {
	if i > len(pn.keys) || i > len(pn.ptrs) {
		return nil, errNodeIdx
	} else if slices.Equal(k, pn.keys[i]) {
		return nil, errNodeUpdate
	}

	pn.ptrs[i] = ptr

	return &pn, nil
}

func (pn PointerNode) Delete(i int) (*PointerNode, error) {
	if i > len(pn.keys) || i > len(pn.ptrs) {
		return nil, errNodeIdx
	}

	k := pn.keys[i]

	pn.keys = slices.Delete(pn.keys, i, i)
	pn.ptrs = slices.Delete(pn.ptrs, i, i)
	pn.size -= len(k) - 8 - 2

	return &pn, nil
}
