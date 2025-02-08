package bpt

import "github.com/pavosql/pavosql/internal/bpt/node"

const pageSize = 4096

type pager interface {
	ReadPage(off uint64) ([]byte, error)
	Commit() error
	Rollback() error
}

type Tree struct {
	root  uint64
	pager pager
}

func (tree *Tree) Get(k []byte) ([]byte, error) {
	cur, err := tree.pager.ReadPage(tree.root)
	if err != nil {
		return nil, err
	}
	// TODO: add defered recover and custom errors
	for {
		i, exists := node.Search(cur, k)
		switch node.TypeOf(cur) {
		case node.TypePointer:
			cur, err = tree.pager.ReadPage(node.Pointer(cur, i))
			if err != nil {
				return nil, err
			}
		case node.TypeLeaf:
			if !exists {
				return nil, nil
			}
			return node.Value(cur, i), nil
		default:
			return nil, nil
		}
	}
}

func (tree *Tree) Insert(k []byte, v []byte) error {
	return nil
}

func (tree *Tree) Update(k []byte, v []byte) error {
	return nil
}

func (tree *Tree) Delete(k []byte) error {
	return nil
}
