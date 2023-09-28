package backend

import (
	"bytes"
	"testing"
)

func TestLeafNodeDecode(t *testing.T) {
	cases := []struct {
		name        string
		input       []byte
		expected    *leafNode
		expectedErr error
	}{
		{
			name: "Successful decoding",
			input: []byte{
				0x00, 0x01, // 1 is representation of lfNode nodeType
				0x00, 0x03, // nKeys is equal to 3
				0x00, 0x04, 0x00, 0x04, 'k', 'e', 'y', '1', 'v', 'a', 'l', '1', // first entry
				0x00, 0x04, 0x00, 0x04, 'k', 'e', 'y', '2', 'v', 'a', 'l', '2', // second entry
				0x00, 0x04, 0x00, 0x04, 'k', 'e', 'y', '3', 'v', 'a', 'l', '3', // third entry
			},
			expected: &leafNode{
				keys: [][]byte{{'k', 'e', 'y', '1'}, {'k', 'e', 'y', '2'}, {'k', 'e', 'y', '3'}},
				vals: [][]byte{{'v', 'a', 'l', '1'}, {'v', 'a', 'l', '2'}, {'v', 'a', 'l', '3'}},
			},
			expectedErr: nil,
		},
		{
			name: "Failed decoding due to wrong nodeType bytes",
			input: []byte{
				0x00, 0x02, // 2 is not the representation of lfNode nodeType
				0x00, 0x03,
				0x00, 0x04, 0x00, 0x04, 'k', 'e', 'y', '1', 'v', 'a', 'l', '1',
				0x00, 0x04, 0x00, 0x04, 'k', 'e', 'y', '2', 'v', 'a', 'l', '2',
				0x00, 0x04, 0x00, 0x04, 'k', 'e', 'y', '3', 'v', 'a', 'l', '3',
			},
			expected: &leafNode{
				keys: [][]byte{},
				vals: [][]byte{},
			},
			expectedErr: errNodeDecode,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ln := &leafNode{}

			err := ln.decode(c.input)

			if err != c.expectedErr {
				t.Errorf("Expected error %v, but got %v", c.expectedErr, err)
			}

			if len(ln.keys) != len(c.expected.keys) {
				t.Errorf("Expected %v keys, but got %v", len(c.expected.keys), len(ln.keys))

				for i, exp := range c.expected.keys {
					if !bytes.Equal(ln.keys[i], exp) {
						t.Errorf("Expected key %v at index %v, but got %v", exp, i, ln.keys[i])
					}
				}
			}

			if len(ln.vals) != len(c.expected.vals) {
				t.Errorf("Expected %v vals, but got %v", len(c.expected.vals), len(ln.vals))

				for i, exp := range c.expected.vals {
					if !bytes.Equal(ln.vals[i], exp) {
						t.Errorf("Expected val %v at index %v, but got %v", exp, i, ln.vals[i])
					}
				}
			}
		})
	}
}

func TestLeafNodeTyp(t *testing.T) {
	ln := &leafNode{}

	if ln.typ() != lfNode {
		t.Errorf("Expected type %v, but got %v", lfNode, ln.typ())
	}
}

