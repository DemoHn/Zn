package value

import (
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type Bool struct {
	value bool
	*r.ElementModel
}

// NewBool - new bool ctx.Value Object from raw bool
func NewBool(value bool) *Bool {
	elem := &Bool{value, r.NewElementModel()}

	// register values
	elem.RegisterGetter("文本", elem.boolGetText)
	return elem
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

//// getters & setters & methods
// getters
// get text representation of current Bool value
func (b *Bool) boolGetText(c *r.Context) (r.Element, error) {
	return NewString(b.String()), nil
}
