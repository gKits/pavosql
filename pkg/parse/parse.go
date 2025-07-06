package parse

import (
	"io"

	"github.com/gkits/pavosql/pkg/ast"
)

func Parse(r io.Reader) ([]ast.Stmnt, error) {
	var (
		stmt ast.Stmnt
		err  error
	)

	stmts := []ast.Stmnt{}
	toks := readTokens(r)
	for tok := range toks {
		switch tok.Type {
		case Select:
			stmt, err = parseSelectStmt(toks)
			if err != nil {
				return nil, err
			}
		case Delete:
		case Create:
		case Update:
		case Insert:
		default:
		}
		stmts = append(stmts, stmt)
	}

	return stmts, nil
}

func parseSelectStmt(toks <-chan Token) (ast.SelectStmt, error) {
	tok := <-toks

	switch tok.Type {
	case Asterisk:
	case Ident:
	}
	return ast.SelectStmt{}, nil
}

func parseDeleteStmt() {}

func parseCreateStmt() {}

func parseUpdateStmt() {}

func parseInsertStmt() {}

func parseFieldSelectList(toks <-chan Token) {
	for tok := range toks {
		_ = tok
	}
}

func readTokens(r io.Reader) <-chan Token {
	toks := make(chan Token)
	go func() {
		for _, tok := range tokenize(r) {
			toks <- tok
		}
	}()
	return toks
}
