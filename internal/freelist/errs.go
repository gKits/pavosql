package freelist

import "errors"

var (
	errNodeDecode = errors.New("freelist: cannot decode, type does not match")
)
