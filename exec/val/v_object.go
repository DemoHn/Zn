package val

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

// Object - 对象型
type Object struct {
	propList map[string]ctx.Value
	ref      ClassRef
}

// ClassRef - aka. Class Definition Reference
// It defines the structure of a class, including compPropList, methodList and propList.
// All instances created from this class MUST inherits from those configurations.
type ClassRef struct {
	// Name - class name
	Name string
	// Constructor defines default logic (mostly for initialization) when a new instance
	// is created by "x 成为 C：P，Q，R"
	Constructor funcExecutor
	// PropList defines all property name of a class, each item COULD NOT BE neither append nor removed
	PropList []string
	// CompPropList - CompProp stands for "Computed Property", which means the value is get or set
	// from a pre-defined function. Computed property offers more extensions for manipulations
	// of properties.
	CompPropList map[string]ClosureRef
	// MethodList - stores all available methods defintion of class
	MethodList map[string]ClosureRef
}

// NewObject -
func NewObject(ref ClassRef) *Object {
	return &Object{
		propList: map[string]ctx.Value{},
		ref:      ref,
	}
}

// NewClassRef - create new empty ctx.ClassRef
func NewClassRef(name string) ClassRef {
	return ClassRef{
		Name:         name,
		Constructor:  nil,
		PropList:     []string{},
		CompPropList: map[string]ClosureRef{},
		MethodList:   map[string]ClosureRef{},
	}
}

// Construct - yield new instance of this class
func (cr *ClassRef) Construct(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	return cr.Constructor(c, params)
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
	sp := c.ShiftChildScope()
	sp.SetThisValue(zo)
	return cprop.Exec(c, []ctx.Value{})
}

// SetProperty -
func (zo *Object) SetProperty(c *ctx.Context, name string, value ctx.Value) *error.Error {
	if _, ok := zo.propList[name]; ok {
		zo.propList[name] = value
		return nil
	}
	// execute computed properites
	if cprop, ok2 := zo.ref.CompPropList[name]; ok2 {
		_, err := cprop.Exec(c, []ctx.Value{})
		return err
	}
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (zo *Object) ExecMethod(c *ctx.Context, name string, values []ctx.Value) (ctx.Value, *error.Error) {
	if method, ok := zo.ref.MethodList[name]; ok {
		return method.Exec(c, values)
	}
	return nil, error.MethodNotFound(name)
}
