package syntax

import (
	"fmt"
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
	whileLoopCasesOK,
	logicExprCasesOK,
}

var testFailSuites = []string{}

const logicExprCasesOK = `
========
1. low -> high precedence
--------
此{A且B或C且D等于E且F为100}为0
--------
$PG($BK(
	$IS(
		L=($OR(
				L=($AND(L=($ID(A)) R=($ID(B))))
				R=($AND(					
					L=($AND(
						L=($ID(C))
						R=($EQ(
							L=($ID(D))
							R=($ID(E))
						))
					))
					R=($VA(
						target=($ID(F))
						assign=($NUM(100))
					))
				))
		))
		R=($NUM(0))
	)
))
`

const whileLoopCasesOK = `
========
1. one line block
--------
每当1：
	令A为B
--------
$PG($BK(
	$WL(
		expr=($NUM(1))
		block=($BK($VD(
				vars[]=($ID(A))
				expr[]=($ID(B))
		)))
	)
))

========
2. nested while loop statement
--------
每当1：
	A为B
	每当2：
		C为D
		E为F
	每当3：
		100
	G为H
	K为L

M为N
--------
$PG($BK(
	$WL(
		expr=($NUM(1))
		block=($BK(
			$VA(target=($ID(A)) assign=($ID(B)))
			$WL(
				expr=($NUM(2))
				block=($BK(
					$VA(target=($ID(C)) assign=($ID(D)))
					$VA(target=($ID(E)) assign=($ID(F)))
				))
			)
			$WL(
				expr=($NUM(3))
				block=($BK($NUM(100)))
			)
			$VA(target=($ID(G)) assign=($ID(H)))
			$VA(target=($ID(K)) assign=($ID(L)))
		))
	)

	$VA(target=($ID(M)) assign=($ID(N)))
))
`

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

type astFailCase struct {
	name     string
	input    string
	failInfo string
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
			//
		})
	}
}

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
				got := fmt.Sprintf("code=%d line=%d col=%d", err.GetCode(), cursor.LineNum, cursor.ColNum)
				if tt.failInfo != got {
					t.Errorf("failInfo compare:\nexpect ->\n%s\ngot ->\n%s", tt.failInfo, got)
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
