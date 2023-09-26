package core

type dbType uint32

const (
	type_error = iota
	type_bytes
	type_int64
)

type value struct {
	typ dbType
}

type row struct {
	cols []string
	vals []value
}
