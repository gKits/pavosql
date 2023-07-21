package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexIdentOrKeyword(t *testing.T) {
	cases := []struct {
		name  string
		input string
		tok   *token
		ok    bool
	}{
		{
			name:  "select keyword with whitespace at the end",
			input: "select ",
			tok:   &token{"select", key, 0, 0},
			ok:    true,
		},
		{
			name:  "from keyword in all caps",
			input: "FROM",
			tok:   &token{"from", key, 0, 0},
			ok:    true,
		},
		{
			name:  "create keyword with line break at the end",
			input: "Create\n",
			tok:   &token{"create", key, 0, 0},
			ok:    true,
		},
		{
			name:  "as keyword with random capitalization",
			input: "aS",
			tok:   &token{"as", key, 0, 0},
			ok:    true,
		},
		{
			name:  "true boolean keyword",
			input: "true",
			tok:   &token{"true", bol, 0, 0},
			ok:    true,
		},
		{
			name:  "false boolean keyword in all caps",
			input: "FALSE",
			tok:   &token{"false", bol, 0, 0},
			ok:    true,
		},
		{
			name:  "random identifier with underscores and numbers",
			input: "test_ident4",
			tok:   &token{"test_ident4", ide, 0, 0},
			ok:    true,
		},
		{
			name:  "identifier ending with comma",
			input: "test,",
			tok:   &token{"test", ide, 0, 0},
			ok:    true,
		},
		{
			name:  "whitespace at beginnig",
			input: " select",
			tok:   nil,
			ok:    false,
		},
		{
			name:  "number at beginning",
			input: "5select",
			tok:   nil,
			ok:    false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lex := Lexer{qry: c.input}
			tok, ok := lex.lexIdentOrKeyword()
			assert.Equal(t, c.ok, ok)
			assert.Equal(t, c.tok, tok)
		})
	}

}

func TestLexNumeric(t *testing.T) {
	cases := []struct {
		name  string
		input string
		tok   *token
		ok    bool
	}{
		{
			name:  "integer",
			input: "123",
			tok:   &token{"123", num, 0, 0},
			ok:    true,
		},
		{
			name:  "integer with whitespace at the end",
			input: "123 ",
			tok:   &token{"123", num, 0, 0},
			ok:    true,
		},
		{
			name:  "integer value with linebreak at the end",
			input: "123\n",
			tok:   &token{"123", num, 0, 0},
			ok:    true,
		},
		{
			name:  "float without digits behind period",
			input: "123.",
			tok:   &token{"123.", num, 0, 0},
			ok:    true,
		},
		{
			name:  "float with digits behind period",
			input: "123.456",
			tok:   &token{"123.456", num, 0, 0},
			ok:    true,
		},
		{
			name:  "float with period at the beginning",
			input: ".456",
			tok:   &token{".456", num, 0, 0},
			ok:    true,
		},
		{
			name:  "float followed by symbol",
			input: "1.234,",
			tok:   &token{"1.234", num, 0, 0},
			ok:    true,
		},
		{
			name:  "integer with letter",
			input: "123A",
			tok:   nil,
			ok:    false,
		},
		{
			name:  "float with double period",
			input: "123.45.6",
			tok:   nil,
			ok:    false,
		},
		{
			name:  "integer with whitespace at the beginnig",
			input: " 123",
			tok:   nil,
			ok:    false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lex := Lexer{qry: c.input}
			tok, ok := lex.lexNumeric()
			assert.Equal(t, c.ok, ok)
			assert.Equal(t, c.tok, tok)
		})
	}
}

func TestLexString(t *testing.T) {
	cases := []struct {
		name  string
		input string
		tok   *token
		ok    bool
	}{
		{
			name:  "string with quotes and whitespaces",
			input: "\"test 123, : select\"",
			tok:   &token{"test 123, : select", str, 0, 0},
			ok:    true,
		},
		{
			name:  "string without quote at the beginning",
			input: "test",
			tok:   nil,
			ok:    false,
		},
		{
			name:  "string with line break in between",
			input: "\"test\n\"",
			tok:   nil,
			ok:    false,
		},
		{
			name:  "string without closing quote",
			input: "\"test",
			tok:   nil,
			ok:    false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lex := Lexer{qry: c.input}
			tok, ok := lex.lexString()
			assert.Equal(t, c.ok, ok)
			assert.Equal(t, c.tok, tok)
		})
	}
}

func TestLexSymbol(t *testing.T) {
	cases := []struct {
		name  string
		input string
		tok   *token
		ok    bool
	}{
		{
			name:  "greater equal symbol with whitespace at the end",
			input: ">= ",
			tok:   &token{">=", sym, 0, 0},
			ok:    true,
		},
		{
			name:  "comma symbol with linebreak at the end",
			input: ",\n",
			tok:   &token{",", sym, 0, 0},
			ok:    true,
		},
		{
			name:  "symbol followed by a non symbol",
			input: ">a",
			tok:   &token{">", sym, 0, 0},
			ok:    true,
		},
		{
			name:  "symbol followed by another symbol unknow combination",
			input: "),",
			tok:   &token{")", sym, 0, 0},
			ok:    true,
		},
		{
			name:  "does not start with a symbol",
			input: "t<",
			tok:   nil,
			ok:    false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lex := Lexer{qry: c.input}
			tok, ok := lex.lexSymbol()
			assert.Equal(t, c.ok, ok)
			assert.Equal(t, c.tok, tok)
		})
	}
}

