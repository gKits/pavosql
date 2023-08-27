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
// split. Existing keys will be updated.
func (bt *BTree) bTreeInsert(p Page, c Cell) (Page, error) {
	i, exists, err := p.BinSearchKeyIdx(c.Key())
	if err != nil {
		return nil, err
	}

	inserted := Page{}
	switch p.Type() {

	case LEAF_PAGE:
		if !exists {
			inserted, err = p.InsertCell(i, c)
			if err != nil {
				return nil, err
			}
		} else {
			inserted, err = p.UpdateCell(i, c)
			if err != nil {
				return nil, err
			}
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
		return nil, fmt.Errorf("btree: page header malformed")

	}

	return inserted, nil
}

func (bt *BTree) bTreeDelete(p Page, k []byte) (Page, error) {
	i, exists, err := p.BinSearchKeyIdx(k)
	if err != nil {
		return nil, err
	}

	deleted := Page{}

	switch p.Type() {
	case LEAF_PAGE:
		if !exists {
			return nil, fmt.Errorf("btree: cannot delete non existing key")
		}

		deleted, err = p.DeleteCell(i)

	case INTERNAL_PAGE:
		ptrCell, _ := p.GetCell(i)
		childPtr, _ := ptrCell.GetChildPtr()
		child := bt.get(childPtr)

		deleted, err = bt.bTreeDelete(child, k)
		if err != nil {
			return nil, err
		}

		if deleted.Size() > PAGE_SIZE/4 {
			deleted, err = bt.mergeChildPtr(i, p, deleted)
			if err != nil {
				return nil, err
			}
		} else {
			deletedPtr := bt.alloc(deleted)
			deleted, err = p.UpdateInternalCell(i, deletedPtr)
			if err != nil {
				return nil, err
			}
		}
		bt.free(childPtr)

	default:
		return nil, fmt.Errorf("btree: page header malformed")
	}

	return deleted, nil
}

// Splits the child page and the pointer to it stored in the parent page into
// two and returns a new parent page with the two pointers pointing to each of
// the halfs.
func (bt *BTree) splitChildPtr(i uint16, parent, child Page) (Page, error) {
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

// Merges the child page, if possible, to it's neighbors and returns the new
// parent page with the pointers replaced by the new one to the merged page.
func (bt *BTree) mergeChildPtr(i uint16, parent, child Page) (Page, error) {
	var err error
	var siblingPtr uint64
	var sibling Page
	var merged bool = false

	if i < parent.NCells()-1 {
		siblingPtr, err = parent.GetInternalCell(i + 1)
		if err != nil {
			return nil, err
		}
		sibling = bt.get(siblingPtr)
		if sibling.Size()+child.Size() < PAGE_SIZE {
			child, err = child.Merge(sibling)
			if err != nil {
				return nil, err
			}
		}
		bt.free(siblingPtr)
		merged = true
	}

	if i != 0 {
		siblingPtr, err = parent.GetInternalCell(i - 1)
		if err != nil {
			return nil, err
		}
		sibling = bt.get(siblingPtr)
		if sibling.Size()+child.Size() < PAGE_SIZE {
			child, err = sibling.Merge(child)
			if err != nil {
				return nil, err
			}
		}
		bt.free(siblingPtr)
		merged = true
	}

	if merged {
		ptr := bt.alloc(child)
		parent, err = parent.UpdateInternalCell(i, ptr)
		if err != nil {
			return nil, err
		}
	}

	return parent, nil
}
