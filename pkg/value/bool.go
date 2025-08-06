package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type boolGetterFunc func(*Bool) (r.Element, error)

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
	if !b.value {
		data = "假"
	}
	return data
}

// GetValue -
func (b *Bool) GetValue() bool {
	return b.value
}

// GetProperty -
func (b *Bool) GetProperty(c *r.Context, name string) (r.Element, error) {
	boolGetterMap := map[string]boolGetterFunc{
		"文本": boolGetText,
	}
	if fn, ok := boolGetterMap[name]; ok {
		return fn(b)
	}
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (b *Bool) SetProperty(c *r.Context, name string, value r.Element) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (b *Bool) ExecMethod(c *r.Context, name string, values []r.Element) (r.Element, error) {
	return nil, zerr.MethodNotFound(name)
}

// getters & setters & methods
// getters
// get text representation of current Bool value
func boolGetText(b *Bool) (r.Element, error) {
	return NewString(b.String()), nil
}