func TestSkipWhitespace(t *testing.T) {
	cases := []struct {
		name  string
		input string
		ok    bool
		ln    int
		col   int
	}{
		{
			name:  "just some whitespaces",
			input: "   ",
			ok:    true,
			ln:    0,
			col:   1,
		},
		{
			name:  "newline at the beginning",
			input: "\n  ",
			ok:    true,
			ln:    1,
			col:   0,
		},
		{
			name:  "tab at the beginning",
			input: "\t  ",
			ok:    true,
			ln:    0,
			col:   1,
		},
		{
			name:  "non space char at the beginning",
			input: "b  ",
			ok:    false,
			ln:    0,
			col:   0,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lex := Lexer{qry: c.input}
			ok := lex.skipWhitespace()
			assert.Equal(t, c.ok, ok)
			assert.Equal(t, c.ln, lex.ln)
			assert.Equal(t, c.col, lex.col)
		})
	}
}

func TestLex(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		tokens []token
		err    error
	}{
		{
			name:  "select a",
			input: "select a",
			tokens: []token{
				{
					lit: "select",
					typ: key,
					ln:  0,
					col: 0,
				},
				{
					lit: "a",
					typ: ide,
					ln:  0,
					col: 7,
				},
			},
		},
		{
			name:  "create table test",
			input: "create table test (\n\tid int primary key not null auto_increment,\n\tname text,\n\theight double\n);",
			tokens: []token{
				{
					lit: "create",
					typ: key,
					ln:  0,
					col: 0,
				},
				{
					lit: "table",
					typ: key,
					ln:  0,
					col: 7,
				},
				{
					lit: "test",
					typ: ide,
					ln:  0,
					col: 13,
				},
				{
					lit: "(",
					typ: sym,
					ln:  0,
					col: 18,
				},
				{
					lit: "id",
					typ: ide,
					ln:  1,
					col: 1,
				},
				{
					lit: "int",
					typ: key,
					ln:  1,
					col: 4,
				},
				{
					lit: "primary",
					typ: key,
					ln:  1,
					col: 8,
				},
				{
					lit: "key",
					typ: key,
					ln:  1,
					col: 16,
				},
				{
					lit: "not",
					typ: key,
					ln:  1,
					col: 20,
				},
				{
					lit: "null",
					typ: key,
					ln:  1,
					col: 24,
				},
				{
					lit: "auto_increment",
					typ: key,
					ln:  1,
					col: 29,
				},
				{
					lit: ",",
					typ: sym,
					ln:  1,
					col: 43,
				},
				{
					lit: "name",
					typ: ide,
					ln:  2,
					col: 1,
				},
				{
					lit: "text",
					typ: key,
					ln:  2,
					col: 6,
				},
				{
					lit: ",",
					typ: sym,
					ln:  2,
					col: 10,
				},
				{
					lit: "height",
					typ: ide,
					ln:  3,
					col: 1,
				},
				{
					lit: "double",
					typ: key,
					ln:  3,
					col: 8,
				},
				{
					lit: ")",
					typ: sym,
					ln:  4,
					col: 0,
				},
				{
					lit: ";",
					typ: sym,
					ln:  4,
					col: 1,
				},
			},
		},
		{
			name:  "insert into table",
			input: "insert into table test (name, height)\nvalues (\"bob\", 1.84);",
			tokens: []token{
				{
					lit: "insert",
					typ: key,
					ln:  0,
					col: 0,
				},
				{
					lit: "into",
					typ: key,
					ln:  0,
					col: 7,
				},
				{
					lit: "table",
					typ: key,
					ln:  0,
					col: 12,
				},
				{
					lit: "test",
					typ: ide,
					ln:  0,
					col: 18,
				},
				{
					lit: "(",
					typ: sym,
					ln:  0,
					col: 23,
				},
				{
					lit: "name",
					typ: ide,
					ln:  0,
					col: 24,
				},
				{
					lit: ",",
					typ: sym,
					ln:  0,
					col: 28,
				},
				{
					lit: "height",
					typ: ide,
					ln:  0,
					col: 30,
				},
				{
					lit: ")",
					typ: sym,
					ln:  0,
					col: 36,
				},
				{
					lit: "values",
					typ: key,
					ln:  1,
					col: 0,
				},
				{
					lit: "(",
					typ: sym,
					ln:  1,
					col: 7,
				},
				{
					lit: "bob",
					typ: str,
					ln:  1,
					col: 8,
				},
				{
					lit: ",",
					typ: sym,
					ln:  1,
					col: 13,
				},
				{
					lit: "1.84",
					typ: num,
					ln:  1,
					col: 15,
				},
				{
					lit: ")",
					typ: sym,
					ln:  1,
					col: 19,
				},
				{
					lit: ";",
					typ: sym,
					ln:  1,
					col: 20,
				},
			},
		},
		{
			name:  "cannot lex broken number",
			input: "123A",
			err:   fmt.Errorf("cannot lex token at [0:0]"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lex := Lexer{qry: c.input}
			tokens, err := lex.Lex()
			assert.Equal(t, c.err, err)
			assert.Equal(t, len(c.tokens), len(tokens))
			for i, tok := range tokens {
				assert.Equal(t, c.tokens[i], tok)
			}
		})
	}
}

func TestNewLexer(t *testing.T) {
	lex := NewLexer("test")
	assert.NotNil(t, lex)
	assert.Equal(t, "test", lex.qry, "lexer query should be %s but is %s instead", "test", lex.qry)
	assert.Equal(t, 0, lex.ptr, "lexer pointer should be %d but is %d instead", 0, lex.ptr)
	assert.Equal(t, 0, lex.ln, "lexer line should be %d but is %d instead", 0, lex.ln)
	assert.Equal(t, 0, lex.col, "lexer column should be %d but is %d instead", 0, lex.col)
}
