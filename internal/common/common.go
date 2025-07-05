package common

const PageSize = 8192

type PageType uint8

const (
	PointerPage PageType = iota + 1
	LeafPage
)
