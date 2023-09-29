package store

import (
	"bytes"
	"encoding/binary"
	"slices"
)

type leafNode struct {
	keys [][]byte
	vals [][]byte
}

func (ln *leafNode) Decode(d []byte) error {
	if nodeType(binary.BigEndian.Uint16(d[0:2])) != lfNode {
		return errNodeDecode
	}

	nKeys := binary.BigEndian.Uint16(d[2:4])
	ln.keys = make([][]byte, nKeys)
	ln.vals = make([][]byte, nKeys)

	off := uint16(4)
	for i := 0; uint16(i) < nKeys; i++ {
		kSize := binary.BigEndian.Uint16(d[off : off+2])
		vSize := binary.BigEndian.Uint16(d[off+2 : off+4])

		ln.keys[i] = d[off+4 : off+4+kSize]
		ln.vals[i] = d[off+4+kSize : off+4+kSize+vSize]

		off += 4 + kSize + vSize
	}

	return nil
}

// node interface methods

func (ln leafNode) Type() nodeType {
	return lfNode
}

func (ln leafNode) Total() int {
	return len(ln.keys)
}

func (ln leafNode) Encode() []byte {
	var b []byte

	b = binary.BigEndian.AppendUint16(b, uint16(lfNode))
	b = binary.BigEndian.AppendUint16(b, uint16(len(ln.keys)))
	for i, k := range ln.keys {
		v := ln.vals[i]
		b = binary.BigEndian.AppendUint16(b, uint16(len(k)))
		b = binary.BigEndian.AppendUint16(b, uint16(len(v)))
		b = append(b, k...)
		b = append(b, v...)
	}

	return b
}

func (ln leafNode) Size() int {
	size := 4
	for i, k := range ln.keys {
		v := ln.vals[i]
		size += 4 + len(k) + len(v)
	}
	return size
}

func (ln leafNode) Key(i int) ([]byte, error) {
	if i < 0 || i >= len(ln.keys) {
		return nil, errNodeIdx
	}
	return ln.keys[i], nil
}

func (ln leafNode) Search(k []byte) (int, bool) {
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

func (ln leafNode) Merge(toMerge node) (node, error) {
	right, ok := toMerge.(*leafNode)
	if !ok {
		return nil, errNodeMerge
	}

	if bytes.Compare(ln.keys[len(ln.keys)-1], right.keys[0]) >= 0 {
		return nil, errNodeMerge
	}

	ln.keys = append(ln.keys, right.keys...)
	ln.vals = append(ln.vals, right.vals...)

	return &ln, nil
}

func (ln leafNode) Split() (node, node) {
	var half int
	var size int = 0
	lnSize := ln.Size()

	for i, k := range ln.keys {
		v := ln.vals[i]
		size += 4 + len(k) + len(v)
		if size > lnSize/2 {
			half = i
			size -= 4 - len(k) - len(v)
			break
		}
	}

	split := leafNode{
		keys: ln.keys[half:],
		vals: ln.vals[half:],
	}

	ln.keys = ln.keys[:half]
	ln.vals = ln.vals[:half]

	return &ln, &split
}

// leafNode specific methods

func (ln *leafNode) Val(i int) ([]byte, error) {
	if i < 0 || i >= len(ln.vals) {
		return nil, errNodeIdx
	}
	return ln.vals[i], nil
}

func (ln *leafNode) KeyVal(i int) ([]byte, []byte, error) {
	if i < 0 || i >= len(ln.keys) || i >= len(ln.vals) {
		return nil, nil, errNodeIdx
	}
	return ln.keys[i], ln.vals[i], nil
}

func (ln leafNode) Insert(i int, k, v []byte) (leafNode, error) {
	if i < 0 || i > len(ln.keys) || i > len(ln.vals) {
		return leafNode{}, errNodeIdx
	}

	ln.keys = slices.Insert(ln.keys, i, k)
	ln.vals = slices.Insert(ln.vals, i, k)

	return ln, nil
}

func (ln leafNode) Update(i int, k, v []byte) (leafNode, error) {
	if i < 0 || i > len(ln.keys) || i > len(ln.vals) {
		return leafNode{}, errNodeIdx
	} else if slices.Equal(k, ln.keys[i]) {
		return leafNode{}, errNodeUpdate
	}

	ln.vals[i] = v

	return ln, nil
}

func (ln leafNode) Delete(i int) (leafNode, error) {
	if i < 0 || i > len(ln.keys) || i > len(ln.vals) {
		return leafNode{}, errNodeIdx
	}

	ln.keys = slices.Delete(ln.keys, i, i)
	ln.vals = slices.Delete(ln.vals, i, i)

	return ln, nil
}
