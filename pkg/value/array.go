package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)
// TODO: add methods
// Array - represents for Zn's 数组型
type Array struct {
	value []r.Value
}

// NewArray - new array ctx.Value Object
func NewArray(value []r.Value) *Array {
	return &Array{value}
}

// GetValue -
func (ar *Array) GetValue() []r.Value {
	return ar.value
}

// AppendValue -
func (ar *Array) AppendValue(value r.Value) {
	ar.value = append(ar.value, value)
}

// GetProperty -
func (ar *Array) GetProperty(c *r.Context, name string) (r.Value, error) {
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (ar *Array) SetProperty(c *r.Context, name string, value r.Value) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (ar *Array) ExecMethod(c *r.Context, name string, values []r.Value) (r.Value, error) {
	return nil, zerr.MethodNotFound(name)
}