package exec

import (
	"github.com/DemoHn/Zn/error"
)

// Object - 对象型
type Object struct {
	propList map[string]Value
	ref      ClassRef
}

// NewObject -
func NewObject(ref ClassRef) *Object {
	return &Object{
		propList: map[string]Value{},
		ref:      ref,
	}
}

// GetProperty -
func (zo *Object) GetProperty(ctx *Context, name string) (Value, *error.Error) {
	// internal properties
	switch name {
	case "自身":
		return zo, nil
	}

	if prop, ok := zo.propList[name]; ok {
		return prop, nil
	}
	// look up computed properties
	cprop, ok2 := zo.ref.CompPropList[name]
	if !ok2 {
		return nil, error.PropertyNotFound(name)
	}
	// execute computed props to get property result
	fctx := ctx.DuplicateNewScope()
	fctx.scope.thisValue = zo
	return cprop.Exec(fctx, []Value{})
}

// SetProperty -
func (zo *Object) SetProperty(ctx *Context, name string, value Value) *error.Error {
	if _, ok := zo.propList[name]; ok {
		zo.propList[name] = value
		return nil
	}
	// execute computed properites
	if cprop, ok2 := zo.ref.CompPropList[name]; ok2 {
		_, err := cprop.Exec(ctx, []Value{})
		return err
	}
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (zo *Object) ExecMethod(ctx *Context, name string, values []Value) (Value, *error.Error) {
	if method, ok := zo.ref.MethodList[name]; ok {
		return method.Exec(ctx, values)
	}
	return nil, error.MethodNotFound(name)
}
