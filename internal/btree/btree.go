package btree

import (
	"errors"
	"fmt"
)

type getFunc func(uint64) ([]byte, error)
type pullFunc func(uint64) ([]byte, error)
type allocFunc func([]byte) (uint64, error)
type freeFunc func(uint64) error

type BTree struct {
	Root     uint64
	pgSize   int
	get      func(uint64) (node, error)
	pull     func(uint64) (node, error)
	alloc    func(node) (uint64, error)
	free     func(uint64) error
	readOnly bool
}

func NewReadOnly(root uint64, pgSize int, get getFunc) BTree {
	return BTree{
		Root:   root,
		pgSize: pgSize,
		get: func(ptr uint64) (node, error) {
			d, err := get(ptr)
			if err != nil {
				return nil, err
			}
			return decodeNode(d)
		},
		readOnly: true,
	}

}

func New(
	root uint64, pgSize int,
	get getFunc, pull pullFunc, alloc allocFunc, free freeFunc,
) BTree {
	return BTree{
		Root:   root,
		pgSize: pgSize,
		get: func(ptr uint64) (node, error) {
			d, err := get(ptr)
			if err != nil {
				return nil, err
			}
			return decodeNode(d)
		},
		pull: func(ptr uint64) (node, error) {
			d, err := pull(ptr)
			if err != nil {
				return nil, err
			}
			return decodeNode(d)
		},
		alloc: func(n node) (uint64, error) {
			return alloc(n.Encode())
		},
		free:     free,
		readOnly: false,
	}
}

func (bt *BTree) Get(k []byte) ([]byte, error) {
	errMsg := "btree: cannot get key: %v"

	if bt.Root == 0 {
		return nil, fmt.Errorf(errMsg, "tree is empty")
	}

	root, err := bt.get(bt.Root)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	_, v, err := bt.bTreeGet(root, k)
	return v, err
}

func (bt *BTree) Set(k, v []byte) (err error) {
	if bt.readOnly {
		return fmt.Errorf("btree: set operation not allow on read only tree")
	}
	errMsg := "btree: cannot set key: %v"

	if bt.Root == 0 {
		root := leafNode{}
		root, err := root.Insert(0, k, v)
		if err != nil {
			return fmt.Errorf(errMsg, err)
		}

		bt.Root, err = bt.alloc(root)
		if err != nil {
			return fmt.Errorf(errMsg, err)
		}

		return nil
	}

	root, err := bt.pull(bt.Root)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	inserted, err := bt.bTreeInsert(root, k, v)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	if inserted.Size() > bt.pgSize {
		insertedPtr, err := bt.alloc(inserted)
		if err != nil {
			return fmt.Errorf(errMsg, err)
		}

		k, err := inserted.Key(0)
		if err != nil {
			return fmt.Errorf(errMsg, err)
		}

		t := pointerNode{}
		t.Insert(0, k, insertedPtr)

		inserted, err = bt.splitChildPtr(0, t, inserted)
		if err != nil {
			return fmt.Errorf(errMsg, err)
		}

	}

	bt.Root, err = bt.alloc(inserted)
	return fmt.Errorf(errMsg, err)
}

func (bt *BTree) Delete(k []byte) (bool, error) {
	if bt.readOnly {
		return false, fmt.Errorf("btree: delete operation not allow on read only tree")
	}
	errMsg := "btree: cannot delete key: %v"

	if bt.Root == 0 {
		return false, fmt.Errorf(errMsg, "tree is empty")
	}

	root, err := bt.pull(bt.Root)
	if err != nil {
		return false, fmt.Errorf(errMsg, err)
	}

	var deleted bool
	root, deleted, err = bt.bTreeDelete(root, k)
	if err != nil {
		return false, fmt.Errorf(errMsg, err)
	}

	if !deleted {
		return deleted, nil
	}

	if root.Type() == btreePointer && root.Total() == 1 {
		pntrRoot := root.(pointerNode)
		bt.Root, _ = pntrRoot.Ptr(0)
	} else {
		bt.Root, err = bt.alloc(root)
		if err != nil {
			return false, fmt.Errorf(errMsg, err)
		}
	}

	return true, nil
}

func (bt *BTree) bTreeGet(n node, k []byte) (node, []byte, error) {
	i, exists := n.Search(k)

	switch n.Type() {
	case btreeLeaf:
		leafN := n.(leafNode)

		if !exists {
			return nil, nil, errors.New("key does not exist")
		}

		v, err := leafN.Val(i)
		if err != nil {
			return nil, nil, err
		}
		return n, v, nil

	case btreePointer:
		pntrN := n.(pointerNode)

		ptr, _ := pntrN.Ptr(i)
		child, err := bt.get(ptr)
		if err != nil {
			return nil, nil, err
		}

		return bt.bTreeGet(child, k)
	default:
		return nil, nil, errors.New("invalid node type")
	}
}

func (bt *BTree) bTreeInsert(n node, k, v []byte) (node, error) {
	i, exists := n.Search(k)

	var err error
	var inserted node

	switch n.Type() {
	case btreeLeaf:
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

	case btreePointer:
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

		inserted, err = bt.splitChildPtr(i, pntrN, inserted)
		if err != nil {
			return nil, err
		}

		inPtr, err := bt.alloc(inserted)
		if err != nil {
			return nil, err
		}

		ptrKey, _ := pntrN.Key(i)
		inserted, err = pntrN.Update(i, ptrKey, inPtr)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("invalid node type")
	}

	return inserted, nil
}

