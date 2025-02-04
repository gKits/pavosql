package bpt

import "io"

const pageSize = 4096

type pager interface {
	io.ReaderAt
	io.WriterAt
}

type Tree struct {
	pager pager
}
