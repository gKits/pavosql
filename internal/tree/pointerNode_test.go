package tree

import (
	"bytes"
	"slices"
	"testing"
)

func TestPointerDecode(t *testing.T) {
	cases := []struct {
		name string
		in   []byte
		res  *PointerNode
		err  error
	}{
		{
			name: "invalid type",
			in:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			res:  &PointerNode{},
			err:  errInvalNodeType,
		},
		{
			name: "decode pointer node",
			in:   []byte{0x00, 0x64, 0x00, 0x00, 0x00, 0x00},
			res:  &PointerNode{keys: [][]byte{}, ptrs: []uint64{}},
			err:  nil,
		},
		{
			name: "decode pointer node with one k-p pair",
			in: []byte{
				0x00, 0x64,
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x01,
				'k',
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
			},
			res: &PointerNode{keys: [][]byte{{'k'}}, ptrs: []uint64{72340172838076673}},
			err: nil,
		},
		{
			name: "decode pointer node with multiple k-p pairs",
			in: []byte{
				0x00, 0x64,
				0x00, 0x00, 0x00, 0x04,
				0x00, 0x01,
				'a',
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x00, 0x01,
				'b',
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff,
				0x00, 0x02,
				'b', 'a',
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x03,
				'c', 'b', 'a',
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
			res: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'b', 'a'}, {'c', 'b', 'a'}},
				ptrs: []uint64{
					72340172838076673,
					255,
					1,
					18446744073709551615,
				},
			},
			err: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := &PointerNode{}
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
				} else if !slices.Equal(res.ptrs, c.res.ptrs) {
					t.Errorf("expected vals %v,  got %v", c.res.ptrs, res.ptrs)
				}
			}
		})
	}
}

func TestPointerEncode(t *testing.T) {
	cases := []struct {
		name string
		in   *PointerNode
		res  []byte
	}{
		{
			name: "decode pointer node",
			in:   &PointerNode{keys: [][]byte{}, ptrs: []uint64{}},
			res:  []byte{0x00, 0x64, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name: "decode pointer node with one k-p pair",
			in:   &PointerNode{keys: [][]byte{{'k'}}, ptrs: []uint64{72340172838076673}},
			res: []byte{
				0x00, 0x64,
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x01,
				'k',
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
			},
		},
		{
			name: "decode pointer node with multiple k-p pairs",
			in: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'b', 'a'}, {'c', 'b', 'a'}},
				ptrs: []uint64{
					72340172838076673,
					255,
					1,
					18446744073709551615,
				},
			},
			res: []byte{
				0x00, 0x64,
				0x00, 0x00, 0x00, 0x04,
				0x00, 0x01,
				'a',
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x00, 0x01,
				'b',
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff,
				0x00, 0x02,
				'b', 'a',
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x03,
				'c', 'b', 'a',
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
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

func TestPointerPtr(t *testing.T) {
	cases := []struct {
		name string
		k    []byte
		res  uint64
		err  error
	}{
		{
			name: "successfully read key",
			k:    []byte{'b'},
			res:  1,
		},
		{
			name: "failed read non existing key",
			k:    []byte{'e'},
			res:  0,
			err:  errKeyNotExists,
		},
	}

	ptr := PointerNode{
		keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
		ptrs: []uint64{0, 1, 2, 3},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := ptr.Ptr(c.k)

			if err != c.err {
				t.Errorf("expected error %v,  got %v", c.err, err)
			} else if res != c.res {
				t.Errorf("expected pointer %v, got %v", c.res, res)
			}
		})
	}
}

