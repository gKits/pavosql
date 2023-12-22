package tree

import (
	"bytes"
	"slices"
	"testing"
)

func TestLeafDecode(t *testing.T) {
	cases := []struct {
		name string
		in   []byte
		res  *LeafNode
		err  error
	}{
		{
			name: "invalid type",
			in:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			res:  &LeafNode{},
			err:  errInvalNodeType,
		},
		{
			name: "decode leaf node",
			in:   []byte{0x00, 0x65, 0x00, 0x00, 0x00, 0x00},
			res:  &LeafNode{keys: [][]byte{}, vals: [][]byte{}},
			err:  nil,
		},
		{
			name: "decode leaf node with one k-v pair",
			in:   []byte{0x00, 0x65, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 'k', 'v'},
			res:  &LeafNode{keys: [][]byte{{'k'}}, vals: [][]byte{{'v'}}},
			err:  nil,
		},
		{
			name: "decode pointer node with multiple k-v pairs",
			in: []byte{
				0x00, 0x65,
				0x00, 0x00, 0x00, 0x04,
				0x00, 0x01,
				0x00, 0x05,
				'a',
				0x01, 0x01, 0x01, 0x01, 0x01,
				0x00, 0x01,
				0x00, 0x0A,
				'b',
				'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a',
				0x00, 0x02,
				0x00, 0x02,
				'b', 'a',
				'b', 'a',
				0x00, 0x03,
				0x00, 0x04,
				'c', 'b', 'a',
				0x01, 0x02, 0x03, 0x04,
			},
			res: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'b', 'a'}, {'c', 'b', 'a'}},
				vals: [][]byte{
					{0x01, 0x01, 0x01, 0x01, 0x01},
					{'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'},
					{'b', 'a'},
					{0x01, 0x02, 0x03, 0x04},
				},
			},
			err: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := &LeafNode{}
			err := res.Decode(c.in)
			if err != c.err {
				t.Errorf("expected error %v, got %v", c.err, err)
			} else if res == nil {
				if c.res != nil {
					t.Errorf("expected result %v, got %v", c.res, res)
				}
			} else {
				if c.res == nil {
					t.Errorf("expected result %v, got %v", c.res, res)
				} else if !slices.EqualFunc(res.keys, c.res.keys, bytes.Equal) {
					t.Errorf("expected keys %v,  got %v", c.res.keys, res.keys)
				} else if !slices.EqualFunc(res.vals, c.res.vals, bytes.Equal) {
					t.Errorf("expected vals %v,  got %v", c.res.vals, res.vals)
				}
			}
		})
	}
}

func TestLeafEncode(t *testing.T) {
	cases := []struct {
		name string
		in   *LeafNode
		res  []byte
	}{
		{
			name: "encode empty leaf node",
			in:   &LeafNode{keys: [][]byte{}, vals: [][]byte{}},
			res:  []byte{0x00, 0x65, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name: "encode leaf node with one k-v pair",
			in:   &LeafNode{keys: [][]byte{{'k'}}, vals: [][]byte{{'v'}}},
			res:  []byte{0x00, 0x65, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 'k', 'v'},
		},
		{
			name: "encode leaf node with multiple k-v pairs",
			in: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'b', 'a'}, {'c', 'b', 'a'}},
				vals: [][]byte{
					{0x01, 0x01, 0x01, 0x01, 0x01},
					{'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'},
					{'b', 'a'},
					{0x01, 0x02, 0x03, 0x04},
				},
			},
			res: []byte{
				0x00, 0x65,
				0x00, 0x00, 0x00, 0x04,
				0x00, 0x01,
				0x00, 0x05,
				'a',
				0x01, 0x01, 0x01, 0x01, 0x01,
				0x00, 0x01,
				0x00, 0x0A,
				'b',
				'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a',
				0x00, 0x02,
				0x00, 0x02,
				'b', 'a',
				'b', 'a',
				0x00, 0x03,
				0x00, 0x04,
				'c', 'b', 'a',
				0x01, 0x02, 0x03, 0x04,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := c.in.Encode()
			if !bytes.Equal(res, c.res) {
				t.Errorf("expected result %v, got %v", c.res, res)
			}
		})
	}
}

func TestLeafVal(t *testing.T) {
	cases := []struct {
		name string
		k    []byte
		res  []byte
		err  error
	}{
		{
			name: "successfully read key",
			k:    []byte{'b'},
			res:  []byte{'b'},
		},
		{
			name: "failed read non existing key",
			k:    []byte{'e'},
			res:  nil,
			err:  errKeyNotExists,
		},
	}

	leaf := LeafNode{
		keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
		vals: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := leaf.Val(c.k)

			if err != c.err {
				t.Errorf("expected error %v,  got %v", c.err, err)
			} else if !bytes.Equal(res, c.res) {
				t.Errorf("expected pointer %v, got %v", c.res, res)
			}
		})
	}
}

