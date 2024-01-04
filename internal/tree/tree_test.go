package tree

import (
	"bytes"
	"errors"
	"maps"
	"slices"
	"testing"
)

func mockTree() (Tree, map[uint64]Node) {
	mock := map[uint64]Node{
		0:  &PointerNode{keys: [][]byte{{'a'}, {'j'}, {'s'}}, ptrs: []uint64{1, 2, 3}},
		1:  &PointerNode{keys: [][]byte{{'a'}, {'d'}, {'g'}}, ptrs: []uint64{4, 5, 6}},
		2:  &PointerNode{keys: [][]byte{{'j'}, {'m'}, {'p'}}, ptrs: []uint64{7, 8, 9}},
		3:  &PointerNode{keys: [][]byte{{'s'}, {'v'}, {'y'}}, ptrs: []uint64{10, 11, 12}},
		4:  &LeafNode{keys: [][]byte{{'a'}, {'b'}, {'c'}}, vals: [][]byte{{'a'}, {'b'}, {'c'}}},
		5:  &LeafNode{keys: [][]byte{{'d'}, {'e'}, {'f'}}, vals: [][]byte{{'d'}, {'e'}, {'f'}}},
		6:  &LeafNode{keys: [][]byte{{'g'}, {'h'}, {'i'}}, vals: [][]byte{{'g'}, {'h'}, {'i'}}},
		7:  &LeafNode{keys: [][]byte{{'j'}, {'k'}, {'l'}}, vals: [][]byte{{'j'}, {'k'}, {'l'}}},
		8:  &LeafNode{keys: [][]byte{{'m'}, {'n'}, {'o'}}, vals: [][]byte{{'m'}, {'n'}, {'o'}}},
		9:  &LeafNode{keys: [][]byte{{'p'}, {'q'}, {'r'}}, vals: [][]byte{{'p'}, {'q'}, {'r'}}},
		10: &LeafNode{keys: [][]byte{{'s'}, {'t'}, {'u'}}, vals: [][]byte{{'s'}, {'t'}, {'u'}}},
		11: &LeafNode{keys: [][]byte{{'v'}, {'w'}, {'x'}}, vals: [][]byte{{'v'}, {'w'}, {'x'}}},
		12: &LeafNode{keys: [][]byte{{'y'}, {'z'}, {'{'}}, vals: [][]byte{{'y'}, {'z'}, {'{'}}},
	}

	next := uint64(13)

	return Tree{
		root: 0,
		read: func(ptr uint64) (Node, error) {
			n, ok := (mock)[ptr]
			if !ok {
				return nil, errors.New("mock get")
			}
			return n, nil
		},
		alloc: func(n Node) (uint64, error) {
			next++
			(mock)[next-1] = n
			return next - 1, nil
		},
		free: func(ptr uint64) error {
			delete(mock, ptr)
			return nil
		},
		maxNodeSize: 39,
	}, mock
}

