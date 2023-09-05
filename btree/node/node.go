package node

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"

	"github.com/gKits/PavoSQL/btree/cell"
	"github.com/gKits/PavoSQL/btree/utils"
)

const (
	NODE_HEAD = 3
)

type Node []byte

func NewNode(t utils.Type, cells ...cell.Cell) Node {
	n := make([]byte, 3+2*len(cells))

	n[0] = byte(t)
	binary.BigEndian.PutUint16(n[1:3], uint16(len(cells)))

	for i, c := range cells {
		if c.Type() != t {
			panic("node: non matching cell and node type")
		}

		binary.BigEndian.PutUint16(n[3+2*i:5+2*i], uint16(len(n)))
		n = append(n, c...)
	}

	return n
}

func (n Node) Type() utils.Type {
	return utils.Type(n[0])
}

func (n Node) NCells() uint16 {
	return binary.BigEndian.Uint16(n[1:3])
}

func (n Node) setNCells(nCells uint16) {
	binary.BigEndian.PutUint16(n[1:3], nCells)
}

func (n Node) Size() uint16 {
	return uint16(len(n))
}

func (n Node) getCellPtr(i uint16) uint16 {
	return binary.BigEndian.Uint16(n[NODE_HEAD+i*2 : NODE_HEAD+2+i*2])
}

func (n Node) setCellPtr(i, ptr uint16) {
	binary.BigEndian.PutUint16(n[NODE_HEAD+i*2:NODE_HEAD+2+i*2], ptr)
}

func (n Node) GetCell(i uint16) cell.Cell {
	if i >= n.NCells() {
		panic("page: index out of n-cell range")
	}
	kSize := binary.BigEndian.Uint16(n[n.getCellPtr(i):])
	vSize := binary.BigEndian.Uint16(n[n.getCellPtr(i)+2:])

	return cell.Cell(n[n.getCellPtr(i) : n.getCellPtr(i)+4+kSize+vSize])
}

func (n Node) InsertCell(i uint16, c cell.Cell) Node {
	if i > n.NCells() {
		panic("page: index out of n-cell range")
	}

	inserted := Node{}
	copy(inserted, n)

	inserted = slices.Insert(inserted, int(n.getCellPtr(i)), c...)

	inserted.setNCells(n.NCells() + 1)

	for i++; i < inserted.NCells(); i++ {
		cur := inserted.getCellPtr(i)
		inserted.setCellPtr(i, cur+c.Size())
	}

	return inserted
}

func (n Node) DeleteCell(i uint16) Node {
	if i >= n.NCells() {
		panic("page: index out of n-cell range")
	}

	c := n.GetCell(i)

	deleted := Node{}
	copy(deleted, n)

	deleted = slices.Delete(
		deleted,
		int(n.getCellPtr(i)),
		int(n.getCellPtr(i)+c.Size()),
	)

	deleted.setNCells(n.NCells() - 1)

	for i++; i < deleted.NCells(); i++ {
		cur := deleted.getCellPtr(i)
		deleted.setCellPtr(i, cur-c.Size())
	}

	return deleted
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
