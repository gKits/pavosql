package parse

import "fmt"

type tokenType uint

const (
	tokErr tokenType = iota
	tokEOF
	tokComment
	tokBool
	tokNumber
	tokString
	tokIdent

	tokKey // Keyword delim
	tokAlter
	tokCreate
	tokDelete
	tokDrop
	tokInsert
	tokSelect
	tokTruncate
	tokUpdate
	tokColumn
	tokDatabase
	tokIndex
	tokTable
	tokDistinct
	tokTop
	tokAs
	tokFrom
	tokSet
	tokLeft
	tokRight
	tokFull
	tokJoin
	tokWhere
	tokHaving
	tokInto
	tokValues
	tokAnd
	tokIs
	tokNot
	tokNull
	tokOr
	tokLike
	tokIn
	tokBetween
	tokCount
	tokSum
	tokAvg
	tokMin
	tokMax
	tokGroup
	tokOrder
	tokBy
	tokDesc
)

var key = map[string]tokenType{
	"alter":    tokAlter,
	"create":   tokCreate,
	"delete":   tokDelete,
	"drop":     tokDrop,
	"insert":   tokInsert,
	"select":   tokSelect,
	"truncate": tokTruncate,
	"update":   tokUpdate,
	"column":   tokColumn,
	"database": tokDatabase,
	"index":    tokIndex,
	"table":    tokTable,
	"distinct": tokDistinct,
	"top":      tokTop,
	"as":       tokAs,
	"from":     tokFrom,
	"set":      tokSet,
}

type token struct {
	typ  tokenType
	val  string
	line int
}

func (t token) String() string {
	switch {
	case t.typ == tokEOF:
		return "EOF"
	case t.typ == tokErr:
		return t.val
	case t.typ > tokKey:
		return fmt.Sprintf("<KW_%s>", t.val)
	case len(t.val) > 13:
		return fmt.Sprintf("%.10q...", t.val)
	}
	return t.val
}
