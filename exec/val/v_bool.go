package val

import "github.com/DemoHn/Zn/error"

// Bool - represents for Zn's 二象型
type Bool struct {
	value bool
}

// NewBool - new bool Value Object from raw bool
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

// GetProperty -
func (b *Bool) GetProperty(ctx *Context, name string) (Value, *error.Error) {
	switch name {
	case "文本*":
		return NewString(b.String()), nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (b *Bool) SetProperty(ctx *Context, name string, value Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (b *Bool) ExecMethod(ctx *Context, name string, values []Value) (Value, *error.Error) {
	return nil, error.MethodNotFound(name)
}
