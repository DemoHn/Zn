package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
)

// Element is the base unit to present a value (aka. variable) - including number, string, array, function, object...
// All kinds of values in Zn language SHOULD implement this interface.
//
// Basically there are 3 methods:
//
// 1. GetProperty - fetch the value from property list of a specific name
// 2. SetProperty - set the value of some property
// 3. ExecMethod - execute one method from method list
type Element interface {
	GetProperty(*Context, string) (Element, error)
	SetProperty(*Context, string, Element) error
	ExecMethod(*Context, string, []Element) (Element, error)
}

/////
// ElementModel is a helper to simplify element initial process
/////

type getterClosure = func(*Context) (Element, error)
type setterClosure = func(*Context, Element) error
type methodClosure = func(*Context, []Element) (Element, error)

type ElementModel struct {
	getterFuncMap map[string]getterClosure
	setterFuncMap map[string]setterClosure
	methodFuncMap map[string]methodClosure
}

func NewElementModel() *ElementModel {
	return &ElementModel{
		getterFuncMap: make(map[string]getterClosure),
		setterFuncMap: make(map[string]setterClosure),
		methodFuncMap: make(map[string]methodClosure),
	}
}

//// register getters & setters & methods
func (em *ElementModel) RegisterGetter(name string, closure getterClosure) {
	em.getterFuncMap[name] = closure
}

func (em *ElementModel) RegisterSetter(name string, closure setterClosure) {
	em.setterFuncMap[name] = closure
}

func (em *ElementModel) RegisterMethod(name string, closure methodClosure) {
	em.methodFuncMap[name] = closure
}

//// impl Element interface
func (em *ElementModel) GetProperty(c *Context, name string) (Element, error) {
	if fn, ok := em.getterFuncMap[name]; ok {
		return fn(c)
	}
	return nil, zerr.PropertyNotFound(name)
}

func (em *ElementModel) SetProperty(c *Context, name string, el Element) error {
	if fn, ok := em.setterFuncMap[name]; ok {
		return fn(c, el)
	}
	return zerr.PropertyNotFound(name)
}

func (em *ElementModel) ExecMethod(c *Context, name string, elems []Element) (Element, error) {
	if fn, ok := em.methodFuncMap[name]; ok {
		return fn(c, elems)
	}
	return nil, zerr.MethodNotFound(name)
}
