package exec

import (
	"github.com/DemoHn/Zn/error"
)

//// General Value types

// ZnValue - general value interface
type ZnValue interface {
	String() string
	GetProperty(string) (ZnValue, *error.Error)
	SetProperty(string, ZnValue) *error.Error
	GetMethod(string) (*ClosureRef, *error.Error)
	FindGetter(string) (bool, *ClosureRef)
}

//////// ZnObject Definition

// ZnObject -
type ZnObject struct {
	// defines all properties (readable and writable)
	PropList map[string]ZnValue
	*ClassRef
}

//////// Primitive Types Definition

// ZnNull - Zn null type - a special marker indicates that
// this value has neither type nor value
type ZnNull struct {
	*ZnObject
}

//////// Variable Type Implementation

func (zo *ZnObject) String() string {
	return "[Object]"
}

// GetProperty -
func (zo *ZnObject) GetProperty(name string) (ZnValue, *error.Error) {
	prop, ok := zo.PropList[name]
	if !ok {
		return nil, error.PropertyNotFound(name)
	}
	return prop, nil
}

// SetProperty -
func (zo *ZnObject) SetProperty(name string, value ZnValue) *error.Error {
	_, ok := zo.PropList[name]
	if !ok {
		return error.PropertyNotFound(name)
	}
	zo.PropList[name] = value
	return nil
}

// GetMethod -
func (zo *ZnObject) GetMethod(name string) (*ClosureRef, *error.Error) {
	methodRef, ok := zo.MethodList[name]
	if !ok {
		return nil, error.MethodNotFound(name)
	}
	return methodRef, nil
}

// FindGetter -
func (zo *ZnObject) FindGetter(name string) (bool, *ClosureRef) {
	getterRef, ok := zo.GetterList[name]
	if !ok {
		return false, nil
	}
	return true, getterRef
}

func (zn *ZnNull) String() string {
	return "空"
}

//////// New[Type] Constructors

// NewZnNull - null value
func NewZnNull() *ZnNull {
	t := &ZnNull{}
	return t
}

// NewZnObject -
func NewZnObject(classRef *ClassRef) *ZnObject {
	return &ZnObject{
		PropList: map[string]ZnValue{},
		ClassRef: classRef,
	}
}
