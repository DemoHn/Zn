//// evalFunctionCall() is a core procedure, but the logic is very complicated. Thus we put the logic into a separate file.

package exec

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/value"
)

// compileFunction - create a Function object (with default param handler logic)
// from Zn code (*syntax.BlockStmt). It's the constructor of 如何XX or (anoymous function in the future)
func compileFunction(upperCtx *r.Context, node *syntax.FunctionDeclareStmt) *value.Function {
	var mainLogicHandler = func(c *r.Context, params []r.Element) (r.Element, error) {
		// 1. handle params (convert params -> varInput map inside the function scope)
		// 1.1 check param length
		if len(params) != len(node.ExecBlock.InputBlock) {
			return nil, zerr.MismatchParamLengthError(len(node.Name.GetLiteral()), len(params))
		}

		varInputMap := map[string]r.Element{}
		// 1.2 match param list with param name
		for idx, inputV := range node.ExecBlock.InputBlock {
			inputName, err := MatchIDName(inputV.ID)
			if err != nil {
				return nil, err
			}
			inputParam := inputName.GetLiteral()
			varInputMap[inputParam] = params[idx]
		}

		// 2. do eval exec block
		if err := evalExecBlock(c, node.ExecBlock, varInputMap); err != nil {
			return nil, err
		}

		return c.GetCurrentScope().GetReturnValue(), nil
	}

	return value.NewFunction(upperCtx.GetCurrentScope(), mainLogicHandler)
}

// （显示：A、B、C），得到D
func evalFunctionCall(c *r.Context, expr *syntax.FuncCallExpr) (r.Element, error) {
	// match & get funcName
	funcName, err := MatchIDName(expr.FuncName)
	if err != nil {
		return nil, err
	}

	// exec params
	params, err := exprsToValues(c, expr.Params)
	if err != nil {
		return nil, err
	}

	resultVal, err := execDirectFunction(c, funcName, params)
	if err != nil {
		return nil, err
	}

	// if exec function call succeed, then the non-nil `resultVal` will be exported.
	// However, if `得到 [someVar]` semi statement is defined, we will bind the `resultVal` to `someVar` first before ending the procedure.
	if expr.YieldResult != nil {
		// add yield result
		ytag, err := MatchIDName(expr.YieldResult)
		if err != nil {
			return nil, err
		}
		// bind yield result
		if err := c.BindScopeSymbolDecl(c.GetCurrentScope(), ytag, resultVal); err != nil {
			return nil, err
		}
	}

	// return result
	return resultVal, nil
}

func execMethodFunction(c *r.Context, root r.Element, funcName *r.IDName, params []r.Element) (r.Element, error) {
	pushCallstack := false

	if robj, ok := root.(*value.Object); ok {
		pushCallstack = true
		refModule := robj.GetRefModule()

		if refModule != nil {
			// append callInfo
			c.PushCallStack()
			c.SetCurrentRefModule(refModule)
		}
	}

	// create a new scope to denote a new 'thisValue'
	newScope := c.PushScope()
	defer c.PopScope()

	newScope.SetThisValue(root)
	// exec method
	elem, err := root.ExecMethod(c, funcName.GetLiteral(), params)
	// pop callInfo only when function execution succeed
	if err == nil && pushCallstack {
		c.PopCallStack()
	}
	return elem, err
}

// direct function: defined as standalone function instead of the method of
// a model
func execDirectFunction(c *r.Context, funcName *r.IDName, params []r.Element) (r.Element, error) {
	sym, err := c.FindSymbol(funcName)
	if err != nil {
		return nil, err
	}
	if sym.GetModule() != nil {
		// push callInfo
		c.PushCallStack()
		c.SetCurrentRefModule(sym.GetModule())
	}

	// assert value is function type
	fn, ok := sym.GetValue().(*value.Function)
	if !ok {
		return nil, zerr.InvalidFuncVariable(funcName.GetLiteral())
	}

	elem, err := fn.Exec(c, nil, params)
	if err == nil {
		c.PopCallStack()
	}
	return elem, err
}
