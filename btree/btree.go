package btree

import (
	"errors"

	"github.com/gKits/PavoSQL/btree/cell"
	"github.com/gKits/PavoSQL/btree/node"
	"github.com/gKits/PavoSQL/btree/utils"
)

const (
	PAGE_SIZE = 4096
)

type BTree struct {
	root  uint64
	get   func(uint64) node.Node
	alloc func(node.Node) uint64
	free  func(uint64)
}

func NewBTree(
	root uint64,
	get func(uint64) node.Node,
	alloc func(node.Node) uint64,
	free func(uint64),
) BTree {
	return BTree{root, get, alloc, free}
}

// Inserts the cell c into the BTree.
func (bt *BTree) Insert(c cell.Cell) error {
	if bt.root == 0 {
		root, _ := node.NewNode(utils.LEAF, c)
		bt.root = bt.alloc(root)
		return nil
	}

	root := bt.get(bt.root)

	inserted, err := bt.bTreeInsert(root, c)
	if err != nil {
		return err
	}

	if inserted.Size() > PAGE_SIZE {
		insertedPtr := bt.alloc(inserted)
		c, err := inserted.GetCell(0)
		if err != nil {
			return err
		}
		c, _ = cell.NewCell(utils.INTERN, c.Key(), insertedPtr)

		root, _ = node.NewNode(utils.INTERN, c)

		inserted, err = bt.splitChildPtr(0, root, inserted)
		if err != nil {
			return err
		}
		bt.root = bt.alloc(inserted)
	} else {
		bt.root = bt.alloc(inserted)
	}

	return nil
}

// Deletes the cell the key k from the BTree.
func (bt *BTree) Delete(k []byte) error {
	if bt.root == 0 {
		return errors.New("cannot delete from empty tree")
	}
	root := bt.get(bt.root)

	deleted, err := bt.bTreeDelete(root, k)
	if err != nil {
		return err
	}

	bt.free(bt.root)
	if deleted.Type() == utils.INTERN && deleted.NCells() == 1 {
		bt.root, _ = deleted.GetInternalCell(0)
	} else {
		bt.root = bt.alloc(deleted)
	}

	return nil
}

func (bt *BTree) bTreeInsert(n node.Node, c cell.Cell) (node.Node, error) {
	i, exists := n.BinSearchKeyIdx(c.Key())

	var err error
	inserted := node.Node{}

	switch n.Type() {
	case utils.LEAF:
		if !exists {
			inserted, err = n.InsertCell(i, c)
			if err != nil {
				return nil, err
			}
		} else {
			inserted, err = n.UpdateCell(i, c)
			if err != nil {
				return nil, err
			}
		}

	case utils.INTERN:
		ptrCell, _ := n.GetCell(i)
		childPtr, _ := ptrCell.GetChildPtr()
		child := bt.get(childPtr)

		inserted, err = bt.bTreeInsert(child, c)
		if err != nil {
			return nil, err
		}

		if inserted.Size() > PAGE_SIZE {
			inserted, err = bt.splitChildPtr(i, n, inserted)
			if err != nil {
				return nil, err
			}
		} else {
			insertedPtr := bt.alloc(inserted)
			inserted, err = n.UpdateInternalCell(i, insertedPtr)
			if err != nil {
				return nil, err
			}
		}
		bt.free(childPtr)

	default:
		return nil, errors.New("btree: page header malformed")

	}

	return inserted, nil
}

func (bt *BTree) bTreeDelete(n node.Node, k []byte) (node.Node, error) {
	i, exists := n.BinSearchKeyIdx(k)

	var err error
	var deleted node.Node

	switch n.Type() {
	case utils.LEAF:
		if !exists {
			return nil, errors.New("btree: cannot delete non existing key")
		}

		deleted, err = n.DeleteCell(i)

	case utils.INTERN:
		ptrCell, _ := n.GetCell(i)
		childPtr, _ := ptrCell.GetChildPtr()
		child := bt.get(childPtr)

		deleted, err = bt.bTreeDelete(child, k)
		if err != nil {
			return nil, err
		}

		if deleted.Size() > PAGE_SIZE/4 {
			deleted, err = bt.mergeChildPtr(i, n, deleted)
			if err != nil {
				return nil, err
			}
		} else {
			deletedPtr := bt.alloc(deleted)
			deleted, err = n.UpdateInternalCell(i, deletedPtr)
			if err != nil {
				return nil, err
			}
		}
		bt.free(childPtr)

	default:
		return nil, errors.New("btree: page header malformed")
	}

	return deleted, nil
}

func (bt *BTree) splitChildPtr(i uint16, parent, child node.Node) (node.Node, error) {
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

func (bt *BTree) mergeChildPtr(i uint16, parent, child node.Node) (node.Node, error) {
	var err error
	var siblingPtr uint64
	var sibling node.Node
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
