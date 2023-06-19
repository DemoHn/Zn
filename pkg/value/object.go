package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

// Object - 对象型
type Object struct {
	propList map[string]r.Element
	model    *ClassModel
}

// ClassModel - aka. Class Definition Reference
// It defines the structure of a class, including compPropList, methodList and propList.
// All instances created from this class MUST inherits from those configurations.
type ClassModel struct {
	// Name - class name
	Name string
	// Constructor defines default logic (mostly for initialization) when a new instance
	// is created by "x 成为 C：P，Q，R"
	Constructor FuncExecutor
	// PropList defines all property name of a class, each item COULD NOT BE neither append nor removed
	PropList []string
	// CompPropList - CompProp stands for "Computed Property", which means the value is get or set
	// from a pre-defined function. Computed property offers more extensions for manipulations
	// of properties.
	CompPropList map[string]*Function
	// MethodList - stores all available methods definition of class
	MethodList map[string]*Function
}

// NewObject -
func NewObject(model *ClassModel) *Object {
	return &Object{
		propList: map[string]r.Element{},
		model:    model,
	}
}

// NewClassModel - create new empty r.ClassRef
func NewClassModel(name string) *ClassModel {
	return &ClassModel{
		Name:         name,
		Constructor:  nil,
		PropList:     []string{},
		CompPropList: map[string]*Function{},
		MethodList:   map[string]*Function{},
	}
}

// GetPropList -
func (zo *Object) GetPropList() map[string]r.Element {
	return zo.propList
}

// SetPropList -
func (zo *Object) SetPropList(propList map[string]r.Element) {
	zo.propList = propList
}

// GetModel -
func (zo *Object) GetModel() *ClassModel {
	return zo.model
}

// GetProperty -
func (zo *Object) GetProperty(c *r.Context, name string) (r.Element, error) {
	// internal properties
	switch name {
	case "自身":
		return zo, nil
	}

	if prop, ok := zo.propList[name]; ok {
		return prop, nil
	}
	// look up computed properties
	cprop, ok2 := zo.model.CompPropList[name]
	if !ok2 {
		return nil, zerr.PropertyNotFound(name)
	}
	// execute computed props to get property result
	return cprop.Exec(c, zo, []r.Element{})
}

// SetProperty -
func (zo *Object) SetProperty(c *r.Context, name string, value r.Element) error {
	if _, ok := zo.propList[name]; ok {
		zo.propList[name] = value
		return nil
	}
	// execute computed properties
	if cprop, ok2 := zo.model.CompPropList[name]; ok2 {
		_, err := cprop.Exec(c, zo, []r.Element{})
		return err
	}
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (zo *Object) ExecMethod(c *r.Context, name string, values []r.Element) (r.Element, error) {
	if method, ok := zo.model.MethodList[name]; ok {
		return method.Exec(c, zo, values)
	}
	return nil, zerr.MethodNotFound(name)
}

//// NOTE: ClassRef is also a type of Value
// Construct - yield new instance of this class
func (cr *ClassModel) Construct(c *r.Context, params []r.Element) (r.Element, error) {
	c.PushScope()
	defer c.PopScope()

	if cr.Constructor != nil {
		return cr.Constructor(c, params)
	}

	return nil, nil
}

// GetProperty - currently there's NO any property inside classRef Value
func (cr *ClassModel) GetProperty(c *r.Context, name string) (r.Element, error) {
	return nil, zerr.PropertyNotFound(name)
}

func (cr *ClassModel) SetProperty(c *r.Context, name string, value r.Element) error {
	return zerr.PropertyNotFound(name)
}

func (cr *ClassModel) ExecMethod(c *r.Context, name string, values []r.Element) (r.Element, error) {
	return nil, zerr.MethodNotFound(name)
}
