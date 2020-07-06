package syntax

import (
	"regexp"
	"strings"
	"testing"

	"github.com/DemoHn/Zn/lex"
)

var testSuccessSuites = []string{
	varDeclCasesOK,
	whileLoopCasesOK,
	logicExprCasesOK,
	arrayListCasesOK,
	funcCallCasesOK,
	branchStmtCasesOK,
	stmtLineBreakCasesOK,
	memberExprCasesOK,
}

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
		expr[]=($STR(「何处相思明月楼」))
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
const funcCallCasesOK = `
========
1. success func call with no param
--------
（显示当前时间）
--------
$PG($BK(
	$FN(name=($ID(显示当前时间)) params=())
))

========
2. success func call with no param (varquote)
--------
（·显示当前之时间·）
--------
$PG($BK(
	$FN(name=($ID(显示当前之时间)) params=())
))

========
3. success func call with 1 parameter
--------
（显示当前时间：「今天」）
--------
$PG($BK(
	$FN(name=($ID(显示当前时间)) params=($STR(「今天」)))
))

========
4. success func call with 2 parameters
--------
（显示当前时间：「今天」，「15:30」）
--------
$PG($BK(
	$FN(name=($ID(显示当前时间)) params=($STR(「今天」) $STR(「15:30」)))
))

========
5. success func call with mutliple parameters
--------
（显示当前时间：「今天」，「15:30」，200，3000）
--------
$PG($BK(
	$FN(name=($ID(显示当前时间)) params=($STR(「今天」) $STR(「15:30」) $NUM(200) $NUM(3000)))
))

========
6. nested functions
--------
（显示当前时间：「今天」，「15:30」，（显示时刻））
--------
$PG($BK(
	$FN(name=($ID(显示当前时间)) params=(
		$STR(「今天」)
		$STR(「15:30」) 
		$FN(name=($ID(显示时刻)) params=())
	))
))
`

const branchStmtCasesOK = `
========
1. if-block only
--------
如果真：
    （X+Y：20，30）
--------
$PG($BK(
	$IF(
		ifExpr=($ID(真))
		ifBlock=($BK(
			$FN(
				name=($ID(X+Y))
				params=($NUM(20) $NUM(30))
			)
		))
	)
))

========
2. if-block and else-block
--------
如果真：
    （X+Y：20，30）
否则：
    （X-Y：20，30）
--------
$PG($BK(
	$IF(
		ifExpr=($ID(真))
		ifBlock=($BK(			
			$FN(
				name=($ID(X+Y))
				params=($NUM(20) $NUM(30))
			)
		))
		elseBlock=($BK(			
			$FN(
				name=($ID(X-Y))
				params=($NUM(20) $NUM(30))
			)
		))
	)
))

========
3. if-block & elseif blocks
--------
如果真：
    （X+Y：20，30）
再如A等于200：
    （X-Y：20，30）
再如A等于300：
    B为10；
注：「
				‘这是一个多行注释’
」
如果此C不为真：
    （ASDF）
--------
$PG($BK(
	$IF(
		ifExpr=($ID(真))
		ifBlock=($BK(
			$FN(
				name=($ID(X+Y))
				params=($NUM(20) $NUM(30))
			)
		))
		otherExpr[]=($EQ(
			L=($ID(A))
			R=($NUM(200))
		))
		otherBlock[]=($BK(
			$FN(
				name=($ID(X-Y))
				params=($NUM(20) $NUM(30))
			)
		))
		otherExpr[]=($EQ(
			L=($ID(A))
			R=($NUM(300))
		))
		otherBlock[]=($BK(
			$VA(
				target=($ID(B))
				assign=($NUM(10))
			)
		$))
	)
	$
	$IF(
		ifExpr=($ISN(L=($ID(C)) R=($ID(真))))
		ifBlock=($BK(
			$FN(name=($ID(ASDF)) params=())
		))
	)
))

========
4. if-elseif-else block
--------
如果真：
    （X+Y：20，30）
再如此A为100：
    （显示）
否则：
    （X-Y：20，30）
--------
$PG($BK(
	$IF(
		ifExpr=($ID(真))
		ifBlock=($BK(			
			$FN(
				name=($ID(X+Y))
				params=($NUM(20) $NUM(30))
			)
		))
		elseBlock=($BK(			
			$FN(
				name=($ID(X-Y))
				params=($NUM(20) $NUM(30))
			)
		))
		otherExpr[]=(
			$IS(L=($ID(A)) R=($NUM(100)))
		)
		otherBlock[]=($BK(
			$FN(
				name=($ID(显示))
				params=()
			)
		))		
	)
))
`

