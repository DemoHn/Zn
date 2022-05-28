package runtime

import "github.com/DemoHn/Zn/pkg/syntax"

// Scope represents a local namespace (aka. environment) of current execution
// including local variables map and current "return" value.
// NOTE: Scope is a doubly linked list
type Scope struct {
	parent *Scope
	// symbolMap - stores current scope stored symbols
	symbolMap map[string]SymbolInfo
	// module - get current module
	module *Module
	// execCursor - record current execution line
	execCursor *ExecCursor
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

type ExecCursor struct {
	*syntax.Lexer
	currentLine int
}


func (s SymbolInfo) GetValue() Value {
	return s.value
}

func NewScope(module *Module, lexer *syntax.Lexer) *Scope {
	return &Scope{
		parent: 	nil,
		module:      module,
		symbolMap:   map[string]SymbolInfo{},
		thisValue:   nil,
		returnValue: nil,
		execCursor: &ExecCursor{
			Lexer:       lexer,
			currentLine: 0,
		},
	}
}

func NewChildScope(sp *Scope) *Scope {
	return &Scope{
		parent: 	 sp,
		module:      sp.module,
		symbolMap:   map[string]SymbolInfo{},
		thisValue:   nil,
		returnValue: nil,
		execCursor: &ExecCursor{
			Lexer:       sp.execCursor.Lexer,
			currentLine: sp.execCursor.currentLine,
		},
	}
}

func (sp *Scope) FindParentScope() *Scope {
	return sp.parent
}

func (sp *Scope) GetModule() *Module {
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
	sp.execCursor.currentLine = line
}