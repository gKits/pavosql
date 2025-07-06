package fmt_test

import (
	"io"
	"testing"

	"github.com/gkits/pavosql/pkg/fmt"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		r       io.Reader
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := fmt.Format(tt.r)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Format() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Format() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