func TestLeafInsert(t *testing.T) {
	cases := []struct {
		name string
		k    []byte
		v    []byte
		res  *LeafNode
		err  error
	}{
		{
			name: "insert in middle",
			k:    []byte{'b', 'a'},
			v:    []byte{'b', 'a'},
			res: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'b', 'a'}, {'c'}, {'d'}},
				vals: [][]byte{{'a'}, {'b'}, {'b', 'a'}, {'c'}, {'d'}},
			},
		},
		{
			name: "insert last",
			k:    []byte{'e'},
			v:    []byte{'e'},
			res: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}},
				vals: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}},
			},
		},
		{
			name: "failed insert existing key",
			k:    []byte{'a'},
			v:    []byte{'a'},
			res: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
				vals: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
			},
			err: errKeyExists,
		},
	}

	leaf := &LeafNode{
		keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
		vals: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := *leaf
			err := res.Insert(c.k, c.v)

			if err != c.err {
				t.Errorf("expected error %v,  got %v", c.err, err)
			}
			if !slices.EqualFunc(res.keys, c.res.keys, bytes.Equal) {
				t.Errorf("expected keys %v,  got %v", c.res.keys, res.keys)
			}
			if !slices.EqualFunc(res.vals, c.res.vals, bytes.Equal) {
				t.Errorf("expected vals %v,  got %v", c.res.vals, res.vals)
			}
		})
	}
}

func TestLeafUpdate(t *testing.T) {
	cases := []struct {
		name string
		k    []byte
		v    []byte
		res  *LeafNode
		err  error
	}{
		{
			name: "update first",
			k:    []byte{'a'},
			v:    []byte{'0', 'a'},
			res: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
				vals: [][]byte{{'0', 'a'}, {'b'}, {'c'}, {'d'}},
			},
		},
		{
			name: "update last",
			k:    []byte{'d'},
			v:    []byte{'3', 'd'},
			res: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
				vals: [][]byte{{'a'}, {'b'}, {'c'}, {'3', 'd'}},
			},
		},
		{
			name: "update middle",
			k:    []byte{'c'},
			v:    []byte{'2', 'c'},
			res: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
				vals: [][]byte{{'a'}, {'b'}, {'2', 'c'}, {'d'}},
			},
		},
		{
			name: "failed update non existing key",
			k:    []byte{'e'},
			v:    []byte{'e'},
			res: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
				vals: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
			},
			err: errKeyNotExists,
		},
	}

	leaf := func() *LeafNode {
		return &LeafNode{
			keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
			vals: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := leaf()
			err := res.Update(c.k, c.v)

			if err != c.err {
				t.Errorf("expected error %v,  got %v", c.err, err)
			}
			if !slices.EqualFunc(res.keys, c.res.keys, bytes.Equal) {
				t.Errorf("expected keys %v,  got %v", c.res.keys, res.keys)
			}
			if !slices.EqualFunc(res.vals, c.res.vals, bytes.Equal) {
				t.Errorf("expected vals %v,  got %v", c.res.vals, res.vals)
			}
		})
	}
}

func TestLeafDelete(t *testing.T) {
	cases := []struct {
		name string
		k    []byte
		res  *LeafNode
		err  error
	}{
		{
			name: "delete first",
			k:    []byte{'a'},
			res: &LeafNode{
				keys: [][]byte{{'b'}, {'c'}, {'d'}},
				vals: [][]byte{{'b'}, {'c'}, {'d'}},
			},
		},
		{
			name: "delete last",
			k:    []byte{'d'},
			res: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}},
				vals: [][]byte{{'a'}, {'b'}, {'c'}},
			},
		},
		{
			name: "delete middle",
			k:    []byte{'c'},
			res: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'d'}},
				vals: [][]byte{{'a'}, {'b'}, {'d'}},
			},
		},
		{
			name: "failed delete non existing key",
			k:    []byte{'e'},
			res: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
				vals: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
			},
			err: errKeyNotExists,
		},
	}

	leaf := func() *LeafNode {
		return &LeafNode{
			keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
			vals: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := leaf()
			err := res.Delete(c.k)

			if err != c.err {
				t.Errorf("expected error %v,  got %v", c.err, err)
			}
			if !slices.EqualFunc(res.keys, c.res.keys, bytes.Equal) {
				t.Errorf("expected keys %v,  got %v", c.res.keys, res.keys)
			}
			if !slices.EqualFunc(res.vals, c.res.vals, bytes.Equal) {
				t.Errorf("expected vals %v,  got %v", c.res.vals, res.vals)
			}
		})
	}
}

