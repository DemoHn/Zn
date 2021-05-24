package val

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/syntax"
)

// Function - 方法类
type Function struct {
	value ClosureRef
}

// NewFunction - new Zn native function
func NewFunction(name string, executor funcExecutor) *Function {
	closureRef := NewClosure(nil, executor)
	return &Function{closureRef}
}

// BuildFunctionFromNode -
func BuildFunctionFromNode(node *syntax.FunctionDeclareStmt) *Function {
	closureRef := BuildClosureFromNode(node.ParamList, node.ExecBlock)
	return &Function{closureRef}
}

// GetProperty -
func (fn *Function) GetProperty(ctx *Context, name string) (Value, *error.Error) {
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (fn *Function) SetProperty(ctx *Context, name string, value Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (fn *Function) ExecMethod(ctx *Context, name string, values []Value) (Value, *error.Error) {
	return nil, error.MethodNotFound(name)
}
