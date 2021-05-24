package val

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

// Object - 对象型
type Object struct {
	propList map[string]ctx.Value
	ref      ctx.ClassRef
}

// NewObject -
func NewObject(ref ctx.ClassRef) *Object {
	return &Object{
		propList: map[string]ctx.Value{},
		ref:      ref,
	}
}

// GetProperty -
func (zo *Object) GetProperty(c *ctx.Context, name string) (ctx.Value, *error.Error) {
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
	return cprop.Exec(fctx, []ctx.Value{})
}

// SetProperty -
func (zo *Object) SetProperty(c *ctx.Context, name string, value ctx.Value) *error.Error {
	if _, ok := zo.propList[name]; ok {
		zo.propList[name] = value
		return nil
	}
	// execute computed properites
	if cprop, ok2 := zo.ref.CompPropList[name]; ok2 {
		_, err := cprop.Exec(ctx, []ctx.Value{})
		return err
	}
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (zo *Object) ExecMethod(c *ctx.Context, name string, values []ctx.Value) (ctx.Value, *error.Error) {
	if method, ok := zo.ref.MethodList[name]; ok {
		return method.Exec(ctx, values)
	}
	return nil, error.MethodNotFound(name)
}
