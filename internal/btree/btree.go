package btree

import (
	"errors"

	"github.com/gKits/PavoSQL/internal/node"
)

const (
	PAGE_SIZE = 4096
)

type BTree struct {
	Root  uint64
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

func (bt *BTree) SetCallbacks(
	get func(uint64) node.Node, alloc func(node.Node) uint64, free func(uint64),
) {
	bt.get = get
	bt.alloc = alloc
	bt.free = free
}

func (bt *BTree) Insert(k, v []byte) error {
	if bt.Root == 0 {
		root := &node.LeafNode{}
		root, err := root.Insert(0, k, v)
		if err != nil {
			return err
		}

		bt.Root = bt.alloc(root)
		return nil
	}

	root := bt.get(bt.Root)

	inserted, err := bt.bTreeInsert(root, k, v)
	if err != nil {
		return err
	}

	if inserted.Size() > PAGE_SIZE {
		insertedPtr := bt.alloc(inserted)
		k, err := inserted.Key(0)
		if err != nil {
			return err
		}

		t := node.PointerNode{}
		t.Insert(0, k, insertedPtr)

		inserted, err = bt.splitChildPtr(0, t, inserted)
		if err != nil {
			return err
		}
		bt.Root = bt.alloc(inserted)
	} else {
		bt.Root = bt.alloc(inserted)
	}

	return nil
}

func (bt *BTree) Delete(k []byte) error {
	if bt.Root == 0 {
		return errors.New("cannot delete from empty tree")
	}
	root := bt.get(bt.Root)

	root, err := bt.bTreeDelete(root, k)
	if err != nil {
		return err
	}

	bt.free(bt.Root)

	if root.Type() == node.PNTR_NODE && root.NKeys() == 1 {
		pntrRoot := root.(*node.PointerNode)
		bt.Root, _ = pntrRoot.Ptr(0)
	} else {
		bt.Root = bt.alloc(root)
	}

	return nil
}

func (bt *BTree) bTreeInsert(n node.Node, k, v []byte) (node.Node, error) {
	i, exists := n.Search(k)

	var err error
	var inserted node.Node

	switch n.Type() {
	case node.LEAF_NODE:
		leafN := n.(*node.LeafNode)

		if !exists {
			inserted, err = leafN.Insert(i, k, v)
			if err != nil {
				return nil, err
			}
		} else {
			inserted, err = leafN.Update(i, k, v)
			if err != nil {
				return nil, err
			}
		}

	case node.PNTR_NODE:
		pntrN := n.(*node.PointerNode)

		ptr, _ := pntrN.Ptr(i)
		child := bt.get(ptr)

		inserted, err = bt.bTreeInsert(child, k, v)
		if err != nil {
			return nil, err
		}

		if inserted.Size() > PAGE_SIZE {
			inserted, err = bt.splitChildPtr(i, *pntrN, inserted)
			if err != nil {
				return nil, err
			}
		} else {
			insertedPtr := bt.alloc(inserted)
			ptrKey, _ := pntrN.Key(i)
			inserted, err = pntrN.Update(i, ptrKey, insertedPtr)
			if err != nil {
				return nil, err
			}
		}

		bt.free(ptr)

	default:
		return nil, errBTreeHeader

	}

	return inserted, nil
}

func (bt *BTree) bTreeDelete(n node.Node, k []byte) (node.Node, error) {
	i, exists := n.Search(k)

	var err error
	var deleted node.Node

	switch n.Type() {
	case node.LEAF_NODE:
		if !exists {
			return nil, errors.New("btree: cannot delete non existing key")
		}

		leafN := n.(*node.LeafNode)

		deleted, err = leafN.Delete(i)

	case node.PNTR_NODE:
		n := n.(*node.PointerNode)

		ptr, _ := n.Ptr(i)
		child := bt.get(ptr)

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

			k, _ := deleted.Key(i)
			deleted, err = n.Update(i, k, deletedPtr)
			if err != nil {
				return nil, err
			}
		}
		bt.free(ptr)

	default:
		return nil, errBTreeHeader
	}

	return deleted, nil
}

func (bt *BTree) splitChildPtr(i int, parent node.PointerNode, child node.Node) (node.Node, error) {
	l, r := child.Split()
	lPtr := bt.alloc(l)
	rPtr := bt.alloc(r)

	lKey, err := parent.Key(i)
	if err != nil {
		return nil, err
	}

	split, err := parent.Update(i, lKey, lPtr)
	if err != nil {
		return nil, err
	}

	rKey, err := r.Key(0)
	if err != nil {
		return nil, err
	}

	split, err = parent.Insert(i+1, rKey, rPtr)
	if err != nil {
		return nil, err
	}

	return split, nil
}

func (bt *BTree) mergeChildPtr(i int, parent *node.PointerNode, child node.Node) (node.Node, error) {
	var err error
	var siblingPtr uint64
	var sibling node.Node
	var merged bool = false

	if i < parent.NKeys()-1 {
		siblingPtr, err = parent.Ptr(i + 1)
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
		i--
		siblingPtr, err = parent.Ptr(i)
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
		k, _ := parent.Key(i)
		parent, err = parent.Update(i, k, ptr)
		if err != nil {
			return nil, err
		}
	}

	return parent, nil
}
