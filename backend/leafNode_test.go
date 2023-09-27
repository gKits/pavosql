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
