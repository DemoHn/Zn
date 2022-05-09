package value

import (
	"fmt"
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"strconv"
	"strings"
)

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
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (n *Number) SetProperty(c *r.Context, name string, value r.Value) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (n *Number) ExecMethod(c *r.Context, name string, values []r.Value) (r.Value, error) {
	return nil, zerr.MethodNotFound(name)
}