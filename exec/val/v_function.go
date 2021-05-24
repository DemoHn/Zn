package val

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
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
func (fn *Function) GetProperty(ctx *ctx.Context, name string) (ctx.Value, *error.Error) {
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (fn *Function) SetProperty(ctx *ctx.Context, name string, value ctx.Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (fn *Function) ExecMethod(ctx *ctx.Context, name string, values []ctx.Value) (ctx.Value, *error.Error) {
	return nil, error.MethodNotFound(name)
}
