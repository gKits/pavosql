package tree

import (
	"errors"
	"fmt"

	"github.com/gKits/PavoSQL/pkg/stack"
)

type Iterator struct {
	t      *Tree
	stk    stack.Stack[Node]
	idxStk stack.Stack[int]
}

var EndOfTree = errors.New("end of tree")

func NewIterator(t *Tree, k []byte) *Iterator {
	return &Iterator{
		t: t,
	}
}

func (it *Iterator) Next() ([]byte, []byte, error) {
	for {
		var (
			cur Node
			err error
		)

		idx, err := it.idxStk.Pop()
		if err != nil {
			return nil, nil, EndOfTree
		}

		cur, err = it.stk.Peek()
		if err != nil {
			return nil, nil, err
		}

		fmt.Println(cur, it.idxStk, idx)

		switch cur.Type() {
		case nodeLeaf:
			if idx == -1 {
				idx++
			}

			leaf := cur.(*LeafNode)
			k, err := leaf.Key(idx)
			if err != nil {
				return nil, nil, err
			}

			v, err := leaf.ValAt(idx)
			if err != nil {
				return nil, nil, err
			}

			if idx+1 >= cur.NKeys() {
				_, err := it.stk.Pop()
				if err != nil {
					return nil, nil, err
				}
			} else {
				it.idxStk.Push(idx + 1)
			}

			return k, v, nil

		case nodePointer:
			if idx+1 >= cur.NKeys() {
				it.stk.Pop()
				continue
			} else {
				idx++
				it.idxStk.Push(idx)
			}

			pointer := cur.(*PointerNode)

			ptr, err := pointer.PtrAt(idx)
			if err != nil {
				return nil, nil, err
			}

			next, err := it.t.read(ptr)
			if err != nil {
				return nil, nil, err
			}

			it.stk.Push(next)
			it.idxStk.Push(-1)

		default:
			return nil, nil, errInvalNodeType
		}
	}
}
