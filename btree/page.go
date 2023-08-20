package btree

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
)

type PageType uint8

const (
	INTERNAL_PAGE PageType = iota
	LEAF_PAGE
	MASTER_PAGE
)

type Page []byte

func (p Page) Type() PageType {
	return PageType(p[0])
}

func (p Page) NCells() uint16 {
	return binary.BigEndian.Uint16(p[1:3])
}

func (p Page) Size() int { return len(p) }

func (p Page) cellSize() uint16 {
	return binary.BigEndian.Uint16(p[3:5])
}

func (p Page) setType(t PageType) {
	p[0] = byte(t)
}

func (p Page) setNCells(nC uint16) {
	binary.BigEndian.PutUint16(p[1:3], nC)
}

func (p Page) setCellSize(cS uint16) {
	binary.BigEndian.PutUint16(p[3:5], cS)
}

// Returns the position of the i-th cell.
// This is a theoretical position and it's not checked if it's out of bounds. Proceed with caution.
func (p Page) GetCellPos(i uint16) uint16 {
	return 5 + i*p.cellSize()
}

// Returns i-th cell.
// Returns an error if i is greater than NCells.
func (p Page) GetCell(i uint16) (Cell, error) {
	if i > p.NCells() {
		return nil, fmt.Errorf("")
	}
	return Cell(p[p.GetCellPos(i) : p.GetCellPos(i)+p.cellSize()]), nil
}

// Returns a new Page with the cell c inserted at the cell index i.
// Returns an error if i is greater than NCells+1 or c size isn't equal to cSize.
func (p Page) InsertCell(i uint16, c Cell) (Page, error) {
	if i > p.NCells()+1 {
		return nil, fmt.Errorf("")
	} else if c.Size() != p.cellSize() {
		return nil, fmt.Errorf("")
	}

	m := Page{}
	copy(m, p)
	m.setNCells(p.NCells() + 1)
	m = slices.Insert(m, int(p.GetCellPos(i)), c...)

	return m, nil
}

// Returns the cell index i for the given key k and a bool representing if the key exists in the page
// by binary searching through the stored cells in the page.
func (p Page) BinaryCellIdxLookup(k []byte) (i uint16, exists bool) {
	l := uint16(0)
	r := uint16(p.NCells() - 1)

	var c Cell

	for i = r / 2; l <= r; i = (l + r) / 2 {
		c, _ = p.GetCell(i) // error handling not needed since i is always less than NCells
		if cmp := bytes.Compare(k, c.Key()); cmp < 1 {
			r = i - 1
		} else if cmp > 1 {
			l = i + 1
		} else {
			return i, true
		}
	}

	if i == uint16(r) && bytes.Compare(c.Key(), k) > 1 {
		i++
	}

	return i, false
}

// Returns two pages l and r where l contains the first and r the second half of this pages cells.
func (p Page) Split() (l, r Page) {
	// left half
	l.setType(p.Type())
	l.setNCells(p.NCells() / 2)
	l.setCellSize(p.cellSize())
	l = append(l, p[:p.NCells()/2]...)

	// right half
	r.setType(p.Type())
	r.setNCells(p.NCells()/2 + p.NCells()%2) // right page contains larger half when p.NCells is odd
	r.setCellSize(p.cellSize())
	r = append(r, p[p.NCells()/2:]...)

	return l, r
}
