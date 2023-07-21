package parser

import (
	"fmt"
	"strings"
	"unicode"
)

type Lexer struct {
	qry string
	ptr int
	ln  int
	col int
}

func NewLexer(qry string) *Lexer {
	return &Lexer{qry: qry}
}

func (lex *Lexer) Lex() ([]token, error) {
	tokens := []token{}

	for lex.ptr < len(lex.qry) {
		var tok *token
		var ok bool

		if skiped := lex.skipWhitespace(); skiped {
			continue
		} else if tok, ok = lex.lexIdentOrKeyword(); ok {
			goto append
		} else if tok, ok = lex.lexNumeric(); ok {
			goto append
		} else if tok, ok = lex.lexString(); ok {
			goto append
		} else if tok, ok = lex.lexSymbol(); ok {
			goto append
		}

	append:
		if tok != nil {
			tokens = append(tokens, *tok)
			continue
		}
		return nil, fmt.Errorf("cannot lex token at [%d:%d]", lex.ln, lex.col)
	}
	return tokens, nil
}

func (lex *Lexer) lexIdentOrKeyword() (*token, bool) {
	if !unicode.IsLetter(lex.curRune()) {
		return nil, false
	}
	bLn, bCol := lex.ln, lex.col
	ln, col := lex.ln, lex.col

	var ptr int
	for ptr = lex.ptr + 1; ptr < len(lex.qry); ptr++ {
		col++
		ch := lex.runeAt(ptr)
		if unicode.IsSpace(ch) {
			break
		}
		if !unicode.IsLetter(ch) && !unicode.IsNumber(ch) && ch != '_' {
			break
		}
	}

	lit := strings.ToLower(lex.qry[lex.ptr:ptr])
	typ := ide
	if _, ok := keywords[lit]; ok {
		typ = key
		if lit == "true" || lit == "false" {
			typ = bol
		}
	}

	lex.updatePos(ptr, ln, col)
	return &token{lit, typ, bLn, bCol}, true
}

func (lex *Lexer) lexNumeric() (*token, bool) {
	if !unicode.IsNumber(lex.curRune()) && lex.curRune() != '.' {
		return nil, false
	}
	bLn, bCol := lex.ln, lex.col
	ln, col := lex.ln, lex.col

	hasPeriod := false

	var ptr int
	for ptr = lex.ptr + 1; ptr < len(lex.qry); ptr++ {
		col++
		ch := lex.runeAt(ptr)
		if unicode.IsSpace(ch) || IsSymbolChar(ch) {
			break
		} else if ch == '.' {
			if hasPeriod {
				return nil, false
			}
			hasPeriod = true
		} else if !unicode.IsNumber(ch) {
			return nil, false
		}
	}

	lit := lex.qry[lex.ptr:ptr]
	lex.updatePos(ptr, ln, col)
	return &token{lit, num, bLn, bCol}, true
}

func (lex *Lexer) lexString() (*token, bool) {
	if lex.curRune() != '"' {
		return nil, false
	}
	bLn, bCol := lex.ln, lex.col
	ln, col := lex.ln, lex.col

	quoteClosed := false

	var ptr int
	for ptr = lex.ptr + 1; ptr < len(lex.qry); ptr++ {
		col++
		ch := lex.runeAt(ptr)
		if ch == '\n' {
			return nil, false
		} else if ch == '"' {
			quoteClosed = true
			ptr++
			col++
			break
		}
	}

	if !quoteClosed {
		return nil, false
	}

	lit := lex.qry[lex.ptr+1 : ptr-1]
	lex.updatePos(ptr, ln, col)
	return &token{lit, str, bLn, bCol}, true
}

func (lex *Lexer) lexSymbol() (*token, bool) {
	bLn, bCol := lex.ln, lex.col
	ln, col := lex.ln, lex.col

	var lit string
	var ptr int
	for ptr = lex.ptr; ptr < len(lex.qry); ptr++ {
		if !IsSymbol(lex.qry[lex.ptr : ptr+1]) {
			break
		}
		col++
		lit = lex.qry[lex.ptr : ptr+1]
	}
	if lit == "" {
		return nil, false
	}

	lex.updatePos(ptr, ln, col)
	return &token{lit, sym, bLn, bCol}, true
}

func (lex *Lexer) skipWhitespace() bool {
	ch := lex.curRune()
	if !unicode.IsSpace(ch) {
		return false
	}
	if ch == '\n' {
		lex.col = 0
		lex.ln++
	} else {
		lex.col++
	}

	lex.ptr++
	return true
}

func (lex *Lexer) runeAt(i int) rune     { return rune(lex.qry[i]) }
func (lex *Lexer) curRune() rune         { return rune(lex.qry[lex.ptr]) }
func (lex *Lexer) updatePos(p, l, c int) { lex.ptr = p; lex.ln = l; lex.col = c }
