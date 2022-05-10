package value

import (
	"fmt"
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"strconv"
	"strings"
)

type numGetterFunc func(*Number, *r.Context) (r.Value, error)
type numMethodFunc func(*Number, *r.Context, []r.Value) (r.Value, error)

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
	return fmt.Sprintf("%f", n.value)
}

// GetValue -
func (n *Number) GetValue() float64 {
	return n.value
}

// GetProperty -
func (n *Number) GetProperty(c *r.Context, name string) (r.Value, error) {
	numGetterMap := map[string]numGetterFunc{
		"文本": numGetText,
		"平方": numGetSquare,
		"立方": numGetCube,
	}
	if fn, ok := numGetterMap[name]; ok {
		return fn(n, c)
	}
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (n *Number) SetProperty(c *r.Context, name string, value r.Value) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (n *Number) ExecMethod(c *r.Context, name string, values []r.Value) (r.Value, error) {
	numMethodMap := map[string]numMethodFunc{
		"加": numExecAdd,
		"减": numExecSub,
		"乘": numExecMul,
		"除": numExecDiv,
		"自增": numExecSelfAdd,
		"自减": numExecSelfSub,
	}
	if fn, ok := numMethodMap[name]; ok {
		return fn(n, c, values)
	}
	return nil, zerr.MethodNotFound(name)
}

//// getters, setters and methods

// getters
func numGetText(n *Number, c *r.Context) (r.Value, error) {
	return NewString(n.String()), nil
}

func numGetSquare(n *Number, c *r.Context) (r.Value, error) {
	res := n.value * n.value
	return NewNumber(res), nil
}

func numGetCube(n *Number, c *r.Context) (r.Value, error) {
	res := n.value * n.value * n.value
	return NewNumber(res), nil
}

// methods
func numExecAdd(n *Number, c *r.Context, values []r.Value) (r.Value, error) {
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

func numExecSub(n *Number, c *r.Context, values []r.Value) (r.Value, error) {
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

func numExecMul(n *Number, c *r.Context, values []r.Value) (r.Value, error) {
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

func numExecDiv(n *Number, c *r.Context, values []r.Value) (r.Value, error) {
	if err := ValidateAllParams(values, "number"); err != nil {
		return nil, err
	}

	sum := n.value
	for _, v := range values {
		vr, _ := v.(*Number)
		if vr.value == 0 {
			return nil, zerr.NewErrorSLOT("被除数不能为0")
		}
		sum /= vr.value
	}

	return NewNumber(sum), nil
}

func numExecSelfAdd(n *Number, c *r.Context, values []r.Value) (r.Value, error) {
	if err := ValidateExactParams(values, "number"); err != nil {
		return nil, err
	}

	target, _ := values[0].(*Number)

	n.value = n.value + target.value

	return n, nil
}

func numExecSelfSub(n *Number, c *r.Context, values []r.Value) (r.Value, error) {
	if err := ValidateExactParams(values, "number"); err != nil {
		return nil, err
	}

	target, _ := values[0].(*Number)

	n.value = n.value - target.value

	return n, nil
}