package exec

import (
	"strings"

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
	switch name {
	case "新增", "添加":
		if err := validateExactParams(values, "any", "decimal"); err != nil {
			return nil, err
		}
		v := values[1].(*Decimal)
		idx, err := v.asInteger()
		if err != nil {
			return nil, err
		}
		ar.value = insertArrayValue(ar.value, idx, values[0])

		return ar, nil
	case "前增":
		if err := validateExactParams(values, "any"); err != nil {
			return nil, err
		}
		ar.value = insertArrayValue(ar.value, 0, values[0])
		return ar, nil
	case "后增":
		if err := validateExactParams(values, "any"); err != nil {
			return nil, err
		}
		ar.value = insertArrayValue(ar.value, len(ar.value), values[0])
		return ar, nil
	case "左移":
		v, newData := shiftArrayValue(ar.value, true)
		ar.value = newData
		return v, nil
	case "右移":
		v, newData := shiftArrayValue(ar.value, false)
		ar.value = newData
		return v, nil
	case "连接":
		// validate input array
		if err := validateAllParams(ar.value, "string"); err != nil {
			return nil, err
		}
		if err := validateExactParams(values, "string"); err != nil {
			return nil, err
		}
		var strArr = []string{}
		for _, v := range ar.value {
			item := v.(*String).value
			strArr = append(strArr, item)
		}

		connector := values[0].(*String).value
		finalStr := strings.Join(strArr, connector)

		return NewString(finalStr), nil
	}
	return nil, error.MethodNotFound(name)
}

////// method handlers
func insertArrayValue(target []Value, idx int, insertItem Value) []Value {
	result := []Value{}

	if idx >= len(target) {
		result = append(target, insertItem)
		return result
	}

	if idx < 0 {
		idx = len(target) + idx
	}
	result = append(result, target[:idx]...)
	result = append(result, insertItem)
	result = append(result, target[idx:]...)
	return result
}

func shiftArrayValue(target []Value, left bool) (Value, []Value) {
	if len(target) == 0 {
		return NewNull(), []Value{}
	}
	if left == true {
		return target[0], target[1:]
	}
	// shift right
	lastIdx := len(target) - 1
	return target[lastIdx], target[:lastIdx]
}
