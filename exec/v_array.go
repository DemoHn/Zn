package exec

import (
	"github.com/DemoHn/Zn/error"
)

// Array - represents for Zn's 数组型
type Array struct {
	value []Value
}

// NewArray - new array Value Object
func NewArray(value []Value) *Array {
	return &Array{value}
}

// GetProperty -
func (ar *Array) GetProperty(ctx *Context, name string) (Value, *error.Error) {
	switch name {
	case "和":

	case "差":
	case "积":
	case "商":
	case "首":
		if len(ar.value) == 0 {
			return NewNull(), nil
		}
		return ar.value[0], nil
	case "尾":
		if len(ar.value) == 0 {
			return NewNull(), nil
		}
		return ar.value[len(ar.value)-1], nil

	case "数目", "长度":
		l := len(ar.value)
		return NewDecimalFromInt(l, 0), nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (ar *Array) SetProperty(ctx *Context, name string, value Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (ar *Array) ExecMethod(ctx *Context, name string, values []Value) (Value, *error.Error) {
	return nil, error.MethodNotFound(name)
}
