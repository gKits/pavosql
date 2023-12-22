package parse

type AST interface{}

type Select struct {
	Fields []string
	Table  string
}

type CreateTable struct {
	Name   string
	Fields []string
	Types  []string
}
