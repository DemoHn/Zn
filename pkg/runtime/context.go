package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/syntax"
)

// Context is a global variable that stores current execution
// states, global configurations
type Context struct {
	// globals - stores all global variables
	globals map[string]Value

	*DependencyTree
	// ScopeStack - trace scopes
	ScopeStack []*Scope
}

// NewContext - create new Zn Context. Notice through the life-cycle
// of one code execution, there's only one running context to store all states.
func NewContext(globalsMap map[string]Value) *Context {
	return &Context{
		globals:        globalsMap,
		DependencyTree: NewDependencyTree(),
		ScopeStack:     []*Scope{},
	}
}

func (ctx *Context) GetCurrentScope() *Scope {
	stackLen := len(ctx.ScopeStack)
	if stackLen == 0 {
		return nil
	}

	return ctx.ScopeStack[stackLen-1]
}

func (ctx *Context) PushScope(module *Module, lexer *syntax.Lexer) *Scope {
	scope := NewScope(module, lexer)
	// push scope into ScopeStack
	ctx.ScopeStack = append(ctx.ScopeStack, scope)

	return ctx.GetCurrentScope()
}

// PushChildScope - create new scope with same module from parent scope
func (ctx *Context) PushChildScope() *Scope {
	sp := ctx.GetCurrentScope()
	if sp == nil {
		return nil
	}
	childScope := NewChildScope(sp)
	// push scope into ScopeStack
	ctx.ScopeStack = append(ctx.ScopeStack, childScope)

	return ctx.GetCurrentScope()
}

func (ctx *Context) PopScope() {
	stackLen := len(ctx.ScopeStack)
	if stackLen == 0 {
		return
	}

	// pop last element
	ctx.ScopeStack = ctx.ScopeStack[:stackLen-1]
}

// SetCurrentLine - set lineIdx to current running scope
func (ctx *Context) SetCurrentLine(line int) {
	if sp := ctx.GetCurrentScope(); sp != nil {
		sp.SetExecLineIdx(line)
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
	// ...then in symbols
	sp := ctx.GetCurrentScope()
	for sp != nil {
		// 1. look up from current module's import map
		if module := sp.GetModule(); module != nil {
			if val, err := module.GetSymbol(name); err == nil {
				return val, nil
			}
		}
		// 2. look up from scope's symbol map
		sym, ok := sp.symbolMap[name]
		if ok {
			return sym.value, nil
		}

		sp = sp.parent
	}
	return nil, zerr.NameNotDefined(name)
}

// SetSymbol -
func (ctx *Context) SetSymbol(name string, value Value) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	// ...then in symbols
	sp := ctx.GetCurrentScope()
	for sp != nil {
		sym, ok := sp.symbolMap[name]
		if ok {
			if sym.isConst {
				return zerr.AssignToConstant()
			}
			sp.symbolMap[name] = SymbolInfo{value, false}
			return nil
		}

		sp = sp.parent
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