package parse

type Statement interface{}

type GetStatement struct {
	Table string
}
