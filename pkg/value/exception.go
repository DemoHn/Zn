package value

import (
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type Exception struct {
	Message string
	*r.ElementModel
}

func NewException(message string) *Exception {
	return &Exception{
		message,
		r.NewElementModel(),
	}
}

func (e *Exception) GetMessage() string {
	return e.Message
}
