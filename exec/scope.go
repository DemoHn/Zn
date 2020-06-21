package exec

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// Scope - tmp Scope solution TODO: will move in the future!
type Scope interface {
	// GetValue - get variable name from current scope
	GetValue(ctx *Context, name string) (ZnValue, *error.Error)
	// SetValue - set variable value from current scope
	SetValue(ctx *Context, name string, value ZnValue) *error.Error
	// BindValue - bind value to current scope
	BindValue(ctx *Context, name string, value ZnValue) *error.Error
	// create new (nested) scope from current scope
	// fails if return scope is nil
	NewScope(ctx *Context, sType string) Scope
	// set current execution line
	SetCurrentLine(line int)
}

//// implementations

// RootScope - as named, this is the root scope for execution one program.
// usually it contains all active variables, scopes, etc
type RootScope struct {
	//// lexical scope
	// file - current execution file directory
	file string
	// currentLine - current exeuction line
	currentLine int
	// lineStack - lexical info of (parsed) current file
	lineStack *lex.LineStack
	//// lastValue - get last valid value even if there's no return statement
	lastValue ZnValue
}

// NewRootScope - create a rootScope from existing Lexer that
// derives from a program file, a piece of code, etc.
//
// That implies a program has one and only one RootScope.
//
// NOTE: When a program file "requires" another one, another RootScope is created
// for that "required" program file.
func NewRootScope() *RootScope {
	return &RootScope{}
}

// Init - init rootScope using new Lexer
func (rs *RootScope) Init(l *lex.Lexer) {
	rs.file = l.InputStream.Scope
	rs.currentLine = 0
	rs.lineStack = l.LineStack
	rs.lastValue = NewZnNull()
}

// SetCurrentLine -
func (rs *RootScope) SetCurrentLine(line int) {
	rs.currentLine = line
}

// GetValue - get variable name from current scope
func (rs *RootScope) GetValue(ctx *Context, name string) (ZnValue, *error.Error) {
	// find on globals first
	if symVal, inGlobals := ctx.globals[name]; inGlobals {
		return symVal, nil
	}
	// ...then in symbols
	symArr, ok := ctx.symbols[name]
	if !ok {
		return nil, error.NameNotDefined(name)
	}

	// find the nearest level of value
	if symArr == nil || len(symArr) == 0 {
		return nil, error.NameNotDefined(name)
	}
	return symArr[0].Value, nil
}

// SetValue - set variable value from current scope
func (rs *RootScope) SetValue(ctx *Context, name string, value ZnValue) *error.Error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}

	symArr, ok := ctx.symbols[name]
	if !ok {
		return error.NameNotDefined(name)
	}

	if symArr != nil && len(symArr) > 0 {
		symArr[0].Value = value
		if symArr[0].IsConstant {
			return error.AssignToConstant()
		}
		return nil
	}

	return error.NameNotDefined(name)
}

// BindValue - bind value to current scope
func (rs *RootScope) BindValue(ctx *Context, name string, value ZnValue) *error.Error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	newInfo := SymbolInfo{
		NestLevel:  1,
		Value:      value,
		IsConstant: false,
	}

	symArr, ok := ctx.symbols[name]
	if !ok {
		// init symbolInfo array
		ctx.symbols[name] = []SymbolInfo{newInfo}
		return nil
	}

	// check if there's variable re-declaration
	if len(symArr) > 0 && symArr[0].NestLevel == 1 {
		return error.NameRedeclared(name)
	}

	// prepend data
	ctx.symbols[name] = append([]SymbolInfo{newInfo}, ctx.symbols[name]...)
	return nil
}

// NewScope - create new (nested) scope from current scope
// fails if return scope is nil
func (rs *RootScope) NewScope(ctx *Context, sType string) Scope {
	// TODO
	return nil
}
