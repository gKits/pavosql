package tree

import (
	"errors"
)

const PageSize = 4096
const MaxCell = PageSize - NodeHeader

var (
	errCellToLarge      = errors.New("cell size exceeds maximum")
	errMalformedRecurse = errors.New("recursive result is malformed")
)

type Tree struct {
	root        uint64
	read        func(uint64) (Node, error) // callback to read node from backend
	alloc       func(Node) (uint64, error) // callback to alloc node in backend
	free        func(uint64) error         // callback to free node in backend
	maxNodeSize int
}

type recurseResult struct {
	Key []byte
	Ptr uint64
}

func (t *Tree) Get(k []byte) ([]byte, error) {
	cur, err := t.read(t.root)
	if err != nil {
		return nil, err
	}

	for {
		switch cur.Type() {
		case nodePointer:
			pointer, ok := cur.(*PointerNode)
			if !ok {
				return nil, errNodeAssert
			}

			idx, exists := pointer.Find(k)
			if !exists {
				idx--
			}
			ptr, err := pointer.PtrAt(idx)
			if err != nil {
				return nil, err
			}

			cur, err = t.read(ptr)
			if err != nil {
				return nil, err
			}
			continue

		case nodeLeaf:
			leaf, ok := cur.(*LeafNode)
			if !ok {
				return nil, errNodeAssert
			}
			return leaf.Val(k)

		default:
			return nil, errInvalNodeType
		}
	}
}

func (t *Tree) Insert(k, v []byte) error {
	if !t.cellFits(k, v) {
		return errCellToLarge
	}

	res, err := t.recursiveInsert(t.root, k, v)
	if err != nil {
		return err
	}

	switch len(res) {
	case 1:
		t.root = res[0].Ptr
	case 2:
		root := &PointerNode{
			keys: [][]byte{res[0].Key, res[1].Key},
			ptrs: []uint64{res[0].Ptr, res[1].Ptr},
		}

		t.root, err = t.alloc(root)
		if err != nil {
			return err
		}
	default:
		return errMalformedRecurse
	}
	return nil
}

func (t *Tree) recursiveInsert(ptr uint64, k, v []byte) ([]recurseResult, error) {
	cur, err := t.read(ptr)
	if err != nil {
		return nil, err
	}
	if err = t.free(ptr); err != nil {
		return nil, err
	}

	switch cur.Type() {
	case nodePointer:
		pointer, ok := cur.(*PointerNode)
		if !ok {
			return nil, errNodeAssert
		}

		idx, exists := pointer.Find(k)
		if !exists {
			idx--
		}

		ptr, err = pointer.PtrAt(idx)
		if err != nil {
			return nil, err
		}

		res, err := t.recursiveInsert(ptr, k, v)
		if err != nil {
			return nil, err
		}

		if err := pointer.Update(res[0].Key, res[0].Ptr); err != nil {
			return nil, err
		}

		if len(res) > 1 {
			if err := pointer.Insert(res[1].Key, res[1].Ptr); err != nil {
				return nil, err
			}
		}
		cur = pointer

	case nodeLeaf:
		leaf, ok := cur.(*LeafNode)
		if !ok {
			return nil, errNodeAssert
		}

		if err := leaf.Insert(k, v); err != nil {
			return nil, err
		}
		cur = leaf

	default:
		return nil, errInvalNodeType
	}

	if cur.Size() > t.maxNodeSize {
		l, r := cur.Split()

		lPtr, err := t.alloc(l)
		if err != nil {
			return nil, err
		}

		rPtr, err := t.alloc(r)
		if err != nil {
			return nil, err
		}

		lK, err := l.Key(0)
		if err != nil {
			return nil, err
		}

		rK, err := r.Key(0)
		if err != nil {
			return nil, err
		}

		return []recurseResult{{Key: lK, Ptr: lPtr}, {Key: rK, Ptr: rPtr}}, nil
	}

	ptr, err = t.alloc(cur)
	if err != nil {
		return nil, err
	}

	origin, err := cur.Key(0)
	if err != nil {
		return nil, err
	}

	return []recurseResult{{Key: origin, Ptr: ptr}}, nil
}

