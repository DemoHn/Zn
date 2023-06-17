package value

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type Number struct {
	value float64
	*r.ElementModel
}

// NewNumber - create new number object (plain float64)
func NewNumber(value float64) *Number {
	n := &Number{value, r.NewElementModel()}
	// init getters & setters & methods
	n.RegisterGetter("文本", n.numGetText)
	n.RegisterGetter("平方", n.numGetSquare)
	n.RegisterGetter("立方", n.numGetCube)
	n.RegisterGetter("平方根", n.numGetSquareRoot)

	n.RegisterMethod("加", n.numExecAdd)
	n.RegisterMethod("减", n.numExecSub)
	n.RegisterMethod("乘", n.numExecMul)
	n.RegisterMethod("除", n.numExecDiv)
	n.RegisterMethod("自增", n.numExecSelfAdd)
	n.RegisterMethod("自减", n.numExecSelfSub)
	n.RegisterMethod("向下取整", n.numExecFloor)
	n.RegisterMethod("向上取整", n.numExecCeil)

	return n
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

// GetValue -
func (n *Number) GetValue() float64 {
	return n.value
}

//// getters, setters and methods

// getters
func (n *Number) numGetText(c *r.Context) (r.Element, error) {
	return NewString(n.String()), nil
}

func (n *Number) numGetSquare(c *r.Context) (r.Element, error) {
	res := n.value * n.value
	return NewNumber(res), nil
}

func (n *Number) numGetCube(c *r.Context) (r.Element, error) {
	res := n.value * n.value * n.value
	return NewNumber(res), nil
}

func (n *Number) numGetSquareRoot(c *r.Context) (r.Element, error) {
	if n.value <= 0 {
		return nil, zerr.ArithRootLessThanZero()
	}
	res := math.Sqrt(n.value)
	return NewNumber(res), nil
}

// methods
func (n *Number) numExecAdd(c *r.Context, values []r.Element) (r.Element, error) {
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

func (n *Number) numExecSub(c *r.Context, values []r.Element) (r.Element, error) {
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

func (n *Number) numExecMul(c *r.Context, values []r.Element) (r.Element, error) {
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

func (n *Number) numExecDiv(c *r.Context, values []r.Element) (r.Element, error) {
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

func (n *Number) numExecSelfAdd(c *r.Context, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "number"); err != nil {
		return nil, err
	}

	target, _ := values[0].(*Number)

	n.value = n.value + target.value

	return n, nil
}

func (n *Number) numExecSelfSub(c *r.Context, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "number"); err != nil {
		return nil, err
	}

	target, _ := values[0].(*Number)

	n.value = n.value - target.value

	return n, nil
}

func (n *Number) numExecFloor(c *r.Context, values []r.Element) (r.Element, error) {
	return NewNumber(math.Floor(n.value)), nil
}

func (n *Number) numExecCeil(c *r.Context, values []r.Element) (r.Element, error) {
	return NewNumber(math.Ceil(n.value)), nil
}
