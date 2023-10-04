package store

import (
	"encoding/binary"
)

type freelistNode struct {
	next uint64
	ptrs []uint64
}

func (fn *freelistNode) Decode(d []byte) error {
	if nodeType(binary.BigEndian.Uint16(d[0:2])) != flNode {
		return errNodeDecode
	}

	nPtrs := binary.BigEndian.Uint16(d[2:4])
	fn.next = binary.BigEndian.Uint64(d[4:12])
	for i := uint16(0); i < nPtrs; i++ {
		fn.ptrs = append(fn.ptrs, binary.BigEndian.Uint64(d[12+i*8:]))
	}

	return nil
}

// node interface methods

func (fn freelistNode) Type() nodeType {
	return flNode
}

func (fn freelistNode) Encode() []byte {
	var b []byte

	binary.BigEndian.AppendUint16(b, uint16(flNode))
	binary.BigEndian.AppendUint16(b, uint16(len(fn.ptrs)))
	binary.BigEndian.AppendUint64(b, fn.next)
	for _, ptr := range fn.ptrs {
		binary.BigEndian.AppendUint64(b, ptr)
	}

	return b
}

func (fn freelistNode) Total() int {
	return len(fn.ptrs)
}

func (fn freelistNode) Size() int {
	return 12 + len(fn.ptrs)*8
}

// Useless methods to implement node interface

// Do not use! This method exists for interface purposes only.
func (fn freelistNode) Key(i int) ([]byte, error) {
	return nil, errNodeUseless
}

// Do not use! This method exists for interface purposes only.
func (fn freelistNode) Search(k []byte) (int, bool) {
	return -1, false
}

// Do not use! This method exists for interface purposes only.
func (fn freelistNode) Merge(n node) (node, error) {
	return nil, errNodeUseless
}

// Do not use! This method exists for interface purposes only.
func (fn freelistNode) Split() (node, node) {
	return nil, nil
}

// freelistNode specific methods

func (fn freelistNode) Pop() (uint64, freelistNode) {
	last := fn.ptrs[fn.Total()-1]
	fn.ptrs = fn.ptrs[:fn.Total()-1]
	return last, fn
}

func (fn freelistNode) Push(ptr uint64) freelistNode {
	fn.ptrs = append(fn.ptrs, ptr)
	return fn
}
