package freelist

import (
	"encoding/binary"
	"errors"

	"github.com/gKits/PavoSQL/internal/node"
)

type Freelist struct {
	next uint64
	ptrs []uint64
}

func (fl *Freelist) Type() node.NodeType {
	return node.FLST_NODE
}

func (fl *Freelist) Encode() []byte {
	var b []byte

	binary.BigEndian.AppendUint16(b, uint16(node.FLST_NODE))
	binary.BigEndian.AppendUint16(b, uint16(len(fl.ptrs)))
	binary.BigEndian.AppendUint64(b, fl.next)
	for _, ptr := range fl.ptrs {
		binary.BigEndian.AppendUint64(b, ptr)
	}

	return b
}

func (fl *Freelist) Decode(d []byte) error {
	if node.NodeType(binary.BigEndian.Uint16(d[0:2])) != node.FLST_NODE {
		return errFLDecode
	}

	nPtrs := binary.BigEndian.Uint16(d[2:4])
	fl.next = binary.BigEndian.Uint64(d[4:12])
	for i := uint16(0); i < nPtrs; i++ {
		fl.ptrs = append(fl.ptrs, binary.BigEndian.Uint64(d[12+i*8:]))
	}

	return nil
}

func (fl *Freelist) Size() int {
	return 12 + len(fl.ptrs)*8
}

func (fl *Freelist) Pop() uint64 {
	last := fl.ptrs[fl.Total()-1]
	fl.ptrs = fl.ptrs[:fl.Total()-1]
	return last
}

func (fl *Freelist) Push(ptr uint64) {
	fl.ptrs = append(fl.ptrs, ptr)
}

func (fl *Freelist) Total() uint {
	return uint(len(fl.ptrs))
}
