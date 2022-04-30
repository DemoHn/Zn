package runtime

// Scope represents a local namespace (aka. environment) of current execution
// including local variables map and current "return" value.
// NOTE: Scope is a doubly linked list
type Scope struct {
	// parent - parent Scope
	parent *Scope
	// child - child Scope
	child *Scope
	// symbolMap - symbolMap
	symbolMap map[string]SymbolInfo
	// returnValue of current scopes
	returnValue Value
}

// SymbolInfo - a wrapper of symbol's value with additional properties.
type SymbolInfo struct {
	// value -
	value Value
	// isConst - if an symbol is const
	isConst bool
}

