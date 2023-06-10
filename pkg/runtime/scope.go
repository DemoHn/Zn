package runtime

import "github.com/DemoHn/Zn/pkg/syntax"

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

func (sp *Scope) FindParentScope() *Scope {
	return sp.parent
}

func (sp *Scope) GetModule() *ModuleOLD {
	return sp.module
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

func (sp *Scope) SetExecLineIdx(line int) {
	sp.execCursor.CurrentLine = line
}

func (sp *Scope) SetLexer(l *syntax.Lexer) {
	sp.execCursor.Lexer = l
}

func (sp *Scope) GetExecCursor() ExecCursor {
	return sp.execCursor
}
