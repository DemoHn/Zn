package syntax

import (
	"testing"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

type preTokenCase struct {
	name    string
	input   string
	astTree string
}

func TestPreTokenParser_OK(t *testing.T) {
	cases := []preTokenCase{
		{
			name:    "preToken test #1",
			input:   "1024",
			astTree: "$BK($NUM(1024))",
		},
		{
			name:    "preToken test #2 (var-decl stmt)",
			input:   "令A为B",
			astTree: "$BK($VD(vars[]=($ID(A)) expr[]=($ID(B))))",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil {
					e, _ := r.(*error.Error)
					t.Errorf("expect no error, got error: %s", e.Display())
				}
			}()

			var err *error.Error
			var tk *lex.Token
			var block Node
			// Parse Round I
			in := lex.NewTextStream(tt.input)
			lexI := lex.NewLexer(in)

			tokenList := []*lex.Token{}

			for {
				tk, err = lexI.NextToken()
				if err != nil {
					panic(err)
				}
				tokenList = append(tokenList, tk)
				if tk.Type == lex.TypeEOF {
					break
				}
			}
			// Parse Round II

			lexII := lex.NewPreTokenLexer(lexI.LineStack, tokenList)
			parserII := NewParser(lexII)
			block, err = parserII.Parse()
			if err != nil {
				panic(err)
			}

			expect := StringifyAST(block)
			got := formatASTstr(tt.astTree)
			if expect != got {
				t.Errorf("AST compare:\nexpect ->\n%s\ngot ->\n%s", expect, got)
			}
		})
	}
}
