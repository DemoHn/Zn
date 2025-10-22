package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

// GoValue wraps golang internal value (like a struct, slice, map etc.) into a "zinc Element"
type GoValue struct {
	tag   string
	value interface{}
}

func NewGoValue(tag string, value interface{}) *GoValue {
	return &GoValue{tag, value}
}

func (gv *GoValue) GetTag() string {
	return gv.tag
}

func (gv *GoValue) GetValue() interface{} {
	return gv.value
}

// GetProperty -
func (gv *GoValue) GetProperty(name string) (r.Element, error) {
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (gv *GoValue) SetProperty(name string, value r.Element) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (gv *GoValue) ExecMethod(name string, values []r.Element) (r.Element, error) {
	return nil, zerr.MethodNotFound(name)
}