func TestPointerInsert(t *testing.T) {
	cases := []struct {
		name string
		k    []byte
		p    uint64
		res  *PointerNode
		err  error
	}{
		{
			name: "insert in middle",
			k:    []byte{'b', 'a'},
			p:    4,
			res: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'b', 'a'}, {'c'}, {'d'}},
				ptrs: []uint64{0, 1, 4, 2, 3},
			},
		},
		{
			name: "insert last",
			k:    []byte{'e'},
			p:    4,
			res: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}},
				ptrs: []uint64{0, 1, 2, 3, 4},
			},
		},
		{
			name: "failed insert existing key",
			k:    []byte{'a'},
			p:    4,
			res: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
				ptrs: []uint64{0, 1, 2, 3},
			},
			err: errKeyExists,
		},
	}

	ptr := &PointerNode{
		keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
		ptrs: []uint64{0, 1, 2, 3},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := *ptr
			err := res.Insert(c.k, c.p)

			if err != c.err {
				t.Errorf("expected error %v,  got %v", c.err, err)
			}
			if !slices.EqualFunc(res.keys, c.res.keys, bytes.Equal) {
				t.Errorf("expected keys %v,  got %v", c.res.keys, res.keys)
			}
			if !slices.Equal(res.ptrs, c.res.ptrs) {
				t.Errorf("expected ptrs %v,  got %v", c.res.ptrs, res.ptrs)
			}
		})
	}
}

func TestPointerUpdate(t *testing.T) {
	cases := []struct {
		name string
		k    []byte
		p    uint64
		res  *PointerNode
		err  error
	}{
		{
			name: "update first",
			k:    []byte{'a'},
			p:    4,
			res: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
				ptrs: []uint64{4, 1, 2, 3},
			},
		},
		{
			name: "update last",
			k:    []byte{'d'},
			p:    4,
			res: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
				ptrs: []uint64{0, 1, 2, 4},
			},
		},
		{
			name: "update middle",
			k:    []byte{'c'},
			p:    4,
			res: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
				ptrs: []uint64{0, 1, 4, 3},
			},
		},
		{
			name: "failed update non existing key",
			k:    []byte{'e'},
			p:    4,
			res: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
				ptrs: []uint64{0, 1, 2, 3},
			},
			err: errKeyNotExists,
		},
	}

	ptr := func() *PointerNode {
		return &PointerNode{
			keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
			ptrs: []uint64{0, 1, 2, 3},
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := ptr()
			err := res.Update(c.k, c.p)

			if err != c.err {
				t.Errorf("expected error %v,  got %v", c.err, err)
			}
			if !slices.EqualFunc(res.keys, c.res.keys, bytes.Equal) {
				t.Errorf("expected keys %v,  got %v", c.res.keys, res.keys)
			}
			if !slices.Equal(res.ptrs, c.res.ptrs) {
				t.Errorf("expected ptrs %v,  got %v", c.res.ptrs, res.ptrs)
			}
		})
	}
}

func TestPointerDelete(t *testing.T) {
	cases := []struct {
		name string
		k    []byte
		res  *PointerNode
		err  error
	}{
		{
			name: "delete first",
			k:    []byte{'a'},
			res: &PointerNode{
				keys: [][]byte{{'b'}, {'c'}, {'d'}},
				ptrs: []uint64{1, 2, 3},
			},
		},
		{
			name: "delete last",
			k:    []byte{'d'},
			res: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}},
				ptrs: []uint64{0, 1, 2},
			},
		},
		{
			name: "delete middle",
			k:    []byte{'c'},
			res: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'d'}},
				ptrs: []uint64{0, 1, 3},
			},
		},
		{
			name: "failed delete non existing key",
			k:    []byte{'e'},
			res: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
				ptrs: []uint64{0, 1, 2, 3},
			},
			err: errKeyNotExists,
		},
	}

	leaf := func() *PointerNode {
		return &PointerNode{
			keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}},
			ptrs: []uint64{0, 1, 2, 3},
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
			if !slices.Equal(res.ptrs, c.res.ptrs) {
				t.Errorf("expected ptrs %v,  got %v", c.res.ptrs, res.ptrs)
			}
		})
	}
}

