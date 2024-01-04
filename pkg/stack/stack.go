package stack

import "errors"

type Stack[T any] []T

var ErrStackEmpty = errors.New("cannot pop, stack is empty")

func (s Stack[T]) Len() int {
	return len(s)
}

func (s *Stack[T]) Push(t T) {
	*s = append(*s, t)
}

func (s *Stack[T]) Pop() (res T, err error) {
	if s.Len() < 1 {
		return res, ErrStackEmpty
	}
	res = (*s)[s.Len()-1]
	*s = (*s)[:s.Len()-1]
	return res, nil
}

func (s Stack[T]) Peek() (res T, err error) {
	if s.Len() < 1 {
		return res, ErrStackEmpty
	}
	return s[s.Len()-1], nil
}
