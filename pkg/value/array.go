package value

import (
	"math"
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type arrayGetterFunc func(*Array, *r.Context) (r.Element, error)
type arraySetterFunc func(*Array, *r.Context, r.Element) error
type arrayMethodFunc func(*Array, *r.Context, []r.Element) (r.Element, error)

// Array - represents for Zn's 数组型
type Array struct {
	value []r.Element
	*r.ElementModel
}

// NewArray - new array r.Value Object
func NewArray(value []r.Element) *Array {
	arr := &Array{value, r.NewElementModel()}

	//// init getters & setters & methods
	arr.RegisterGetter("文本", arr.arrayGetText)
	arr.RegisterGetter("首项", arr.arrayGetFirstItem)
	arr.RegisterGetter("末项", arr.arrayGetLastItem)
	arr.RegisterGetter("数目", arr.arrayGetLength)
	arr.RegisterGetter("长度", arr.arrayGetLength)
	arr.RegisterGetter("逆序", arr.arrayGetReverse)

	arr.RegisterSetter("首项", arr.arraySetFirstItem)
	arr.RegisterSetter("末项", arr.arraySetLastItem)

	arr.RegisterMethod("新增", arr.arrayExecInsert)
	arr.RegisterMethod("添加", arr.arrayExecInsert)
	arr.RegisterMethod("前增", arr.arrayExecPrepend)
	arr.RegisterMethod("后增", arr.arrayExecAppend)
	arr.RegisterMethod("左移", arr.arrayExecShift)
	arr.RegisterMethod("右移", arr.arrayExecPop)
	arr.RegisterMethod("拼接", arr.arrayExecJoin)
	arr.RegisterMethod("合并", arr.arrayExecMerge)
	arr.RegisterMethod("包含", arr.arrayExecContains)
	arr.RegisterMethod("寻找", arr.arrayExecFind)
	arr.RegisterMethod("交换", arr.arrayExecSwap)
	return arr
}

// GetValue -
func (ar *Array) GetValue() []r.Element {
	return ar.value
}

// AppendValue -
func (ar *Array) AppendValue(value r.Element) {
	ar.value = append(ar.value, value)
}

//// getters, setters & methods

// getters
func (ar *Array) arrayGetText(c *r.Context) (r.Element, error) {
	return NewString(StringifyValue(ar)), nil
}

func (ar *Array) arrayGetFirstItem(c *r.Context) (r.Element, error) {
	if len(ar.value) == 0 {
		return NewNull(), nil
	}
	return ar.value[0], nil
}

func (ar *Array) arrayGetLastItem(c *r.Context) (r.Element, error) {
	if len(ar.value) == 0 {
		return NewNull(), nil
	}
	return ar.value[len(ar.value)-1], nil
}

func (ar *Array) arrayGetLength(c *r.Context) (r.Element, error) {
	l := len(ar.value)
	return NewNumber(float64(l)), nil
}

func (ar *Array) arrayGetReverse(c *r.Context) (r.Element, error) {
	var result []r.Element
	l := len(ar.value)
	for i := 0; i < l; i++ {
		result = append(result, ar.value[l-1-i])
	}

	return NewArray(result), nil
}

func (ar *Array) arrayGetAddResult(c *r.Context) (r.Element, error) {
	if err := ValidateLeastParams(ar.value, "number+"); err != nil {
		return nil, err
	}

	var sum float64 = 0
	// validate types
	for _, param := range ar.value {
		vparam := param.(*Number)
		sum = sum + vparam.value
	}

	return NewNumber(sum), nil
}

func (ar *Array) arrayGetSubResult(c *r.Context) (r.Element, error) {
	if err := ValidateLeastParams(ar.value, "number+"); err != nil {
		return nil, err
	}

	var sum float64 = 0

	// validate types
	for idx, param := range ar.value {
		vparam := param.(*Number)
		if idx == 0 {
			sum = vparam.value
		} else {
			sum = sum - vparam.value
		}
	}

	return NewNumber(sum), nil
}

func (ar *Array) arrayGetMulResult(c *r.Context) (r.Element, error) {
	if err := ValidateLeastParams(ar.value, "number+"); err != nil {
		return nil, err
	}

	var sum float64 = 0

	// validate types
	for idx, param := range ar.value {
		vparam := param.(*Number)
		if idx == 0 {
			sum = vparam.value
		} else {
			sum = sum * vparam.value
		}
	}

	return NewNumber(sum), nil
}

func (ar *Array) arrayGetDivResult(c *r.Context) (r.Element, error) {
	if err := ValidateLeastParams(ar.value, "number+"); err != nil {
		return nil, err
	}

	var sum float64 = 0

	// validate types
	for idx, param := range ar.value {
		vparam := param.(*Number)
		if idx == 0 {
			sum = vparam.value
		} else {
			if vparam.value == 0 {
				return nil, zerr.ArithDivZero()
			}
			sum = sum / vparam.value
		}
	}

	return NewNumber(sum), nil
}

// setters
func (ar *Array) arraySetFirstItem(c *r.Context, value r.Element) error {
	if len(ar.value) == 0 {
		result := []r.Element{value}
		ar.value = result
		return nil
	}
	ar.value[0] = value
	return nil
}

func (ar *Array) arraySetLastItem(c *r.Context, value r.Element) error {
	if len(ar.value) == 0 {
		result := []r.Element{value}
		ar.value = result
		return nil
	}
	ar.value[len(ar.value)-1] = value
	return nil
}

// methods
func (ar *Array) arrayExecInsert(c *r.Context, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "any", "number"); err != nil {
		return nil, err
	}
	v := values[1].(*Number)
	ar.value = insertArrayValue(ar.value, int(v.value), values[0])

	return ar, nil
}

func (ar *Array) arrayExecPrepend(c *r.Context, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "any"); err != nil {
		return nil, err
	}
	ar.value = insertArrayValue(ar.value, 0, values[0])
	return ar, nil
}

func (ar *Array) arrayExecAppend(c *r.Context, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "any"); err != nil {
		return nil, err
	}
	ar.value = insertArrayValue(ar.value, len(ar.value), values[0])
	return ar, nil
}

func (ar *Array) arrayExecShift(c *r.Context, values []r.Element) (r.Element, error) {
	v, newData := shiftArrayValue(ar.value, true)
	ar.value = newData
	return v, nil
}

func (ar *Array) arrayExecPop(c *r.Context, values []r.Element) (r.Element, error) {
	v, newData := shiftArrayValue(ar.value, false)
	ar.value = newData
	return v, nil
}

func (ar *Array) arrayExecJoin(c *r.Context, values []r.Element) (r.Element, error) {
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

func (ar *Array) arrayExecMerge(c *r.Context, values []r.Element) (r.Element, error) {
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

func (ar *Array) arrayExecContains(c *r.Context, values []r.Element) (r.Element, error) {
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

func (ar *Array) arrayExecFind(c *r.Context, values []r.Element) (r.Element, error) {
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

func (ar *Array) arrayExecSwap(c *r.Context, values []r.Element) (r.Element, error) {
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

////// method handlers
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
