package btree

import "fmt"

const (
	PAGE_SIZE = 4096
)

type BTree struct {
	root uint64
}

// Returns the page at the given pointer.
func (bt *BTree) GetPage(ptr uint64) Page { return nil }

// Recursively walks through the BTree and inserts the cell into the correct page and returns it.
// Returns an error if the page is malformed or the cell does not fit the pages requirements like
// key uniquenes, key placement or cell size.
func (bt *BTree) Insert(p Page, c Cell) (newPage Page, err error) {
	i, exists, err := p.BinSearchKeyIdx(c.Key())
	if err != nil {
		return nil, fmt.Errorf("btree: could not insert cell into page: %e", err)
	}
	if exists {
		return nil, fmt.Errorf("btree: could not insert cell into page, key must be unique")
	}

	switch p.Type() {
	case INTERNAL_PAGE:
		internalC, err := p.GetCell(i)
		if err != nil {
			return nil, fmt.Errorf("btree: could not get cell: %e", err)
		}

		childPtr, err := internalC.GetChildPtr()
		if err != nil {
			return nil, fmt.Errorf("btree: could not get child pointer from cell value: %e", err)
		}

		newPage, err = bt.Insert(bt.GetPage(childPtr), c)
		if err != nil {
			return nil, err
		}

	case LEAF_PAGE:
		newPage, err = p.InsertCell(i, c)
		if err != nil {
			return nil, fmt.Errorf("btree: could not insert new cell into page: %e", err)
		}

	default:
		return nil, fmt.Errorf("btree: bad page type")
	}

	return newPage, nil
}
