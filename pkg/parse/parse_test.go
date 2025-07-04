package parse_test

import (
	"io"
	"testing"

	"github.com/pavosql/pavosql/pkg/ast"
	"github.com/pavosql/pavosql/pkg/parse"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		r       io.Reader
		want    []ast.Stmnt
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := parse.Parse(tt.r)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Parse() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Parse() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