func TestLeafNodeEncode(t *testing.T) {
	cases := []struct {
		name     string
		input    *leafNode
		expected []byte
	}{
		{
			name: "Successful encoding",
			input: &leafNode{
				keys: [][]byte{{'k', 'e', 'y', '1'}, {'k', 'e', 'y', '2'}, {'k', 'e', 'y', '3'}},
				vals: [][]byte{{'v', 'a', 'l', '1'}, {'v', 'a', 'l', '2'}, {'v', 'a', 'l', '3'}},
			},
			expected: []byte{
				0x00, 0x01, // 1 is representation of lfNode nodeType
				0x00, 0x03, // nKeys is equal to 3
				0x00, 0x04, 0x00, 0x04, 'k', 'e', 'y', '1', 'v', 'a', 'l', '1', // first entry
				0x00, 0x04, 0x00, 0x04, 'k', 'e', 'y', '2', 'v', 'a', 'l', '2', // second entry
				0x00, 0x04, 0x00, 0x04, 'k', 'e', 'y', '3', 'v', 'a', 'l', '3', // third entry
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

func TestLeafNodeSize(t *testing.T) {
	cases := []struct {
		name     string
		input    *leafNode
		expected int
	}{
		{
			name: "Size calculation",
			input: &leafNode{
				keys: [][]byte{{'k', 'e', 'y', '1'}, {'k', 'e', 'y', '2'}, {'k', 'e', 'y', '3'}},
				vals: [][]byte{{'v', 'a', 'l', '1'}, {'v', 'a', 'l', '2'}, {'v', 'a', 'l', '3'}},
			},
			expected: 40,
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

func TestLeafNodeKey(t *testing.T) {
	ln := &leafNode{
		keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}, {'f'}},
		vals: [][]byte{{'1'}, {'2'}, {'3'}, {'4'}, {'5'}, {'6'}},
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
			res, err := ln.key(c.input)

			if err != c.expectedErr {
				t.Errorf("Expected error %v, but got %v", c.expectedErr, err)
			}

			if !bytes.Equal(res, c.expected) {
				t.Errorf("Expected key %v, but got %v", c.expected, res)
			}
		})
	}
}

func TestLeafNodeSearch(t *testing.T) {
	cases := []struct {
		name           string
		input          []byte
		ln             *leafNode
		expected       int
		expectedExists bool
	}{
		{
			name:  "Search before first key in odd amount of keys",
			input: []byte{'a'},
			ln: &leafNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}},
				vals: [][]byte{{}, {}, {}, {}, {}},
			},
			expected:       0,
			expectedExists: false,
		},
		{
			name:  "Search before first key in even amount of keys",
			input: []byte{'a'},
			ln: &leafNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}, {'l'}},
				vals: [][]byte{{}, {}, {}, {}, {}, {}},
			},
			expected:       0,
			expectedExists: false,
		},
		{
			name:  "Search after last key in odd amount of keys",
			input: []byte{'k'},
			ln: &leafNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}},
				vals: [][]byte{{}, {}, {}, {}, {}},
			},
			expected:       5,
			expectedExists: false,
		},
		{
			name:  "Search after last key in even amount of keys",
			input: []byte{'m'},
			ln: &leafNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}, {'l'}},
				vals: [][]byte{{}, {}, {}, {}, {}, {}},
			},
			expected:       6,
			expectedExists: false,
		},
		{
			name:  "Search existing key in odd ammount of keys",
			input: []byte{'h'},
			ln: &leafNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}},
				vals: [][]byte{{}, {}, {}, {}, {}},
			},
			expected:       3,
			expectedExists: true,
		},
		{
			name:  "Search existing key in even ammount of keys",
			input: []byte{'j'},
			ln: &leafNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}, {'l'}},
				vals: [][]byte{{}, {}, {}, {}, {}, {}},
			},
			expected:       4,
			expectedExists: true,
		},
		{
			name:  "Search non-existing key in odd ammount of keys",
			input: []byte{'c'},
			ln: &leafNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}},
				vals: [][]byte{{}, {}, {}, {}, {}},
			},
			expected:       1,
			expectedExists: false,
		},
		{
			name:  "Search non-existing key in even ammount of keys",
			input: []byte{'c'},
			ln: &leafNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}, {'l'}},
				vals: [][]byte{{}, {}, {}, {}, {}, {}},
			},
			expected:       1,
			expectedExists: false,
		},
		{
			name:  "Search first key in odd ammount of keys",
			input: []byte{'b'},
			ln: &leafNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}},
				vals: [][]byte{{}, {}, {}, {}, {}},
			},
			expected:       0,
			expectedExists: true,
		},
		{
			name:  "Search first key in even ammount of keys",
			input: []byte{'b'},
			ln: &leafNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}, {'l'}},
				vals: [][]byte{{}, {}, {}, {}, {}, {}},
			},
			expected:       0,
			expectedExists: true,
		},
		{
			name:  "Search last key in odd ammount of keys",
			input: []byte{'j'},
			ln: &leafNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}},
				vals: [][]byte{{}, {}, {}, {}, {}},
			},
			expected:       4,
			expectedExists: true,
		},
		{
			name:  "Search last key in even ammount of keys",
			input: []byte{'l'},
			ln: &leafNode{
				keys: [][]byte{{'b'}, {'d'}, {'f'}, {'h'}, {'j'}, {'l'}},
				vals: [][]byte{{}, {}, {}, {}, {}, {}},
			},
			expected:       5,
			expectedExists: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, exists := c.ln.search(c.input)

			if exists != c.expectedExists {
				t.Errorf("Expected the key existing %v, but got %v", c.expectedExists, exists)
			}

			if res != c.expected {
				t.Errorf("Expected index %v, but got %v", c.expected, res)
			}
		})
	}
}

