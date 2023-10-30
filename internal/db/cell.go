package db

import (
	"encoding/binary"
	"fmt"
)

const errCellMsg = "cell: cannot interpret '%x' as type '%s'"

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
	if c.typ != dbBytes || int(c.b[0]) != len(c.b[2:]) {
		return nil, fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	return c.b[2:], nil
}

func (c *Cell) Bool() (bool, error) {
	if c.typ != dbBool || len(c.b) != 1 {
		return false, fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	n := uint8(c.b[0])
	return n == 1, nil
}

func (c *Cell) String() (string, error) {
	if c.typ != dbString || int(c.b[0]) != len(c.b[2:]) {
		return "", fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	return string(c.b[2 : 2+int(c.b[0])]), nil
}

func (c *Cell) RawString() (string, error) {
	if c.typ != dbRawString || int(c.b[0]) != len(c.b[2:]) {
		return "", fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	return string(c.b[2 : 2+int(c.b[0])]), nil
}

func (c *Cell) Int8() (int8, error) {
	if c.typ != dbInt8 {
		return 0, fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	return int8(c.b[0]), nil
}

func (c *Cell) Uint8() (uint8, error) {
	if c.typ != dbUint8 {
		return 0, fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	return c.b[0], nil
}

func (c *Cell) Int16() (int16, error) {
	if c.typ != dbInt16 {
		return 0, fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	n := binary.BigEndian.Uint16(c.b)
	return int16(n), nil
}

func (c *Cell) Uint16() (uint16, error) {
	if c.typ != dbUint16 {
		return 0, fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	n := binary.BigEndian.Uint16(c.b)
	return n, nil
}

func (c *Cell) Int32() (int32, error) {
	if c.typ != dbInt32 {
		return 0, fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	n := binary.BigEndian.Uint32(c.b)
	return int32(n), nil
}

func (c *Cell) Uint32() (uint32, error) {
	if c.typ != dbUint32 {
		return 0, fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	n := binary.BigEndian.Uint32(c.b)
	return n, nil
}

func (c *Cell) Int64() (int64, error) {
	if c.typ != dbInt64 {
		return 0, fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	n := binary.BigEndian.Uint64(c.b)
	return int64(n), nil
}

func (c *Cell) Uint64() (uint64, error) {
	if c.typ != dbUint64 {
		return 0, fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	n := binary.BigEndian.Uint64(c.b)
	return n, nil
}

func (c *Cell) Float32() (float32, error) {
	if c.typ != dbFloat32 {
		return 0, fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	n := binary.BigEndian.Uint32(c.b)
	return float32(n), nil
}

func (c *Cell) Float64() (float64, error) {
	if c.typ != dbFloat64 {
		return 0, fmt.Errorf(errCellMsg, c.b, c.typ)
	}

	n := binary.BigEndian.Uint64(c.b)
	return float64(n), nil
}
