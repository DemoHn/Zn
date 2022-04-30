package runtime

// Value is the base unit to present a value (aka. variable) - including number, string, array, function, object...
// All kinds of values in Zn language SHOULD implement this interface.
//
// Basically there are 3 methods:
//
// 1. GetProperty - fetch the value from property list of a specific name
// 2. SetProperty - set the value of some property
// 3. ExecMethod - execute one method from method list
type Value interface {
	GetProperty(*Context, string) (Value, error)
	SetProperty(*Context, string, Value) error
	ExecMethod(*Context, string, []Value) error
}
