package btree

import "errors"

var (
	errBTreeHeader      = errors.New("btree: cannot decode page, header is malformed")
	errBTreeDeleteEmpty = errors.New("btree: cannot delete key from empty btree")
	errBTreeDeleteKey   = errors.New("btree: cannot delete non existing key")
)
