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
	var executor = func(c *r.Context, params []r.Value) (r.Value, error) {
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

	var paramHandler = func(c *r.Context, params []r.Value) (r.Value, error) {
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
func evalFunctionCall(c *r.Context, expr *syntax.FuncCallExpr) (r.Value, error) {
	var zval *value.Function
	vtag := expr.FuncName.GetLiteral()

	// for a function call, if thisValue NOT FOUND, that means the target closure is a FUNCTION
	// instead of a METHOD (which is defined on class definition statement)
	//
	// If thisValue != nil, we will attempt to find closure from its method list;
	// then look up from scope's values.
	//
	// If thisValue == nil, we will look up target closure from scope's values directly.
	thisValue, _ := c.FindThisValue()

	// exec params first
	params, err := exprsToValues(c, expr.Params)
	if err != nil {
		return nil, err
	}

	// if thisValue exists, find ID from its method list
	if thisValue != nil {
		v, err := thisValue.ExecMethod(c, vtag, params)
		if err == nil {
			if expr.YieldResult != nil {
				// add yield result
				vtag := expr.YieldResult.GetLiteral()
				// bind yield result
				if err := c.BindScopeSymbolDecl(c.GetCurrentScope(), vtag, v); err != nil {
					return nil, err
				}
			}
			// return result
			return v, nil
		}

		if errX, ok := err.(*zerr.RuntimeError); ok {
			if errX.Code != zerr.ErrMethodNotFound {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// if function value not found from object scope, look up from local scope

	// find function definition
	v, err := c.FindSymbol(vtag)
	if err != nil {
		return nil, err
	}
	// assert value
	zval, ok := v.(*value.Function)
	if !ok {
		return nil, zerr.InvalidFuncVariable(vtag)
	}

	// exec function call via its ClosureRef
	v2, err := zval.Exec(c, thisValue, params)
	if err != nil {
		return nil, err
	}

	if expr.YieldResult != nil {
		// add yield result
		ytag := expr.YieldResult.GetLiteral()
		// bind yield result
		if err := c.BindScopeSymbolDecl(c.GetCurrentScope(), ytag, v2); err != nil {
			return nil, err
		}
	}

	// return result
	return v2, nil
}
