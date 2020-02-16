package syntax

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DemoHn/Zn/lex"
)

var testFailSuites = []string{
	varDeclCasesFAIL,
}

const varDeclCasesFAIL = `
========
1. non-identifiers as assigner (InvalidSyntax)
--------
注：第一行留给度娘

令某变量，另一变量，1240为1000
--------
code=2250 line=3 col=10

========
2. incomplete statement (additional comma) (InvalidSyntax)
--------
    
令某变量，另一变量，
【A，B】为1
    
--------
code=2250 line=3 col=0

========
3. incomplete statement (InvalidSyntax)
--------
    
令某变量，另一变量
    【A，B】为100
    
--------
code=2252 line=2 col=5

========
4. invalid token (lexError)
--------
令锅为「锅」

令#$x为100
    
--------
code=2024 line=3 col=2
`

func TestAST_FAIL(t *testing.T) {
	astCases := []astFailCase{}

	for _, suData := range testFailSuites {
		suites := splitTestSuites(suData)
		for _, suite := range suites {
			astCases = append(astCases, astFailCase{
				name:     suite[0],
				input:    suite[1],
				failInfo: suite[2],
			})
		}
	}

	// TODO: filter
	// after filtering...
	for _, tt := range astCases {
		t.Run(tt.name, func(t *testing.T) {
			in := lex.NewTextStream(tt.input)
			l := lex.NewLexer(in)
			p := NewParser(l)

			_, err := p.Parse()

			if err == nil {
				t.Errorf("expect error, got no error found")
			} else {
				// compare with error code
				cursor := err.GetCursor()
				got := fmt.Sprintf("code=%x line=%d col=%d", err.GetCode(), cursor.LineNum, cursor.ColNum)
				failInfof := strings.TrimSpace(tt.failInfo)
				if failInfof != got {
					t.Errorf("failInfo compare:\nexpect ->\n%s\ngot ->\n%s", tt.failInfo, got)
				}
			}
		})
	}
}
