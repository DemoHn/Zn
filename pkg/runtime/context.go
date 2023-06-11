package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
)

// Context is a global variable that stores current execution
// states, global configurations
type Context struct {
	// globals - stores all global variables
	globals map[string]Value
	// hasPrinted - if stdout has been used to output message before program end, set `hasPrinted` -> true; so that after message is done
	hasPrinted bool
	*DependencyTree

	currentModule *Module
	// callStack - get current call module & line for traceback
	callStack []CallInfo
}

type CallInfo struct {
	*Module
	LastLineIdx int
}

// NewContext - create new Zn Context. Notice through the life-cycle
// of one code execution, there's only one running context to store all states.
func NewContext(globalsMap map[string]Value) *Context {
	return &Context{
		globals:        globalsMap,
		hasPrinted:     false,
		DependencyTree: NewDependencyTree(),
		currentModule:  nil,
		callStack:      []CallInfo{},
	}
}

func (ctx *Context) GetCurrentScope() *Scope {
	if ctx.currentModule != nil {
		return ctx.currentModule.GetCurrentScope()
	}
	return nil
}

// PushChildScope - create new scope with same module from parent scope
func (ctx *Context) PushChildScope() *Scope {
	sp := ctx.GetCurrentScope()
	if sp == nil {
		return nil
	}

	return ctx.currentModule.PushScope()
}

func (ctx *Context) PopScope() {
	if ctx.currentModule != nil {
		ctx.currentModule.PopScope()
	}
}

// SetCurrentLine - set lineIdx to current running scope
func (ctx *Context) SetCurrentLine(line int) {
	if ctx.currentModule != nil {
		ctx.currentModule.SetCurrentLine(line)
	}
}

//// scope symbols getters / setters

// FindSymbol - find symbol in the context from current scope
// up to its root scope
func (ctx *Context) FindSymbol(name string) (Value, error) {
	// find on globals first
	if symVal, inGlobals := ctx.globals[name]; inGlobals {
		return symVal, nil
	}

	if ctx.currentModule != nil {
		return ctx.currentModule.FindScopeValue(name)
	}
	return nil, zerr.UnexpectedNilModule()
}

// SetSymbol -
func (ctx *Context) SetSymbol(name string, value Value) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	// ...then in symbols
	if ctx.currentModule != nil {
		return ctx.currentModule.SetScopeValue(name, value)
	}
	return zerr.UnexpectedNilModule()
}

// BindSymbol - bind non-const value with re-declaration check on same scope
func (ctx *Context) BindSymbol(name string, value Value) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	if ctx.currentModule != nil {
		return ctx.currentModule.BindSymbol(name, SymbolInfo{value, false}, true)
	}
	return zerr.UnexpectedNilModule()
}

// BindSymbolDecl - bind value for declaration statement - that variables could be re-bind.
func (ctx *Context) BindSymbolDecl(name string, value Value, isConst bool) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	if ctx.currentModule != nil {
		return ctx.currentModule.BindSymbol(name, SymbolInfo{value, isConst}, false)
	}
	return zerr.UnexpectedNilModule()
}

// BindScopeSymbolDecl - bind value for declaration statement - that variables could be re-bind.
func (ctx *Context) BindScopeSymbolDecl(scope *Scope, name string, value Value) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	if scope != nil {
		scope.SetSymbolValue(name, false, value)
	}
	return nil
}

// FindThisValue -
func (ctx *Context) FindThisValue() (Value, error) {
	sp := ctx.GetCurrentScope()
	for sp != nil {
		thisValue := sp.thisValue
		if thisValue != nil {
			return thisValue, nil
		}

		sp = sp.parent
	}

	return nil, zerr.PropertyNotFound("thisValue")
}

// MarkHasPrinted - called by `显示` function only
func (ctx *Context) MarkHasPrinted() {
	ctx.hasPrinted = true
}

func (ctx *Context) GetHasPrinted() bool {
	return ctx.hasPrinted
}
