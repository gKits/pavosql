package db

import (
	"encoding/binary"
	"fmt"
)

const (
	errCellWrongTypeMsg = "cell: expected type '%v', but got '%v'"
	errCellWrongLenMsg  = "cell: expected length '%v', but got '%v'"
)

type Cell struct {
	typ dbType
	b   []byte
}

func (c *Cell) Type() dbType {
	return c.typ
}

func (c *Cell) Raw() []byte {
	return c.b
}

func (c *Cell) Bytes() ([]byte, error) {
	if c.typ != dbBytes {
		return nil, fmt.Errorf(errCellWrongTypeMsg, dbBytes, c.typ)
	}

	l := int(binary.BigEndian.Uint16(c.b))

	if l != len(c.b[2:]) {
		return nil, fmt.Errorf(errCellWrongLenMsg, l, len(c.b[2:]))
	}

	return c.b[2:], nil
}

func (c *Cell) Bool() (bool, error) {
	if c.typ != dbBool {
		return false, fmt.Errorf(errCellWrongTypeMsg, dbBool, c.typ)
	}

	if len(c.b) != 1 {
		return false, fmt.Errorf(errCellWrongLenMsg, 1, len(c.b))
	}

	n := uint8(c.b[0])
	if n > 1 {
		return false, nil
	}

	return n == 1, nil
}

func (c *Cell) String() (string, error) {
	if c.typ != dbString {
		return "", fmt.Errorf(errCellWrongTypeMsg, dbString, c.typ)
	}

	l := int(c.b[0])

	if l != len(c.b[2:]) {
		return "", fmt.Errorf(errCellWrongLenMsg, l, len(c.b[2:]))
	}

	return string(c.b[2 : 2+l]), nil
}

func (c *Cell) RawString() (string, error) {
	if c.typ != dbRawString {
		return "", fmt.Errorf(errCellWrongTypeMsg, dbRawString, c.typ)
	}

	l := int(c.b[0])

	if l != len(c.b[2:]) {
		return "", fmt.Errorf(errCellWrongLenMsg, l, len(c.b[2:]))
	}

	return string(c.b[2 : 2+l]), nil
}

func (c *Cell) Int8() (int8, error) {
	if c.typ != dbInt8 {
		return 0, fmt.Errorf(errCellWrongTypeMsg, dbInt8, c.typ)
	}

	n := c.b[0]
	return int8(n), nil
}

func (c *Cell) Uint8() (uint8, error) {
	if c.typ != dbUint8 {
		return 0, fmt.Errorf(errCellWrongTypeMsg, dbUint8, c.typ)
	}

	n := c.b[0]
	return n, nil
}

func (c *Cell) Int16() (int16, error) {
	if c.typ != dbInt16 {
		return 0, fmt.Errorf(errCellWrongTypeMsg, dbInt16, c.typ)
	}

	n := binary.BigEndian.Uint16(c.b)
	return int16(n), nil
}

func (c *Cell) Uint16() (uint16, error) {
	if c.typ != dbUint16 {
		return 0, fmt.Errorf(errCellWrongTypeMsg, dbUint16, c.typ)
	}

	n := binary.BigEndian.Uint16(c.b)
	return n, nil
}

func (c *Cell) Int32() (int32, error) {
	if c.typ != dbInt32 {
		return 0, fmt.Errorf(errCellWrongTypeMsg, dbInt32, c.typ)
	}

	n := binary.BigEndian.Uint32(c.b)
	return int32(n), nil
}

func (c *Cell) Uint32() (uint32, error) {
	if c.typ != dbUint32 {
		return 0, fmt.Errorf(errCellWrongTypeMsg, dbUint32, c.typ)
	}

	n := binary.BigEndian.Uint32(c.b)
	return n, nil
}

func (c *Cell) Int64() (int64, error) {
	if c.typ != dbInt64 {
		return 0, fmt.Errorf(errCellWrongTypeMsg, dbInt64, c.typ)
	}

	n := binary.BigEndian.Uint64(c.b)
	return int64(n), nil
}

func (c *Cell) Uint64() (uint64, error) {
	if c.typ != dbUint64 {
		return 0, fmt.Errorf(errCellWrongTypeMsg, dbUint64, c.typ)
	}

	n := binary.BigEndian.Uint64(c.b)
	return n, nil
}

func (c *Cell) Float32() (float32, error) {
	if c.typ != dbFloat32 {
		return 0, fmt.Errorf(errCellWrongTypeMsg, dbFloat32, c.typ)
	}

	n := binary.BigEndian.Uint32(c.b)
	return float32(n), nil
}

func (c *Cell) Float64() (float64, error) {
	if c.typ != dbFloat64 {
		return 0, fmt.Errorf(errCellWrongTypeMsg, dbFloat64, c.typ)
	}

	n := binary.BigEndian.Uint64(c.b)
	return float64(n), nil
}
