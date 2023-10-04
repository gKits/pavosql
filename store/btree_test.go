package store

import (
	"bytes"
	"testing"
)

var mockStore map[uint64]node
var lastPtr uint64 = 0

func mockGetPage(ptr uint64) (node, error) {
	n, ok := mockStore[ptr]
	if !ok {
		return nil, errKVBadPtr
	}
	return n, nil
}

func mockPullPage(ptr uint64) (node, error) {
	n, err := mockGetPage(ptr)
	if err != nil {
		return nil, err
	}
	if err := mockFreePage(ptr); err != nil {
		return nil, err
	}

	return n, nil
}

func mockAllocPage(n node) (uint64, error) {
	lastPtr++
	mockStore[lastPtr] = n

	return lastPtr, nil
}

func mockFreePage(ptr uint64) error {
	if _, ok := mockStore[ptr]; !ok {
		return errKVBadPtr
	}

	delete(mockStore, ptr)
	return nil
}

func TestBTreeGet(t *testing.T) {
	cases := []struct {
		name        string
		bt          bTree
		input       []byte
		expected    []byte
		expectedErr error
	}{}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := c.bt.Get(c.input)

			if err != c.expectedErr {

			}

			if !bytes.Equal(c.expected, res) {

			}
		})
	}
}
