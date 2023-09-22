package backend

import (
	"bytes"
	"encoding/binary"
	"slices"
)

type pointerNode struct {
	keys [][]byte
	ptrs []uint64
}

func (pn *pointerNode) decode(d []byte) error {
	if nodeType(binary.BigEndian.Uint16(d[0:2])) != ptrNode {
		return errNodeDecode
	}

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

// node interface methods

func (pn pointerNode) typ() nodeType {
	return ptrNode
}

func (pn pointerNode) total() int {
	return len(pn.keys)
}

func (pn pointerNode) encode() []byte {
	var b []byte

	binary.BigEndian.AppendUint16(b, uint16(ptrNode))
	binary.BigEndian.AppendUint16(b, uint16(len(pn.keys)))
	for i, k := range pn.keys {
		ptr := pn.ptrs[i]
		binary.BigEndian.AppendUint16(b, uint16(len(k)))
		b = append(b, k...)
		binary.BigEndian.AppendUint64(b, ptr)
	}

	return b
}

func (pn pointerNode) size() int {
	size := 4
	for _, k := range pn.keys {
		size += 4 + len(k) + 8
	}
	return size
}

func (pn pointerNode) key(i int) ([]byte, error) {
	if i >= len(pn.keys) {
		return nil, errNodeIdx
	}
	return pn.keys[i], nil
}

func (ln pointerNode) search(k []byte) (int, bool) {
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

func (pn pointerNode) merge(toMerge node) (node, error) {
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

func (pn pointerNode) split() (node, node) {
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

func (pn *pointerNode) ptr(i int) (uint64, error) {
	if i >= len(pn.ptrs) {
		return 0, errNodeIdx
	}
	return pn.ptrs[i], nil
}

func (pn *pointerNode) keyPtr(i int) ([]byte, uint64, error) {
	if i >= len(pn.keys) || i >= len(pn.ptrs) {
		return nil, 0, errNodeIdx
	}
	return pn.keys[i], pn.ptrs[i], nil
}

func (pn pointerNode) insert(i int, k []byte, ptr uint64) (pointerNode, error) {
	if i > len(pn.keys) || i > len(pn.ptrs) {
		return pointerNode{}, errNodeIdx
	}

	pn.keys = slices.Insert(pn.keys, i, k)
	pn.ptrs = slices.Insert(pn.ptrs, i, ptr)

	return pn, nil
}

func (pn pointerNode) update(i int, k []byte, ptr uint64) (pointerNode, error) {
	if i > len(pn.keys) || i > len(pn.ptrs) {
		return pointerNode{}, errNodeIdx
	} else if slices.Equal(k, pn.keys[i]) {
		return pointerNode{}, errNodeUpdate
	}

	pn.ptrs[i] = ptr

	return pn, nil
}

func (pn pointerNode) delete(i int) (pointerNode, error) {
	if i > len(pn.keys) || i > len(pn.ptrs) {
		return pointerNode{}, errNodeIdx
	}

	pn.keys = slices.Delete(pn.keys, i, i)
	pn.ptrs = slices.Delete(pn.ptrs, i, i)

	return pn, nil
}
