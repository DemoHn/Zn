package syntax

import (
	"regexp"
	"strings"
	"testing"

	"github.com/DemoHn/Zn/lex"
)

/**
# Introduction of ast testcases

TODO
*/
var testSuccessSuites = []string{
	varDeclCasesOK,
}

const varDeclCasesOK = `
========
1. inline one var
--------
令某变量为100
--------
$PG($BK(
	$VD(
		vars[]=($ID(某变量))
		expr[]=($NUM(100))
	)
))

========
2. two variables
--------
令变量1，变量2为100
--------
$PG($BK(
	$VD(
		vars[]=($ID(变量1) $ID(变量2))
		expr[]=($NUM(100))
	)
))

========
3. paired variables inline
--------
令小A，小B为100，小C为「何处相思明月楼」，D，E，F为B
--------
$PG($BK(
	$VD(
		vars[]=($ID(小A) $ID(小B))
		expr[]=($NUM(100))
		vars[]=($ID(小C))
		expr[]=($STR(何处相思明月楼))
		vars[]=($ID(D) $ID(E) $ID(F))
		expr[]=($ID(B))
	)
))

========
4. with varquotes
--------
令小A，·先令·为200
--------
$PG($BK(
	$VD(
		vars[]=($ID(小A) $ID(先令))
		expr[]=($NUM(200))
	)
))

========
5. A -> B -> C
--------
令A为B为C
--------
$PG($BK(
	$VD(
		vars[]=($ID(A))
		expr[]=(
			$VA(
				target=($ID(B))
				assign=($ID(C))
			)
		)
	)
))

========
6. block var declare
--------
令：
	A为1，B为2，C，D为3，
	E，F为4

令G为5
--------
$PG($BK(
	$VD(
		vars[]=($ID(A))		expr[]=($NUM(1))
		vars[]=($ID(B))		expr[]=($NUM(2))
		vars[]=($ID(C) $ID(D))		expr[]=($NUM(3))
		vars[]=($ID(E) $ID(F))		expr[]=($NUM(4))
	)
	$VD(
		vars[]=($ID(G))
		expr[]=($NUM(5))
	)
))
`

type astSuccessCase struct {
	name    string
	input   string
	astTree string
}

func TestAST_OK(t *testing.T) {
	astCases := []astSuccessCase{}

	for _, suData := range testSuccessSuites {
		suites := splitTestSuites(suData)
		for _, suite := range suites {
			astCases = append(astCases, astSuccessCase{
				name:    suite[0],
				input:   suite[1],
				astTree: suite[2],
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

			node, err := p.Parse()
			if err != nil {
				t.Errorf("expect no error, got error: %s", err.Error())
			} else {
				// compare with ast
				expect := StringifyAST(node)
				got := formatASTstr(tt.astTree)

				if expect != got {
					t.Errorf("AST compare:\nexpect ->\n%s\ngot ->\n%s", expect, got)
				}
			}
		})

	}
}

func splitTestSuites(source string) [][3]string {
	result := [][3]string{}

	source = strings.Replace(source, "\r\n", "\n", -1)
	sourceArr := strings.Split(source, "\n")

	const (
		sInit    = 0
		sPartI   = 1
		sPartII  = 2
		sPartIII = 3
	)
	var state = sInit
	l1 := []string{}
	l2 := []string{}
	l3 := []string{}
	for _, line := range sourceArr {
		if strings.HasPrefix(line, "========") {
			// push old data
			if state == sPartIII {
				result = append(result, [3]string{
					strings.Join(l1, "\n"),
					strings.Join(l2, "\n"),
					strings.Join(l3, "\n"),
				})
			}
			state = sPartI
			// clear buffer
			l1 = []string{}
			l2 = []string{}
			l3 = []string{}
			continue
		}
		if strings.HasPrefix(line, "--------") {
			if state == sPartI {
				state = sPartII
			} else if state == sPartII {
				state = sPartIII
			}
			continue
		}

		switch state {
		case sPartI:
			l1 = append(l1, line)
		case sPartII:
			l2 = append(l2, line)
		case sPartIII:
			l3 = append(l3, line)
		}
	}

	// tail append
	if state == sPartIII {
		result = append(result, [3]string{
			strings.Join(l1, "\n"),
			strings.Join(l2, "\n"),
			strings.Join(l3, "\n"),
		})
	}
	return result
}

func formatASTstr(input string) string {
	reL := regexp.MustCompile(`\((\s)+`)
	reR := regexp.MustCompile(`(\s)+\)`)
	reS := regexp.MustCompile(`(\s)+`)

	input = reL.ReplaceAllString(input, "(")
	input = reR.ReplaceAllString(input, ")")
	input = reS.ReplaceAllString(input, " ")

	return strings.TrimSpace(input)
}
