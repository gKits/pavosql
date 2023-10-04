package store

import (
	"errors"
)

const PageSize = 4096

type getFunc func(uint64) (node, error)
type pullFunc func(uint64) (node, error)
type allocFunc func(node) (uint64, error)
type freeFunc func(uint64) error

type bTree struct {
	root  uint64
	get   getFunc
	pull  pullFunc
	alloc allocFunc
	free  freeFunc
}

var (
	errBTreeHeader      = errors.New("btree: cannot decode page, header is malformed")
	errBTreeDeleteEmpty = errors.New("btree: cannot delete key from empty btree")
	errBTreeDeleteKey   = errors.New("btree: cannot delete non existing key")
	errBTreeGetEmpty    = errors.New("btree: cannot read from empty btree")
	errBTreeGetKey      = errors.New("btree: cannot read non existing key")
)

func (bt *bTree) Get(k []byte) ([]byte, error) {
	if bt.root == 0 {
		return nil, errBTreeGetEmpty
	}

	root, err := bt.get(bt.root)
	if err != nil {
		return nil, err
	}

	_, v, err := bt.bTreeGet(root, k)
	return v, err
}

func (bt *bTree) Set(k, v []byte) error {
	if bt.root == 0 {
		root := leafNode{}
		root, err := root.Insert(0, k, v)
		if err != nil {
			return err
		}

		bt.root, err = bt.alloc(root)
		if err != nil {
			return err
		}

		return nil
	}

	root, err := bt.pull(bt.root)
	if err != nil {
		return err
	}

	inserted, err := bt.bTreeInsert(root, k, v)
	if err != nil {
		return err
	}

	if inserted.Size() > PageSize {
		insertedPtr, err := bt.alloc(inserted)
		if err != nil {
			return err
		}

		k, err := inserted.Key(0)
		if err != nil {
			return err
		}

		t := pointerNode{}
		t.Insert(0, k, insertedPtr)

		inserted, err = bt.splitChildPtr(0, t, inserted)
		if err != nil {
			return err
		}

	}

	bt.root, err = bt.alloc(inserted)
	return err
}

func (bt *bTree) Delete(k []byte) (bool, error) {
	if bt.root == 0 {
		return false, errBTreeDeleteEmpty
	}

	root, err := bt.pull(bt.root)
	if err != nil {
		return false, err
	}

	var deleted bool
	root, deleted, err = bt.bTreeDelete(root, k)
	if err != nil {
		return false, err
	}

	if !deleted {
		return deleted, nil
	}

	if root.Type() == ptrNode && root.Total() == 1 {
		pntrRoot := root.(pointerNode)
		bt.root, _ = pntrRoot.Ptr(0)
	} else {
		bt.root, err = bt.alloc(root)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (bt *bTree) bTreeGet(n node, k []byte) (node, []byte, error) {
	i, exists := n.Search(k)

	switch n.Type() {
	case lfNode:
		leafN := n.(leafNode)

		if !exists {
			return nil, nil, errBTreeGetKey
		}

		v, err := leafN.Val(i)
		if err != nil {
			return nil, nil, err
		}
		return n, v, nil

	case ptrNode:
		pntrN := n.(pointerNode)

		ptr, _ := pntrN.Ptr(i)
		child, err := bt.get(ptr)
		if err != nil {
			return nil, nil, err
		}

		return bt.bTreeGet(child, k)
	default:
		return nil, nil, errBTreeHeader
	}
}

func (bt *bTree) bTreeInsert(n node, k, v []byte) (node, error) {
	i, exists := n.Search(k)

	var err error
	var inserted node

	switch n.Type() {
	case lfNode:
		leafN := n.(leafNode)

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

	case ptrNode:
		pntrN := n.(pointerNode)

		ptr, _ := pntrN.Ptr(i)
		child, err := bt.get(ptr)
		if err != nil {
			return nil, err
		}

		inserted, err = bt.bTreeInsert(child, k, v)
		if err != nil {
			return nil, err
		}

		if inserted.Size() > PageSize {
			inserted, err = bt.splitChildPtr(i, pntrN, inserted)
			if err != nil {
				return nil, err
			}
		} else {
			insertedPtr, err := bt.alloc(inserted)
			if err != nil {
				return nil, err
			}

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

func (bt *bTree) bTreeDelete(n node, k []byte) (node, bool, error) {
	i, exists := n.Search(k)

	var err error
	var deleted node

	switch n.Type() {
	case lfNode:
		if !exists {
			return nil, false, errBTreeDeleteKey
		}

		leafN := n.(leafNode)
		deleted, err = leafN.Delete(i)
		if err != nil {
			return nil, false, err
		}

	case ptrNode:
		n := n.(pointerNode)

		ptr, _ := n.Ptr(i)
		child, err := bt.pull(ptr)
		if err != nil {
			return nil, false, err
		}

		deleted, _, err = bt.bTreeDelete(child, k)
		if err != nil {
			return nil, false, err
		}

		if deleted.Size() > PageSize/4 {
			deleted, err = bt.mergeChildPtr(i, n, deleted)
			if err != nil {
				return nil, false, err
			}
		} else {
			deletedPtr, err := bt.alloc(deleted)
			if err != nil {
				return nil, false, err
			}

			k, _ := deleted.Key(i)
			deleted, err = n.Update(i, k, deletedPtr)
			if err != nil {
				return nil, false, err
			}
		}

	default:
		return nil, false, errBTreeHeader
	}

	return deleted, true, nil
}

func (bt *bTree) splitChildPtr(i int, parent pointerNode, child node) (node, error) {
	l, r := child.Split()

	lPtr, err := bt.alloc(l)
	if err != nil {
		return nil, err
	}

	rPtr, err := bt.alloc(r)
	if err != nil {
		return nil, err
	}

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

	split, err = split.Insert(i+1, rKey, rPtr)
	if err != nil {
		return nil, err
	}

	return split, nil
}

func (bt *bTree) mergeChildPtr(i int, parent pointerNode, child node) (node, error) {
	var err error
	var siblingPtr uint64
	var sibling node
	var merged bool = false

	if i < parent.Total()-1 {
		siblingPtr, err = parent.Ptr(i + 1)
		if err != nil {
			return nil, err
		}

		sibling, err = bt.pull(siblingPtr)
		if err != nil {
			return nil, err
		}

		if sibling.Size()+child.Size() < PageSize {
			child, err = child.Merge(sibling)
			if err != nil {
				return nil, err
			}
		}

		merged = true
	}

	if i != 0 {
		i--
		siblingPtr, err = parent.Ptr(i)
		if err != nil {
			return nil, err
		}

		sibling, err = bt.pull(siblingPtr)
		if err != nil {
			return nil, err
		}

		if sibling.Size()+child.Size() < PageSize {
			child, err = sibling.Merge(child)
			if err != nil {
				return nil, err
			}
		}

		merged = true
	}

	if merged {
		ptr, err := bt.alloc(child)
		if err != nil {
			return nil, err
		}

		k, _ := parent.Key(i)

		parent, err = parent.Update(i, k, ptr)
		if err != nil {
			return nil, err
		}
	}

	return parent, nil
}
