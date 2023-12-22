package parse

import "testing"

func TestLex(t *testing.T) {
	q := `
	/*
		Multiline
		comment
	*/

	select test From test
	WHERE test = "foo bar" and
	// Inline comment

	baz > 1 or baz < 2.0;
	`

	lex := NewLexer(q)

	expected := []Token{
		{Type: tokComment, Text: "/* Multiline comment */"},
		{Type: tokSelect, Text: "select"},
		{Type: tokIdent, Text: "test"},
		{Type: tokFrom, Text: "from"},
		{Type: tokIdent, Text: "test"},
		{Type: tokWhere, Text: "where"},
		{Type: tokIdent, Text: "test"},
		{Type: tokEqual, Text: "="},
		{Type: tokString, Text: "foo bar"},
		{Type: tokAnd, Text: "and"},
		{Type: tokComment, Text: "// Inline comment"},
		{Type: tokIdent, Text: "baz"},
		{Type: tokGreater, Text: ">"},
		{Type: tokInt, Text: "1"},
		{Type: tokOr, Text: "or"},
		{Type: tokIdent, Text: "baz"},
		{Type: tokLess, Text: "<"},
		{Type: tokFloat, Text: "2.0"},
		{Type: tokSemicol, Text: ";"},
	}

	i := 0
	for tok := lex.Lex(); tok.Type != tokEOF; tok = lex.Lex() {
		if i >= len(expected) {
			t.Fatalf("expected only %d tokens, got more", len(expected))
		} else if tok.Type != expected[i].Type {
			t.Fatalf("expected type %v tokens, got %v: %s", expected[i].Type, tok.Type, tok.Text)
		} else if tok.Text != expected[i].Text {
			t.Fatalf("expected text %v tokens, got %v", expected[i].Text, tok.Text)
		}
		i++
	}
}
