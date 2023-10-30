package db

import "fmt"

type dbType uint8

const (
	dbBytes dbType = iota
	dbBool
	dbString
	dbRawString
	dbInt8
	dbUint8
	dbInt16
	dbUint16
	dbInt32
	dbUint32
	dbInt64
	dbUint64
	dbFloat32
	dbFloat64
)

func (dT dbType) String() string {
	switch dT {
	case dbBytes:
		return "Bytes"
	case dbBool:
		return "Bool"
	case dbString:
		return "String"
	case dbRawString:
		return "RawString"
	case dbInt8:
		return "Int8"
	case dbUint8:
		return "Uint8"
	case dbInt16:
		return "Int16"
	case dbUint16:
		return "Uint16"
	case dbInt32:
		return "Int32"
	case dbUint32:
		return "Uint32"
	case dbInt64:
		return "Int64"
	case dbUint64:
		return "Uint64"
	case dbFloat32:
		return "Float32"
	case dbFloat64:
		return "Float64"
	default:
		return fmt.Sprintf("%d", dT)
	}
}
