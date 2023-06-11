package runtime

// Scope represents a local namespace (aka. environment) of current execution
// including local variables map and current "return" value.
// NOTE: Scope is a doubly linked list
type Scope struct {
	// symbolMap - stores current scope stored symbols
	symbolMap map[string]SymbolInfo
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

func (s SymbolInfo) GetValue() Value {
	return s.value
}

func NewScope() *Scope {
	return &Scope{
		symbolMap:   map[string]SymbolInfo{},
		thisValue:   nil,
		returnValue: nil,
	}
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

func (sp *Scope) SetSymbolValue(name string, isConst bool, v Value) {
	sp.symbolMap[name] = SymbolInfo{
		isConst: isConst,
		value:   v,
	}
}

func (sp *Scope) GetSymbolValue(name string) (bool, Value) {
	if info, ok := sp.symbolMap[name]; ok {
		return true, info.value
	}
	return false, nil
}

func (sp *Scope) GetSymbol(name string) (bool, SymbolInfo) {
	if info, ok := sp.symbolMap[name]; ok {
		return true, info
	}
	return false, SymbolInfo{}
}
