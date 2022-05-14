package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/syntax"
)

// Context is a global variable that stores current execution
// states, global configurations
type Context struct {
	astBuilder syntax.ASTBuilder
	// globals - stores all global variables
	globals map[string]Value
	// scopeStack - trace scopes
	scopeStack []*Scope
}

// NewContext - create new Zn Context. Notice through the life-cycle
// of one code execution, there's only one running context to store all states.
func NewContext(globalsMap map[string]Value, builder syntax.ASTBuilder) *Context {
	return &Context{
		astBuilder: builder,
		globals: globalsMap,
		scopeStack: []*Scope{},
	}
}

func (ctx *Context) GetBuilder() syntax.ASTBuilder {
	return ctx.astBuilder
}

func (ctx *Context) GetCurrentScope() *Scope {
	stackLen := len(ctx.scopeStack)
	if stackLen == 0 {
		return nil
	}

	return ctx.scopeStack[stackLen-1]
}

func (ctx *Context) PushScope(module *Module) *Scope {
	scope := NewScope(module)
	// push scope into scopeStack
	ctx.scopeStack = append(ctx.scopeStack, scope)

	return ctx.GetCurrentScope()
}

// PushChildScope - create new scope with same module from parent scope
func (ctx *Context) PushChildScope() *Scope {
	sp := ctx.GetCurrentScope()
	if sp == nil {
		return nil
	}
	childScope := NewChildScope(sp)
	// push scope into scopeStack
	ctx.scopeStack = append(ctx.scopeStack, childScope)

	return ctx.GetCurrentScope()
}

func (ctx *Context) PopScope() {
	stackLen := len(ctx.scopeStack)
	if stackLen == 0 {
		return
	}

	// pop last element
	ctx.scopeStack = ctx.scopeStack[:stackLen-1]
}


//// scope symbols getters / setters

// FindSymbol - find symbol in the context from current scope
// up to its root scope
func (ctx *Context) FindSymbol(name string) (Value, error) {
	// find on globals first
	if symVal, inGlobals := ctx.globals[name]; inGlobals {
		return symVal, nil
	}
	// ...then in symbols
	for i := len(ctx.scopeStack)-1; i >= 0; i-- {
		sp := ctx.scopeStack[i]
		sym, ok := sp.symbolMap[name]
		if ok {
			return sym.value, nil
		}
		// if not found, search in prev scope
	}
	return nil, zerr.NameNotDefined(name)
}

// SetSymbol -
func (ctx *Context) SetSymbol(name string, value Value) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	// ...then in symbols
	for i := len(ctx.scopeStack)-1; i >= 0; i-- {
		sp := ctx.scopeStack[i]
		sym, ok := sp.symbolMap[name]
		if ok {
			if sym.isConst {
				return zerr.AssignToConstant()
			}
			sp.symbolMap[name] = SymbolInfo{value, false}
			return nil
		}
		// if not found, search in previous scope
	}
	return zerr.NameNotDefined(name)
}

// BindSymbol - bind non-const value with re-declaration check on same scope
func (ctx *Context) BindSymbol(name string, value Value) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	// bind directly
	sp := ctx.GetCurrentScope()
	if sp != nil {
		if _, ok := sp.symbolMap[name]; ok {
			return zerr.NameRedeclared(name)
		}
		// set value
		sp.symbolMap[name] = SymbolInfo{value, false}
	}
	return nil
}

// BindSymbolDecl - bind value for declaration statement - that variables could be re-bind.
func (ctx *Context) BindSymbolDecl(name string, value Value, isConst bool) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	sp := ctx.GetCurrentScope()
	if sp != nil {
		sp.symbolMap[name] = SymbolInfo{value, isConst}
	}
	return nil
}

// BindScopeSymbolDecl - bind value for declaration statement - that variables could be re-bind.
func (ctx *Context) BindScopeSymbolDecl(scope *Scope, name string, value Value) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	if scope != nil {
		scope.symbolMap[name] = SymbolInfo{value, false}
	}
	return nil
}

// FindThisValue -
func (ctx *Context) FindThisValue() (Value, error) {
	for i := len(ctx.scopeStack)-1; i >= 0; i-- {
		sp := ctx.scopeStack[i]
		thisValue := sp.thisValue
		if thisValue != nil {
			return thisValue, nil
		}

		// otherwise, find thisValue from previous scope
	}

	return nil, zerr.PropertyNotFound("thisValue")
}