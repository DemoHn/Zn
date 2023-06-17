package value

import (
	r "github.com/DemoHn/Zn/pkg/runtime"
)

// Null -
type Null struct {
	*r.ElementModel
}

// NewNull -
func NewNull() *Null {
	return &Null{r.NewElementModel()}
}
