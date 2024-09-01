package exec

import (
	"fmt"
	"testing"

	"github.com/DemoHn/Zn/pkg/io"
	"github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/syntax/zh"
	"github.com/DemoHn/Zn/pkg/value"
)

func setupMockContext() *runtime.Context {
	// init an empty context with init module
	// in this case, the initModule's name = main, lexer = nil (i.e no source code context)
	initModule := runtime.NewModule("main", nil)
	return runtime.NewContext(globalValues, initModule)
}

// a helper function to digest statement object from source code
// for writing testcases easier
func setupStmtFromCode(text string) (syntax.Statement, error) {
	in := io.NewByteStream([]byte(text))
	source, _ := in.ReadAll()
	p := syntax.NewParser(source, zh.NewParserZH())
	program, pErr := p.Parse()
	if pErr != nil {
		return nil, fmt.Errorf("syntax error on init:%v", pErr)
	}

	stmts := program.Content.Children
	if len(stmts) > 0 {
		// get first child statement ONLY
		return stmts[0], nil
	} else {
		return nil, fmt.Errorf("no suitable statement")
	}
}

func injectValuesToRootScope(c *runtime.Context, nameMap map[string]runtime.Element) {
	for k, v := range nameMap {
		c.BindSymbolDecl(runtime.NewIDName(k), v, false)
	}
}

func TestEvalWhileLoopStmt_OKCases(t *testing.T) {
	cases := []struct {
		name        string
		code        string
		initValue   map[string]runtime.Element
		expectLogic func(*runtime.Context, *testing.T)
	}{
		{
			name: "normal loop (5 times)",
			code: "每当真且A小于5：\n    A = A + 1\n    B = B + 10",
			initValue: map[string]runtime.Element{
				"A": value.NewNumber(1),
				"B": value.NewNumber(10),
			},
			expectLogic: func(ctx *runtime.Context, t *testing.T) {
				sym, _ := ctx.FindElement(runtime.NewIDName("B"))
				assertB := 50
				if sym.(*value.Number).GetValue() != float64(assertB) {
					t.Errorf("expect B (in root scope) = %f, got %f", sym.(*value.Number).GetValue(), float64(assertB))
				}
			},
		},
		{
			name: "break the loop via '结束循环' (before finish 4 times)",
			code: "每当A小于5：\n    A = A + 1；计数 = 计数 + 1\n    如果A == 4：\n        结束循环",
			initValue: map[string]runtime.Element{
				"A":  value.NewNumber(1),
				"计数": value.NewNumber(0),
			},
			expectLogic: func(ctx *runtime.Context, t *testing.T) {
				sym, _ := ctx.FindElement(runtime.NewIDName("计数"))
				assertB := 3
				if sym.(*value.Number).GetValue() != float64(assertB) {
					t.Errorf("expect B (in root scope) = %f, got %f", float64(assertB), sym.(*value.Number).GetValue())
				}
			},
		},
		{
			name: "break the loop via '结束循环' in inner ifs",
			code: `
每当A小于5：
    A = A + 1
    计数 = 计数 + 1
    如果A >= 3：
        如果A == 4：
            如果A > 0：
                结束循环`,
			initValue: map[string]runtime.Element{
				"A":  value.NewNumber(1),
				"计数": value.NewNumber(0),
			},
			expectLogic: func(ctx *runtime.Context, t *testing.T) {
				sym, _ := ctx.FindElement(runtime.NewIDName("计数"))
				assertB := 3
				if sym.(*value.Number).GetValue() != float64(assertB) {
					t.Errorf("expect B (in root scope) = %f, got %f", float64(assertB), sym.(*value.Number).GetValue())
				}
			},
		},
		{
			name: "continue the loop via '继续循环'",
			code: "每当A小于5：\n    A = A + 1；\n    如果A >= 4：\n        继续循环\n    计数 = 计数 + 1",
			initValue: map[string]runtime.Element{
				"A":  value.NewNumber(1),
				"计数": value.NewNumber(0),
			},
			expectLogic: func(ctx *runtime.Context, t *testing.T) {
				sym, _ := ctx.FindElement(runtime.NewIDName("计数"))
				assertB := 2
				if sym.(*value.Number).GetValue() != float64(assertB) {
					t.Errorf("expect B (in root scope) = %f, got %f", float64(assertB), sym.(*value.Number).GetValue())
				}
			},
		},
		{
			name: "'结束循环' only jumps out inner loop, not jump outer loop",
			code: `
每当A 小于 5：
	A = A + 1
	B = 0
	每当真：
		B = B + 1
		S = S + 1

		如果 B > 3：
			结束循环`,
			initValue: map[string]runtime.Element{
				"A": value.NewNumber(0),
				"B": value.NewNumber(0),
				"S": value.NewNumber(0),
			},
			expectLogic: func(ctx *runtime.Context, t *testing.T) {
				sym, _ := ctx.FindElement(runtime.NewIDName("S"))
				assertB := 20 /* loop for 4 * 5 = 20 times*/
				if sym.(*value.Number).GetValue() != float64(assertB) {
					t.Errorf("expect B (in root scope) = %f, got %f", float64(assertB), sym.(*value.Number).GetValue())
				}
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := setupMockContext()

			injectValuesToRootScope(ctx, tt.initValue)

			ss, err := setupStmtFromCode(tt.code)
			if err != nil {
				t.Errorf("FATAL got error:%v", err)
				return
			}

			// run the core function: evalWhileLoopStmt
			if err := evalWhileLoopStmt(ctx, ss.(*syntax.WhileLoopStmt)); err != nil {
				t.Errorf("expect OK, but got error: %v", err)
				return
			}

			tt.expectLogic(ctx, t)
		})
	}
}
