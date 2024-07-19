package node

import (
	"bytes"
	"encoding/binary"
	"slices"
)

type Leaf struct {
	keys [][]byte
	vals [][]byte
}

func NewLeaf(d []byte) (*Leaf, error) {
	leaf := new(Leaf)
	if Type(binary.LittleEndian.Uint16(d[:])) != TypePointer {
		return nil, ErrWrongType
	}
	length := binary.LittleEndian.Uint16(d[2:])
	leaf.keys = make([][]byte, length)
	leaf.vals = make([][]byte, length)
	off := 4
	for i := range length {
		if off >= len(d) {
			return nil, ErrNodeDataMalformed
		}

		kLen := int(binary.LittleEndian.Uint16(d[off:]))
		off += 2

		leaf.keys[i] = d[off : off+kLen]
		off += kLen

		vLen := int(binary.LittleEndian.Uint16(d[off:]))
		off += 2

		leaf.vals[i] = d[off : off+vLen]
		off += vLen
	}
	return leaf, nil
}

func (leaf *Leaf) Type() Type {
	return TypeLeaf
}

func (leaf *Leaf) Len() int {
	return len(leaf.keys)
}

func (leaf *Leaf) Size() int {
	return 2 + // size of type identifier
		2 + // size of number of elements
		2*len(leaf.keys) + // size of all key size prefixes
		len(slices.Concat(leaf.keys...)) + // size of all keys
		len(leaf.vals) + // size of all value size prefixes
		len(slices.Concat(leaf.vals...)) // size of all values
}

func (leaf *Leaf) Key(i int) ([]byte, error) {
	if i >= len(leaf.keys) || i < 0 {
		return nil, ErrIndexOutOfBounds
	}
	return leaf.keys[i], nil
}

func (leaf *Leaf) Val(i int) ([]byte, error) {
	if i >= len(leaf.vals) || i < 0 {
		return nil, ErrIndexOutOfBounds
	}
	return leaf.vals[i], nil
}

func (leaf *Leaf) KeyVal(i int) ([]byte, []byte, error) {
	if i >= len(leaf.keys) || i >= len(leaf.vals) || i < 0 {
		return nil, nil, ErrIndexOutOfBounds
	}
	return leaf.keys[i], leaf.vals[i], nil
}

func (leaf *Leaf) Search(key []byte) (int, bool) {
	return slices.BinarySearchFunc(leaf.keys, key, bytes.Compare)
}

func (leaf *Leaf) Insert(i int, key, val []byte) error {
	if i > len(leaf.keys) || i > len(leaf.vals) || i < 0 {
		return ErrIndexOutOfBounds
	}
	leaf.keys = slices.Insert(leaf.keys, i, key)
	leaf.vals = slices.Insert(leaf.vals, i, val)
	return nil
}

func (leaf *Leaf) Update(i int, val []byte) error {
	if i >= len(leaf.vals) || i < 0 {
		return ErrIndexOutOfBounds
	}
	leaf.vals[i] = val
	return nil
}

func (leaf *Leaf) Delete(i int) error {
	if i >= len(leaf.keys) || i >= len(leaf.vals) || i < 0 {
		return ErrIndexOutOfBounds
	}
	leaf.keys = slices.Delete(leaf.keys, i, i)
	leaf.vals = slices.Delete(leaf.vals, i, i)
	return nil
}

func (leaf *Leaf) Split() (Leaf, Leaf, error) {
	if len(leaf.keys) <= 1 || len(leaf.vals) <= 1 {
		return Leaf{}, Leaf{}, ErrCannotSplit
	}
	left := Leaf{
		keys: leaf.keys[:len(leaf.keys)/2],
		vals: leaf.vals[:len(leaf.vals)/2],
	}
	right := Leaf{
		keys: leaf.keys[len(leaf.keys)/2:],
		vals: leaf.vals[len(leaf.vals)/2:],
	}

	return left, right, nil
}

func (leaf *Leaf) Append(toAdd Leaf) error {
	if bytes.Compare(leaf.keys[0], toAdd.keys[0]) >= 0 {
		return ErrCannotAppend
	}
	leaf.keys = append(leaf.keys, toAdd.keys...)
	leaf.vals = append(leaf.vals, toAdd.vals...)
	return nil
}

func (leaf *Leaf) Bytes() ([]byte, error) {
	b := make([]byte, leaf.Size())

	binary.LittleEndian.PutUint16(b[:], uint16(TypeLeaf))
	binary.LittleEndian.PutUint16(b[2:], uint16(len(leaf.keys)))
	off := 4
	for i, key := range leaf.keys {
		val := leaf.vals[i]

		binary.LittleEndian.PutUint16(b[off:], uint16(len(key)))
		off += 2
		copy(b[off:], key)
		off += len(key)

		binary.LittleEndian.PutUint16(b[off:], uint16(len(val)))
		off += 2
		copy(b[off:], val)
		off += len(val)

	}
	return b, nil
}
