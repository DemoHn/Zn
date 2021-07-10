package syntax

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DemoHn/Zn/lex"
)

var testFailSuites = []string{
	varDeclCasesFAIL,
	whileLoopCasesFAIL,
	funcCallCasesFAIL,
	arrayListCasesFAIL,
	funcDeclareCasesFAIL,
	importStmtCasesFAIL,
}

const varDeclCasesFAIL = `
========
1. non-identifiers as assigner (InvalidSyntax)
--------
注：第一行留给度娘

令某变量、另一变量、1240为1000
--------
code=2250 line=3 col=9

========

========
3. incomplete statement (InvalidSyntax)
--------
    
令某变量、另一变量
    【A，B】为100
    
--------
code=2250 line=2 col=5

========
4. invalid token (lexError)
--------
令锅为「锅」

令#'x为100
    
--------
code=2024 line=3 col=2

========
5. block indent unexpected
--------
    
令：
A为B，
B为C
    
--------
code=2250 line=2 col=1

========
6. block indent unexpected #2
--------
    
令：
        A为B，
        B为C
    
--------
code=2250 line=2 col=1

========
8. keyword only
--------
令
--------
code=2250 line=1 col=0
`

const whileLoopCasesFAIL = `
========
1. keyword only
--------
每当
--------
code=2250 line=1 col=2

========
2. keyword only #2
--------
每当：
--------
code=2250 line=1 col=2

========
3. missing true blocks
--------
每当真：
--------
code=2250 line=1 col=4

========
4. unexpected indents
--------
每当真：
每当又是真：
    （显示：「每当」）
--------
code=2250 line=2 col=0

========
5. trueExpr <- var declare stmt
--------
每当令变量为真：
    （显示：「变量为真」）
--------
code=2250 line=1 col=2


========
6. block statement fail
--------
每当变量为真：
    令数组为【【233】
--------
code=2250 line=2 col=10
`

const funcCallCasesFAIL = `
========
1. missing right paren
--------
（显示代码 等
--------
code=2250 line=1 col=7

========
2. func name not ID
--------
（{80000}）
--------
code=2250 line=1 col=0

========
3. func name includes otehr expr
--------
（显示时间，「2020」）
--------
code=2250 line=1 col=6

========
4. additional right paren
--------
（显示时间：「2020」））
--------
code=2250 line=1 col=13

========
5. additional pause comma
--------
（显示时间：「2020」、、500）
--------
code=2250 line=1 col=13
`

const arrayListCasesFAIL = `
========
2. missing right brancket
--------
【10，
--------
code=2250 line=1 col=4

========
3. incomplete map mark
--------
【「正定」 == 】
--------
code=2250 line=1 col=9

========
4. incomplete map mark #2
--------
【 == 「正定」 】
--------
code=2250 line=1 col=5

========
5. mixture of hashmap and array
--------
【 100，「正定」== 10 】
--------
code=2255 line=1 col=16
`

const funcDeclareCasesFAIL = `
========
1. no exec blocks
--------
如何搞个大新闻？
A为B
--------
code=2251 line=2 col=0

========
2. multiple 已知 blocks
--------
如何搞个大新闻？
    已知A
    已知B
--------
code=2250 line=3 col=0

========
3. more than one ref mark
--------
如何搞个大新闻？
    已知A、&&B	
--------
code=2250 line=2 col=4
`

const importStmtCasesFAIL = `
========
1. invalid quote
--------
导入“ASDF”
--------
code=2250 line=1 col=2

========
2. incomplete 之
--------
导入《ASDF》 之
令代码为空
--------
code=2250 line=1 col=9
`

type astFailCase struct {
	name     string
	input    string
	failInfo string
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
