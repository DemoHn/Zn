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

	// modulegraph - store module dependency & all preloaded modules
	moduleGraph *ModuleGraph

	// current execution module. Must be NON-EMPTY at initialization
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
// NOTE: initModule DO NOT accept nil value at initialization!!
func NewContext(globalsMap map[string]Value, initModule *Module) *Context {
	// init module dep graph
	graph := NewModuleGraph()
	graph.AddModule(initModule)

	return &Context{
		globals:       globalsMap,
		hasPrinted:    false,
		moduleGraph:   graph,
		currentModule: initModule,
		callStack:     []CallInfo{},
	}
}

func (ctx *Context) GetCurrentScope() *Scope {
	if ctx.currentModule != nil {
		return ctx.currentModule.GetCurrentScope()
	}
	return nil
}

func (ctx *Context) FindParentScope() *Scope {
	if ctx.currentModule != nil {
		sLen := len(ctx.currentModule.scopeStack)

		if sLen > 1 {
			return ctx.currentModule.scopeStack[sLen-2]
		}
	}
	return nil
}

// PushScope - create new scope with same module from parent scope
func (ctx *Context) PushScope() *Scope {
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

//// enter & exist modules
func (ctx *Context) EnterModule(module *Module) {
	m := ctx.currentModule
	// push callstack
	if m != nil {
		ctx.callStack = append(ctx.callStack, CallInfo{
			Module:      m,
			LastLineIdx: m.currentLine,
		})
	}

	// set current module
	ctx.currentModule = module
}

func (ctx *Context) ExitModule() {
	sLen := len(ctx.callStack)
	if sLen > 0 {
		// get last module in callstack
		last := ctx.callStack[sLen-1]
		// pop last one
		ctx.callStack = ctx.callStack[:sLen-1]

		ctx.currentModule = last.Module
	}
}

func (ctx *Context) FindModule(name string) *Module {
	return ctx.moduleGraph.FindModule(name)
}

func (ctx *Context) GetCurrentModule() *Module {
	return ctx.currentModule
}

func (ctx *Context) GetCallStack() []CallInfo {
	return ctx.callStack
}

//// scope symbols getters / setters

// FindSymbol - find symbol in the context from current scope
// up to its root scope
func (ctx *Context) FindSymbol(name string) (Value, error) {
	// find on globals first
	if symVal, inGlobals := ctx.globals[name]; inGlobals {
		return symVal, nil
	}

	return ctx.currentModule.FindScopeValue(name)
}

// SetSymbol -
func (ctx *Context) SetSymbol(name string, value Value) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	// ...then in symbols
	return ctx.currentModule.SetScopeValue(name, value)
}

// BindSymbol - bind non-const value with re-declaration check on same scope
func (ctx *Context) BindSymbol(name string, value Value) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}

	return ctx.currentModule.BindSymbol(name, value, false, true)
}

// BindSymbolDecl - bind value for declaration statement - that variables could be re-bind.
func (ctx *Context) BindSymbolDecl(name string, value Value, isConst bool) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}

	return ctx.currentModule.BindSymbol(name, value, isConst, false)
}

// BindScopeSymbolDecl - bind value for declaration statement - that variables could be re-bind.
func (ctx *Context) BindScopeSymbolDecl(scope *Scope, name string, value Value) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	if scope != nil {
		scope.SetSymbolValue(name, value, false)
	}
	return nil
}

// FindThisValue -
func (ctx *Context) FindThisValue() (Value, error) {
	m := ctx.currentModule
	for cursor := len(m.scopeStack) - 1; cursor >= 0; cursor-- {
		sp := m.scopeStack[cursor]
		if sp.thisValue != nil {
			return sp.thisValue, nil
		}
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
