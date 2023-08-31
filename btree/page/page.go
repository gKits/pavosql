package page

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gKits/PavoSQL/btree/cell"
	"slices"
)

type PageType uint8

const (
	INTERNAL_PAGE PageType = iota
	LEAF_PAGE
	MASTER_PAGE
)

type Page []byte

func NewInternalPage(c cell.Cell) Page {
	p := Page{byte(INTERNAL_PAGE)}
	binary.BigEndian.AppendUint16(p, 1)
	p = append(p, c...)
	return p
}

func (p Page) Type() PageType {
	return PageType(p[0])
}

func (p Page) NCells() uint16 {
	return binary.BigEndian.Uint16(p[1:3])
}

func (p Page) Size() uint16 {
	return 5 + p.NCells()*p.cellSize()
}

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

// Returns the position of the i-th cell. This is a theoretical position and
// it's not checked if i exceeds nCells. Proceed with caution.
func (p Page) GetCellPos(i uint16) uint16 {
	return 5 + i*p.cellSize()
}

// Returns i-th cell. Returns an error if i is greater equal than NCells.
func (p Page) GetCell(i uint16) (cell.Cell, error) {
	if i >= p.NCells() {
		return nil, fmt.Errorf("page: index out of n-cell range")
	}
	return cell.Cell(p[p.GetCellPos(i) : p.GetCellPos(i)+p.cellSize()]), nil
}

// Returns a new Page with the cell c inserted after the index i. Returns an
// error if i is greater equal than NCells or c's size isn't equal to cellSize.
func (p Page) InsertCell(i uint16, c cell.Cell) (Page, error) {
	if i >= p.NCells() {
		return nil, fmt.Errorf("page: index out of n-cell range")
	} else if c.Size() != p.cellSize() {
		return nil, fmt.Errorf("page: cell has wrong size for page")
	}

	inserted := Page{}
	copy(inserted, p)
	inserted.setNCells(p.NCells() + 1)
	inserted = slices.Insert(inserted, int(p.GetCellPos(i))+1, c...)

	return inserted, nil
}

// Returns a new Page without the cell at index i. Returns an error if i is out
// of range.
func (p Page) DeleteCell(i uint16) (Page, error) {
	if i >= p.NCells() {
		return nil, fmt.Errorf("page: index out of n-cell range")
	}

	deleted := Page{}
	copy(deleted, p)
	deleted.setNCells(p.NCells() - 1)
	deleted = slices.Delete(
		deleted,
		int(p.GetCellPos(i)),
		int(p.GetCellPos(i)+p.cellSize()),
	)

	return deleted, nil
}

// Returns a new Page with the cell at index i updated to the value of cell c.
// Returns an error if the page does not have a cell at i or c's size isn't
// equal to cSize.
func (p Page) UpdateCell(i uint16, c cell.Cell) (Page, error) {
	if c.Size() != p.cellSize() {
		return nil, fmt.Errorf("page: cell has wrong size for page")
	}

	ogC, err := p.GetCell(i)
	if err != nil {
		return nil, err
	} else if !slices.Equal(ogC.Key(), c.Key()) {
		return nil, fmt.Errorf("page: updated key needs to equal original")
	}

	updated := Page{}
	copy(updated, p)
	copy(updated[p.GetCellPos(i):], c)

	return updated, nil
}

// Returns the child pointer stored in the internal cell at i. Returns an error
// if the page type is not internal or the cell does not exist.
func (p Page) GetInternalCell(i uint16) (uint64, error) {
	if p.Type() != INTERNAL_PAGE {
		return 0, fmt.Errorf("page: ptr are only stored on internal pages")
	}
	c, err := p.GetCell(i)
	if err != nil {
		return 0, err
	}
	return c.GetChildPtr()
}

func (p Page) UpdateInternalCell(i uint16, ptr uint64) (Page, error) {
	if p.Type() != INTERNAL_PAGE {
		return nil, fmt.Errorf("page: ptr are only stored on internal pages")
	}

	c, err := p.GetCell(i)
	if err != nil {
		return nil, nil
	}

	c, err = c.SetChildPtr(ptr)
	if err != nil {
		return nil, nil
	}

	updated := Page{}
	copy(updated, p)
	copy(updated[p.GetCellPos(i):], c)

	return updated, nil
}

// Returns the cell index for the given key k and a bool representing if the key
// exists by binary searching over the pages cells.
func (p Page) BinSearchKeyIdx(k []byte) (uint16, bool) {
	l := uint16(0)
	r := uint16(p.NCells() - 1)

	var i uint16
	var cmp int
	for i = r / 2; l < r; i = (l + r) / 2 {
		c, _ := p.GetCell(i)

		cmp = bytes.Compare(k, c.Key())
		if cmp < 0 {
			r = i - 1
		} else if cmp > 0 {
			l = i + 1
		} else {
			return i, true
		}
	}

	if cmp > 1 {
		i++
	}

	return i, false
}

// Returns two pages l and r where l contains the first and r the second half of
// this pages cells.
func (p Page) Split() (l, r Page) {
	// left half
	l.setType(p.Type())
	l.setNCells(p.NCells() / 2)
	l.setCellSize(p.cellSize())
	l = append(l, p[:p.NCells()/2]...)

	// right half
	r.setType(p.Type())
	r.setNCells(p.NCells()/2 + p.NCells()%2)
	r.setCellSize(p.cellSize())
	r = append(r, p[p.NCells()/2:]...)

	return l, r
}

// Merges page p and toMerge into one and returns it. ToMerge will be merged
// onto p (p left, toMerge right). Returns an error if the page types don't
// match or if p's last key >= toMerge's first key to stay sorted.
func (p Page) Merge(toMerge Page) (Page, error) {
	if p.Type() != toMerge.Type() {
		return nil, fmt.Errorf("page: cannot merge pages of different types")
	}

	pLast, err := p.GetCell(p.NCells() - 1)
	if err != nil {
		return nil, err
	}
	tMFirst, err := toMerge.GetCell(0)
	if err != nil {
		return nil, err
	}
	if cmp := bytes.Compare(pLast.Key(), tMFirst.Key()); cmp >= 0 {
		return nil, fmt.Errorf("page: last key >= merged pages first key")
	}

	merged := Page{}
	merged.setType(p.Type())
	merged.setNCells(p.NCells() + toMerge.NCells())
	merged = append(merged, p[5:]...)
	merged = append(merged, toMerge[5:]...)

	return merged, nil
}
