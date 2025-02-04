package parse

import (
	"strings"
	"testing"
)

func Test_tokenize(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []Token
	}{
		{
			name: "tokenize with keywords and special chars #1",
			in:   `SeLECt name frOM users wHERE name == "john doe"`,
			want: []Token{
				{"SeLECt", Select, 1, 1},
				{"name", Ident, 1, 8},
				{"frOM", From, 1, 13},
				{"users", Ident, 1, 18},
				{"wHERE", Where, 1, 24},
				{"name", Ident, 1, 30},
				{"=", Equal, 1, 35},
				{"=", Equal, 1, 36},
				{"\"john doe\"", String, 1, 38},
			},
		},
		{
			name: "tokenize with keywords and special chars #2",
			in:   `  creaTe TABlE iF exists users ( )`,
			want: []Token{
				{"creaTe", Create, 1, 3},
				{"TABlE", Table, 1, 10},
				{"iF", If, 1, 16},
				{"exists", Exists, 1, 19},
				{"users", Ident, 1, 26},
				{"(", LParen, 1, 32},
				{")", RParen, 1, 34},
			},
		},
		{
			name: "tokenize with comments #1",
			in: `"hello" // this is an inline comment
                    /* This
                        is a multiline comment
                    */
                . = [] {} 'xxxx'
            `,
			want: []Token{
				{"\"hello\"", String, 1, 1},
				{".", Period, 5, 17},
				{"=", Equal, 5, 19},
				{"[", LBracket, 5, 21},
				{"]", RBracket, 5, 22},
				{"{", LBrace, 5, 24},
				{"}", RBrace, 5, 25},
				{"'xxxx'", String, 5, 27},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var i int
			for got := range tokenize(strings.NewReader(c.in)) {
				if i >= len(c.want) {
					t.Fatalf("want %d tokens, got at least %d", len(c.want), i+1)
				}
				if got != c.want[i] {
					t.Fatalf("want token %v, got %v", c.want[i], got)
				}
				i++
			}
		})
	}
}
