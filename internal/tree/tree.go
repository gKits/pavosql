package tree

import (
	"errors"
	"fmt"

	"github.com/pavosql/pavosql/internal/common"
	"github.com/pavosql/pavosql/internal/tree/node"
)

type pager interface {
	ReadPage(off uint64) ([common.PageSize]byte, error)
	Commit() error
	Rollback() error
}

type Tree struct {
	root  uint64
	pager pager
}

func New() *Tree {
	return &Tree{}
}

func (t *Tree) Get(k []byte) ([]byte, error) {
	page, err := t.pager.ReadPage(t.root)
	if err != nil {
		return nil, err
	}
	cur := node.Node(page)

	for {
		i, exists := cur.Search(k)

		switch cur.Type() {
		case common.PointerPage:
			ptr := cur.Pointer(i)
			page, err = t.pager.ReadPage(ptr)
			if err != nil {
				return nil, fmt.Errorf("tree: failed to read page: %w", err)
			}
			cur = node.Node(page)
		case common.LeafPage:
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
	return nil
}

func (t *Tree) Delete(k []byte) error {
	return nil
}
