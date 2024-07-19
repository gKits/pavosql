package parse

import (
	"fmt"
	"io"
	"text/scanner"
)

type parseState func(p *parser) (parseState, error)

type parser struct {
	q    Statement
	scan *scanner.Scanner
}

func Parse(query io.Reader) ([]Statement, error) {
	p := parser{
		scan: new(scanner.Scanner).Init(query),
	}

	var (
		state = parseOperator
		err   error
	)
	for state != nil {
		if state, err = state(&p); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func parseOperator(p *parser) (parseState, error) {
	tok := p.scan.Scan()
	if tok == scanner.Ident {
		if state, ok := operators[p.scan.TokenText()]; ok {
			return state, nil
		}
	}
	return nil, fmt.Errorf("expected operator, got %s", p.scan.TokenText())
}

func parseGet(p *parser) (parseState, error) {
	tok := p.scan.Scan()
	if tok != '(' {
		return nil, fmt.Errorf("expected '(', got %s", p.scan.TokenText())
	}

	tok = p.scan.Scan()
	if tok != scanner.Ident {
		return nil, fmt.Errorf("expected table name identifier, got %s", p.scan.TokenText())
	}

	tok = p.scan.Scan()
	if tok != ',' {
		return nil, fmt.Errorf("expected ',', got %s", p.scan.TokenText())
	}

	return nil, nil
}

func parseInsert(p *parser) (parseState, error) {
	return nil, nil
}

func parseUpdate(p *parser) (parseState, error) {
	return nil, nil
}

func parseDelete(p *parser) (parseState, error) {
	return nil, nil
}

func parseCreate(p *parser) (parseState, error) {
	return nil, nil
}

func parseCondition(p *parser) (parseState, error) {
	for tok := p.scan.Scan(); tok != scanner.EOF; tok = p.scan.Scan() {
		break
	}
	return nil, nil
}
