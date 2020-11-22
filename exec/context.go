package exec

import (
	"github.com/DemoHn/Zn/debug"
	"github.com/DemoHn/Zn/lex"
)

const arithPrecision = 8

// Context is a global variable that stores current execution
// states, global configurations
type Context struct {
	// globals - stores all global variables
	globals map[string]Value
	// arith - an arith object that manages all arith calculations
	arith *Arith

	// a seperate map to store inner debug data
	// usage: call （__probe：「tagName」，variable）
	// it will record all logs (including variable value, curernt scope, etc.)
	// the value is deep-copied so don't worry - the value logged won't be changed
	_probe *debug.Probe
	scope  *Scope
}

// Scope represents a local namespace (aka. environment) of current execution
// including local variables map and current "return" value.
type Scope struct {
	// fileInfo -
	fileInfo *FileInfo
	// classRefMap stores class definition template (reference) within the scope
	// this item only exists on RootScope since class defition block IS allowed
	// ONLY in root block
	classRefMap map[string]ClassRef
	// parent - parent Scope
	parent    *Scope
	symbolMap map[string]SymbolInfo
	// sgValue - scope variable
	sgValue Value
	// thisValue - "this" variable of the scope
	thisValue Value
	// retrunValue - return value of scope
	returnValue Value
}

// FileInfo records current file info, usually for displaying error
type FileInfo struct {
	//// lexical scope
	// file - current execution file directory
	file string
	// currentLine - current exeuction line
	currentLine int
	// lineStack - lexical info of (parsed) current file
	lineStack *lex.LineStack
}

// SymbolInfo - a wrapper of symbol's value with additional properties.
type SymbolInfo struct {
	value Value
	// Constant - if an symbol is const
	isConst bool
}

// NewContext - create new Zn Context. Notice through the life-cycle
// of one code execution, there's only one running context to store all states.
func NewContext() *Context {
	return &Context{
		globals: map[string]Value{}, // TODO
		arith:   NewArith(arithPrecision),
		_probe:  debug.NewProbe(),
		scope:   nil,
	}
}

// InitScope - init root scope
func (ctx *Context) InitScope(l *lex.Lexer) {
	fileInfo := &FileInfo{
		file:        l.InputStream.GetFile(),
		currentLine: 0,
		lineStack:   l.LineStack,
	}
	newScope := &Scope{
		fileInfo:    fileInfo,
		classRefMap: map[string]ClassRef{},
		parent:      nil,
		symbolMap:   map[string]SymbolInfo{},
		sgValue:     nil,
		thisValue:   nil,
		returnValue: nil,
	}
	ctx.scope = newScope
}

// DuplicateNewScope - create a new Context with new scope which parent points to duplicator's scope.
func (ctx *Context) DuplicateNewScope() *Context {
	newContext := *ctx
	newContext.scope = createChildScope(ctx.scope)

	return &newContext
}

//// helpers
func createChildScope(old *Scope) *Scope {
	newScope := &Scope{
		fileInfo:    old.fileInfo,
		classRefMap: old.classRefMap,
		parent:      old,
		symbolMap:   map[string]SymbolInfo{},
		sgValue:     nil,
		returnValue: nil,
	}

	return newScope
}
