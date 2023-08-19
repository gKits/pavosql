package btree

import "encoding/binary"

type CellType uint8

const (
	INTERNAL_CELL CellType = iota
	DATA_CELL
)

type Cell []byte

func (c Cell) Type() CellType { return CellType(c[0]) }
func (c Cell) Size() uint16   { return 5 + c.kSize() + c.vSize() }
func (c Cell) Key() []byte    { return c[5 : 5+c.kSize()] }
func (c Cell) Val() []byte    { return c[5+c.kSize() : 5+c.kSize()+c.vSize()] }
func (c Cell) kSize() uint16  { return binary.BigEndian.Uint16(c[1:3]) }
func (c Cell) vSize() uint16  { return binary.BigEndian.Uint16(c[3:5]) }
