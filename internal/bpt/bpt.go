package bpt

import "io"

const pageSize = 4096

type pager interface {
	io.ReaderAt
	io.WriterAt
	Commit() error
	Rollback() error
}

type Tree struct {
	pager pager
}
