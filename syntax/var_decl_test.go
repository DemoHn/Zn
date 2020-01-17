package syntax

import (
	"reflect"
	"testing"

	"github.com/DemoHn/Zn/lex"
)

func TestVarDecl_OK(t *testing.T) {
	testcases := []struct {
		name     string
		tokens   string
		declVars []string
		exprType reflect.Type
	}{
		{
			name:     "one var decl",
			tokens:   "$40[令] $101[他] $41[为] $100[1200]",
			declVars: []string{"他"},
			exprType: reflect.TypeOf(&Number{}),
		},
		{
			name:     "two var decl (with var quote)",
			tokens:   "$40[令] $101[变量] $11[，] $91[此之代码] $41[为] $100[1200]",
			declVars: []string{"变量", "此之代码"},
			exprType: reflect.TypeOf(&Number{}),
		},
		{
			name:     "three var decl",
			tokens:   "$40[令] $101[变量] $11[，] $101[大新闻] $11[，] $101[名字] $41[为] $101[空]",
			declVars: []string{"变量", "大新闻", "名字"},
			exprType: reflect.TypeOf(&ID{}),
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			tkList := lex.ParseTokenStr(tt.tokens)
			p := NewParser(nil)
			p.InitMockToken(tkList)

			pg, err := p.Parse()
			if err != nil {
				t.Errorf("Parse() error! should have no error, got error")
				t.Error(err)
				return
			}

			// assert programNode
			if len(pg.Children) == 0 {
				t.Errorf("Parsed programNode should have at least 1 stmt, got 0!")
				return
			}
			stmt, ok := pg.Children[0].(*VarDeclareStmt)
			if !ok {
				t.Errorf("Parsed first item should be a *VarDeclareStmt!")
				return
			}

			// assert data
			vars := []string{}
			for _, item := range stmt.Variables {
				vars = append(vars, item.literal)
			}

			if !reflect.DeepEqual(vars, tt.declVars) {
				t.Errorf("DeclVars not same! expect: %v, got: %v", tt.declVars, vars)
			}

			// assert assignExpr
			if reflect.TypeOf(stmt.AssignExpr) != tt.exprType {
				t.Errorf("AssignExpr not same! expect: %v, got: %v", tt.exprType, reflect.TypeOf(stmt.AssignExpr))
			}
		})
	}
}
