package exec

import "github.com/DemoHn/Zn/error"

// Null -
type Null struct{}

// NewNull -
func NewNull() *Null {
	return &Null{}
}

// GetProperty -
func (nl *Null) GetProperty(ctx *Context, name string) (Value, *error.Error) {
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (nl *Null) SetProperty(ctx *Context, name string, value Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (nl *Null) ExecMethod(ctx *Context, name string, values []Value) (Value, *error.Error) {
	return nil, error.MethodNotFound(name)
}
