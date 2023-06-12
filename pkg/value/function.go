package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type FuncExecutor = func(*r.Context, []r.Value) (r.Value, error)

type Function struct {
	// closureScope - when compiling a function, we will create an exclusive scope for this function to store values created inside the function.
	closureScope  *r.Scope
	paramHandler  FuncExecutor
	logicHandlers []FuncExecutor
	// exceptionHandler - when an exception raise up, run this handler to catch the exception (like try...catch...)
	// TODO -
	exceptionHandler FuncExecutor
}

func NewFunction(executor FuncExecutor) *Function {
	logicHandlers := []FuncExecutor{}
	if executor != nil {
		logicHandlers = append(logicHandlers, executor)
	}

	return &Function{
		closureScope:     r.NewScope(),
		paramHandler:     nil,
		logicHandlers:    logicHandlers,
		exceptionHandler: nil,
	}
}

// setters
func (fn *Function) SetParamHandler(handler FuncExecutor) {
	fn.paramHandler = handler
}

// add logic handler
func (fn *Function) AddLogicHandler(handler FuncExecutor) {
	fn.logicHandlers = append(fn.logicHandlers, handler)
}

//// core function: Exec
// Exec - execute a closure - accepts input params, execute from closure executor and
// yields final result
func (fn *Function) Exec(c *r.Context, thisValue r.Value, params []r.Value) (r.Value, error) {
	// init scope
	fnScope := fn.closureScope
	module := c.GetCurrentModule()

	// add the pre-defined closure scope to current module
	module.AddScope(fnScope)
	defer c.PopScope()

	fnScope.SetThisValue(thisValue)

	if fn.paramHandler != nil {
		if _, err := fn.paramHandler(c, params); err != nil {
			return nil, err
		}
	}
	if len(fn.logicHandlers) == 0 {
		return nil, zerr.UnexpectedEmptyExecLogic()
	}
	// do execution
	var lastValue r.Value
	var execError error
	for _, handler := range fn.logicHandlers {
		lastValue, execError = handler(c, params)
		if execError != nil {
			return nil, execError
		}
	}
	return lastValue, nil
}

// impl Value interface
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
