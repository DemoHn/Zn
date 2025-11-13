package exec

import (
	"fmt"

	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/syntax/zh"
)

func evalVarAssignBlockText(vm *r.VM, blockText string) (r.ElementMap, error) {
	// #0. for empty string, skip parsing directly
	if len(blockText) == 0 {
		return r.ElementMap{}, nil
	}

	// #1. convert blockText into rune slice
	in := io.NewByteStream([]byte(blockText))

	runeStr, err := in.ReadAll()
	if err != nil {
		return nil, zerr.ReadVarInputError(err)
	}

	// #2. parse the rune slice into an AST
	parser := syntax.NewParser(runeStr, zh.NewParserZH())

	programAST, err := parser.Parse()
	if err != nil {
		return nil, zerr.NewErrorSLOT(fmt.Sprintf("解析 (varInput) 出现错误：‘%s’", blockText))
	}

	// #3. assert AST contains a bunch of varAssignExpr
	varAssignExprs, assertOK := assertASTIsVarAssignBlock(programAST)
	if !assertOK {
		return nil, zerr.NewErrorSLOT(fmt.Sprintf("解析 (varInput) 出现错误：‘%s’", blockText))
	}

	// #4. evaluate each varAssignExpr
	varInputMap := make(map[string]r.Element)
	for _, vpair := range varAssignExprs {
		// ensure vpair.TargetVar is a single ID
		idExpr, ok := vpair.TargetVar.(*syntax.ID)
		if !ok {
			return nil, zerr.NewErrorSLOT("目标变量必须是一个标识符")
		}

		// evaluate the assigned value
		evalResult, err := evalExpression(vm, vpair.AssignExpr)
		if err != nil {
			return nil, err
		}

		// add the result into varInputMap
		varInputMap[idExpr.GetLiteral()] = evalResult
	}

	return varInputMap, nil
}

// evalExpressionString - evaluate a valid expression string into runtime.Element
// e.g.: "10 + 8 * 3" -> value.Number(34)
func evalExpressionText(vm *r.VM, exprStr string) (r.Element, error) {
	// #1. convert exprStr into rune slice
	in := io.NewByteStream([]byte(exprStr))

	runeStr, err := in.ReadAll()
	if err != nil {
		return nil, zerr.ReadVarInputError(err)
	}
	// #2. parse the rune slice into an AST
	parser := syntax.NewParser(runeStr, zh.NewParserZH())

	programAST, err := parser.Parse()
	if err != nil {
		return nil, zerr.NewErrorSLOT(fmt.Sprintf("解析表达式出现错误：‘%s’", exprStr))
	}

	// #3. assert AST to be a single expression
	exprAST, assertOK := assertASTIsSingleExpr(programAST)
	if !assertOK {
		return nil, zerr.NewErrorSLOT(fmt.Sprintf("表达式‘%s’必须是一个单一表达式", exprStr))
	}

	// #4. evaluate expression
	return evalExpression(vm, exprAST)
}

func assertASTIsSingleExpr(ast *syntax.Program) (syntax.Expression, bool) {
	if ast.ExecBlock == nil {
		return nil, false
	}

	if ast.ExecBlock.StmtBlock == nil {
		return nil, false
	}

	if len(ast.ExecBlock.StmtBlock.Children) != 1 {
		return nil, false
	}

	expr := ast.ExecBlock.StmtBlock.Children[0]
	exprAST, ok := expr.(syntax.Expression)
	return exprAST, ok
}

func assertASTIsVarAssignBlock(ast *syntax.Program) ([]*syntax.VarAssignExpr, bool) {
	if ast.ExecBlock == nil {
		return nil, false
	}

	if ast.ExecBlock.StmtBlock == nil {
		return nil, false
	}

	resultExprs := make([]*syntax.VarAssignExpr, 0)
	for _, stmt := range ast.ExecBlock.StmtBlock.Children {
		switch expr := stmt.(type) {
		case *syntax.VarAssignExpr:
			resultExprs = append(resultExprs, expr)
		case *syntax.EmptyStmt:
			continue
		default:
			return nil, false
		}
	}

	return resultExprs, true
}

// ExecVarInputText implements a special logic that initialize a value map
// before the main program starts. We name the predefined value map as "variable inputs" ("varInputs")
//
// The grammer of one "varInput" is a simplified version of the "varDeclareStmt", that is:
//
//	a. ‹symbol› = ‹value›   /  ‹symbol› 设为 ‹value›
//	b. ‹symbol› 成为 ‹type›：‹param1›、‹param2› ...
//
// for example:
//
// [varInput]
// 客单价 = 28
//
// [mainCode]
// 输入客单价
// 令销量 = 300
// 输出客单价 * 销量  ->  8400
func ExecVarInputText(source string) (r.ElementMap, error) {
	vm := r.InitVM(globalValues)

	return evalVarAssignBlockText(vm, source)
}

func ExecExpressionInputText(exprStrMap map[string]string) (r.ElementMap, error) {
	vm := r.InitVM(globalValues)
	result := make(map[string]r.Element)
	for k, v := range exprStrMap {
		evalResult, err := evalExpressionText(vm, v)
		if err != nil {
			return nil, err
		}
		result[k] = evalResult
	}

	return result, nil
}
