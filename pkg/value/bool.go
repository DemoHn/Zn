package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type boolGetterFunc func(*Bool, *r.Context) (r.Value, error)

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
func (b *Bool) GetProperty(c *r.Context, name string) (r.Value, error) {
	boolGetterMap := map[string]boolGetterFunc{
		"文本": boolGetText,
	}
	if fn, ok := boolGetterMap[name]; ok {
		return fn(b, c)
	}
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (b *Bool) SetProperty(c *r.Context, name string, value r.Value) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (b *Bool) ExecMethod(c *r.Context, name string, values []r.Value) (r.Value, error) {
	return nil, zerr.MethodNotFound(name)
}

//// getters & setters & methods
// getters
// get text representation of current Bool value
func boolGetText(b *Bool, c *r.Context) (r.Value, error) {
	return NewString(b.String()), nil
}