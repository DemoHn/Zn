package value

import (
	"math"
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type arrayGetterFunc func(*Array) (r.Element, error)
type arraySetterFunc func(*Array, r.Element) error
type arrayMethodFunc func(*Array, []r.Element) (r.Element, error)

// Array - represents for Zn's 数组型
type Array struct {
	value []r.Element
}

// NewArray - new array r.Value Object
func NewArray(value []r.Element) *Array {
	return &Array{value}
}

func NewEmptyArray() *Array {
	return &Array{
		value: []r.Element{},
	}
}

func (ar *Array) Length() int {
	return len(ar.value)
}

// GetValue -
func (ar *Array) GetValue() []r.Element {
	return ar.value
}

// AppendValue -
func (ar *Array) AppendValue(value r.Element) {
	ar.value = append(ar.value, value)
}

// GetProperty -
func (ar *Array) GetProperty(name string) (r.Element, error) {
	arrayGetterMap := map[string]arrayGetterFunc{
		"文本": arrayGetText,
		"首项": arrayGetFirstItem,
		"末项": arrayGetLastItem,
		"数目": arrayGetLength,
		"长度": arrayGetLength,
		"逆序": arrayGetReverse,
	}

	if fn, ok := arrayGetterMap[name]; ok {
		return fn(ar)
	}
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (ar *Array) SetProperty(name string, value r.Element) error {
	arraySetterMap := map[string]arraySetterFunc{
		"首项": arraySetFirstItem,
		"末项": arraySetLastItem,
	}

	if fn, ok := arraySetterMap[name]; ok {
		return fn(ar, value)
	}
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (ar *Array) ExecMethod(name string, values []r.Element) (r.Element, error) {
	arrayMethodMap := map[string]arrayMethodFunc{
		"新增": arrayExecInsert,
		"添加": arrayExecInsert,
		"前增": arrayExecPrepend,
		"后增": arrayExecAppend,
		"左移": arrayExecShift,
		"右移": arrayExecPop,
		"拼接": arrayExecJoin,
		"合并": arrayExecMerge,
		"包含": arrayExecContains,
		"寻找": arrayExecFind,
		"交换": arrayExecSwap,
	}

	if fn, ok := arrayMethodMap[name]; ok {
		return fn(ar, values)
	}
	return nil, zerr.MethodNotFound(name)
}

//// getters, setters & methods

// getters
func arrayGetText(ar *Array) (r.Element, error) {
	return NewString(StringifyValue(ar)), nil
}

func arrayGetFirstItem(ar *Array) (r.Element, error) {
	if len(ar.value) == 0 {
		return NewNull(), nil
	}
	return ar.value[0], nil
}

func arrayGetLastItem(ar *Array) (r.Element, error) {
	if len(ar.value) == 0 {
		return NewNull(), nil
	}
	return ar.value[len(ar.value)-1], nil
}

func arrayGetLength(ar *Array) (r.Element, error) {
	l := len(ar.value)
	return NewNumber(float64(l)), nil
}

func arrayGetReverse(ar *Array) (r.Element, error) {
	var result []r.Element
	l := len(ar.value)
	for i := 0; i < l; i++ {
		result = append(result, ar.value[l-1-i])
	}

	return NewArray(result), nil
}

// setters
func arraySetFirstItem(ar *Array, value r.Element) error {
	if len(ar.value) == 0 {
		result := []r.Element{value}
		ar.value = result
		return nil
	}
	ar.value[0] = value
	return nil
}

func arraySetLastItem(ar *Array, value r.Element) error {
	if len(ar.value) == 0 {
		result := []r.Element{value}
		ar.value = result
		return nil
	}
	ar.value[len(ar.value)-1] = value
	return nil
}

// methods
func arrayExecInsert(ar *Array, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "any", "number"); err != nil {
		return nil, err
	}
	v := values[1].(*Number)
	ar.value = insertArrayValue(ar.value, int(v.value), values[0])

	return ar, nil
}

func arrayExecPrepend(ar *Array, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "any"); err != nil {
		return nil, err
	}
	ar.value = insertArrayValue(ar.value, 0, values[0])
	return ar, nil
}

func arrayExecAppend(ar *Array, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "any"); err != nil {
		return nil, err
	}
	ar.value = insertArrayValue(ar.value, len(ar.value), values[0])
	return ar, nil
}

func arrayExecShift(ar *Array, values []r.Element) (r.Element, error) {
	v, newData := shiftArrayValue(ar.value, true)
	ar.value = newData
	return v, nil
}

func arrayExecPop(ar *Array, values []r.Element) (r.Element, error) {
	v, newData := shiftArrayValue(ar.value, false)
	ar.value = newData
	return v, nil
}

func arrayExecJoin(ar *Array, values []r.Element) (r.Element, error) {
	// validate input array
	if err := ValidateAllParams(ar.value, "string"); err != nil {
		return nil, err
	}
	if err := ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	var strArr []string
	for _, v := range ar.value {
		item := v.(*String).value
		strArr = append(strArr, item)
	}

	connector := values[0].(*String).value
	finalStr := strings.Join(strArr, connector)

	return NewString(finalStr), nil
}

func arrayExecMerge(ar *Array, values []r.Element) (r.Element, error) {
	if err := ValidateAllParams(values, "array"); err != nil {
		return nil, err
	}

	var result []r.Element
	result = append(result, ar.value...)
	for _, v := range values {
		varr := v.(*Array).value
		result = append(result, varr...)
	}
	// update new array
	ar.value = result

	return NewArray(result), nil
}

func arrayExecContains(ar *Array, values []r.Element) (r.Element, error) {
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
}

func arrayExecFind(ar *Array, values []r.Element) (r.Element, error) {
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
	return NewNumber(float64(idx)), nil
}

func arrayExecSwap(ar *Array, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "number", "number"); err != nil {
		return nil, err
	}
	// check if all indexes in the range of array
	l := len(ar.value)

	cursor0 := int(math.Floor(values[0].(*Number).GetValue()) - 1)
	cursor1 := int(math.Floor(values[1].(*Number).GetValue()) - 1)

	if cursor0 < 0 || cursor0 >= l {
		return nil, zerr.IndexOutOfRange()
	}
	if cursor1 < 0 || cursor1 >= l {
		return nil, zerr.IndexOutOfRange()
	}
	// swap item
	tmp := ar.value[cursor0]
	ar.value[cursor0] = ar.value[cursor1]
	ar.value[cursor1] = tmp

	return ar, nil
}

// //// method handlers
func insertArrayValue(target []r.Element, idx int, insertItem r.Element) []r.Element {
	var result []r.Element

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

func shiftArrayValue(target []r.Element, left bool) (r.Element, []r.Element) {
	if len(target) == 0 {
		return NewNull(), []r.Element{}
	}
	if left {
		return target[0], target[1:]
	}
	// shift right
	lastIdx := len(target) - 1
	return target[lastIdx], target[:lastIdx]
}
