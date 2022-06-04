package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type funcExecutor = func(*r.Context, []r.Value) (r.Value, error)

// Function - 方法类
type Function struct {
	value *ClosureRef
}

// ClosureRef - aka. Closure Execution Reference
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
func NewFunctionFromClosure(closure *ClosureRef) *Function {
	return &Function{closure}
}

// NewClosure - wraps a closure from native code (Golang code)
func NewClosure(paramHandler funcExecutor, executor funcExecutor) *ClosureRef {
	return &ClosureRef{
		ParamHandler: paramHandler,
		Executor:     executor,
	}
}

// Exec - execute a closure - accepts input params, execute from closure executor and
// yields final result
func (cs *ClosureRef) Exec(c *r.Context, thisValue r.Value, params []r.Value) (r.Value, error) {
	// init scope
	newScope := c.PushChildScope()
	defer c.PopScope()
	newScope.SetThisValue(thisValue)

	if cs.ParamHandler != nil {
		if _, err := cs.ParamHandler(c, params); err != nil {
			return nil, err
		}
	}
	if cs.Executor == nil {
		return nil, zerr.NewErrorSLOT("执行逻辑不能为空")
	}
	// do execution
	val, err := cs.Executor(c, params)
	if err != nil {
		return nil, err
	}

	return val, nil
}

// GetValue -
func (fn *Function) GetValue() *ClosureRef {
	return fn.value
}

// GetProperty -
func (fn *Function) GetProperty(c *r.Context, name string) (r.Value, error) {
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (fn *Function) SetProperty(c *r.Context, name string, value r.Value) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (fn *Function) ExecMethod(c *r.Context, name string, values []r.Value) (r.Value, error) {
	return nil, zerr.MethodNotFound(name)
}