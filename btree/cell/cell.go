package cell

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/gKits/PavoSQL/btree/utils"
)

type Cell []byte

func NewCell(t utils.Type, k []byte, v any) (Cell, error) {
	c := Cell{byte(t)}
	binary.BigEndian.AppendUint16(c, uint16(len(k)))

	switch t {
	case utils.INTERN:
		ptr, ok := v.(uint64)
		if !ok {
			return nil, errors.New("cell: internal cell requires uint64")
		}

		binary.BigEndian.AppendUint16(c, 8)
		c = append(c, k...)
		binary.BigEndian.AppendUint64(c, ptr)
		break

	case utils.LEAF:
		val, ok := v.([]byte)
		if !ok {
			return nil, errors.New("cell: leaf cell requires bytes")
		}

		binary.BigEndian.AppendUint16(c, uint16(len(val)))
		c = append(c, k...)
		c = append(c, val...)
		break

	default:
		return nil, errors.New("cell: unknown cell type")

	}
	return c, nil
}

func (c Cell) Type() utils.Type {
	return utils.Type(c[0])
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

// Returns a new cell with the value set to v. Returns an error if the size of v
// doesn't equal vSize of original cell c.
func (c Cell) SetVal(v []byte) (Cell, error) {
	if uint16(len(v)) != c.vSize() {
		return nil, errors.New("cell: value has wrong size")
	}
	updated := Cell{}
	copy(updated, c)
	copy(updated[5+c.kSize():5+c.kSize()+c.vSize()], v)

	return updated, nil
}

// Returns a new internal cell with a new child ptr. Returns an error if the
// type of the cell is not internal.
func (c Cell) SetChildPtr(ptr uint64) (Cell, error) {
	if c.Type() != utils.INTERN {
		return nil, errors.New("cell: cannot store ptr in non internal cell")
	}

	updated := Cell{}
	copy(updated, c)
	binary.BigEndian.PutUint64(updated[5+updated.kSize():], ptr)

	return updated, nil
}

// Returns the pointer to the child page stored in the internal cell. Returns an
// error if the type of the cell is not internal.
func (c Cell) GetChildPtr() (uint64, error) {
	if c.Type() != utils.INTERN {
		return 0, errors.New("cell: cannot get ptr from non internal cell")
	}
	return binary.BigEndian.Uint64(c[5+c.kSize() : 5+c.kSize()+c.vSize()]), nil
}