const arrayListCasesOK = `
========
1. empty array
--------
【】
--------
$PG($BK($ARR()))

========
2. empty hashmap
--------
【==】
--------
$PG($BK($HM()))

========
3. mixed string and decimal array
--------
【「MacBook Air 12"」 ， 2080， 3000】
--------
$PG($BK($ARR($STR(「MacBook Air 12"」) $NUM(2080) $NUM(3000))))

========
4. array with newline
--------
【
    「MacBook Air 12"」，
    2080， 
    3000
】
--------
$PG($BK($ARR($STR(「MacBook Air 12"」) $NUM(2080) $NUM(3000))))

========
5. array nest with array
--------
【
    「MacBook Air 12"」，
    2080， 
    【100，200，300】
】
--------
$PG($BK($ARR(
	$STR(「MacBook Air 12"」) 
	$NUM(2080) 
	$ARR($NUM(100) $NUM(200) $NUM(300))
)))

========
6. array nest with array nest with array
--------
【
    「MacBook Air 12"」，
    2080， 
    【100，200，300，
        【
            10000
        】
    】
】
--------
$PG($BK($ARR(
	$STR(「MacBook Air 12"」) 
	$NUM(2080) 
	$ARR($NUM(100) $NUM(200) $NUM(300) $ARR($NUM(10000)))
)))

========
7. a simple hashmap
--------
【
		「数学」 == 80，
		「语文」 == 90
】
--------
$PG($BK($HM(
	key[]=($STR(「数学」)) value[]=($NUM(80)) 
	key[]=($STR(「语文」)) value[]=($NUM(90))
)))

========
8. a hashmap nest with hashmap
--------
【
		「数学」 == 80，
		「语文」 == 【
				「阅读」 == 20，
				「听力」 == 30.5，
				「比例」 == 0.12345
		】
】
--------
$PG($BK($HM(
	key[]=($STR(「数学」)) value[]=($NUM(80)) 
	key[]=($STR(「语文」)) value[]=($HM(
		key[]=($STR(「阅读」)) value[]=($NUM(20))
		key[]=($STR(「听力」)) value[]=($NUM(30.5))
		key[]=($STR(「比例」)) value[]=($NUM(0.12345))
	))
)))
`

const stmtLineBreakCasesOK = `
========
1. a statement in oneline
--------
令香港记者为记者名为「张宝华」
--------
$PG($BK($VD(vars[]=($ID(香港记者)) expr[]=($VA(target=($ID(记者名)) assign=($STR(「张宝华」)))))))

========
2. a complete statement with comma list - 3 lines
--------
令树叶为「绿」，鲜花为「红」，
    雪花为「白」，
        墨水为「黑」
--------
$PG($BK($VD(vars[]=($ID(树叶)) expr[]=($STR(「绿」)) vars[]=($ID(鲜花)) expr[]=($STR(「红」)) vars[]=($ID(雪花)) expr[]=($STR(「白」)) vars[]=($ID(墨水)) expr[]=($STR(「黑」)))))

========
3. nested function calls with multiple lines
--------
（显示：
    「1」，（调用参数：200，300，
        4000，5000））
--------
$PG($BK($FN(name=($ID(显示)) params=($STR(「1」) $FN(name=($ID(调用参数)) params=($NUM(200) $NUM(300) $NUM(4000) $NUM(5000)))))))

========
4. multi-line hashmap
--------
令对象表为【
		1 == 「象」，
		2 == 「士」，
		3 == 「车」
】
--------
$PG($BK($VD(vars[]=($ID(对象表)) expr[]=($HM(key[]=($NUM(1)) value[]=($STR(「象」)) key[]=($NUM(2)) value[]=($STR(「士」)) key[]=($NUM(3)) value[]=($STR(「车」)))))))
`

const memberExprCasesOK = `
========
1. normal dot member
--------
天之涯
--------
$PG($BK(
	$MB(root=($ID(天)) type=(mID) object=($ID(涯)))
))

========
2. normal dot member (nested)
--------
雪花之天涯之海角
--------
$PG($BK(
	$MB(
		root=(
			$MB(root=($ID(雪花)) type=(mID) object=($ID(天涯)))
		)
		type=(mID)
		object=($ID(海角))
	)
))

========
3. call method
--------
雪花之（执行方法：A，B，C）
--------
$PG($BK(
	$MB(
		root=($ID(雪花))
		type=(mMethod)
		object=($FN(
			name=($ID(执行方法))
			params=($ID(A) $ID(B) $ID(C))
		))
	)
))

========
4. array index
--------
Array#123
--------
$PG($BK(
	$MB(root=($ID(Array)) type=(mIndex) object=($NUM(123)))
))

========
5. array index (using {})
--------
Array#{天之涯}
--------
$PG($BK(
	$MB(root=($ID(Array)) type=(mIndex) object=(
		$MB(root=($ID(天)) type=(mID) object=($ID(涯)))
	))
))

========
6. array index (nested)
--------
Array#20#30#{QR}
--------
$PG($BK(
	$MB(
		root=(
			$MB(
				root=(
					$MB(
						root=($ID(Array))
						type=(mIndex)
						object=($NUM(20))
					)
				)
				type=(mIndex)
				object=($NUM(30))
			)
		)
		type=(mIndex)
		object=($ID(QR))
	)
))

========
7. mix methods & members & indexes
--------
Array#10之首之（执行）
--------
$PG($BK(
	$MB(
		root=(
			$MB(
				root=(
					$MB(
						root=($ID(Array))
						type=(mIndex)
						object=($NUM(10))
					)
				)
				type=(mID)
				object=($ID(首))
			)
		)
		type=(mMethod)
		object=($FN(
			name=($ID(执行))
			params=()
		))
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

	for _, tt := range astCases {
		t.Run(tt.name, func(t *testing.T) {
			in := lex.NewTextStream(tt.input)
			l := lex.NewLexer(in)
			p := NewParser(l)

			block, err := p.Parse()
			pg := new(Program)
			pg.Content = block
			if err != nil {
				t.Errorf("expect no error, got error: %s", err.Display())
			} else {
				// compare with ast
				expect := StringifyAST(pg)
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
