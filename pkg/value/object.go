package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

// Object - 对象型
type Object struct {
	propList map[string]r.Value
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
	// MethodList - stores all available methods definition of class
	MethodList map[string]ClosureRef
}

// NewObject -
func NewObject(ref ClassRef) *Object {
	return &Object{
		propList: map[string]r.Value{},
		ref:      ref,
	}
}

// NewClassRef - create new empty r.ClassRef
func NewClassRef(name string) ClassRef {
	return ClassRef{
		Name:         name,
		Constructor:  nil,
		PropList:     []string{},
		CompPropList: map[string]ClosureRef{},
		MethodList:   map[string]ClosureRef{},
	}
}

// GetPropList -
func (zo *Object) GetPropList() map[string]r.Value {
	return zo.propList
}

// SetPropList -
func (zo *Object) SetPropList(propList map[string]r.Value) {
	zo.propList = propList
}

// GetRef -
func (zo *Object) GetRef() ClassRef {
	return zo.ref
}

// GetProperty -
func (zo *Object) GetProperty(c *r.Context, name string) (r.Value, error) {
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
		return nil, zerr.PropertyNotFound(name)
	}
	// execute computed props to get property result
	return cprop.Exec(c, zo, []r.Value{})
}

// SetProperty -
func (zo *Object) SetProperty(c *r.Context, name string, value r.Value) error {
	if _, ok := zo.propList[name]; ok {
		zo.propList[name] = value
		return nil
	}
	// execute computed properties
	if cprop, ok2 := zo.ref.CompPropList[name]; ok2 {
		_, err := cprop.Exec(c, zo, []r.Value{})
		return err
	}
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (zo *Object) ExecMethod(c *r.Context, name string, values []r.Value) (r.Value, error) {
	if method, ok := zo.ref.MethodList[name]; ok {
		return method.Exec(c, zo, values)
	}
	return nil, zerr.MethodNotFound(name)
}

//// NOTE: ClassRef is also a type of Value
// Construct - yield new instance of this class
func (cr *ClassRef) Construct(c *r.Context, params []r.Value) (r.Value, error) {
	c.PushChildScope()

	if cr.Constructor != nil {
		val, err := cr.Constructor(c, params)
		if err != nil {
			return nil, err
		}
		c.PopScope()
		return val, nil
	}

	c.PopScope()
	return nil, nil
}

// GetProperty - currently there's NO any property inside classRef Value
func (cr *ClassRef) GetProperty(c *r.Context, name string) (r.Value, error) {
	return nil, zerr.PropertyNotFound(name)
}

func (cr *ClassRef) SetProperty(c *r.Context, name string, value r.Value) error {
	return zerr.PropertyNotFound(name)
}

func (cr *ClassRef) ExecMethod(c *r.Context, name string, values []r.Value) (r.Value, error) {
	return nil, zerr.MethodNotFound(name)
}
