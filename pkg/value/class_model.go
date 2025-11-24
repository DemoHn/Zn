package value

import (
	"fmt"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

// ClassModel - aka. Class Definition Reference
// It defines the structure of a class, including compPropList, methodList and propList.
// All instances created from this class MUST inherits from those configurations.
type ClassModel struct {
	// Name - class name
	name string

	// Constructor defines default logic (mostly for initialization) when a new instance
	// is created by "x = (新建C：P，Q，R)"
	constructor r.FuncExecutor

	// PropList defines all property name & default value of the class, each property CANNOT be appended or removed
	propList map[string]r.Element

	// CompPropList - CompProp stands for "Computed Property", which means the value is get or set
	// from a pre-defined function. Computed property offers more extensions for manipulations
	// of properties.
	compPropList map[string]*Function

	// methodList - stores all available methods definition of class
	methodList map[string]*Function
}

// NewClassModel - create new empty r.ClassRef
func NewClassModel(name string) *ClassModel {
	model := &ClassModel{
		name:         name,
		constructor:  nil,
		propList:     map[string]r.Element{},
		compPropList: map[string]*Function{},
		methodList:   map[string]*Function{},
	}

	// set default constructor - NO ACTION AT ALL
	model.constructor = func(instance r.Element, params []r.Element) (r.Element, error) {
		return instance, nil
	}
	return model
}

func (cm *ClassModel) String() string {
	return fmt.Sprintf("‹类型·%s›", cm.name)
}

// Construct - yield new instance of this class
func (cm *ClassModel) Construct(params []r.Element) (r.Element, error) {
	// initialize a new object - an instance of class with no props set
	instance := NewObject(cm, r.ElementMap{})
	return cm.constructor(instance, params)
}

// //// GETTERS //////
func (cm *ClassModel) GetName() string {
	return cm.name
}

// GetPropList - list all defined properties to help duplicate initial properties to new Object
func (cm *ClassModel) GetPropList() map[string]r.Element {
	return cm.propList
}

func (cm *ClassModel) FindCompProp(name string) (*Function, bool) {
	cprop, ok := cm.compPropList[name]
	return cprop, ok
}

func (cm *ClassModel) FindMethod(name string) (*Function, bool) {
	method, ok := cm.methodList[name]
	return method, ok
}

// //// SETTERS //////
func (cm *ClassModel) SetConstructor(fn r.FuncExecutor) *ClassModel {
	cm.constructor = fn

	return cm
}

// DefineProperty - define property of model and set the defaultValue
func (cm *ClassModel) DefineProperty(name string, defaultValue r.Element) *ClassModel {
	cm.propList[name] = defaultValue

	return cm
}

func (cm *ClassModel) DefineCompProperty(name string, compFunc *Function) *ClassModel {
	cm.compPropList[name] = compFunc

	return cm
}

func (cm *ClassModel) DefineMethod(name string, methodFunc *Function) *ClassModel {
	cm.methodList[name] = methodFunc

	return cm
}

// // impl methods as a "Element"
// GetProperty - currently there's NO any property inside classRef Value
func (cr *ClassModel) GetProperty(name string) (r.Element, error) {
	return nil, zerr.PropertyNotFound(name)
}

func (cr *ClassModel) SetProperty(name string, value r.Element) error {
	return zerr.PropertyNotFound(name)
}

func (cr *ClassModel) ExecMethod(name string, values []r.Element) (r.Element, error) {
	return nil, zerr.MethodNotFound(name)
}

func (cr *ClassModel) Exportable() bool {
	return true
}
