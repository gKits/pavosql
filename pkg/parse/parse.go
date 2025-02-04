package parse

import (
	"io"
	"sync"

	"github.com/pavosql/pavosql/pkg/ast"
)

func Parse(r io.Reader) []ast.Stmnt {
	var wg sync.WaitGroup
	wg.Add(2)

	toks := make(chan Token)
	go func() {
		defer wg.Done()
		for tok := range tokenize(r) {
			toks <- tok
		}
	}()

	go func() {
		defer wg.Done()
		for tok := range toks {
			_ = tok
			// TODO: Implement parsing logic
		}
	}()

	wg.Wait()

	return nil
}
