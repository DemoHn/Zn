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

const (
	// RootLevel - RootScope nest level
	RootLevel = 1
	//
	sTypeRoot = "scopeRoot"
	sTypeFunc = "scopeFunc"
)

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
	return &RootScope{
		lastValue: NewZnNull(),
	}
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
	return getValue(ctx, name)
}

// SetValue - set variable value from current scope
func (rs *RootScope) SetValue(ctx *Context, name string, value ZnValue) *error.Error {
	return setValue(ctx, name, value)
}

// BindValue - bind value to current scope
func (rs *RootScope) BindValue(ctx *Context, name string, value ZnValue) *error.Error {
	return bindValue(ctx, RootLevel, name, value)
}

// NewScope - create new (nested) scope from current scope
// fails if return scope is nil
func (rs *RootScope) NewScope(ctx *Context, sType string) Scope {
	if sType == sTypeFunc {
		return &FuncScope{
			returnValue: NewZnNull(),
			root:        rs,
			parent:      rs,
			nestLevel:   RootLevel + 1,
		}
	}
	return nil
}

// SetLastValue - set last value
func (rs *RootScope) SetLastValue(value ZnValue) {
	rs.lastValue = value
}

// GetLastValue -
func (rs *RootScope) GetLastValue() ZnValue {
	return rs.lastValue
}

// FuncScope - function scope
type FuncScope struct {
	//// returnValue - the final return value of current scope
	returnValue ZnValue
	root        *RootScope
	parent      Scope
	nestLevel   int
	returnFlag  bool
}

// GetValue - get variable name from current scope
func (fs *FuncScope) GetValue(ctx *Context, name string) (ZnValue, *error.Error) {
	return getValue(ctx, name)
}

// SetValue - set variable value from current scope
func (fs *FuncScope) SetValue(ctx *Context, name string, value ZnValue) *error.Error {
	return setValue(ctx, name, value)
}

// BindValue - bind value to current scope
func (fs *FuncScope) BindValue(ctx *Context, name string, value ZnValue) *error.Error {
	return bindValue(ctx, fs.nestLevel, name, value)
}

// NewScope - create new (nested) scope from current scope
// fails if return scope is nil
func (fs *FuncScope) NewScope(ctx *Context, sType string) Scope {
	if sType == sTypeFunc {
		return &FuncScope{
			returnValue: NewZnNull(),
			root:        fs.root,
			parent:      fs,
			nestLevel:   fs.nestLevel + 1,
		}
	}
	return nil
}

// SetCurrentLine - set current execution line
func (fs *FuncScope) SetCurrentLine(line int) {
	fs.root.SetCurrentLine(line)
}

// SetReturnValue - set last value
func (fs *FuncScope) SetReturnValue(value ZnValue) {
	fs.returnValue = value
}

// GetReturnValue -
func (fs *FuncScope) GetReturnValue() ZnValue {
	return fs.returnValue
}

// SetReturnFlag - set last value
func (fs *FuncScope) SetReturnFlag(flag bool) {
	fs.returnFlag = flag
}

// GetReturnFlag -
func (fs *FuncScope) GetReturnFlag() bool {
	return fs.returnFlag
}

//// helpers
func bindValue(ctx *Context, nestLevel int, name string, value ZnValue) *error.Error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	newInfo := SymbolInfo{
		NestLevel:  nestLevel,
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
	if len(symArr) > 0 && symArr[0].NestLevel == nestLevel {
		return error.NameRedeclared(name)
	}

	// prepend data
	ctx.symbols[name] = append([]SymbolInfo{newInfo}, ctx.symbols[name]...)
	return nil
}

func setValue(ctx *Context, name string, value ZnValue) *error.Error {
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

func getValue(ctx *Context, name string) (ZnValue, *error.Error) {
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
