package btree

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
)

type leafNode struct {
	keys [][]byte
	vals [][]byte
}

func DecodeLeaf(d []byte) (leafNode, error) {
	ln := leafNode{}

	if nodeType(binary.BigEndian.Uint16(d[0:2])) != btreeLeaf {
		return leafNode{}, fmt.Errorf("leafNode: cannot decode to leaf, wrong type identifier")
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

	return ln, nil

}

func (ln leafNode) Type() nodeType {
	return btreeLeaf
}

func (ln leafNode) Total() int {
	return len(ln.keys)
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
		return nil, fmt.Errorf("leafNode: key at index '%d' does not exist", i)
	}
	return ln.keys[i], nil
}

func (ln *leafNode) Val(i int) ([]byte, error) {
	if i < 0 || i >= len(ln.vals) {
		return nil, fmt.Errorf("leafNode: val at index '%d' does not exist", i)
	}
	return ln.vals[i], nil
}

func (ln leafNode) Insert(i int, k, v []byte) (newLn leafNode, err error) {
	if i < 0 || i > len(ln.keys) || i > len(ln.vals) {
		return leafNode{}, fmt.Errorf("leafNode: cannot insert at non existing index '%d'", i)
	}

	newLn.keys = slices.Insert(ln.keys, i, k)
	newLn.vals = slices.Insert(ln.vals, i, k)

	return ln, nil
}

func (ln leafNode) Update(i int, k, v []byte) (newLn leafNode, err error) {
	if i < 0 || i > len(ln.keys) || i > len(ln.vals) {
		return leafNode{}, fmt.Errorf("leafNode: cannot update at non existing index '%d'", i)
	}

	ln.keys[i] = k
	ln.vals[i] = v

	return ln, nil
}

func (ln leafNode) Delete(i int) (leafNode, error) {
	if i < 0 || i > len(ln.keys) || i > len(ln.vals) {
		return leafNode{}, fmt.Errorf("leafNode: cannot delete at non existing index '%d'", i)
	}

	ln.keys = slices.Delete(ln.keys, i, i)
	ln.vals = slices.Delete(ln.vals, i, i)

	return ln, nil
}

func (ln leafNode) Search(k []byte) (int, bool) {
	return slices.BinarySearchFunc(ln.keys, k, bytes.Compare)
}

func (ln leafNode) Merge(right leafNode) (leafNode, error) {
	if bytes.Compare(ln.keys[len(ln.keys)-1], right.keys[0]) >= 0 {
		return leafNode{}, fmt.Errorf("leafNode: cannot merge, last key of left is GE first key of right node")
	}

	ln.keys = append(ln.keys, right.keys...)
	ln.vals = append(ln.vals, right.vals...)

	return ln, nil
}

func (ln leafNode) Split() (leafNode, leafNode) {
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

	return ln, split
}

func (ln leafNode) Encode() []byte {
	var b []byte

	b = binary.BigEndian.AppendUint16(b, uint16(btreeLeaf))
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
