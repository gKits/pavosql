package backend

import (
	"errors"
)

const pageSize = 4096

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

func (bt *bTree) Insert(k, v []byte) error {
	if bt.root == 0 {
		root := leafNode{}
		root, err := root.insert(0, k, v)
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

	if inserted.size() > pageSize {
		insertedPtr, err := bt.alloc(inserted)
		if err != nil {
			return err
		}

		k, err := inserted.key(0)
		if err != nil {
			return err
		}

		t := pointerNode{}
		t.insert(0, k, insertedPtr)

		inserted, err = bt.splitChildPtr(0, t, inserted)
		if err != nil {
			return err
		}

	}

	bt.root, err = bt.alloc(inserted)
	return err
}

func (bt *bTree) Delete(k []byte) error {
	if bt.root == 0 {
		return errBTreeDeleteEmpty
	}

	root, err := bt.pull(bt.root)
	if err != nil {
		return err
	}

	root, err = bt.bTreeDelete(root, k)
	if err != nil {
		return err
	}

	if root.typ() == ptrNode && root.total() == 1 {
		pntrRoot := root.(pointerNode)
		bt.root, _ = pntrRoot.ptr(0)
	} else {
		bt.root, err = bt.alloc(root)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bt *bTree) bTreeGet(n node, k []byte) (node, []byte, error) {
	i, exists := n.search(k)

	switch n.typ() {
	case lfNode:
		leafN := n.(leafNode)

		if !exists {
			return nil, nil, errBTreeGetKey
		}

		v, err := leafN.val(i)
		if err != nil {
			return nil, nil, err
		}
		return n, v, nil

	case ptrNode:
		pntrN := n.(pointerNode)

		ptr, _ := pntrN.ptr(i)
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
	i, exists := n.search(k)

	var err error
	var inserted node

	switch n.typ() {
	case lfNode:
		leafN := n.(leafNode)

		if !exists {
			inserted, err = leafN.insert(i, k, v)
			if err != nil {
				return nil, err
			}
		} else {
			inserted, err = leafN.update(i, k, v)
			if err != nil {
				return nil, err
			}
		}

	case ptrNode:
		pntrN := n.(pointerNode)

		ptr, _ := pntrN.ptr(i)
		child, err := bt.get(ptr)
		if err != nil {
			return nil, err
		}

		inserted, err = bt.bTreeInsert(child, k, v)
		if err != nil {
			return nil, err
		}

		if inserted.size() > pageSize {
			inserted, err = bt.splitChildPtr(i, pntrN, inserted)
			if err != nil {
				return nil, err
			}
		} else {
			insertedPtr, err := bt.alloc(inserted)
			if err != nil {
				return nil, err
			}

			ptrKey, _ := pntrN.key(i)
			inserted, err = pntrN.update(i, ptrKey, insertedPtr)
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

func (bt *bTree) bTreeDelete(n node, k []byte) (node, error) {
	i, exists := n.search(k)

	var err error
	var deleted node

	switch n.typ() {
	case lfNode:
		if !exists {
			return nil, errBTreeDeleteKey
		}

		leafN := n.(leafNode)
		deleted, err = leafN.delete(i)
		if err != nil {
			return nil, err
		}

	case ptrNode:
		n := n.(pointerNode)

		ptr, _ := n.ptr(i)
		child, err := bt.pull(ptr)
		if err != nil {
			return nil, err
		}

		deleted, err = bt.bTreeDelete(child, k)
		if err != nil {
			return nil, err
		}

		if deleted.size() > pageSize/4 {
			deleted, err = bt.mergeChildPtr(i, n, deleted)
			if err != nil {
				return nil, err
			}
		} else {
			deletedPtr, err := bt.alloc(deleted)
			if err != nil {
				return nil, err
			}

			k, _ := deleted.key(i)
			deleted, err = n.update(i, k, deletedPtr)
			if err != nil {
				return nil, err
			}
		}

	default:
		return nil, errBTreeHeader
	}

	return deleted, nil
}

func (bt *bTree) splitChildPtr(i int, parent pointerNode, child node) (node, error) {
	l, r := child.split()

	lPtr, err := bt.alloc(l)
	if err != nil {
		return nil, err
	}

	rPtr, err := bt.alloc(r)
	if err != nil {
		return nil, err
	}

	lKey, err := parent.key(i)
	if err != nil {
		return nil, err
	}

	split, err := parent.update(i, lKey, lPtr)
	if err != nil {
		return nil, err
	}

	rKey, err := r.key(0)
	if err != nil {
		return nil, err
	}

	split, err = parent.insert(i+1, rKey, rPtr)
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

	if i < parent.total()-1 {
		siblingPtr, err = parent.ptr(i + 1)
		if err != nil {
			return nil, err
		}

		sibling, err = bt.pull(siblingPtr)
		if err != nil {
			return nil, err
		}

		if sibling.size()+child.size() < pageSize {
			child, err = child.merge(sibling)
			if err != nil {
				return nil, err
			}
		}

		merged = true
	}

	if i != 0 {
		i--
		siblingPtr, err = parent.ptr(i)
		if err != nil {
			return nil, err
		}

		sibling, err = bt.pull(siblingPtr)
		if err != nil {
			return nil, err
		}

		if sibling.size()+child.size() < pageSize {
			child, err = sibling.merge(child)
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

		k, _ := parent.key(i)

		parent, err = parent.update(i, k, ptr)
		if err != nil {
			return nil, err
		}
	}

	return parent, nil
}
