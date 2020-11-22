package exec

import (
	"github.com/DemoHn/Zn/debug"
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
	"github.com/DemoHn/Zn/syntax"
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

// Result - context execution result structure
// NOTICE: when HasError = true, Value = nil, while execution yields error
//         when HasError = false, Error = nil, Value = <result Value>
//
// Currently only one value is supported as return argument.
type Result struct {
	HasError bool
	Value    Value
	Error    *error.Error
}

// NewContext - create new Zn Context. Notice through the life-cycle
// of one code execution, there's only one running context to store all states.
func NewContext() *Context {
	return &Context{
		globals: globalValues,
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

// ExecuteCode - execute program from input Zn code (whether from file or REPL)
func (ctx *Context) ExecuteCode(in *lex.InputStream) Result {
	l := lex.NewLexer(in)
	p := syntax.NewParser(l)
	// start
	block, err := p.Parse()
	if err != nil {
		return Result{true, nil, err}
	}

	// init scope
	ctx.InitScope(l)

	// construct root (program) node
	program := syntax.NewProgramNode(block)

	// eval program
	if err := evalProgram(ctx, program); err != nil {
		wrapError(ctx, err)
		return Result{true, nil, err}
	}
	return Result{false, ctx.scope.returnValue, nil}
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

// wrapError if lineInfo is missing (mostly for non-syntax errors)
// If lineInfo missing, then we will add current execution line and hide some part to
// display errors properly.
func wrapError(ctx *Context, err *error.Error) {
	cursor := err.GetCursor()

	if cursor.LineNum == 0 {
		fileInfo := ctx.scope.fileInfo
		newCursor := error.Cursor{
			File:    fileInfo.file,
			LineNum: fileInfo.currentLine,
			Text:    fileInfo.lineStack.GetLineText(fileInfo.currentLine, false),
		}
		err.SetCursor(newCursor)
	}
}
