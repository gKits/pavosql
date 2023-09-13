package node

import (
	"bytes"
	"encoding/binary"
	"slices"
)

type LeafNode struct {
	keys [][]byte
	vals [][]byte
	size int
}

func NewLeafNode() LeafNode {
	return LeafNode{}
}

// Node interface functions

func (ln *LeafNode) Type() NodeType {
	return LEAF_NODE
}

func (ln LeafNode) NKeys() int {
	return len(ln.keys)
}

func (ln *LeafNode) Encode() []byte {
	var b []byte

	b = binary.BigEndian.AppendUint16(b, uint16(LEAF_NODE))
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

func (ln *LeafNode) Decode(d []byte) error {
	if NodeType(binary.BigEndian.Uint16(d[0:2])) != LEAF_NODE {
		return errNodeDecode
	}

	ln.size = len(d)

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

func (ln *LeafNode) Size() int {
	return ln.size
}

func (ln *LeafNode) Key(i int) ([]byte, error) {
	if i >= len(ln.keys) {
		return nil, errNodeIdx
	}
	return ln.keys[i], nil
}

func (ln *LeafNode) Search(k []byte) (int, bool) {
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

func (ln LeafNode) Merge(toMerge Node) (Node, error) {
	right, ok := toMerge.(*LeafNode)
	if !ok {
		return nil, errNodeMerge
	}

	if bytes.Compare(ln.keys[len(ln.keys)-1], right.keys[0]) >= 0 {
		return nil, errNodeMerge
	}

	ln.keys = append(ln.keys, right.keys...)
	ln.vals = append(ln.vals, right.vals...)
	ln.size += right.size

	return &ln, nil
}

func (ln LeafNode) Split() (Node, Node) {
	var half int
	var size int = 0

	for i, k := range ln.keys {
		v := ln.vals[i]
		size += 4 + len(k) + len(v)
		if size > ln.size/2 {
			half = i
			size -= 4 - len(k) - len(v)
			break
		}
	}

	split := LeafNode{
		keys: ln.keys[half:],
		vals: ln.vals[half:],
		size: ln.size - size,
	}

	ln.keys = ln.keys[:half]
	ln.vals = ln.vals[:half]
	ln.size = size

	return &ln, &split
}

// LeafNode specific functions

func (ln *LeafNode) Val(i int) ([]byte, error) {
	if i >= len(ln.vals) {
		return nil, errNodeIdx
	}
	return ln.vals[i], nil
}

func (ln *LeafNode) KeyVal(i int) ([]byte, []byte, error) {
	if i >= len(ln.keys) || i >= len(ln.vals) {
		return nil, nil, errNodeIdx
	}
	return ln.keys[i], ln.vals[i], nil
}

func (ln LeafNode) Insert(i int, k, v []byte) (*LeafNode, error) {
	if i > len(ln.keys) || i > len(ln.vals) {
		return nil, errNodeIdx
	}

	ln.keys = slices.Insert(ln.keys, i, k)
	ln.vals = slices.Insert(ln.vals, i, k)
	ln.size += len(k) + len(v) + 4

	return &ln, nil
}

func (ln LeafNode) Update(i int, k, v []byte) (*LeafNode, error) {
	if i > len(ln.keys) || i > len(ln.vals) {
		return nil, errNodeIdx
	} else if slices.Equal(k, ln.keys[i]) {
		return nil, errNodeUpdate
	}

	og := ln.vals[i]
	ln.vals[i] = v
	ln.size += len(v) - len(og)

	return &ln, nil
}

func (ln LeafNode) Delete(i int) (*LeafNode, error) {
	if i > len(ln.keys) || i > len(ln.vals) {
		return nil, errNodeIdx
	}

	k := ln.keys[i]
	v := ln.vals[i]

	ln.keys = slices.Delete(ln.keys, i, i)
	ln.vals = slices.Delete(ln.vals, i, i)
	ln.size -= len(k) - len(v) - 4

	return &ln, nil
}
