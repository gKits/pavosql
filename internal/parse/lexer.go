package parse

import (
	"fmt"
	"strings"
	"text/scanner"
)

type Lexer struct {
	scan *scanner.Scanner
}

func NewLexer(text string) *Lexer {
	scan := &scanner.Scanner{}
	scan.Init(strings.NewReader(text))
	scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments
	scan.Filename = "at"

	return &Lexer{scan: scan}
}

func (lex *Lexer) Lex() Token {
	r := lex.scan.Scan()

	tok := Token{}
	text := lex.scan.TokenText()

	switch r {
	case scanner.EOF:
		tok.Type = tokEOF

	case scanner.Int:
		tok.Type = tokInt
		tok.Text = text

	case scanner.Float:
		tok.Type = tokFloat
		tok.Text = text

	case scanner.String, scanner.RawString:
		tok.Type = tokString
		tok.Text = text[1 : len(text)-1]

	case scanner.Ident:
		text = strings.ToLower(text)

		var ok bool
		if tok.Type, ok = keywords[text]; !ok {
			tok.Type = tokIdent
		}
		tok.Text = text

	case scanner.Comment:
		tok.Text = strings.Join(strings.Fields(text), " ")
		tok.Type = tokComment

	default:
		tok.Text = text
		var ok bool
		if tok.Type, ok = keywords[text]; !ok {
			tok.Type = tokError
			tok.Text = fmt.Sprintf("lex: %s: invalid token '%s'", lex.scan.Position, text)
		}

	}
	return tok
}
