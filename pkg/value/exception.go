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

func (e *Exception) GetMessage() string {
	return e.GetMessage()
}

// GetProperty -
func (e *Exception) GetProperty(c *r.Context, name string) (r.Value, error) {
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (e *Exception) SetProperty(c *r.Context, name string, value r.Value) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (e *Exception) ExecMethod(c *r.Context, name string, values []r.Value) (r.Value, error) {
	return nil, zerr.MethodNotFound(name)
}