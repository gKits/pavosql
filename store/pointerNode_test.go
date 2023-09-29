package store

import (
	"bytes"
	"testing"
)

func TestPointerNodeDecode(t *testing.T) {
	cases := []struct {
		name        string
		input       []byte
		expected    *pointerNode
		expectedErr error
	}{
		{
			name: "Successful decoding",
			input: []byte{
				0x00, 0x00, // 1 is representation of ptrNode nodeType
				0x00, 0x03, // nKeys is equal to 3
				0x00, 0x04, 'k', 'e', 'y', '1', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // first entry
				0x00, 0x04, 'k', 'e', 'y', '2', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // second entry
				0x00, 0x04, 'k', 'e', 'y', '3', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // third entry
			},
			expected: &pointerNode{
				keys: [][]byte{{'k', 'e', 'y', '1'}, {'k', 'e', 'y', '2'}, {'k', 'e', 'y', '3'}},
				ptrs: []uint64{0, 1, 2},
			},
			expectedErr: nil,
		},
		{
			name: "Failed decoding due to wrong nodeType bytes",
			input: []byte{
				0x00, 0x02, // 2 is not the representation of ptrNode nodeType
				0x00, 0x03,
				0x00, 0x04, 'k', 'e', 'y', '1', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // first entry
				0x00, 0x04, 'k', 'e', 'y', '2', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // second entry
				0x00, 0x04, 'k', 'e', 'y', '3', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // third entry
			},
			expected: &pointerNode{
				keys: [][]byte{},
				ptrs: []uint64{},
			},
			expectedErr: errNodeDecode,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pn := &pointerNode{}

			err := pn.decode(c.input)

			if err != c.expectedErr {
				t.Errorf("Expected error %v, but got %v", c.expectedErr, err)
			}

			if len(pn.keys) != len(c.expected.keys) {
				t.Errorf("Expected %v keys, but got %v", len(c.expected.keys), len(pn.keys))

				for i, exp := range c.expected.keys {
					if !bytes.Equal(pn.keys[i], exp) {
						t.Errorf("Expected key %v at index %v, but got %v", exp, i, pn.keys[i])
					}
				}
			}

			if len(pn.ptrs) != len(c.expected.ptrs) {
				t.Errorf("Expected %v ptrs, but got %v", len(c.expected.ptrs), len(pn.ptrs))

				for i, exp := range c.expected.ptrs {
					if pn.ptrs[i] != exp {
						t.Errorf("Expected val %v at index %v, but got %v", exp, i, pn.ptrs[i])
					}
				}
			}
		})
	}
}

func TestPointerNodeTyp(t *testing.T) {
	pn := &pointerNode{}

	if pn.typ() != ptrNode {
		t.Errorf("Expected type %v, but got %v", ptrNode, pn.typ())
	}
}

func TestPointerNodeEncode(t *testing.T) {
	cases := []struct {
		name     string
		input    *pointerNode
		expected []byte
	}{
		{
			name: "Successful encoding",
			input: &pointerNode{
				keys: [][]byte{{'k', 'e', 'y', '1'}, {'k', 'e', 'y', '2'}, {'k', 'e', 'y', '3'}},
				ptrs: []uint64{0, 1, 2},
			},
			expected: []byte{
				0x00, 0x00, // 0 is representation of ptrNode nodeType
				0x00, 0x03, // nKeys is equal to 3
				0x00, 0x04, 'k', 'e', 'y', '1', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // first entry
				0x00, 0x04, 'k', 'e', 'y', '2', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, // second entry
				0x00, 0x04, 'k', 'e', 'y', '3', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // third entry
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := c.input.encode()

			if !bytes.Equal(res, c.expected) {
				t.Errorf("Expected %v, but got %v", c.expected, res)
			}
		})
	}
}

func TestPointerNodeSize(t *testing.T) {
	cases := []struct {
		name     string
		input    *pointerNode
		expected int
	}{
		{
			name: "Size calculation",
			input: &pointerNode{
				keys: [][]byte{{'k', 'e', 'y', '1'}, {'k', 'e', 'y', '2'}, {'k', 'e', 'y', '3'}},
				ptrs: []uint64{0, 1, 2},
			},
			expected: 46,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := c.input.size()

			if res != c.expected {
				t.Errorf("Expected node size %v, but got %v", c.expected, res)
			}
		})
	}
}

