package backend

import (
	"crypto/rand"
	"slices"
	"testing"
)

func randBytes() []byte {
	b := make([]byte, 10)
	rand.Read(b)
	return b
}

func TestSize(t *testing.T) {

}

func TestDecode(t *testing.T) {
	cases := []struct {
		name string
		d    []byte
		leaf LeafNode
		err  error
	}{
		{
			"wrong type number throws error",
			[]byte{0, 0},
			LeafNode{},
			errNodeDecode,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			leaf := LeafNode{}
			err := leaf.Decode(c.d)

			if c.err != err {
				t.Errorf("unexpected error result: %s != %s", c.err, err)
			} else {

			}
		})
	}

}

func TestEncodeDecodeSize(t *testing.T) {
	val := []byte("test")

	cases := []struct {
		leaf LeafNode
	}{
		{
			LeafNode{
				[][]byte{val, val, val},
				[][]byte{val, val, val},
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			d := c.leaf.Encode()
			n := LeafNode{}
			n.Decode(d)

			nSize := n.Size()
			if len(d) != nSize {
				t.Errorf("sizes are not equal: %d != %d", len(d), nSize)
			}

			if !slices.EqualFunc(n.keys, c.leaf.keys, func(b1, b2 []byte) bool {
				return string(b1) == string(b2)
			}) {
				t.Errorf("keys are not equal: %v != %v", c.leaf.keys, n.keys)
			}

			if !slices.EqualFunc(n.vals, c.leaf.vals, func(b1, b2 []byte) bool {
				return string(b1) == string(b2)
			}) {
				t.Errorf("vals are not equal: %v != %v", c.leaf.vals, n.vals)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	val := []byte("test")

	keysEven := [][]byte{
		[]byte("b"), []byte("d"), []byte("f"), []byte("h"),
		[]byte("j"), []byte("l"), []byte("n"), []byte("p"),
		[]byte("r"), []byte("t"), []byte("v"), []byte("x"),
	}

	keysOdd := [][]byte{
		[]byte("b"), []byte("d"), []byte("f"), []byte("h"),
		[]byte("j"), []byte("l"), []byte("n"), []byte("p"),
		[]byte("r"), []byte("t"), []byte("v"),
	}

	valsEven := [][]byte{
		val, val, val, val,
		val, val, val, val,
		val, val, val, val,
	}

	valsOdd := [][]byte{
		val, val, val, val,
		val, val, val, val,
		val, val, val,
	}

	cases := []struct {
		name   string
		leaf   LeafNode
		search []byte
		idx    int
		exists bool
	}{
		{
			"Even numbers of keys + search key exists",
			LeafNode{
				keysEven,
				valsEven,
			},
			[]byte("d"),
			1,
			true,
		},
		{
			"Even numbers of keys + search key does not exist",
			LeafNode{
				keysEven,
				valsEven,
			},
			[]byte("e"),
			2,
			false,
		},
		{
			"Odd numbers of keys + search key exists",
			LeafNode{
				keysOdd,
				valsOdd,
			},
			[]byte("p"),
			7,
			true,
		},
		{
			"Odd numbers of keys + search key does not exist",
			LeafNode{
				keysOdd,
				valsOdd,
			},
			[]byte("u"),
			10,
			false,
		},
		{
			"key before first index",
			LeafNode{
				keysEven,
				valsEven,
			},
			[]byte("a"),
			0,
			false,
		},
		{
			"key after last index",
			LeafNode{
				keysOdd,
				valsOdd,
			},
			[]byte("z"),
			11,
			false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			idx, exists := c.leaf.Search(c.search)

			if c.idx != idx {
				t.Errorf("result idx does not match expected: %d != %d", c.idx, idx)
			}
			if c.exists != exists {
				t.Errorf("different existing status expected: %v != %v", c.exists, exists)
			}
		})
	}
}

func TestMerge(t *testing.T) {}

func TestSplit(t *testing.T) {}

func TestInsert(t *testing.T) {}

func TestUpdate(t *testing.T) {}

func TestDelete(t *testing.T) {}
