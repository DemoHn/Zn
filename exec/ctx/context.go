package ctx

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

// DuplicateNewScope - create a new Context with new scope which parent points to duplicator's scope.
func (ctx *Context) DuplicateNewScope() *Context {
	newContext := *ctx
	newContext.scope = createChildScope(ctx.scope)

	return &newContext
}

// ExecuteCode - execute program from input Zn code (whether from file or REPL)
func (ctx *Context) ExecuteCode(in *lex.InputStream) (Value, *error.Error) {
	program, err := ctx.parseCode(in)
	if err != nil {
		return nil, err
	}
	// init scope
	ctx.initScope(program.Lexer)

	// eval program
	return ctx.execProgram(program)
}

// parseCode - lex & parse code text
func (ctx *Context) parseCode(in *lex.InputStream) (*syntax.Program, *error.Error) {
	l := lex.NewLexer(in)
	p := syntax.NewParser(l)
	// start
	block, err := p.Parse()
	if err != nil {
		return nil, err
	}

	return syntax.NewProgramNode(block, l), nil
}

func (ctx *Context) execProgram(program *syntax.Program) (Value, *error.Error) {
	err := evalProgram(ctx, program)
	if err != nil {
		cursor := err.GetCursor()

		// wrapError if lineInfo is missing (mostly for non-syntax errors)
		// If lineInfo missing, then we will add current execution line and hide some part to
		// display errors properly.
		if cursor.LineNum == 0 {
			fileInfo := ctx.scope.fileInfo
			newCursor := error.Cursor{
				File:    fileInfo.file,
				LineNum: fileInfo.currentLine,
				Text:    fileInfo.lineStack.GetLineText(fileInfo.currentLine, false),
			}
			err.SetCursor(newCursor)
		}
		return nil, err
	}
	return ctx.scope.returnValue, nil
}

// InitScope - init root scope
func (ctx *Context) initScope(l *lex.Lexer) {
	fileInfo := &FileInfo{
		file:        l.InputStream.GetFile(),
		currentLine: 0,
		lineStack:   l.LineStack,
	}
	if ctx.scope == nil {
		newScope := &Scope{
			fileInfo:    fileInfo,
			classRefMap: map[string]ClassRef{},
			parent:      nil,
			symbolMap:   map[string]SymbolInfo{},
			sgValue:     nil,
			thisValue:   nil,
			returnValue: NewNull(),
		}
		ctx.scope = newScope
	} else {
		// refresh scope fileInfo
		ctx.scope.fileInfo = fileInfo
	}
}

// resetScopeValue - reset classRefMap, symbolMap, sgValue, thisValue, returnValue
// to initial value
func (ctx *Context) resetScopeValue() {
	if ctx.scope != nil {
		// preserve fileInfo
		fileInfo := ctx.scope.fileInfo
		ctx.scope = &Scope{
			fileInfo:    fileInfo,
			classRefMap: map[string]ClassRef{},
			parent:      nil,
			symbolMap:   map[string]SymbolInfo{},
			sgValue:     nil,
			thisValue:   nil,
			returnValue: NewNull(),
		}
	}
}

//// helpers
func createChildScope(old *Scope) *Scope {
	newScope := &Scope{
		fileInfo:    old.fileInfo,
		classRefMap: old.classRefMap,
		parent:      old,
		symbolMap:   map[string]SymbolInfo{},
		sgValue:     nil,
		thisValue:   nil,
		returnValue: nil,
	}

	return newScope
}
