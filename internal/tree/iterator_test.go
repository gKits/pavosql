package tree

import (
	"bytes"
	"testing"

	"github.com/gKits/PavoSQL/pkg/stack"
)

func TestIterator(t *testing.T) {
	tree, m := mockTree()

	cases := []struct {
		name string
		it   Iterator
		num  int
		res  [][]byte
		err  error
	}{
		{
			name: "succefully iterate over one adjacent nodes",
			it: Iterator{
				stk:    stack.Stack[Node]{m[0], m[1], m[4]},
				idxStk: stack.Stack[int]{0, 0, 0},
				t:      &tree,
			},
			num: 5,
			res: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}},
		},
		{
			name: "succefully iterate multiple adjacent nodes",
			it: Iterator{
				stk:    stack.Stack[Node]{m[0], m[1], m[4]},
				idxStk: stack.Stack[int]{0, 0, 0},
				t:      &tree,
			},
			num: 7,
			res: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}, {'f'}, {'g'}},
		},
		{
			name: "succefully iterate multiple levels",
			it: Iterator{
				stk:    stack.Stack[Node]{m[0], m[1], m[4]},
				idxStk: stack.Stack[int]{0, 0, 0},
				t:      &tree,
			},
			num: 10,
			res: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}, {'f'}, {'g'}, {'h'}, {'i'}, {'j'}},
		},
		{
			name: "succefully iterate over whole tree",
			it: Iterator{
				stk:    stack.Stack[Node]{m[0], m[1], m[4]},
				idxStk: stack.Stack[int]{0, 0, 0},
				t:      &tree,
			},
			num: 30,
			res: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}, {'f'}, {'g'}, {'h'}, {'i'}, {'j'}, {'k'}, {'l'}, {'m'}, {'n'}, {'o'}, {'p'}, {'q'}, {'r'}, {'s'}, {'t'}, {'u'}, {'v'}, {'w'}, {'x'}, {'y'}, {'z'}, {'{'}},
			err: EndOfTree,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for i := 0; i < c.num; i++ {
				k, v, err := c.it.Next()
				if err == nil {
					if !bytes.Equal(k, c.res[i]) || !bytes.Equal(v, c.res[i]) {
						t.Log(c.it)
						t.Errorf("expected k-v %v, got %v, %v", c.res[i], k, v)
					}
				} else if err != c.err {
					t.Errorf("expected error %v, got %v", c.err, err)
				}
			}
		})
	}
}
