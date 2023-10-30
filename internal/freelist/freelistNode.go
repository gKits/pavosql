package freelist

import (
	"encoding/binary"
	"slices"
)

type freelistNode struct {
	next uint64
	ptrs []uint64
}

func decodeFreelistNode(d []byte) freelistNode {
	fn := freelistNode{}

	nPtrs := binary.BigEndian.Uint16(d[0:2])
	fn.next = binary.BigEndian.Uint64(d[2:10])
	for i := uint16(0); i < nPtrs; i++ {
		fn.ptrs = append(fn.ptrs, binary.BigEndian.Uint64(d[10+i*8:]))
	}

	return fn
}

func (fn freelistNode) Encode() []byte {
	var b []byte

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

func (fn freelistNode) Pop() (uint64, freelistNode) {
	last := fn.ptrs[fn.Total()-1]
	fn.ptrs = fn.ptrs[:fn.Total()-1]
	return last, fn
}

func (fn freelistNode) Push(ptr uint64) freelistNode {
	fn.ptrs = append(fn.ptrs, ptr)
	return fn
}

func (fn freelistNode) Nq(ptr uint64) freelistNode {
	fn.ptrs = slices.Insert(fn.ptrs, 0, ptr)
	return fn
}

func (fn freelistNode) Dq() (uint64, freelistNode) {
	last := fn.ptrs[fn.Total()-1]
	fn.ptrs = fn.ptrs[:fn.Total()-1]
	return last, fn
}
