package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type Exception struct {
	Message string
}

func NewException(message string) *Exception {
	return &Exception{Message: message}
}

func (e *Exception) Error() string {
	return e.Message
}

// GetProperty -
func (e *Exception) GetProperty(c *r.Context, name string) (r.Element, error) {
	if name == "内容" {
		return NewString(e.Message), nil
	}
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (e *Exception) SetProperty(c *r.Context, name string, value r.Element) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (e *Exception) ExecMethod(c *r.Context, name string, values []r.Element) (r.Element, error) {
	return nil, zerr.MethodNotFound(name)
}
