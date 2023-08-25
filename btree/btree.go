package btree

import "fmt"

const (
	PAGE_SIZE uint16 = 4096
)

type BTree struct {
	root  uint64
	get   func(ptr uint64) Page
	alloc func(p Page) uint64
	free  func(ptr uint64)
}

// Inserts the cell into the correct leaf by recursively walking down the B-Tree
// and updating the internal cells pointers along the way. If a page exceeds the
// page size limit it and the pointer to it stored in it's parent page will be
// split.
func (bt *BTree) bTreeInsert(p Page, c Cell) (Page, error) {
	i, exists, err := p.BinSearchKeyIdx(c.Key())
	if err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("btree: cannot insert key that already exists")
	}

	inserted := Page{}
	switch p.Type() {

	case LEAF_PAGE:
		inserted, err = p.InsertCell(i, c)
		if err != nil {
			return nil, err
		}

	case INTERNAL_PAGE:
		ptrCell, _ := p.GetCell(i)
		childPtr, _ := ptrCell.GetChildPtr()
		child := bt.get(childPtr)

		inserted, err = bt.bTreeInsert(child, c)
		if err != nil {
			return nil, err
		}

		if inserted.Size() > PAGE_SIZE {
			inserted, err = bt.splitChildPtr(i, p, inserted)
			if err != nil {
				return nil, err
			}
		} else {
			insertedPtr := bt.alloc(inserted)
			inserted, err = p.UpdateInternalCell(i, insertedPtr)
			if err != nil {
				return nil, err
			}
		}
		bt.free(childPtr)

	default:

	}

	return inserted, nil
}

// Updates the cell on the correct leaf by recursively walking down the B-Tree
// and updating the internal cells pointers along the way.
func (bt *BTree) bTreeUpdate(p Page, c Cell) (Page, error) {
	i, exists, err := p.BinSearchKeyIdx(c.Key())
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("btree: cannot update key that does not exist")
	}

	updated := Page{}
	switch p.Type() {

	case LEAF_PAGE:
		updated, err = p.UpdateCell(i, c)
		if err != nil {
			return nil, err
		}

	case INTERNAL_PAGE:
		ptrCell, _ := p.GetCell(i)
		childPtr, _ := ptrCell.GetChildPtr()
		child := bt.get(childPtr)

		updated, err = bt.bTreeUpdate(child, c)
		if err != nil {
			return nil, err
		}

		updatedPtr := bt.alloc(updated)

		updated, err = p.UpdateInternalCell(i, updatedPtr)
		if err != nil {
			return nil, err
		}
		bt.free(childPtr)

	default:

	}

	return updated, nil
}

func (bt BTree) splitChildPtr(i uint16, parent, child Page) (Page, error) {
	l, r := child.Split()
	lPtr := bt.alloc(l)
	rPtr := bt.alloc(r)

	split, err := parent.UpdateInternalCell(i, lPtr)
	if err != nil {
		return nil, err
	}

	rCell, err := r.GetCell(0)
	if err != nil {
		return nil, err
	}
	rCell, err = rCell.SetChildPtr(rPtr)
	if err != nil {
		return nil, err
	}

	split, err = parent.InsertCell(i, rCell)
	if err != nil {
		return nil, err
	}

	return split, nil
}
