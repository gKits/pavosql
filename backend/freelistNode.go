package backend

import (
	"encoding/binary"
	"errors"
)

type freelistNode struct {
	next uint64
	ptrs []uint64
}

func (fn *freelistNode) decode(d []byte) error {
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

func (fn freelistNode) typ() nodeType {
	return flNode
}

func (fn freelistNode) encode() []byte {
	var b []byte

	binary.BigEndian.AppendUint16(b, uint16(flNode))
	binary.BigEndian.AppendUint16(b, uint16(len(fn.ptrs)))
	binary.BigEndian.AppendUint64(b, fn.next)
	for _, ptr := range fn.ptrs {
		binary.BigEndian.AppendUint64(b, ptr)
	}

	return b
}

func (fn freelistNode) total() int {
	return len(fn.ptrs)
}

func (fn freelistNode) size() int {
	return 12 + len(fn.ptrs)*8
}

// Useless methods to implement node interface

// Do not use! This method exists for interface purposes only.
func (fn freelistNode) key(i int) ([]byte, error) {
	return nil, errors.New("")
}

// Do not use! This method exists for interface purposes only.
func (fn freelistNode) search(k []byte) (int, bool) {
	return -1, false
}

// Do not use! This method exists for interface purposes only.
func (fn freelistNode) merge(n node) (node, error) {
	return nil, errors.New("")
}

// Do not use! This method exists for interface purposes only.
func (fn freelistNode) split() (node, node) {
	return nil, nil
}

// freelistNode specific methods

func (fn freelistNode) pop() (uint64, freelistNode) {
	last := fn.ptrs[fn.total()-1]
	fn.ptrs = fn.ptrs[:fn.total()-1]
	return last, fn
}

func (fn freelistNode) push(ptr uint64) freelistNode {
	fn.ptrs = append(fn.ptrs, ptr)
	return fn
}