func (t *Tree) Update(k, v []byte) error {
	if !t.cellFits(k, v) {
		return errCellToLarge
	}

	res, err := t.recursiveUpdate(t.root, k, v)
	if err != nil {
		return err
	}

	switch len(res) {
	case 1:
		t.root = res[0].Ptr
	case 2:
		root := &PointerNode{
			keys: [][]byte{res[0].Key, res[1].Key},
			ptrs: []uint64{res[0].Ptr, res[1].Ptr},
		}

		t.root, err = t.alloc(root)
		if err != nil {
			return err
		}
	default:
		return errMalformedRecurse
	}
	return nil
}

func (t *Tree) recursiveUpdate(ptr uint64, k, v []byte) ([]recurseResult, error) {
	cur, err := t.read(ptr)
	if err != nil {
		return nil, err
	}
	if err = t.free(ptr); err != nil {
		return nil, err
	}

	switch cur.Type() {
	case nodePointer:
		pointer, ok := cur.(*PointerNode)
		if !ok {
			return nil, errNodeAssert
		}

		idx, _ := pointer.Find(k)
		ptr, err = pointer.PtrAt(idx)
		if err != nil {
			return nil, err
		}

		res, err := t.recursiveUpdate(ptr, k, v)
		if err != nil {
			return nil, err
		}

		if err := pointer.Update(res[0].Key, res[0].Ptr); err != nil {
			return nil, err
		}

		if len(res) > 1 {
			if err := pointer.Insert(res[1].Key, res[1].Ptr); err != nil {
				return nil, err
			}
		}

	case nodeLeaf:
		leaf, ok := cur.(*LeafNode)
		if !ok {
			return nil, errNodeAssert
		}

		if err := leaf.Update(k, v); err != nil {
			return nil, err
		}
		cur = leaf

	default:
		return nil, errInvalNodeType
	}

	if cur.Size() > t.maxNodeSize {
		l, r := cur.Split()

		lPtr, err := t.alloc(l)
		if err != nil {
			return nil, err
		}

		rPtr, err := t.alloc(r)
		if err != nil {
			return nil, err
		}

		nK, err := r.Key(0)
		if err != nil {
			return nil, err
		}

		return []recurseResult{{Key: k, Ptr: lPtr}, {Key: nK, Ptr: rPtr}}, nil
	}

	ptr, err = t.alloc(cur)
	if err != nil {
		return nil, err
	}

	return []recurseResult{{Key: k, Ptr: ptr}}, nil
}

func (t *Tree) Delete(k []byte) error {
	return nil
}

func (t *Tree) recursiveDelete(ptr uint64, k, v []byte) (Node, error) {
	cur, err := t.read(ptr)
	if err != nil {
		return nil, err
	}
	if err = t.free(ptr); err != nil {
		return nil, err
	}

	switch cur.Type() {
	case nodePointer:
		pointer, ok := cur.(*PointerNode)
		if !ok {
			return nil, errNodeAssert
		}

		idx, _ := pointer.Find(k)
		ptr, err = pointer.PtrAt(idx)
		if err != nil {
			return nil, err
		}

		child, err := t.recursiveDelete(ptr, k, v)
		if err != nil {
			return nil, err
		}

		if err = t.tryMergeNeighbors(pointer, child, idx); err != nil {
			return nil, err
		}

	case nodeLeaf:
		if err := cur.Delete(k); err != nil {
			return nil, err
		}

	default:
		return nil, errInvalNodeType
	}

	ptr, err = t.alloc(cur)
	if err != nil {
		return nil, err
	}

	return cur, nil
}

func (t *Tree) tryMergeNeighbors(parent *PointerNode, child Node, idx int) error {
	if child.Size() < t.maxNodeSize/3 {
		if idx > 0 {
			lPtr, err := parent.PtrAt(idx - 1)
			if err != nil {
				return err
			}

			l, err := t.read(lPtr)
			if err != nil {
				return err
			}

			if child.Size()+l.Size() < t.maxNodeSize {
				if err := l.Merge(child); err != nil {
					return err
				}
				child = l
			}
		}
		if idx < parent.NKeys()-1 {
			rPtr, err := parent.PtrAt(idx + 1)
			if err != nil {
				return err
			}

			r, err := t.read(rPtr)
			if err != nil {
				return err
			}

			if child.Size()+r.Size() < t.maxNodeSize {
				if err := child.Merge(r); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (t *Tree) cellFits(k, v []byte) bool {
	return 4+len(k)+len(v) <= t.maxNodeSize-NodeHeader || 2+len(k)+8 <= t.maxNodeSize-NodeHeader
}
