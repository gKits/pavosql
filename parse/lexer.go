package parse

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type stateFn func(*lexer) stateFn

type lexer struct {
	name    string // name of lexer for debugging
	input   string // input string beeing tokenized
	pos     int    // current position in input string
	start   int    // starting position of current token beeing scanned
	width   int    // width of last rune read
	tok     token  // last emitted token
	line    int    // current line in multt-line statement
	startLn int    // starting line of current token
	opt     lexOptions
}

type lexOptions struct {
	emitComments  bool
	allowComments bool
	leftComment   string
	rightComment  string
	singleComment string
}

var defaultLexOptions lexOptions = lexOptions{
	emitComments:  false,
	allowComments: true,
	leftComment:   "/*",
	rightComment:  "*/",
	singleComment: "//",
}

const eof = -1

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	if r == '\n' {
		l.line++
	}

	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	if l.pos > 0 {
		r, w := utf8.DecodeLastRuneInString(l.input[:l.pos])
		l.pos -= w
		if r == '\n' {
			l.line--
		}
	}
}

func (l *lexer) ignore() {
	l.line += strings.Count(l.input[l.start:l.pos], "\n")
	l.start = l.pos
	l.startLn = l.line
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *lexer) curToken(typ tokenType) token {
	tok := token{typ, l.input[l.start:l.pos], l.startLn}
	l.start = l.pos
	l.startLn = l.line
	return tok
}

func (l *lexer) emit(typ tokenType) stateFn {
	return l.emitToken(l.curToken(typ))

}

func (l *lexer) emitToken(t token) stateFn {
	l.tok = t
	return nil
}

func (l *lexer) run() {
	for state := lexIdent; state != nil; {
		state = nil
	}
}

func lexIdent(l *lexer) stateFn {
	for {
		r := l.next()

		if unicode.IsSpace(r) {
			l.backup()
			word := l.input[l.start:l.pos]

			typ, ok := key[word]
			if !ok {
				typ = tokIdent
			}

			return l.emit(typ)
		}
	}
}

func lexString(l *lexer) stateFn {
loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof {
				break
			}
			fallthrough

		case eof, '\n':
			return l.errorf("unterminated string")

		case '"':
			break loop
		}
	}

	return l.emit(tokString)
}

func lexRawString(l *lexer) stateFn {
	for {
		switch l.next() {
		case '`':

		}
	}
}

func lexNumber(l *lexer) stateFn {
	l.accept("+-")
	digits := "0123456789"

	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}

	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun(digits)
	}

	if unicode.IsLetter(l.peek()) {
		l.next()
		return l.errorf("invalid number syntax: %q", l.input[l.start:l.pos])
	}
	return l.emit(tokNumber)
}

func lexComment(l *lexer) stateFn {
	l.pos += len(l.opt.leftComment)
	closeIdx := strings.Index(l.input[l.pos:], l.opt.rightComment)
	if closeIdx < 0 {
		return l.errorf("unterminated comment")
	}
	l.pos += closeIdx + len(l.opt.rightComment)

	l.curToken(tokComment)

	return nil
}

func lexSingleLineComment(l *lexer) stateFn {
	l.pos += len(l.opt.singleComment)
	newlineIdx := strings.Index(l.input[l.pos:], "\n")
	if newlineIdx < 0 {
	}
	l.pos += newlineIdx
	return nil
}

func lexSpace(l *lexer) stateFn {
	for r := l.peek(); unicode.IsSpace(r); r = l.peek() {
		l.next()
	}
	return nil
}

func (l *lexer) errorf(format string, args ...any) stateFn {
	l.tok = token{tokErr, fmt.Sprintf(format, args...), l.line}
	l.start = 0
	l.pos = 0
	l.input = l.input[:0]
	return nil
}
