package node

import (
	"bytes"
	"encoding/binary"
	"errors"
	"slices"

	"github.com/gKits/PavoSQL/btree/cell"
	"github.com/gKits/PavoSQL/btree/utils"
)

const (
	NODE_HEAD = 3
)

type Node []byte

func NewNode(t utils.Type, cells ...cell.Cell) (Node, error) {
	n := make([]byte, 3+2*len(cells))

	n[0] = byte(t)
	binary.BigEndian.PutUint16(n[1:3], uint16(len(cells)))

	for i, c := range cells {
		if c.Type() != t {
			return nil, errors.New("node: non matching cell and node type")
		}

		binary.BigEndian.PutUint16(n[3+2*i:5+2*i], uint16(len(n)))
		n = append(n, c...)
	}

	return n, nil
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

func (n Node) GetCell(i uint16) (cell.Cell, error) {
	if i >= n.NCells() {
		return nil, errors.New("page: index out of n-cell range")
	}
	kSize := binary.BigEndian.Uint16(n[n.getCellPtr(i):])
	vSize := binary.BigEndian.Uint16(n[n.getCellPtr(i)+2:])

	return cell.Cell(n[n.getCellPtr(i) : n.getCellPtr(i)+4+kSize+vSize]), nil
}

func (n Node) getCells() []cell.Cell {
	cells := []cell.Cell{}
	for i := uint16(0); i < n.NCells(); i++ {
		c, _ := n.GetCell(i)
		cells = append(cells, c)
	}
	return cells
}

func (n Node) InsertCell(i uint16, c cell.Cell) (Node, error) {
	if i > n.NCells() {
		return nil, errors.New("page: index out of n-cell range")
	}

	inserted := Node{}
	copy(inserted, n)

	inserted = slices.Insert(inserted, int(n.getCellPtr(i)), c...)

	inserted.setNCells(n.NCells() + 1)

	for i++; i < inserted.NCells(); i++ {
		cur := inserted.getCellPtr(i)
		inserted.setCellPtr(i, cur+c.Size())
	}

	return inserted, nil
}

func (n Node) DeleteCell(i uint16) (Node, error) {
	c, err := n.GetCell(i)
	if err != nil {
		return nil, err
	}

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

	return deleted, nil
}

func (n Node) UpdateCell(i uint16, c cell.Cell) (Node, error) {
	og, err := n.GetCell(i)
	if err != nil {
		return nil, err
	}

	if !slices.Equal(og.Key(), c.Key()) {
		return nil, errors.New("page: updated key needs to equal original")
	}

	updated := Node{}
	copy(updated, n)
	updated = slices.Replace(updated, int(n.getCellPtr(i)), int(n.getCellPtr(i+1)), c...)

	return updated, nil
}

func (n Node) GetInternalCell(i uint16) (uint64, error) {
	if n.Type() != utils.INTERN {
		return 0, errors.New("page: ptr are only stored on internal pages")
	}
	c, err := n.GetCell(i)
	if err != nil {
		return 0, err
	}

	return c.GetChildPtr()
}

func (n Node) UpdateInternalCell(i uint16, ptr uint64) (Node, error) {
	if n.Type() != utils.INTERN {
		return nil, errors.New("page: ptr are only stored on internal pages")
	}

	c, err := n.GetCell(i)
	if err != nil {
		return nil, err
	}

	c, err = c.SetChildPtr(ptr)
	if err != nil {
		panic(err)
	}

	updated := Node{}
	copy(updated, n)
	updated = slices.Replace(updated, int(n.getCellPtr(i)), int(n.getCellPtr(i+1)), c...)

	return n.UpdateCell(i, c)
}

func (n Node) BinSearchKeyIdx(k []byte) (uint16, bool) {
	l := uint16(0)
	r := uint16(n.NCells() - 1)

	var i uint16
	var cmp int
	for i = r / 2; l < r; i = (l + r) / 2 {
		c, _ := n.GetCell(i)

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

func (n Node) Split() (Node, Node) {
	c := n.getCells()

	l, _ := NewNode(n.Type(), c[:n.NCells()/2]...)
	r, _ := NewNode(n.Type(), c[n.NCells()/2:]...)

	return l, r
}

func (n Node) Merge(toMerge Node) (Node, error) {
	if n.Type() != toMerge.Type() {
		return nil, errors.New("page: cannot merge pages of different types")
	}

	nCells := n.getCells()
	tMCells := n.getCells()

	if cmp := bytes.Compare(nCells[n.NCells()-1].Key(), tMCells[0].Key()); cmp >= 0 {
		return nil, errors.New("page: last key >= merged pages first key")
	}
	nCells = append(nCells, tMCells...)

	merged, _ := NewNode(n.Type(), nCells...)
	return merged, nil
}