func TestLeafNodeMerge(t *testing.T) {
	cases := []struct {
		name        string
		left        *leafNode
		right       *leafNode
		expected    *leafNode
		expectedErr error
	}{
		{
			name: "Successful merge",
			left: &leafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}},
				vals: [][]byte{{}, {}, {}},
			},
			right: &leafNode{
				keys: [][]byte{{'d'}, {'e'}, {'f'}},
				vals: [][]byte{{}, {}, {}},
			},
			expected: &leafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}, {'f'}},
				vals: [][]byte{{}, {}, {}, {}, {}, {}},
			},
			expectedErr: nil,
		},
		{
			name: "Failed merge due to non-ordered keys",
			left: &leafNode{
				keys: [][]byte{{'a'}, {'b'}, {'d'}},
				vals: [][]byte{{}, {}, {}},
			},
			right: &leafNode{
				keys: [][]byte{{'c'}, {'e'}, {'f'}},
				vals: [][]byte{{}, {}, {}},
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
				ln := res.(*leafNode)

				if len(ln.keys) != len(c.expected.keys) {
					t.Errorf("Expected %v keys, but got %v", len(c.expected.keys), len(ln.keys))

					for i, exp := range c.expected.keys {
						if !bytes.Equal(ln.keys[i], exp) {
							t.Errorf("Expected key %v at index %v, but got %v", exp, i, ln.keys[i])
						}
					}
				}

				if len(ln.vals) != len(c.expected.vals) {
					t.Errorf("Expected %v vals, but got %v", len(c.expected.vals), len(ln.vals))

					for i, exp := range c.expected.vals {
						if !bytes.Equal(ln.vals[i], exp) {
							t.Errorf("Expected val %v at index %v, but got %v", exp, i, ln.vals[i])
						}
					}
				}
			} else if res == nil && c.expected != nil || res != nil && c.expected == nil {
				t.Errorf("Expected %v, but got %v", c.expected, res)
			}
		})
	}
}

func TestLeafNodeSplit(t *testing.T) {
	cases := []struct {
		name  string
		ln    *leafNode
		left  *leafNode
		right *leafNode
	}{
		{
			name: "Split even number of same size",
			ln: &leafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}, {'f'}},
				vals: [][]byte{{'1'}, {'2'}, {'3'}, {'4'}, {'5'}, {'6'}},
			},
			left: &leafNode{
				keys: [][]byte{{'a'}, {'b'}, {'c'}},
				vals: [][]byte{{'1'}, {'2'}, {'3'}},
			},
			right: &leafNode{
				keys: [][]byte{{'d'}, {'e'}, {'f'}},
				vals: [][]byte{{'4'}, {'5'}, {'6'}},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resL, resR := c.ln.split()

			left := resL.(*leafNode)
			right := resR.(*leafNode)

			if len(left.keys) != len(c.left.keys) {
				t.Errorf("Expected %v keys, but got %v", len(c.left.keys), len(left.keys))

				for i, exp := range c.left.keys {
					if !bytes.Equal(left.keys[i], exp) {
						t.Errorf("Expected key %v at index %v, but got %v", exp, i, left.keys[i])
					}
				}
			}

			if len(left.vals) != len(c.left.vals) {
				t.Errorf("left %v vals, but got %v", len(c.left.vals), len(left.vals))

				for i, exp := range c.left.vals {
					if !bytes.Equal(left.vals[i], exp) {
						t.Errorf("Expected val %v at index %v, but got %v", exp, i, left.vals[i])
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

			if len(right.vals) != len(c.right.vals) {
				t.Errorf("right %v vals, but got %v", len(c.right.vals), len(right.vals))

				for i, exp := range c.right.vals {
					if !bytes.Equal(right.vals[i], exp) {
						t.Errorf("Expected val %v at index %v, but got %v", exp, i, right.vals[i])
					}
				}
			}
		})
	}
}
