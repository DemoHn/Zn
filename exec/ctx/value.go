package ctx

import (
	"github.com/DemoHn/Zn/error"
)

// Value is the base unit to present a value (aka. variable) - including number, string, array, function, object...
// All kinds of values in Zn language SHOULD implement this interface.
//
// Basically there're 3 methods:
//
// 1. GetProperty - fetch the value from property list of a specific name
// 2. SetProperty - set the value of some property
// 3. ExecMethod - execute one method from method list
type Value interface {
	GetProperty(*Context, string) (Value, *error.Error)
	SetProperty(*Context, string, Value) *error.Error
	ExecMethod(*Context, string, []Value) (Value, *error.Error)
}

// ClosureRef - aka. Closure Exection Reference
// It's the structure of a closure which wraps execution logic.
// The executor could be either a bunch of code or some native code.
type ClosureRef struct {
	ParamHandler funcExecutor
	Executor     funcExecutor // closure execution logic
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
