package parse

type TokenType uint

const (
	tokError TokenType = iota
	tokEOF
	tokComment
	tokNumber
	tokInt
	tokFloat
	tokString
	tokIdent
	tokSelect
	tokAsterisk
	tokLeftPar
	tokRightPar
	tokSemicol
	tokComma
	tokEqual
	tokGreater
	tokLess
	tokGreaterEq
	tokLessEq
	tokNotEq
	tokFrom
	tokWhere
	tokCreate
	tokTable
	tokNot
	tokAnd
	tokOr
)

var keywords = map[string]TokenType{
	"select": tokSelect,
	"from":   tokFrom,
	"where":  tokWhere,
	"create": tokCreate,
	"table":  tokTable,
	",":      tokComma,
	"*":      tokAsterisk,
	";":      tokSemicol,
	"(":      tokLeftPar,
	")":      tokRightPar,
	"=":      tokEqual,
	">":      tokGreater,
	"<":      tokLess,
	">=":     tokGreaterEq,
	"<=":     tokLessEq,
	"!=":     tokNotEq,
	"not":    tokNot,
	"and":    tokAnd,
	"or":     tokOr,
}

type Token struct {
	Type   TokenType
	Text   string
	Line   int
	Column int
}
