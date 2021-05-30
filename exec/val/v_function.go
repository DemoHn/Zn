package val

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

type funcExecutor = func(*ctx.Context, []ctx.Value) (ctx.Value, *error.Error)

// Function - 方法类
type Function struct {
	value ClosureRef
}

// ClosureRef - aka. Closure Exection Reference
// It's the structure of a closure which wraps execution logic.
// The executor could be either a bunch of code or some native code.
type ClosureRef struct {
	ParamHandler funcExecutor
	Executor     funcExecutor // closure execution logic
}

// NewFunction - new Zn native function
func NewFunction(name string, executor funcExecutor) *Function {
	closureRef := NewClosure(nil, executor)
	return &Function{closureRef}
}

// NewFunctionFromClosure -
func NewFunctionFromClosure(closure ClosureRef) *Function {
	return &Function{closure}
}

// NewClosure - wraps a closure from native code (Golang code)
func NewClosure(paramHandler funcExecutor, executor funcExecutor) ClosureRef {
	return ClosureRef{
		ParamHandler: paramHandler,
		Executor:     executor,
	}
}

// Exec - execute a closure - accepts input params, execute from closure exeuctor and
// yields final result
func (cs *ClosureRef) Exec(c *ctx.Context, thisValue ctx.Value, params []ctx.Value) (ctx.Value, *error.Error) {
	// init scope
	currentScope := c.GetScope()
	newScope := currentScope.CreateChildScope()
	newScope.SetThisValue(thisValue)
	// set and revert scope
	c.SetScope(newScope)
	defer c.SetScope(currentScope)

	if cs.ParamHandler != nil {
		if _, err := cs.ParamHandler(c, params); err != nil {
			return nil, err
		}
	}
	if cs.Executor == nil {
		return nil, error.NewErrorSLOT("执行逻辑不能为空")
	}
	// do execution
	return cs.Executor(c, params)
}

// GetValue -
func (fn *Function) GetValue() *ClosureRef {
	return &fn.value
}

// GetProperty -
func (fn *Function) GetProperty(c *ctx.Context, name string) (ctx.Value, *error.Error) {
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (fn *Function) SetProperty(c *ctx.Context, name string, value ctx.Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (fn *Function) ExecMethod(c *ctx.Context, name string, values []ctx.Value) (ctx.Value, *error.Error) {
	return nil, error.MethodNotFound(name)
}