func TestPointerSplit(t *testing.T) {
	cases := []struct {
		name string
		in   *PointerNode
		l    *PointerNode
		r    *PointerNode
	}{
		{
			name: "split even number",
			in: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}},
				ptrs: []uint64{0, 1},
			},
			l: &PointerNode{
				keys: [][]byte{{'a'}},
				ptrs: []uint64{0},
			},
			r: &PointerNode{
				keys: [][]byte{{'b'}},
				ptrs: []uint64{1},
			},
		},
		{
			name: "split odd number",
			in: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}},
				ptrs: []uint64{0, 1, 2},
			},
			l: &PointerNode{
				keys: [][]byte{{'a'}},
				ptrs: []uint64{0},
			},
			r: &PointerNode{
				keys: [][]byte{{'b'}, {'c'}},
				ptrs: []uint64{1, 2},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			outL, outR := c.in.Split()

			if l, ok := outL.(*PointerNode); !ok {
				t.Error("expected successful type assertion")
			} else {
				if !slices.EqualFunc(c.l.keys, l.keys, bytes.Equal) {
					t.Errorf("expected keys %v, got %v", c.l.keys, l.keys)
				} else if !slices.Equal(c.l.ptrs, l.ptrs) {
					t.Errorf("expected ptrs %v, got %v", c.l.ptrs, l.ptrs)
				}
			}

			if r, ok := outR.(*PointerNode); !ok {
				t.Error("expected successful type assertion")
			} else {
				if !slices.EqualFunc(c.r.keys, r.keys, bytes.Equal) {
					t.Errorf("expected keys %v, got %v", c.r.keys, r.keys)
				} else if !slices.Equal(c.r.ptrs, r.ptrs) {
					t.Errorf("expected ptrs %v, got %v", c.r.ptrs, r.ptrs)
				}
			}
		})
	}
}

func TestPointerMerge(t *testing.T) {
	cases := []struct {
		name  string
		in    *PointerNode
		merge Node
		res   *PointerNode
		err   error
	}{
		{
			name: "merge same size",
			in: &PointerNode{
				keys: [][]byte{{'a'}},
				ptrs: []uint64{0},
			},
			merge: &PointerNode{
				keys: [][]byte{{'b'}},
				ptrs: []uint64{1},
			},
			res: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}},
				ptrs: []uint64{0, 1},
			},
		},
		{
			name: "merge not same size",
			in: &PointerNode{
				keys: [][]byte{{'a'}},
				ptrs: []uint64{0},
			},
			merge: &PointerNode{
				keys: [][]byte{{'b'}, {'c'}},
				ptrs: []uint64{1, 2},
			},
			res: &PointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}},
				ptrs: []uint64{0, 1, 2},
			},
		},
		{
			name: "failed merge wrong order",
			in: &PointerNode{
				keys: [][]byte{{'b'}},
				ptrs: []uint64{1},
			},
			merge: &PointerNode{
				keys: [][]byte{{'a'}},
				ptrs: []uint64{0},
			},
			res: &PointerNode{
				keys: [][]byte{{'b'}},
				ptrs: []uint64{1},
			},
			err: errMergeOrder,
		},
		{
			name: "failed merge wrong order equal key",
			in: &PointerNode{
				keys: [][]byte{{'b'}},
				ptrs: []uint64{1},
			},
			merge: &PointerNode{
				keys: [][]byte{{'b'}},
				ptrs: []uint64{0},
			},
			res: &PointerNode{
				keys: [][]byte{{'b'}},
				ptrs: []uint64{1},
			},
			err: errMergeOrder,
		},
		{
			name: "failed merge wrong type",
			in: &PointerNode{
				keys: [][]byte{{'a'}},
				ptrs: []uint64{0},
			},
			merge: &LeafNode{},
			res: &PointerNode{
				keys: [][]byte{{'a'}},
				ptrs: []uint64{0},
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
			} else if !slices.Equal(c.res.ptrs, c.in.ptrs) {
				t.Errorf("expected ptrs %v, got %v", c.res.ptrs, c.in.ptrs)
			}
		})
	}
}
