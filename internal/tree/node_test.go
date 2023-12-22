package tree

import "testing"

func TestNewNode(t *testing.T) {
	cases := []struct {
		name string
		in   []byte
		res  Node
		err  error
	}{
		{
			name: "new pointer node",
			in:   []byte{0x00, 0x64, 0x00, 0x00, 0x00, 0x00},
			res:  &PointerNode{},
		},
		{
			name: "new leaf node",
			in:   []byte{0x00, 0x65, 0x00, 0x00, 0x00, 0x00},
			res:  &LeafNode{},
		},
		{
			name: "invalid node type",
			in:   []byte{0x00, 0x00},
			res:  nil,
			err:  errInvalNodeType,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := NewNode(c.in)

			if err != c.err {
				t.Errorf("expected error %v, got %v", c.err, err)
			} else if c.res != nil {
				if res.Type() != c.res.Type() {
					t.Errorf("expected node type %v, got %v", c.res.Type(), res.Type())
				}
			}
		})
	}
}
