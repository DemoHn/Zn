package exec

import (
	"unicode/utf8"

	"github.com/DemoHn/Zn/error"
)

// String - represents for Zn's 文本型
type String struct {
	value string
}

// NewString - new string Value Object from raw string
func NewString(value string) *String {
	return &String{value}
}

// String - display string value's string
func (s *String) String() string {
	return s.value
}

// GetProperty -
func (s *String) GetProperty(ctx *Context, name string) (Value, *error.Error) {
	switch name {
	case "长度":
		l := utf8.RuneCountInString(s.value)
		return NewDecimalFromInt(l, 0), nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (s *String) SetProperty(ctx *Context, name string, value Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (s *String) ExecMethod(ctx *Context, name string, values []Value) (Value, *error.Error) {
	return nil, error.MethodNotFound(name)
}
