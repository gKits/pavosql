package stack

import (
	"slices"
	"testing"
)

func TestPush(t *testing.T) {
	cases := []struct {
		name string
		s    Stack[int]
		in   int
		res  Stack[int]
	}{
		{
			name: "successful push",
			s:    Stack[int]{0, 1, 2, 3},
			in:   4,
			res:  Stack[int]{0, 1, 2, 3, 4},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.s.Push(c.in)
			if !slices.Equal(c.s, c.res) {
				t.Errorf("expected stack %v, got %v", c.res, c.s)
			}
		})
	}
}

func TestPop(t *testing.T) {
	cases := []struct {
		name string
		in   Stack[int]
		res  int
		out  Stack[int]
		err  error
	}{
		{
			name: "successful pop",
			in:   Stack[int]{0, 1, 2, 3, 4},
			res:  4,
			out:  Stack[int]{0, 1, 2, 3},
		},
		{
			name: "failed pop",
			in:   Stack[int]{},
			res:  0,
			out:  Stack[int]{},
			err:  ErrStackEmpty,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := c.in.Pop()
			if err != c.err {
				t.Errorf("expected error %v, got %v", c.err, err)
			} else if !slices.Equal(c.in, c.out) {
				t.Errorf("expected stack %v, got %v", c.out, c.in)
			} else if res != c.res {
				t.Errorf("expected res %v, got %v", c.res, res)
			}
		})
	}
}
