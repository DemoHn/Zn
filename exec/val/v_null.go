package val

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

// Null -
type Null struct{}

// NewNull -
func NewNull() *Null {
	return &Null{}
}

// GetProperty -
func (nl *Null) GetProperty(ctx *ctx.Context, name string) (ctx.Value, *error.Error) {
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (nl *Null) SetProperty(ctx *ctx.Context, name string, value ctx.Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (nl *Null) ExecMethod(ctx *ctx.Context, name string, values []ctx.Value) (ctx.Value, *error.Error) {
	return nil, error.MethodNotFound(name)
}
