package db

import (
	"encoding/binary"
	"fmt"
)

type table struct {
	Name     string   `json:"name"`
	Cols     []string `json:"columns"`
	Types    []dbType `json:"types"`
	Nullable []bool   `json:"nullable"`
	PKeys    int      `json:"pkeys"`
	Prefix   uint32   `json:"prefix"`
	Indices  []int    `json:"indices"`
}

func (tb *table) encodeKey(k []byte) []byte {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], tb.Prefix)
	return append(k, buf[:]...)
}

func (tb *table) checkRow(row *Row) error {
	set := map[string]struct{}{}
	for _, col := range tb.Cols {
		set[col] = struct{}{}
	}

	for _, col := range row.Cols {
		if _, ok := set[col]; !ok {
			return fmt.Errorf("table: unknow column '%s' in table '%s'", col, tb.Name)
		}
	}

	return nil
}
