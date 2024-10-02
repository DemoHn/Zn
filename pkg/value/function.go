package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

const (
	EVConstThisVariableName = "此"
)

type FuncExecutor = func(*r.Context, []r.Element) (r.Element, error)

type Function struct {
	name string
	// outerScope - when compiling an "inner" function (i.e. a "child" function defined inside another "parent" function), we will copy the outer scope into the function to impl the Closure feature
	outerScope   *r.Scope
	logicHandler FuncExecutor
}

func NewFunction(outerScope *r.Scope, executor FuncExecutor) *Function {
	return &Function{
		name:         "",
		outerScope:   outerScope,
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
	if fn.outerScope != nil {
		module.AddScope(fn.outerScope)
		defer c.PopScope()
	}

	// create new scope for the function ITSELF
	fnScope := c.PushScope()
	defer c.PopScope()

	if thisValue != nil {
		// set thisValue of current scope
		fnScope.SetThisValue(thisValue)

		// add a const variable "此" to represent "$this"
		// usage: 以此（调用某方法：XX、YY）
		if err := c.BindSymbolConst(r.NewIDName(EVConstThisVariableName), thisValue); err != nil {
			return nil, err
		}
	}

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
