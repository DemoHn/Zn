package runtime

// Scope represents a local namespace (aka. environment) of current execution
// including local variables map and current "return" value.
// NOTE: Scope is a doubly linked list
type Scope struct {
	parent *Scope
	// symbolMap - stores current scope stored symbols
	symbolMap map[string]SymbolInfo
	// module - get current module
	module *Module
	// thisValue - "this" variable of the scope
	thisValue Value
	// returnValue - return value of scope
	returnValue Value
}

// SymbolInfo - a wrapper of symbol's value with additional properties.
type SymbolInfo struct {
	// value -
	value Value
	// isConst - if an symbol is const
	isConst bool
}

func NewScope(module *Module) *Scope {
	return &Scope{
		parent: 	nil,
		module:      module,
		symbolMap:   map[string]SymbolInfo{},
		thisValue:   nil,
		returnValue: nil,
	}
}

func NewChildScope(sp *Scope) *Scope {
	return &Scope{
		parent: 	 sp,
		module:      sp.module,
		symbolMap:   map[string]SymbolInfo{},
		thisValue:   nil,
		returnValue: nil,
	}
}

func (sp *Scope) FindParentScope() *Scope {
	return sp.parent
}

// GetThisValue -
func (sp *Scope) GetThisValue() Value {
	return sp.thisValue
}

// SetThisValue -
func (sp *Scope) SetThisValue(v Value) {
	sp.thisValue = v
}

// GetReturnValue -
func (sp *Scope) GetReturnValue() Value {
	return sp.returnValue
}

// SetReturnValue -
func (sp *Scope) SetReturnValue(v Value) {
	sp.returnValue = v
}
