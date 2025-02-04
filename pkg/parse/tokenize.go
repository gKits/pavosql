package parse

import (
	"io"
	"iter"
	"strings"
	"text/scanner"
)

type TokenType int

const (
	LexError TokenType = iota - 1
	String
	RawString
	Int
	Float
	Ident

	SpecialChar // This is a separator => if Token.Type > SpecialChar && Token.Type < Keyword {...}
	Semicolon
	Period
	Equal
	LParen
	RParen
	LBracket
	RBracket
	LBrace
	RBrace
	Plus
	Hyphen
	Asterisk
	Greater
	Less

	Keyword // This is a separator => if Token.Type > Keyword {...}
	Select
	Delete
	Create
	Update
	Insert
	From
	Into
	Table
	Set
	Values
	Where
	If
	Exists
	Not
	And
	Or
)

var keywords map[string]TokenType = map[string]TokenType{
	"select": Select,
	"delete": Delete,
	"create": Create,
	"update": Update,
	"insert": Insert,
	"from":   From,
	"into":   Into,
	"table":  Table,
	"set":    Set,
	"values": Values,
	"where":  Where,
	"if":     If,
	"exists": Exists,
	"not":    Not,
	"and":    And,
	"or":     Or,
}

var specialChars map[string]TokenType = map[string]TokenType{
	";": Semicolon,
	".": Period,
	"=": Equal,
	"(": LParen,
	")": RParen,
	"[": LBracket,
	"]": RBracket,
	"{": LBrace,
	"}": RBrace,
}

type Token struct {
	Val          string
	Type         TokenType
	Line, Column int
}

func tokenize(r io.Reader) iter.Seq[Token] {
	scan := new(scanner.Scanner)
	scan.Init(r)

	return func(yield func(Token) bool) {
		for r := scan.Scan(); r != scanner.EOF; r = scan.Scan() {
			tok := Token{
				Val:    scan.TokenText(),
				Line:   scan.Pos().Line,
				Column: scan.Pos().Column - len(scan.TokenText()),
			}

			switch r {
			case scanner.Int:
				tok.Type = Int
			case scanner.Float:
				tok.Type = Float
			case scanner.String, scanner.Char:
				tok.Type = String
			case scanner.RawString:
				tok.Type = RawString
			case scanner.Ident:
				if keyword, ok := keywords[strings.ToLower(tok.Val)]; ok {
					tok.Type = keyword
					break
				}
				tok.Type = Ident
			default:
				if specCh, ok := specialChars[tok.Val]; ok {
					tok.Type = specCh
					break
				}
				tok.Type = LexError
			}

			if !yield(tok) {
				return
			}
		}
	}
}
