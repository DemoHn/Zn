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
		return addValueExecutor(ctx, ar.value)
	case "差":
		return subValueExecutor(ctx, ar.value)
	case "积":
		return mulValueExecutor(ctx, ar.value)
	case "商":
		return divValueExecutor(ctx, ar.value)
	case "首", "首项", "第一项":
		if len(ar.value) == 0 {
			return NewNull(), nil
		}
		return ar.value[0], nil
	case "尾", "末项", "最后项":
		if len(ar.value) == 0 {
			return NewNull(), nil
		}
		return ar.value[len(ar.value)-1], nil
	case "数目", "长度":
		l := len(ar.value)
		return NewDecimalFromInt(l, 0), nil
	case "文本*":
		valStr := StringifyValue(ar)
		return NewString(valStr), nil
	case "逆":
		result := []Value{}
		l := len(ar.value)
		for i := 0; i < l; i++ {
			result = append(result, result[l-1-i])
		}

		return NewArray(result), nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (ar *Array) SetProperty(ctx *Context, name string, value Value) *error.Error {
	switch name {
	case "首", "首项", "第一项":
		if len(ar.value) == 0 {
			result := []Value{value}
			ar.value = result
			return nil
		}
		ar.value[0] = value
		return nil
	case "尾", "末项", "最后项":
		if len(ar.value) == 0 {
			result := []Value{value}
			ar.value = result
			return nil
		}
		ar.value[len(ar.value)-1] = value
		return nil
	}
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (ar *Array) ExecMethod(ctx *Context, name string, values []Value) (Value, *error.Error) {
	return nil, error.MethodNotFound(name)
}
