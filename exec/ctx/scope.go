package ctx

// Scope represents a local namespace (aka. environment) of current execution
// including local variables map and current "return" value.
// NOTE: Scope is a doubly linked list
type Scope struct {
	// parent - parent Scope
	parent *Scope
	// child - child Scope
	child *Scope
	// symbolMap - stores current scope stored symbols
	symbolMap map[string]SymbolInfo
	// sgValue - scoped global variable
	sgValue Value
	// thisValue - "this" variable of the scope
	thisValue Value
	// retrunValue - return value of scope
	returnValue Value
}

// SymbolInfo - a wrapper of symbol's value with additional properties.
type SymbolInfo struct {
	// value -
	value Value
	// isConst - if an symbol is const
	isConst bool
}

func NewScope() *Scope {
	return &Scope{
		parent:      nil,
		child:       nil,
		symbolMap:   map[string]SymbolInfo{},
		sgValue:     nil,
		thisValue:   nil,
		returnValue: nil,
	}
}

// CreateChildScope -
func (sp *Scope) CreateChildScope() *Scope {
	newScope := &Scope{
		parent:      sp,
		child:       nil,
		symbolMap:   map[string]SymbolInfo{},
		sgValue:     nil,
		thisValue:   nil,
		returnValue: nil,
	}

	sp.child = newScope
	return newScope
}

// FindParentScope - find parent scope
func (sp *Scope) FindParentScope() *Scope {
	return sp.parent
}

// FindChildScope - find child scope
func (sp *Scope) FindChildScope() *Scope {
	return sp.child
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

// SetSgValue -
func (sp *Scope) SetSgValue(v Value) {
	sp.sgValue = v
}

// GetSgValue -
func (sp *Scope) GetSgValue() Value {
	return sp.sgValue
}
