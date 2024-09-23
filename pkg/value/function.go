package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type FuncExecutor = func(*r.Context, []r.Element) (r.Element, error)

type Function struct {
	name string
	// closureScope - when compiling a function, we will create an exclusive scope for this function to store values created inside the function.
	closureScope *r.Scope
	logicHandler FuncExecutor
}

func NewFunction(closureScope *r.Scope, executor FuncExecutor) *Function {
	return &Function{
		name:         "",
		closureScope: closureScope,
		logicHandler: executor,
	}
}

func (fn *Function) SetName(name string) *Function {
	fn.name = name
	return fn
}

// Exec - execute the Function Object - accepts input params, execute from closure executor and
// yields final result
func (fn *Function) Exec(c *r.Context, thisValue r.Element, params []r.Element) (r.Element, error) {
	// init scope
	module := c.GetCurrentModule()

	// add the pre-defined closure scope to current module
	if fn.closureScope != nil {
		module.AddScope(fn.closureScope)
		defer c.PopScope()
	}

	// create new scope for the function ITSELF
	fnScope := c.PushScope()
	defer c.PopScope()

	fnScope.SetThisValue(thisValue)

	fnLogiHandler := fn.logicHandler
	return fnLogiHandler(c, params)
}

// impl Value interface
// GetProperty -
func (fn *Function) GetProperty(c *r.Context, name string) (r.Element, error) {
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (fn *Function) SetProperty(c *r.Context, name string, value r.Element) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (fn *Function) ExecMethod(c *r.Context, name string, values []r.Element) (r.Element, error) {
	return nil, zerr.MethodNotFound(name)
}
