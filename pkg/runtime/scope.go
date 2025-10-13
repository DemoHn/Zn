package runtime

import zerr "github.com/DemoHn/Zn/pkg/error"

// Scope represents a local namespace (aka. environment) of current execution
// including local variables map and current "return" value.
// NOTE: Scope is a doubly linked list
/**
type Scope struct {
	// symbolMap - stores current scope stored symbols
	symbolMap map[string]SymbolInfo
	// thisValue - "this" variable of the scope
	thisValue Element
}

// SymbolInfo - a wrapper of symbol's value with additional properties.
type SymbolInfo struct {
	// value -
	value Element
	// isConst - if an symbol is const
	isConst bool
	// module - get original module (for reference)
	module *Module
}

func (s SymbolInfo) GetValue() Element {
	return s.value
}

func (s SymbolInfo) GetModule() *Module {
	return s.module
}

func NewScope(thisValue Element) *Scope {
	return &Scope{
		symbolMap: map[string]SymbolInfo{},
		thisValue: thisValue,
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

func (sp *Scope) SetSymbolValue(name string, v Element, isConst bool, module *Module) {
	sp.symbolMap[name] = SymbolInfo{
		isConst: isConst,
		value:   v,
		module:  module,
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

**/
////////////////////////////////////////////
////////////////////////////////////////////
type Scope struct {
	locals       []LocalSymbol
	localCount   int
	currentDepth int
	values       []Element
}

type LocalSymbol struct {
	name    string
	depth   int
	isConst bool
}

func NewScope() Scope {
	return Scope{
		locals:       []LocalSymbol{},
		localCount:   0,
		currentDepth: 0,
		values:       []Element{},
	}
}

func (sp *Scope) BeginScope() {
	sp.currentDepth++
}

func (sp *Scope) EndScope() {
	sp.currentDepth--

	// pop all deeper values
	for sp.localCount > 0 && sp.locals[sp.localCount-1].depth > sp.currentDepth {
		sp.localCount--
	}
}

func (sp *Scope) AddValue(name string, value Element) {
	sp.locals = append(sp.locals, LocalSymbol{
		name:    name,
		depth:   sp.currentDepth,
		isConst: false,
	})

	sp.values = append(sp.values, value)
}

func (sp *Scope) AddConstValue(name string, value Element) {
	sp.locals = append(sp.locals, LocalSymbol{
		name:    name,
		depth:   sp.currentDepth,
		isConst: true,
	})

	sp.values = append(sp.values, value)
}

func (sp *Scope) GetValue(name string) Element {
	for i := sp.localCount - 1; i >= 0; i-- {
		if sp.locals[i].name == name {
			return sp.values[i]
		}
	}
	return nil
}

func (sp *Scope) SetValue(name string, value Element) error {
	for i := sp.localCount - 1; i >= 0; i-- {
		if sp.locals[i].name == name {
			if sp.locals[i].isConst {
				// error: cannot change const value
				return zerr.AssignToConstant()
			}
			sp.values[i] = value
			return nil
		}
	}
	return zerr.NameNotDefined(name)
}
