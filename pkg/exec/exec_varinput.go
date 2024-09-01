package exec

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/syntax/zh"
	"github.com/DemoHn/Zn/pkg/value"
)

// ExecVarInputs implements a special logic that initialize a value map
// before the main program starts. We name the predefined value map as "variable inputs" ("varInputs")
//
// The grammer of one "varInput" is a simplified version of the "varDeclareStmt", that is:
//   a. ‹symbol› = ‹value›   /  ‹symbol› 设为 ‹value›
//   b. ‹symbol› 成为 ‹type›：‹param1›、‹param2› ...
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
func ExecVarInputs(source string) (map[string]r.Element, error) {
	in := io.NewByteStream([]byte(source))

	sourceCode, err := in.ReadAll()
	if err != nil {
		return nil, zerr.ReadVarInputError(err)
	}

	// #2.  parse program
	parser := syntax.NewParser(sourceCode, zh.NewParserZH())
	vdStmt, err := parser.ParseVarInputs()
	if err != nil {
		return nil, zerr.NewErrorSLOT("解析预定义变量出现错误")
	}

	return evalVarInputStmt(vdStmt)
}

// evalVarInputStmt - similar to varDeclareStmt, but we add some limitations:
// 1. the assigned value must be basicType (ID1 = “文本” ok; ID1 = ID2 is not allowed)
// 2. 恒为 is same as 为 (since all predefined inputs are consts)
// 3. 成为 is NOT supported (since currently we have no way to fetch object class before actual code starts)
func evalVarInputStmt(node *syntax.VarDeclareStmt) (map[string]r.Element, error) {
	blankCtx := r.NewContext(globalValues, r.NewMainModule(nil))
	varInputMap := make(map[string]r.Element)

	for _, vpair := range node.AssignPair {
		switch vpair.Type {
		case syntax.VDTypeAssign, syntax.VDTypeAssignConst:
			obj, err := evalPrimeExpr(blankCtx, vpair.AssignExpr)
			if err != nil {
				return nil, err
			}

			for _, v := range vpair.Variables {
				vtag := v.GetLiteral()
				varInputMap[vtag] = value.DuplicateValue(obj)
			}
		default:
			return nil, zerr.NewErrorSLOT("不支持的赋值类型")
		}
	}

	return varInputMap, nil
}
