package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

// Null -
type Null struct{}

// NewNull -
func NewNull() *Null {
	return &Null{}
}

func (s *Null) String() string {
	return "空"
}

// GetProperty -
func (s *Null) GetProperty(name string) (r.Element, error) {
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (s *Null) SetProperty(name string, value r.Element) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (s *Null) ExecMethod(name string, values []r.Element) (r.Element, error) {
	return nil, zerr.MethodNotFound(name)
}
