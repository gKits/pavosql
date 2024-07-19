package node

import (
	"encoding/binary"
	"errors"
)

var (
	ErrIndexOutOfBounds  = errors.New("index is out of bounds")
	ErrCannotSplit       = errors.New("cannot split node with less than 2 entries")
	ErrCannotAppend      = errors.New("cannot append node with first key gte last key of current node")
	ErrWrongType         = errors.New("wrong type identifier")
	ErrNodeDataMalformed = errors.New("node data is malformed")
)

type Type uint16

const (
	typeInvalid Type = iota
	TypePointer
	TypeLeaf
)

func TypeOf(b []byte) Type {
	return Type(binary.LittleEndian.Uint16(b))
}
