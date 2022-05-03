package value

import (
	r "github.com/DemoHn/Zn/pkg/runtime"
)

// Null -
type Null struct{
	r.ValueBase
}

// NewNull -
func NewNull() *Null {
	return &Null{}
}