func (bt *BTree) bTreeDelete(n node, k []byte) (node, bool, error) {
	i, exists := n.Search(k)

	var err error
	var deleted node

	switch n.Type() {
	case btreeLeaf:
		if !exists {
			return nil, false, errors.New("key does not exist")
		}

		leafN := n.(leafNode)
		deleted, err = leafN.Delete(i)
		if err != nil {
			return nil, false, err
		}

	case btreePointer:
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

		if deleted.Size() > bt.pgSize/4 {
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
		return nil, false, errors.New("invalid node type")
	}

	return deleted, true, nil
}

func (bt *BTree) splitChildPtr(i int, par pointerNode, child node) (node, error) {
	if !bt.shouldSplit(child) {
		return par, nil
	}

	var (
		lPtr   uint64
		rPtr   uint64
		lKey   []byte
		rKey   []byte
		err    error
		errMsg = "cannot split child ptr: %v"
	)

	switch child.Type() {
	case btreePointer:
		ptrChild := child.(pointerNode)

		l, r := ptrChild.Split()

		lKey = l.keys[0]
		rKey = r.keys[0]

		lPtr, err = bt.alloc(l)
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}

		rPtr, err = bt.alloc(r)
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
		break

	case btreeLeaf:
		leafChild := child.(pointerNode)

		l, r := leafChild.Split()

		lKey = l.keys[0]
		rKey = r.keys[0]

		lPtr, err = bt.alloc(l)
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}

		rPtr, err = bt.alloc(r)
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
		break

	default:
		return nil, fmt.Errorf(errMsg, "child has invalid type")
	}

	par, err = par.Update(i, lKey, lPtr)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	par, err = par.Insert(i+1, rKey, rPtr)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return par, nil
}

func (bt *BTree) mergeChildPtr(i int, par pointerNode, child node) (node, error) {
	var (
		leftSib  node
		rightSib node
		merge    uint8 = 0
		err      error
		errMsg   = "cannot merge child pointer: %v"
	)

	if i > 0 {
		leftSib, err = bt.get(par.ptrs[i-1])
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
		if bt.canMerge2(child, leftSib) {
			merge++
		}
	}

	if i < par.Total()-1 {
		rightSib, err = bt.get(par.ptrs[i+1])
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
		if bt.canMerge2(child, rightSib) {
			merge++
		}
	}

	if merge == 2 {
		if !bt.canMerge3(child, leftSib, rightSib) {
			merge--
		}
	} else if merge == 0 {
		return par, nil
	}

	switch child.Type() {

	case btreePointer:
		pntrChild := child.(pointerNode)

		if merge >= 1 {
			if rightSib != nil {
				pntrRight := rightSib.(pointerNode)
				pntrChild, err = pntrChild.Merge(pntrRight)
				if err != nil {
					return nil, fmt.Errorf(errMsg, err)
				}
			} else {
				pntrLeft := rightSib.(pointerNode)
				pntrChild, err = pntrChild.Merge(pntrLeft)
				if err != nil {
					return nil, fmt.Errorf(errMsg, err)
				}
			}
		}

		if merge == 2 {
			pntrLeft := rightSib.(pointerNode)
			pntrChild, err = pntrChild.Merge(pntrLeft)
			if err != nil {
				return nil, fmt.Errorf(errMsg, err)
			}
		}

		child = pntrChild

	case btreeLeaf:
		leafChild := child.(leafNode)

		if merge >= 1 {
			if rightSib != nil {
				leafRight := rightSib.(leafNode)
				leafChild, err = leafChild.Merge(leafRight)
				if err != nil {
					return nil, fmt.Errorf(errMsg, err)
				}
			} else {
				leafLeft := rightSib.(leafNode)
				leafChild, err = leafChild.Merge(leafLeft)
				if err != nil {
					return nil, fmt.Errorf(errMsg, err)
				}
			}
		}

		if merge == 2 {
			leafLeft := rightSib.(leafNode)
			leafChild, err = leafChild.Merge(leafLeft)
			if err != nil {
				return nil, fmt.Errorf(errMsg, err)
			}
		}

		child = leafChild

	default:
		return nil, fmt.Errorf(errMsg, "child has invalid type")
	}

	first, err := child.Key(0)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	ptr, err := bt.alloc(child)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	var j int
	if merge > 1 {
		j = i - 1
		par, err = par.Delete(i)
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}

		if merge == 3 {
			par, err = par.Delete(i + 1)
			if err != nil {
				return nil, fmt.Errorf(errMsg, err)
			}
		}
	} else {
		j = i
		par, err = par.Delete(i + 1)
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
	}

	par, err = par.Update(j, first, ptr)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return par, nil
}

func (bt *BTree) canMerge2(a, b node) bool {
	return a.Type() == b.Type() && a.Size()+b.Size() <= bt.pgSize
}

func (bt *BTree) canMerge3(a, b, c node) bool {
	return a.Type() == b.Type() && a.Type() == c.Type() && a.Size()+b.Size()+c.Size() <= bt.pgSize
}

func (bt *BTree) shouldSplit(n node) bool {
	return n.Size() > bt.pgSize
}
