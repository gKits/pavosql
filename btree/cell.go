package btree

import (
	"encoding/binary"
	"fmt"
)

type CellType uint8

const (
	INTERNAL_CELL CellType = iota
	DATA_CELL
)

type Cell []byte

func (c Cell) Type() CellType {
	return CellType(c[0])
}

func (c Cell) Size() uint16 {
	return 5 + c.kSize() + c.vSize()
}

func (c Cell) Key() []byte {
	return c[5 : 5+c.kSize()]
}

func (c Cell) Val() []byte {
	return c[5+c.kSize() : 5+c.kSize()+c.vSize()]
}

func (c Cell) kSize() uint16 {
	return binary.BigEndian.Uint16(c[1:3])
}

func (c Cell) vSize() uint16 {
	return binary.BigEndian.Uint16(c[3:5])
}

// Returns a new cell with the value set to v.
// Returns an error if the size of v doesn't equal vSize of original cell c.
func (c Cell) SetVal(v []byte) (Cell, error) {
	if uint16(len(v)) != c.vSize() {
		return nil, fmt.Errorf("cell: size of set value %d needs to equal vSize %d", len(v), c.vSize())
	}
	updated := Cell{}
	copy(updated, c)
	copy(updated[5+c.kSize():5+c.kSize()+c.vSize()], v)

	return updated, nil
}

// Returns a new internal cell with a new child ptr.
// Returns an error if the type of the cell is not internal.
func (c Cell) SetChildPtr(ptr uint64) (Cell, error) {
	if c.Type() != INTERNAL_CELL {
		return nil, fmt.Errorf("cell: non internal cell cannot contain a child pointer")
	}

	updated := Cell{}
	copy(updated, c)
	binary.BigEndian.PutUint64(updated[5+updated.kSize():], ptr)

	return updated, nil
}

// Returns the pointer to the child page stored in the internal cell.
// Returns an error if the type of the cell is not internal.
func (c Cell) GetChildPtr() (uint64, error) {
	if c.Type() != INTERNAL_CELL {
		return 0, fmt.Errorf("cell: non internal cell cannot contain a child pointer")
	}
	return binary.BigEndian.Uint64(c[5+c.kSize() : 5+c.kSize()+c.vSize()]), nil
}