func TestPointerNodeKey(t *testing.T) {
	pn := &pointerNode{
		keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}, {'f'}},
		ptrs: []uint64{0, 1, 2, 3, 4, 5},
	}

	cases := []struct {
		name        string
		input       int
		expected    []byte
		expectedErr error
	}{
		{
			name:        "Key at index 0",
			input:       0,
			expected:    []byte{'a'},
			expectedErr: nil,
		},
		{
			name:        "Key at last index",
			input:       5,
			expected:    []byte{'f'},
			expectedErr: nil,
		},
		{
			name:        "Too large key",
			input:       6,
			expected:    nil,
			expectedErr: errNodeIdx,
		},
		{
			name:        "Negative key",
			input:       -1,
			expected:    nil,
			expectedErr: errNodeIdx,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := pn.key(c.input)

			if err != c.expectedErr {
				t.Errorf("Expected error %v, but got %v", c.expectedErr, err)
			}

			if !bytes.Equal(res, c.expected) {
				t.Errorf("Expected key %v, but got %v", c.expected, res)
			}
		})
	}
}

func TestPointerNodeSearch(t *testing.T) {
	cases := []struct {
		name           string
		input          []byte
		pn             *pointerNode
		expected       int
		expectedExists bool
	}{
		{
			name:  "Search before first key in odd amount of keys",
			input: []byte{'a'},
			pn: &pointerNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}},
				ptrs: []uint64{0, 0, 0, 0, 0},
			},
			expected:       0,
			expectedExists: false,
		},
		{
			name:  "Search before first key in even amount of keys",
			input: []byte{'a'},
			pn: &pointerNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}, {'l'}},
				ptrs: []uint64{0, 0, 0, 0, 0, 0},
			},
			expected:       0,
			expectedExists: false,
		},
		{
			name:  "Search after last key in odd amount of keys",
			input: []byte{'k'},
			pn: &pointerNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}},
				ptrs: []uint64{0, 0, 0, 0, 0, 0},
			},
			expected:       5,
			expectedExists: false,
		},
		{
			name:  "Search after last key in even amount of keys",
			input: []byte{'m'},
			pn: &pointerNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}, {'l'}},
				ptrs: []uint64{0, 0, 0, 0, 0, 0},
			},
			expected:       6,
			expectedExists: false,
		},
		{
			name:  "Search existing key in odd ammount of keys",
			input: []byte{'h'},
			pn: &pointerNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}},
				ptrs: []uint64{0, 0, 0, 0, 0, 0},
			},
			expected:       3,
			expectedExists: true,
		},
		{
			name:  "Search existing key in even ammount of keys",
			input: []byte{'j'},
			pn: &pointerNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}, {'l'}},
				ptrs: []uint64{0, 0, 0, 0, 0, 0},
			},
			expected:       4,
			expectedExists: true,
		},
		{
			name:  "Search non-existing key in odd ammount of keys",
			input: []byte{'c'},
			pn: &pointerNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}},
				ptrs: []uint64{0, 0, 0, 0, 0, 0},
			},
			expected:       1,
			expectedExists: false,
		},
		{
			name:  "Search non-existing key in even ammount of keys",
			input: []byte{'c'},
			pn: &pointerNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}, {'l'}},
				ptrs: []uint64{0, 0, 0, 0, 0, 0},
			},
			expected:       1,
			expectedExists: false,
		},
		{
			name:  "Search first key in odd ammount of keys",
			input: []byte{'b'},
			pn: &pointerNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}},
				ptrs: []uint64{0, 0, 0, 0, 0, 0},
			},
			expected:       0,
			expectedExists: true,
		},
		{
			name:  "Search first key in even ammount of keys",
			input: []byte{'b'},
			pn: &pointerNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}, {'l'}},
				ptrs: []uint64{0, 0, 0, 0, 0, 0},
			},
			expected:       0,
			expectedExists: true,
		},
		{
			name:  "Search last key in odd ammount of keys",
			input: []byte{'j'},
			pn: &pointerNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}},
				ptrs: []uint64{0, 0, 0, 0, 0, 0},
			},
			expected:       4,
			expectedExists: true,
		},
		{
			name:  "Search last key in even ammount of keys",
			input: []byte{'l'},
			pn: &pointerNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}, {'l'}},
				ptrs: []uint64{0, 0, 0, 0, 0, 0},
			},
			expected:       5,
			expectedExists: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, exists := c.pn.search(c.input)

			if exists != c.expectedExists {
				t.Errorf("Expected the key existing %v, but got %v", c.expectedExists, exists)
			}

			if res != c.expected {
				t.Errorf("Expected index %v, but got %v", c.expected, res)
			}
		})
	}
}

