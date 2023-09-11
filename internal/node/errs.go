package node

import "errors"

var (
	errNodeIdx    = errors.New("node: index out of key range")
	errNodeDecode = errors.New("node: cannot decode, type does not match")
	errNodeUpdate = errors.New("node: cannot update k-v, keys must match")
	errNodeMerge  = errors.New("node: cannot merge, lefts last key must be less than rights first key")
)