func TestLeafSplit(t *testing.T) {
	cases := []struct {
		name string
		in   *LeafNode
		l    *LeafNode
		r    *LeafNode
	}{
		{
			name: "split even number",
			in: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}},
				vals: [][]byte{{'a'}, {'b'}},
			},
			l: &LeafNode{
				keys: [][]byte{{'a'}},
				vals: [][]byte{{'a'}},
			},
			r: &LeafNode{
				keys: [][]byte{{'b'}},
				vals: [][]byte{{'b'}},
			},
		},
		{
			name: "split odd number",
			in: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}},
				vals: [][]byte{{'a'}, {'b'}, {'c'}},
			},
			l: &LeafNode{
				keys: [][]byte{{'a'}},
				vals: [][]byte{{'a'}},
			},
			r: &LeafNode{
				keys: [][]byte{{'b'}, {'c'}},
				vals: [][]byte{{'b'}, {'c'}},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			outL, outR := c.in.Split()

			if l, ok := outL.(*LeafNode); !ok {
				t.Error("expected successful type assertion")
			} else {
				if !slices.EqualFunc(c.l.keys, l.keys, bytes.Equal) {
					t.Errorf("expected keys %v, got %v", c.l.keys, l.keys)
				} else if !slices.EqualFunc(c.l.vals, l.vals, bytes.Equal) {
					t.Errorf("expected vals %v, got %v", c.l.vals, l.vals)
				}
			}

			if r, ok := outR.(*LeafNode); !ok {
				t.Error("expected successful type assertion")
			} else {
				if !slices.EqualFunc(c.r.keys, r.keys, bytes.Equal) {
					t.Errorf("expected keys %v, got %v", c.r.keys, r.keys)
				} else if !slices.EqualFunc(c.r.vals, r.vals, bytes.Equal) {
					t.Errorf("expected vals %v, got %v", c.r.vals, r.vals)
				}
			}
		})
	}
}

func TestLeafMerge(t *testing.T) {
	cases := []struct {
		name  string
		in    *LeafNode
		merge Node
		res   *LeafNode
		err   error
	}{
		{
			name: "merge same size",
			in: &LeafNode{
				keys: [][]byte{{'a'}},
				vals: [][]byte{{'a'}},
			},
			merge: &LeafNode{
				keys: [][]byte{{'b'}},
				vals: [][]byte{{'b'}},
			},
			res: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}},
				vals: [][]byte{{'a'}, {'b'}},
			},
		},
		{
			name: "merge not same size",
			in: &LeafNode{
				keys: [][]byte{{'a'}},
				vals: [][]byte{{'a'}},
			},
			merge: &LeafNode{
				keys: [][]byte{{'b'}, {'c'}},
				vals: [][]byte{{'b'}, {'c'}},
			},
			res: &LeafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}},
				vals: [][]byte{{'a'}, {'b'}, {'c'}},
			},
		},
		{
			name: "failed merge wrong order",
			in: &LeafNode{
				keys: [][]byte{{'b'}},
				vals: [][]byte{{'b'}},
			},
			merge: &LeafNode{
				keys: [][]byte{{'a'}},
				vals: [][]byte{{'a'}},
			},
			res: &LeafNode{
				keys: [][]byte{{'b'}},
				vals: [][]byte{{'b'}},
			},
			err: errMergeOrder,
		},
		{
			name: "failed merge wrong order equal key",
			in: &LeafNode{
				keys: [][]byte{{'b'}},
				vals: [][]byte{{'b'}},
			},
			merge: &LeafNode{
				keys: [][]byte{{'b'}},
				vals: [][]byte{{'a'}},
			},
			res: &LeafNode{
				keys: [][]byte{{'b'}},
				vals: [][]byte{{'b'}},
			},
			err: errMergeOrder,
		},
		{
			name: "failed merge wrong type",
			in: &LeafNode{
				keys: [][]byte{{'a'}},
				vals: [][]byte{{'a'}},
			},
			merge: &PointerNode{},
			res: &LeafNode{
				keys: [][]byte{{'a'}},
				vals: [][]byte{{'a'}},
			},
			err: errMergeType,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.in.Merge(c.merge)

			if err != c.err {
				t.Errorf("expected error %v, got %v", c.err, err)
			}

			if !slices.EqualFunc(c.res.keys, c.in.keys, bytes.Equal) {
				t.Errorf("expected keys %v, got %v", c.res.keys, c.in.keys)
			} else if !slices.EqualFunc(c.res.vals, c.in.vals, bytes.Equal) {
				t.Errorf("expected vals %v, got %v", c.res.vals, c.in.vals)
			}
		})
	}
}
