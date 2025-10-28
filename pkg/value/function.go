package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type Function struct {
	name         string
	logicHandler r.FuncExecutor
}

func NewFunction(executor r.FuncExecutor) *Function {
	return &Function{
		name:         "",
		logicHandler: executor,
	}
}

func (fn *Function) SetName(name string) *Function {
	fn.name = name
	return fn
}

// Exec - execute the Function Object - accepts input params, execute from closure executor and
// yields final result
func (fn *Function) Exec(thisValue r.Element, params []r.Element) (r.Element, error) {
	fnLogicHandler := fn.logicHandler
	return fnLogicHandler(params)
}

// impl Value interface
// GetProperty -
func (fn *Function) GetProperty(name string) (r.Element, error) {
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (fn *Function) SetProperty(name string, value r.Element) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (fn *Function) ExecMethod(name string, values []r.Element) (r.Element, error) {
	return nil, zerr.MethodNotFound(name)
}