func TestPointerNodeMerge(t *testing.T) {
	cases := []struct {
		name        string
		left        *pointerNode
		right       *pointerNode
		expected    *pointerNode
		expectedErr error
	}{
		{
			name: "Successful merge",
			left: &pointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}},
				ptrs: []uint64{0, 0, 0},
			},
			right: &pointerNode{
				keys: [][]byte{{'d'}, {'e'}, {'f'}},
				ptrs: []uint64{0, 0, 0},
			},
			expected: &pointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}, {'f'}},
				ptrs: []uint64{0, 0, 0, 0, 0, 0},
			},
			expectedErr: nil,
		},
		{
			name: "Failed merge due to non-ordered keys",
			left: &pointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'d'}},
				ptrs: []uint64{0, 0, 0},
			},
			right: &pointerNode{
				keys: [][]byte{{'c'}, {'e'}, {'f'}},
				ptrs: []uint64{0, 0, 0},
			},
			expected:    nil,
			expectedErr: errNodeMerge,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := c.left.merge(c.right)

			if c.expectedErr != err {
				t.Errorf("Expected error %v, but got %v", c.expectedErr, err)
			}

			if res != nil && c.expected != nil {
				pn := res.(*pointerNode)
				t.Log(pn)

				if len(pn.keys) != len(c.expected.keys) {
					t.Errorf("Expected %v keys, but got %v", len(c.expected.keys), len(pn.keys))

					for i, exp := range c.expected.keys {
						if !bytes.Equal(pn.keys[i], exp) {
							t.Errorf("Expected key %v at index %v, but got %v", exp, i, pn.keys[i])
						}
					}
				}

				if len(pn.ptrs) != len(c.expected.ptrs) {
					t.Errorf("Expected %v ptrs, but got %v", len(c.expected.ptrs), len(pn.ptrs))

					for i, exp := range c.expected.ptrs {
						if pn.ptrs[i] != exp {
							t.Errorf("Expected val %v at index %v, but got %v", exp, i, pn.ptrs[i])
						}
					}
				}
			} else if res == nil && c.expected != nil || res != nil && c.expected == nil {
				t.Errorf("Expected %v, but got %v", c.expected, res)
			}
		})
	}
}

func TestPointerNodeSplit(t *testing.T) {
	cases := []struct {
		name  string
		pn    *pointerNode
		left  *pointerNode
		right *pointerNode
	}{
		{
			name: "Split even number of same size",
			pn: &pointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}, {'f'}},
				ptrs: []uint64{0, 1, 2, 3, 4, 5},
			},
			left: &pointerNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}},
				ptrs: []uint64{0, 1, 2},
			},
			right: &pointerNode{
				keys: [][]byte{{'d'}, {'e'}, {'f'}},
				ptrs: []uint64{3, 4, 5},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resL, resR := c.pn.split()

			left := resL.(*pointerNode)
			right := resR.(*pointerNode)

			if len(left.keys) != len(c.left.keys) {
				t.Errorf("Expected %v keys, but got %v", len(c.left.keys), len(left.keys))

				for i, exp := range c.left.keys {
					if !bytes.Equal(left.keys[i], exp) {
						t.Errorf("Expected key %v at index %v, but got %v", exp, i, left.keys[i])
					}
				}
			}

			if len(left.ptrs) != len(c.left.ptrs) {
				t.Errorf("left %v ptrs, but got %v", len(c.left.ptrs), len(left.ptrs))

				for i, exp := range c.left.ptrs {
					if left.ptrs[i] != exp {
						t.Errorf("Expected val %v at index %v, but got %v", exp, i, left.ptrs[i])
					}
				}
			}

			if len(right.keys) != len(c.right.keys) {
				t.Errorf("Expected %v keys, but got %v", len(c.right.keys), len(right.keys))

				for i, exp := range c.right.keys {
					if !bytes.Equal(right.keys[i], exp) {
						t.Errorf("Expected key %v at index %v, but got %v", exp, i, right.keys[i])
					}
				}
			}

			if len(right.ptrs) != len(c.right.ptrs) {
				t.Errorf("right %v ptrs, but got %v", len(c.right.ptrs), len(right.ptrs))

				for i, exp := range c.right.ptrs {
					if right.ptrs[i] != exp {
						t.Errorf("Expected val %v at index %v, but got %v", exp, i, right.ptrs[i])
					}
				}
			}
		})
	}
}
