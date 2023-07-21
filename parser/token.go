package parser

type token struct {
	lit string
	typ tokenType
	ln  int
	col int
}

type tokenType uint8

const (
	key tokenType = iota
	sym
	ide
	str
	num
	bol
)

const eof = rune(0)

var keywords map[string]struct{} = map[string]struct{}{
	"select":         {},
	"from":           {},
	"create":         {},
	"as":             {},
	"table":          {},
	"insert":         {},
	"into":           {},
	"values":         {},
	"int":            {},
	"text":           {},
	"double":         {},
	"bool":           {},
	"where":          {},
	"and":            {},
	"or":             {},
	"true":           {},
	"false":          {},
	"unique":         {},
	"primary":        {},
	"key":            {},
	"not":            {},
	"null":           {},
	"limit":          {},
	"auto_increment": {},
}

var symbols map[string]struct{} = map[string]struct{}{
	"*":  {},
	";":  {},
	",":  {},
	"(":  {},
	")":  {},
	"=":  {},
	"!=": {},
	"+":  {},
	"-":  {},
	"<":  {},
	">":  {},
	"<=": {},
	">=": {},
}

func IsSymbol(s string) bool   { _, ok := symbols[s]; return ok }
func IsSymbolChar(r rune) bool { _, ok := symbols[string(r)]; return ok }
