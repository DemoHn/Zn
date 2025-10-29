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

// NewObject -
func NewObject(model *ClassModel, initProps map[string]r.Element) *Object {
	objPropList := make(map[string]r.Element)
	for prop, elem := range model.GetPropList() {
		// find value from initial props first
		if initValue, ok := initProps[prop]; ok {
			objPropList[prop] = initValue
		} else {
			// duplicate default prop values
			objPropList[prop] = DuplicateValue(elem)
		}
	}

	return &Object{
		propList: objPropList,
		model:    model,
	}
}

// get object name
func (zo *Object) GetObjectName() string {
	return zo.model.GetName()
}

// GetProperty -
func (zo *Object) GetProperty(name string) (r.Element, error) {
	// internal properties
	switch name {
	case "自身":
		return zo, nil
	}

	if prop, ok := zo.propList[name]; ok {
		return prop, nil
	}
	// look up computed properties
	cprop, ok2 := zo.model.FindCompProp(name)
	if !ok2 {
		return nil, zerr.PropertyNotFound(name)
	}
	// execute computed props to get property result
	return cprop.Exec(zo, []r.Element{})
}

// SetProperty -
func (zo *Object) SetProperty(name string, value r.Element) error {
	if _, ok := zo.propList[name]; ok {
		zo.propList[name] = value
		return nil
	}
	// execute computed properties
	if cprop, ok2 := zo.model.FindCompProp(name); ok2 {
		_, err := cprop.Exec(zo, []r.Element{})
		return err
	}
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (zo *Object) ExecMethod(name string, values []r.Element) (r.Element, error) {
	if method, ok := zo.model.FindMethod(name); ok {
		return method.Exec(zo, values)
	}
	return nil, zerr.MethodNotFound(name)
}
