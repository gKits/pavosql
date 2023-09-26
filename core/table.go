package core

type table struct {
	name          string
	types         []uint32
	cols          []string
	pKeys         int
	prefix        uint32
	indices       [][]string
	indexPrefixes []uint32
}

var TABLE_META = table{
	name:   "@meta",
	prefix: 1,
	types:  []uint32{type_bytes, type_bytes},
	cols:   []string{"key", "val"},
	pKeys:  1,
}

var TABLE_TDEF = table{
	name:   "@tabledef",
	prefix: 2,
	types:  []uint32{type_bytes, type_bytes},
	cols:   []string{"name", "def"},
	pKeys:  1,
}
