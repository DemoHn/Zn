package val

import (
	"strings"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

// Array - represents for Zn's 数组型
type Array struct {
	value []ctx.Value
}

// NewArray - new array ctx.Value Object
func NewArray(value []ctx.Value) *Array {
	return &Array{value}
}

// GetValue -
func (ar *Array) GetValue() []ctx.Value {
	return ar.value
}

// AppendValue -
func (ar *Array) AppendValue(value ctx.Value) {
	ar.value = append(ar.value, value)
}

// GetProperty -
func (ar *Array) GetProperty(c *ctx.Context, name string) (ctx.Value, *error.Error) {
	switch name {
	case "和":
		return AddValueExecutor(c, ar.value)
	case "差":
		return SubValueExecutor(c, ar.value)
	case "积":
		return MulValueExecutor(c, ar.value)
	case "商":
		return DivValueExecutor(c, ar.value)
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
		result := []ctx.Value{}
		l := len(ar.value)
		for i := 0; i < l; i++ {
			result = append(result, result[l-1-i])
		}

		return NewArray(result), nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (ar *Array) SetProperty(c *ctx.Context, name string, value ctx.Value) *error.Error {
	switch name {
	case "首", "首项", "第一项":
		if len(ar.value) == 0 {
			result := []ctx.Value{value}
			ar.value = result
			return nil
		}
		ar.value[0] = value
		return nil
	case "尾", "末项", "最后项":
		if len(ar.value) == 0 {
			result := []ctx.Value{value}
			ar.value = result
			return nil
		}
		ar.value[len(ar.value)-1] = value
		return nil
	}
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (ar *Array) ExecMethod(c *ctx.Context, name string, values []ctx.Value) (ctx.Value, *error.Error) {
	switch name {
	case "新增", "添加":
		if err := ValidateExactParams(values, "any", "decimal"); err != nil {
			return nil, err
		}
		v := values[1].(*Decimal)
		idx, err := v.AsInteger()
		if err != nil {
			return nil, err
		}
		ar.value = insertArrayValue(ar.value, idx, values[0])

		return ar, nil
	case "前增":
		if err := ValidateExactParams(values, "any"); err != nil {
			return nil, err
		}
		ar.value = insertArrayValue(ar.value, 0, values[0])
		return ar, nil
	case "后增":
		if err := ValidateExactParams(values, "any"); err != nil {
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
		if err := ValidateAllParams(ar.value, "string"); err != nil {
			return nil, err
		}
		if err := ValidateExactParams(values, "string"); err != nil {
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
	case "合并":
		if err := ValidateAllParams(values, "array"); err != nil {
			return nil, err
		}

		result := []ctx.Value{}
		result = append(result, ar.value...)
		for _, v := range values {
			varr := v.(*Array).value
			result = append(result, varr...)
		}
		// update new array
		ar.value = result

		return NewArray(result), nil
	case "包含":
		result := false
		if err := ValidateExactParams(values, "any"); err != nil {
			return nil, err
		}
		for _, item := range ar.value {
			if res, err := CompareValues(item, values[0], CmpEq); err != nil {
				return nil, err
			} else if res {
				result = true
				break
			}
		}
		return NewBool(result), nil
	case "寻找":
		idx := -1
		if err := ValidateExactParams(values, "any"); err != nil {
			return nil, err
		}
		for i, item := range ar.value {
			if res, err := CompareValues(item, values[0], CmpEq); err != nil {
				return nil, err
			} else if res {
				idx = i
				break
			}
		}
		return NewDecimalFromInt(idx, 0), nil
	}
	return nil, error.MethodNotFound(name)
}

////// method handlers
func insertArrayValue(target []ctx.Value, idx int, insertItem ctx.Value) []ctx.Value {
	result := []ctx.Value{}

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

func shiftArrayValue(target []ctx.Value, left bool) (ctx.Value, []ctx.Value) {
	if len(target) == 0 {
		return NewNull(), []ctx.Value{}
	}
	if left {
		return target[0], target[1:]
	}
	// shift right
	lastIdx := len(target) - 1
	return target[lastIdx], target[:lastIdx]
}
