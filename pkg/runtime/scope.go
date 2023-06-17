package runtime

// Scope represents a local namespace (aka. environment) of current execution
// including local variables map and current "return" value.
// NOTE: Scope is a doubly linked list
type Scope struct {
	// symbolMap - stores current scope stored symbols
	symbolMap map[string]SymbolInfo
	// thisValue - "this" variable of the scope
	thisValue Element
	// returnValue - return value of scope
	returnValue Element
}

// SymbolInfo - a wrapper of symbol's value with additional properties.
type SymbolInfo struct {
	// value -
	value Element
	// isConst - if an symbol is const
	isConst bool
}

func (s SymbolInfo) GetValue() Element {
	return s.value
}

func MakeSymbolInfo(value Element, isConst bool) SymbolInfo {
	return SymbolInfo{value, isConst}
}

func NewScope() *Scope {
	return &Scope{
		symbolMap:   map[string]SymbolInfo{},
		thisValue:   nil,
		returnValue: nil,
	}
}

// GetThisValue -
func (sp *Scope) GetThisValue() Element {
	return sp.thisValue
}

// SetThisValue -
func (sp *Scope) SetThisValue(v Element) {
	sp.thisValue = v
}

// GetReturnValue -
func (sp *Scope) GetReturnValue() Element {
	return sp.returnValue
}

// SetReturnValue -
func (sp *Scope) SetReturnValue(v Element) {
	sp.returnValue = v
}

func (sp *Scope) SetSymbolValue(name string, v Element, isConst bool) {
	sp.symbolMap[name] = SymbolInfo{
		isConst: isConst,
		value:   v,
	}
}

func (sp *Scope) GetSymbolValue(name string) (bool, Element) {
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
