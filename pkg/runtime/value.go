package runtime

import zerr "github.com/DemoHn/Zn/pkg/error"

// Value is the base unit to present a value (aka. variable) - including number, string, array, function, object...
// All kinds of values in Zn language SHOULD implement this interface.
//
// Basically there are 3 methods:
//
// 1. GetProperty - fetch the value from property list of a specific name
// 2. SetProperty - set the value of some property
// 3. ExecMethod - execute one method from method list
type Value interface {
	GetProperty(*Context, string) (Value, error)
	SetProperty(*Context, string, Value) error
	ExecMethod(*Context, string, []Value) (Value, error)
}

type getterFunc = func(*Context) (Value, error)
type setterFunc = func(*Context, Value) error
type methodFunc = func(*Context, []Value) (Value, error)

type ValueBase struct {
	getters []getterFunc
	setters []setterFunc
	methods []methodFunc
	// index map - to find getter/setter/method from its name
	getterIdxMap map[string]int
	setterIdxMap map[string]int
	methodIdxMap map[string]int
}

//// implement Value methods
// GetProperty -
func (vb ValueBase) GetProperty(ctx *Context, name string) (Value, error) {
	// find index
	if idx, ok := vb.getterIdxMap[name]; ok {
		if idx < len(vb.getters) {
			return vb.getters[idx](ctx)
		}
	}
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (vb ValueBase) SetProperty(ctx *Context, name string, value Value) error {
	// find index
	if idx, ok := vb.setterIdxMap[name]; ok {
		if idx < len(vb.setters) {
			return vb.setters[idx](ctx, value)
		}
	}
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (vb ValueBase) ExecMethod(ctx *Context, name string, values []Value) (Value, error) {
	// find index
	if idx, ok := vb.methodIdxMap[name]; ok {
		if idx < len(vb.methods) {
			return vb.methods[idx](ctx, values)
		}
	}
	return nil, zerr.MethodNotFound(name)
}
