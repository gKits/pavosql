package node

import (
	"bytes"
	"encoding/binary"
	"slices"
)

type Pointer struct {
	keys [][]byte
	ptrs []uint64
}

func NewPointer(d []byte) (*Pointer, error) {
	pointer := new(Pointer)
	if Type(binary.LittleEndian.Uint16(d[:])) != TypePointer {
		return nil, ErrWrongType
	}
	length := binary.LittleEndian.Uint16(d[2:])
	pointer.keys = make([][]byte, length)
	pointer.ptrs = make([]uint64, length)
	off := 4
	for i := range length {
		if off >= len(d) {
			return nil, ErrNodeDataMalformed
		}

		kLen := int(binary.LittleEndian.Uint16(d[off:]))
		off += 2

		pointer.keys[i] = d[off : off+kLen]
		off += kLen

		pointer.ptrs[i] = binary.LittleEndian.Uint64(d[off:])
		off += 8
	}
	return pointer, nil
}

func (pointer *Pointer) Type() Type {
	return TypePointer
}

func (pointer *Pointer) Len() int {
	return len(pointer.keys)
}

func (pointer *Pointer) Size() int {
	return 2 + // size of type identifier
		2 + // size of number of elements
		2*len(pointer.keys) + // size of all key size prefixes
		len(slices.Concat(pointer.keys...)) + // size of all keys
		len(pointer.ptrs)*8 // size of all pointers
}

func (pointer *Pointer) Key(i int) ([]byte, error) {
	if i >= len(pointer.keys) || i < 0 {
		return nil, ErrIndexOutOfBounds
	}
	return pointer.keys[i], nil
}

func (pointer *Pointer) Ptr(i int) (uint64, error) {
	if i >= len(pointer.ptrs) || i < 0 {
		return 0, ErrIndexOutOfBounds
	}
	return pointer.ptrs[i], nil
}

func (pointer *Pointer) KeyPtr(i int) ([]byte, uint64, error) {
	if i >= len(pointer.keys) || i >= len(pointer.ptrs) || i < 0 {
		return nil, 0, ErrIndexOutOfBounds
	}
	return pointer.keys[i], pointer.ptrs[i], nil
}

func (pointer *Pointer) Search(key []byte) (int, bool) {
	return slices.BinarySearchFunc(pointer.keys, key, bytes.Compare)
}

func (pointer *Pointer) Insert(i int, key []byte, ptr uint64) error {
	if i > len(pointer.keys) || i > len(pointer.ptrs) || i < 0 {
		return ErrIndexOutOfBounds
	}
	pointer.keys = slices.Insert(pointer.keys, i, key)
	pointer.ptrs = slices.Insert(pointer.ptrs, i, ptr)
	return nil
}

func (pointer *Pointer) Update(i int, ptr uint64) error {
	if i >= len(pointer.ptrs) || i < 0 {
		return ErrIndexOutOfBounds
	}
	pointer.ptrs[i] = ptr
	return nil
}

func (pointer *Pointer) Delete(i int) error {
	if i >= len(pointer.keys) || i >= len(pointer.ptrs) || i < 0 {
		return ErrIndexOutOfBounds
	}
	pointer.keys = slices.Delete(pointer.keys, i, i)
	pointer.ptrs = slices.Delete(pointer.ptrs, i, i)
	return nil
}

func (pointer *Pointer) Split() (Pointer, Pointer, error) {
	if len(pointer.keys) <= 1 || len(pointer.ptrs) <= 1 {
		return Pointer{}, Pointer{}, ErrCannotSplit
	}
	left := Pointer{
		keys: pointer.keys[:len(pointer.keys)/2],
		ptrs: pointer.ptrs[:len(pointer.ptrs)/2],
	}
	right := Pointer{
		keys: pointer.keys[len(pointer.keys)/2:],
		ptrs: pointer.ptrs[len(pointer.ptrs)/2:],
	}

	return left, right, nil
}

func (pointer *Pointer) Append(toAdd Pointer) error {
	if bytes.Compare(pointer.keys[0], toAdd.keys[0]) >= 0 {
		return ErrCannotAppend
	}
	pointer.keys = append(pointer.keys, toAdd.keys...)
	pointer.ptrs = append(pointer.ptrs, toAdd.ptrs...)
	return nil
}

func (pointer *Pointer) Bytes() ([]byte, error) {
	b := make([]byte, pointer.Size())

	binary.LittleEndian.PutUint16(b[:], uint16(TypePointer))
	binary.LittleEndian.PutUint16(b[2:], uint16(len(pointer.keys)))
	off := 4
	for i, key := range pointer.keys {
		ptr := pointer.ptrs[i]

		binary.LittleEndian.PutUint16(b[off:], uint16(len(key)))
		off += 2
		copy(b[off:], key)
		off += len(key)

		binary.LittleEndian.PutUint64(b[off:], ptr)
		off += 8

	}
	return b, nil
}
