package value

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type numGetterFunc func(*Number) (r.Element, error)
type numMethodFunc func(*Number, []r.Element) (r.Element, error)

type Number struct {
	value float64
}

// NewNumber - create new number object (plain float64)
func NewNumber(value float64) *Number {
	return &Number{value}
}

func NewNumberFromString(value string) (*Number, error) {
	v := strings.ReplaceAll(value, ",", "")
	v = strings.Replace(v, "*^", "", 1)
	v = strings.Replace(v, "*10^", "e", 1)

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, err
	}
	return NewNumber(f), nil
}

// String -
func (n *Number) String() string {
	return fmt.Sprintf("%v", n.value)
}

// Construct - make Number construtable
// TODO: introduce TypeElement
func (n *Number) Construct(params []r.Element) (r.Element, error) {
	if err := ValidateExactParams(params, "number"); err != nil {
		return nil, err
	}

	p0 := params[0].(*Number)
	return p0, nil
}

// GetValue -
func (n *Number) GetValue() float64 {
	return n.value
}

// GetProperty -
func (n *Number) GetProperty(name string) (r.Element, error) {
	numGetterMap := map[string]numGetterFunc{
		"文本":  numGetText,
		"平方":  numGetSquare,
		"立方":  numGetCube,
		"平方根": numGetSquareRoot,
	}
	if fn, ok := numGetterMap[name]; ok {
		return fn(n)
	}
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (n *Number) SetProperty(name string, value r.Element) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (n *Number) ExecMethod(name string, values []r.Element) (r.Element, error) {
	numMethodMap := map[string]numMethodFunc{
		"加":    numExecAdd,
		"减":    numExecSub,
		"乘":    numExecMul,
		"除":    numExecDiv,
		"自增":   numExecSelfAdd,
		"自减":   numExecSelfSub,
		"向下取整": numExecFloor,
		"向上取整": numExecCeil,
	}
	if fn, ok := numMethodMap[name]; ok {
		return fn(n, values)
	}
	return nil, zerr.MethodNotFound(name)
}

//// getters, setters and methods

// getters
func numGetText(n *Number) (r.Element, error) {
	return NewString(n.String()), nil
}

func numGetSquare(n *Number) (r.Element, error) {
	res := n.value * n.value
	return NewNumber(res), nil
}

func numGetCube(n *Number) (r.Element, error) {
	res := n.value * n.value * n.value
	return NewNumber(res), nil
}

func numGetSquareRoot(n *Number) (r.Element, error) {
	if n.value <= 0 {
		return nil, zerr.ArithRootLessThanZero()
	}
	res := math.Sqrt(n.value)
	return NewNumber(res), nil
}

// methods
func numExecAdd(n *Number, values []r.Element) (r.Element, error) {
	if err := ValidateAllParams(values, "number"); err != nil {
		return nil, err
	}

	sum := n.value
	for _, v := range values {
		vr, _ := v.(*Number)
		sum += vr.value
	}

	return NewNumber(sum), nil
}

func numExecSub(n *Number, values []r.Element) (r.Element, error) {
	if err := ValidateAllParams(values, "number"); err != nil {
		return nil, err
	}

	sum := n.value
	for _, v := range values {
		vr, _ := v.(*Number)
		sum -= vr.value
	}

	return NewNumber(sum), nil
}

func numExecMul(n *Number, values []r.Element) (r.Element, error) {
	if err := ValidateAllParams(values, "number"); err != nil {
		return nil, err
	}

	sum := n.value
	for _, v := range values {
		vr, _ := v.(*Number)
		sum *= vr.value
	}

	return NewNumber(sum), nil
}

func numExecDiv(n *Number, values []r.Element) (r.Element, error) {
	if err := ValidateAllParams(values, "number"); err != nil {
		return nil, err
	}

	sum := n.value
	for _, v := range values {
		vr, _ := v.(*Number)
		if vr.value == 0 {
			return nil, zerr.ArithDivZero()
		}
		sum /= vr.value
	}

	return NewNumber(sum), nil
}

func numExecSelfAdd(n *Number, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "number"); err != nil {
		return nil, err
	}

	target, _ := values[0].(*Number)

	n.value = n.value + target.value

	return n, nil
}

func numExecSelfSub(n *Number, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "number"); err != nil {
		return nil, err
	}

	target, _ := values[0].(*Number)

	n.value = n.value - target.value

	return n, nil
}

func numExecFloor(n *Number, values []r.Element) (r.Element, error) {
	return NewNumber(math.Floor(n.value)), nil
}

func numExecCeil(n *Number, values []r.Element) (r.Element, error) {
	return NewNumber(math.Ceil(n.value)), nil
}
