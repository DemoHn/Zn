package val

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

// Bool - represents for Zn's 二象型
type Bool struct {
	value bool
}

// NewBool - new bool ctx.Value Object from raw bool
func NewBool(value bool) *Bool {
	return &Bool{value}
}

// String - show displayed value
func (b *Bool) String() string {
	data := "真"
	if b.value == false {
		data = "假"
	}
	return data
}

// GetValue -
func (b *Bool) GetValue() bool {
	return b.value
}

// GetProperty -
func (b *Bool) GetProperty(c *ctx.Context, name string) (ctx.Value, *error.Error) {
	switch name {
	case "文本*":
		return NewString(b.String()), nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (b *Bool) SetProperty(c *ctx.Context, name string, value ctx.Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (b *Bool) ExecMethod(c *ctx.Context, name string, values []ctx.Value) (ctx.Value, *error.Error) {
	return nil, error.MethodNotFound(name)
}