func TestTreeGet(t *testing.T) {
	cases := []struct {
		name string
		k    []byte
		res  []byte
		err  error
	}{
		{
			name: "get first k-v",
			k:    []byte{'a'},
			res:  []byte{'a'},
		},
		{
			name: "get last k-v",
			k:    []byte{'{'},
			res:  []byte{'{'},
		},
		{
			name: "get middle k-v", k: []byte{'n'},
			res: []byte{'n'},
		},
		{
			name: "failed non existing key",
			k:    []byte{'}'},
			res:  nil,
			err:  errKeyNotExists,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tree, _ := mockTree()
			res, err := tree.Get(c.k)
			if err != c.err {
				t.Errorf("expected error %v, got %v", c.err, err)
			}
			if !bytes.Equal(res, c.res) {
				t.Errorf("expected val %v, got %v", c.res, res)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	cases := []struct {
		name string
		k    []byte
		v    []byte
		res  map[uint64]Node
		err  error
	}{
		{
			name: "insert without split",
			k:    []byte{'n', 'a'},
			v:    []byte{'n', 'a'},
			res: map[uint64]Node{
				1:  &PointerNode{keys: [][]byte{{'a'}, {'d'}, {'g'}}, ptrs: []uint64{4, 5, 6}},
				3:  &PointerNode{keys: [][]byte{{'s'}, {'v'}, {'y'}}, ptrs: []uint64{10, 11, 12}},
				4:  &LeafNode{keys: [][]byte{{'a'}, {'b'}, {'c'}}, vals: [][]byte{{'a'}, {'b'}, {'c'}}},
				5:  &LeafNode{keys: [][]byte{{'d'}, {'e'}, {'f'}}, vals: [][]byte{{'d'}, {'e'}, {'f'}}},
				6:  &LeafNode{keys: [][]byte{{'g'}, {'h'}, {'i'}}, vals: [][]byte{{'g'}, {'h'}, {'i'}}},
				7:  &LeafNode{keys: [][]byte{{'j'}, {'k'}, {'l'}}, vals: [][]byte{{'j'}, {'k'}, {'k'}}},
				9:  &LeafNode{keys: [][]byte{{'p'}, {'q'}, {'r'}}, vals: [][]byte{{'p'}, {'q'}, {'r'}}},
				10: &LeafNode{keys: [][]byte{{'s'}, {'t'}, {'u'}}, vals: [][]byte{{'s'}, {'t'}, {'u'}}},
				11: &LeafNode{keys: [][]byte{{'v'}, {'w'}, {'x'}}, vals: [][]byte{{'v'}, {'w'}, {'x'}}},
				12: &LeafNode{keys: [][]byte{{'y'}, {'z'}, {'{'}}, vals: [][]byte{{'y'}, {'z'}, {'{'}}},
				13: &LeafNode{keys: [][]byte{{'m'}, {'n'}, {'n', 'a'}, {'o'}}, vals: [][]byte{{'m'}, {'n'}, {'n', 'a'}, {'o'}}},
				14: &PointerNode{keys: [][]byte{{'j'}, {'m'}, {'p'}}, ptrs: []uint64{7, 13, 9}},
				15: &PointerNode{keys: [][]byte{{'a'}, {'j'}, {'s'}}, ptrs: []uint64{1, 14, 3}},
			},
		},
		{
			name: "insert with splits",
			k:    []byte{'n', 'a', 'a', 'a', 'a', 'a', 'a'},
			v:    []byte{'n', 'a', 'a', 'a', 'a', 'a', 'a'},
			res: map[uint64]Node{
				1:  &PointerNode{keys: [][]byte{{'a'}, {'d'}, {'g'}}, ptrs: []uint64{4, 5, 6}},
				3:  &PointerNode{keys: [][]byte{{'s'}, {'v'}, {'y'}}, ptrs: []uint64{10, 11, 12}},
				4:  &LeafNode{keys: [][]byte{{'a'}, {'b'}, {'c'}}, vals: [][]byte{{'a'}, {'b'}, {'c'}}},
				5:  &LeafNode{keys: [][]byte{{'d'}, {'e'}, {'f'}}, vals: [][]byte{{'d'}, {'e'}, {'f'}}},
				6:  &LeafNode{keys: [][]byte{{'g'}, {'h'}, {'i'}}, vals: [][]byte{{'g'}, {'h'}, {'i'}}},
				7:  &LeafNode{keys: [][]byte{{'j'}, {'k'}, {'l'}}, vals: [][]byte{{'j'}, {'k'}, {'k'}}},
				9:  &LeafNode{keys: [][]byte{{'p'}, {'q'}, {'r'}}, vals: [][]byte{{'p'}, {'q'}, {'r'}}},
				10: &LeafNode{keys: [][]byte{{'s'}, {'t'}, {'u'}}, vals: [][]byte{{'s'}, {'t'}, {'u'}}},
				11: &LeafNode{keys: [][]byte{{'v'}, {'w'}, {'x'}}, vals: [][]byte{{'v'}, {'w'}, {'x'}}},
				12: &LeafNode{keys: [][]byte{{'y'}, {'z'}, {'{'}}, vals: [][]byte{{'y'}, {'z'}, {'{'}}},
				13: &LeafNode{keys: [][]byte{{'m'}, {'n'}}, vals: [][]byte{{'m'}, {'n'}}},
				14: &LeafNode{keys: [][]byte{{'n', 'a', 'a', 'a', 'a', 'a', 'a'}, {'o'}}, vals: [][]byte{{'n', 'a', 'a', 'a', 'a', 'a', 'a'}, {'o'}}},
				15: &PointerNode{keys: [][]byte{{'j'}, {'m'}}, ptrs: []uint64{7, 13}},
				16: &PointerNode{keys: [][]byte{{'n', 'a', 'a', 'a', 'a', 'a', 'a'}, {'p'}}, ptrs: []uint64{14, 9}},
				17: &PointerNode{keys: [][]byte{{'a'}, {'j'}}, ptrs: []uint64{1, 15}},
				18: &PointerNode{keys: [][]byte{{'n', 'a', 'a', 'a', 'a', 'a', 'a'}, {'s'}}, ptrs: []uint64{16, 3}},
				19: &PointerNode{keys: [][]byte{{'a'}, {'n', 'a', 'a', 'a', 'a', 'a', 'a'}}, ptrs: []uint64{17, 18}},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tree, mock := mockTree()
			err := tree.Insert(c.k, c.v)
			if err != c.err {
				t.Errorf("expected error %v, got %v", c.err, err)
			} else if !maps.EqualFunc(c.res, mock, func(n1, n2 Node) bool {
				if n1.Type() != n2.Type() {
					return false
				}
				switch n1.Type() {
				case nodeLeaf:
					l1 := n1.(*LeafNode)
					l2 := n2.(*LeafNode)
					return slices.EqualFunc(l1.keys, l2.keys, bytes.Equal) && slices.EqualFunc(l1.keys, l2.keys, bytes.Equal)
				case nodePointer:
					p1 := n1.(*PointerNode)
					p2 := n2.(*PointerNode)
					return slices.EqualFunc(p1.keys, p2.keys, bytes.Equal) && slices.EqualFunc(p1.keys, p2.keys, bytes.Equal)
				}
				return false
			}) {
				t.Errorf("expected map %v, got %v", c.res, mock)
			}
		})
	}
}
