package tree

import (
	"errors"
	"fmt"
	"slices"
)

type pager interface {
	ReadPage(int64) ([PageSize]byte, error)
	Alloc([PageSize]byte) (int64, error)
	Free(int64) error
	Commit() error
	Abort() error
}

type Tree struct {
	root     int64
	pager    pager
	readOnly bool
}

func New() *Tree {
	return &Tree{}
}

func (t *Tree) Get(k []byte) ([]byte, error) {
	page, err := t.pager.ReadPage(t.root)
	if err != nil {
		return nil, err
	}
	cur := node(page)

	for {
		i, exists := cur.Search(k)

		switch cur.Type() {
		case PointerPage:
			ptr := cur.Pointer(i)
			page, err = t.pager.ReadPage(ptr)
			if err != nil {
				return nil, fmt.Errorf("tree: failed to read page: %w", err)
			}
			cur = node(page)
		case LeafPage:
			if !exists {
				return nil, errors.New("key does not exists on leaf node")
			}
			return cur.Val(i), nil
		default:
			return nil, errors.New("invalid page type")
		}
	}
}

func (t *Tree) Set(k []byte, v []byte) error {
	if t.readOnly {
		return errors.New("tree: cannot write onto read only tree")
	}
	page, err := t.pager.ReadPage(t.root)
	if err != nil {
		return fmt.Errorf("tree: failed to read root page: %w", err)
	}
	cur := node(page)

	visited := []node{cur}
	for {
		i, exists := cur.Search(k)

		switch cur.Type() {
		case PointerPage:
			ptr := cur.Pointer(i)
			page, err = t.pager.ReadPage(ptr)
			if err != nil {
				return fmt.Errorf("tree: failed to read page: %w", err)
			}
			cur = node(page)
			visited = append(visited, cur)
			continue

		case LeafPage:
			if !exists {
				return errors.New("key does not exists on leaf node")
			}

			if cur.CanSet(k, v) {
				newNode := cur.Set(i, k, v)
				ptr, err := t.pager.Alloc(newNode)
				if err != nil {
					return fmt.Errorf("tree: failed to allocate page: %w", err)
				}
				// TODO: pass ptr to parent node
				_ = ptr
				break
			}

			// TODO: handle splitting
			left, right := cur.Split()
			_, _ = left, right

		default:
			return errors.New("invalid page type")
		}
		break
	}

	for _, n := range slices.Backward(visited) {
		ptr, err := t.pager.Alloc(n)
		if err != nil {
			return err
		}
		// TODO: pass ptr to parent nodes
		_ = ptr
	}

	return nil
}

func (t *Tree) Delete(k []byte) error {
	if t.readOnly {
		return errors.New("tree: cannot write onto read only tree")
	}
	return nil
}
