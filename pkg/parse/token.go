package parse

type tokenType uint

const (
	tokError tokenType = iota
	tokEOF
	tokNumber
	tokString
	tokComment
	tokIdent
	// keywords
	tokDot      // .
	tokComma    // ,
	tokLeftPar  // (
	tokRightPar // )
	// operator keyword
	tokOperator // only used for separation of operators
	tokGet      // get
	tokDelete   // delete
	tokUpdate   // update
	tokInsert   // insert
)

var operators map[string]parseState = map[string]parseState{
	"get":    parseGet,
	"delete": parseDelete,
	"update": parseUpdate,
	"insert": parseInsert,
	"create": parseCreate,
}

var keywords map[string]tokenType = map[string]tokenType{
	".":      tokDot,
	",":      tokComma,
	"(":      tokLeftPar,
	")":      tokRightPar,
	"get":    tokGet,
	"delete": tokDelete,
	"update": tokUpdate,
	"insert": tokInsert,
}

type token struct {
	typ tokenType
	val string
}
