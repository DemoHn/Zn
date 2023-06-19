//// evalFunctionCall() is a core procedure, but the logic is very complicated. Thus we put the logic into a separate file.

package exec

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/value"
)

// compileFunction - create a function (with default param handler logic)
// from Zn code (*syntax.BlockStmt). It's the constructor of 如何XX or (anoymous function in the future)
func compileFunction(upperContext *r.Context, paramTags []*syntax.ParamItem, stmtBlock *syntax.BlockStmt) *value.Function {
	var executor = func(c *r.Context, params []r.Element) (r.Element, error) {
		// iterate block round I - function hoisting
		// NOTE: function hoisting means bind function definitions at the beginning
		// of execution so that even if "function execution" statement is before
		// "function definition" statement.
		for _, stmtI := range stmtBlock.Children {
			if v, ok := stmtI.(*syntax.FunctionDeclareStmt); ok {
				fn := compileFunction(c, v.ParamList, v.ExecBlock)
				if err := c.BindSymbol(v.FuncName.GetLiteral(), fn); err != nil {
					return nil, err
				}
			}
		}
		// iterate block round II - execution of rest code blocks
		for _, stmtII := range stmtBlock.Children {
			if _, ok := stmtII.(*syntax.FunctionDeclareStmt); !ok {
				if err := evalStatement(c, stmtII); err != nil {
					return extractSignalValue(err, zerr.SigTypeReturn)
				}
			}
		}
		return c.GetCurrentScope().GetReturnValue(), nil
	}

	var paramHandler = func(c *r.Context, params []r.Element) (r.Element, error) {
		// check param length
		if len(params) != len(paramTags) {
			return nil, zerr.MismatchParamLengthError(len(paramTags), len(params))
		}

		// bind params (as variable) to function scope
		for idx, paramVal := range params {
			param := paramTags[idx]
			// if param is NOT a reference type, then we need additionally
			// copy its value
			if !param.RefMark {
				paramVal = value.DuplicateValue(paramVal)
			}
			paramName := param.ID.GetLiteral()
			if err := c.BindSymbol(paramName, paramVal); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	fn := value.NewFunction(upperContext.GetCurrentScope(), executor)
	fn.SetParamHandler(paramHandler)
	return fn
}

// （显示：A、B、C），得到D
func evalFunctionCall(c *r.Context, expr *syntax.FuncCallExpr) (r.Element, error) {
	var resultVal r.Element
	var err error

	// eval method
	funcName := expr.FuncName.GetLiteral()

	// exec params
	params, err := exprsToValues(c, expr.Params)
	if err != nil {
		return nil, err
	}
	// for a function call, if thisValue NOT FOUND, that means the target closure is a FUNCTION
	// instead of a METHOD (which is defined on class definition statement)
	//
	// If thisValue != nil, we will attempt to find closure from its method list;
	// then look up from scope's values.
	//
	// If thisValue == nil, we will look up target closure from scope's values directly.
	thisValue := c.GetThisValue()

	// if thisValue exists, find ID from its method list
	/* example:
	如何外部方法？
		输出「这是外部方法」

	定义示例类：
		如何内部类方法？
			输出「内部类方法」

		如何方法B？
			（内部类方法）  //  等价于 `以 [某示例类对象]（内部类方法）`
			（外部方法）   //  1. 先示例类中寻找「外部方法」，如同调用 `以 [某示例类对象]（内部类方法）` 2. 寻找无果（抛出 zerr.ErrMethodNotFound 错误）后再去全局作用域寻找「外部方法」的方法对象并调用其逻辑
	*/
	if thisValue != nil {
		resultVal, err = execMethodFunction(c, thisValue, funcName, params)
		if err != nil {
			if errX, ok := err.(*zerr.RuntimeError); ok {
				if errX.Code == zerr.ErrMethodNotFound {
					// fallback to execute direct function
					resultVal, err = execDirectFunction(c, funcName, params)
				}
			}
		}
	} else {
		// no parent object denoted, execute function directly
		resultVal, err = execDirectFunction(c, funcName, params)
	}

	if err != nil {
		return nil, err
	}

	// if exec function call succeed, then the non-nil `resultVal` will be exported.
	// However, if `得到 [someVar]` semi statement is defined, we will bind the `resultVal` to `someVar` first before ending the procedure.
	if expr.YieldResult != nil {
		// add yield result
		ytag := expr.YieldResult.GetLiteral()
		// bind yield result
		if err := c.BindScopeSymbolDecl(c.GetCurrentScope(), ytag, resultVal); err != nil {
			return nil, err
		}
	}

	// return result
	return resultVal, nil
}

func execMethodFunction(c *r.Context, root r.Element, funcName string, params []r.Element) (r.Element, error) {
	// create a new scope to denote a new 'thisValue'
	newScope := c.PushScope()
	defer c.PopScope()

	newScope.SetThisValue(root)
	// exec method
	return root.ExecMethod(c, funcName, params)
}

// direct function: defined as standalone function instead of the method of
// a model
func execDirectFunction(c *r.Context, funcName string, params []r.Element) (r.Element, error) {
	v, err := c.FindSymbol(funcName)
	if err != nil {
		return nil, err
	}
	// assert value is function type
	funcVal, ok := v.(*value.Function)
	if !ok {
		return nil, zerr.InvalidFuncVariable(funcName)
	}

	return funcVal.Exec(c, nil, params)
}
