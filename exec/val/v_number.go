package val

import (
	"fmt"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

type Number struct {
	value float64
}

// NewNumber - create new number object (plain float64)
func NewNumber(value float64) *Number {
	return &Number{value}
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
func (n *Number) GetProperty(c *ctx.Context, name string) (ctx.Value, *error.Error) {
	switch name {
	case "文本":
		return NewString(n.String()), nil
	case "+1":
		return NewNumber(n.value + 1), nil
	case "-1":
		return NewNumber(n.value - 1), nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (n *Number) SetProperty(c *ctx.Context, name string, value ctx.Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (n *Number) ExecMethod(c *ctx.Context, name string, values []ctx.Value) (ctx.Value, *error.Error) {
	switch name {
	case "+1":
		n.value = n.value + 1
		return n, nil
	case "-1":
		n.value = n.value - 1
		return n, nil
	}
	return nil, error.MethodNotFound(name)
}

//// arithmetic calculations
// Add -
func (n *Number) Add(others ...*Number) *Number {
	var result = n.value
	for _, item := range others {
		result += item.value
	}

	return NewNumber(result)
}

// Sub -
func (n *Number) Sub(others ...*Number) *Number {
	var result = n.value
	for _, item := range others {
		result -= item.value
	}

	return NewNumber(result)
}

// Mul - TODO
func (n *Number) Mul(others ...*Number) *Number {
	return NewNumber(0)
}

// Div - TODO
func (n *Number) Div(others ...*Number) *Number {
	return NewNumber(0)
}
