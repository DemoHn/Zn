package val

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/syntax"
)

type funcExecutor func(*ctx.Context, []ctx.Value) (ctx.Value, *error.Error)

// BuildClosureFromNode - create a closure (with default param handler logic)
// from Zn code (*syntax.BlockStmt). It's the constructor of 如何XX or (anoymous function in the future)
func BuildClosureFromNode(paramTags []*syntax.ParamItem, stmtBlock *syntax.BlockStmt) ctx.ClosureRef {
	var executor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
		// iterate block round I - function hoisting
		// NOTE: function hoisting means bind function definitions at the beginning
		// of execution so that even if "function execution" statement is before
		// "function definition" statement.
		for _, stmtI := range stmtBlock.Children {
			if v, ok := stmtI.(*syntax.FunctionDeclareStmt); ok {
				fn := BuildFunctionFromNode(v)
				if err := bindValue(ctx, v.FuncName.GetLiteral(), fn); err != nil {
					return nil, err
				}
			}
		}
		// iterate block round II - execution of rest code blocks
		for _, stmtII := range stmtBlock.Children {
			if _, ok := stmtII.(*syntax.FunctionDeclareStmt); !ok {
				if err := evalStatement(ctx, stmtII); err != nil {
					// if recv breaks
					if err.GetCode() == error.ReturnBreakSignal {
						if extra, ok := err.GetExtra().(ctx.Value); ok {
							return extra, nil
						}
					}
					return nil, err
				}
			}
		}
		return ctx.scope.returnValue, nil
	}

	var paramHandler = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
		// check param length
		if len(params) != len(paramTags) {
			return nil, error.MismatchParamLengthError(len(paramTags), len(params))
		}

		// bind params (as variable) to function scope
		for idx, paramVal := range params {
			param := paramTags[idx]
			// if param is NOT a reference type, then we need additionally
			// copy its value
			if !param.RefMark {
				paramVal = duplicateValue(paramVal)
			}
			paramName := param.ID.GetLiteral()
			if err := bindValue(ctx, paramName, paramVal); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	return ctx.ClosureRef{
		ParamHandler: paramHandler,
		Executor:     executor,
	}
}

// NewClosure - wraps a closure from native code (Golang code)
func NewClosure(paramHandler funcExecutor, executor funcExecutor) ctx.ClosureRef {
	return ctx.ClosureRef{
		ParamHandler: paramHandler,
		Executor:     executor,
	}
}

// Exec - execute a closure - accepts input params, execute from closure exeuctor and
// yields final result
func (cs *ctx.ClosureRef) Exec(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	if cs.ParamHandler != nil {
		if _, err := cs.ParamHandler(ctx, params); err != nil {
			return nil, err
		}
	}
	if cs.Executor == nil {
		return nil, error.NewErrorSLOT("执行逻辑不能为空")
	}
	// do execution
	return cs.Executor(ctx, params)
}
