package parse

type Parser struct {
	lex *Lexer
}

func NewParser(query string) *Parser {
	return &Parser{lex: NewLexer(query)}
}

func (p *Parser) Parse() (AST, error) {
	return nil, nil
}
