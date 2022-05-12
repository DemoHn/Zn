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
	// import - imported value from stdlib or elsewhere
	imports map[string]Value
	// scopeStack - trace scopes
	scopeStack []Scope
}

// NewContext - create new Zn Context. Notice through the life-cycle
// of one code execution, there's only one running context to store all states.
func NewContext(globalsMap map[string]Value) *Context {
	return &Context{
		globals: globalsMap,
		imports: map[string]Value{},
		scopeStack: []Scope{},
	}
}

// ImportModule -
func (ctx *Context) ImportModule(moduleName string, l *syntax.Lexer) {
	module := NewModule(moduleName, l)
	scope := NewScope(module)

	// push scope into scopeStack
	ctx.scopeStack = append(ctx.scopeStack, scope)
}

func (ctx *Context) GetCurrentScope() *Scope {
	stackLen := len(ctx.scopeStack)
	if stackLen == 0 {
		return nil
	}

	lastScope := ctx.scopeStack[stackLen-1]
	return &lastScope
}

func (ctx *Context) PushNewScope() *Scope {

}

// GetScope -
func (ctx *Context) GetScope() *Scope {
	return ctx.scope
}

// SetScope -
func (ctx *Context) SetScope(sp *Scope) {
	ctx.scope = sp
}


//// scope symbols getters / setters

// FindSymbol - find symbol in the context from current scope
// up to its root scope
func (ctx *Context) FindSymbol(name string) (Value, error) {
	// find on globals first
	if symVal, inGlobals := ctx.globals[name]; inGlobals {
		return symVal, nil
	}
	// next in imports
	if imVal, inImports := ctx.imports[name]; inImports {
		return imVal, nil
	}
	// ...then in symbols
	sp := ctx.scope
	for sp != nil {
		sym, ok := sp.symbolMap[name]
		if ok {
			return sym.value, nil
		}
		// if not found, search its parent
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
	sp := ctx.scope
	for sp != nil {
		sym, ok := sp.symbolMap[name]
		if ok {
			if sym.isConst {
				return zerr.AssignToConstant()
			}
			sp.symbolMap[name] = SymbolInfo{value, false}
			return nil
		}
		// if not found, search its parent
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
	if ctx.scope != nil {
		if _, ok := ctx.scope.symbolMap[name]; ok {
			return zerr.NameRedeclared(name)
		}
		// set value
		ctx.scope.symbolMap[name] = SymbolInfo{value, false}
	}
	return nil
}

// BindSymbolDecl - bind value for declaration statement - that variables could be re-bind.
func (ctx *Context) BindSymbolDecl(name string, value Value, isConst bool) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	if ctx.scope != nil {
		ctx.scope.symbolMap[name] = SymbolInfo{value, isConst}
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
	sp := ctx.scope
	for sp != nil {
		thisValue := sp.thisValue
		if thisValue != nil {
			return thisValue, nil
		}

		// otherwise, find thisValue from parent scope
		sp = sp.parent
	}

	return nil, zerr.PropertyNotFound("thisValue")
}

// fetch from imports
// GetImportValue -
func (ctx *Context) GetImportValue(name string) (Value, error) {
	// find on globals first
	if symVal, inImports := ctx.imports[name]; inImports {
		return symVal, nil
	}

	return nil, zerr.NameNotDefined(name)
}

// SetImportValue -
func (ctx *Context) SetImportValue(name string, value Value) error {
	if _, inImports := ctx.imports[name]; inImports {
		return zerr.NameRedeclared(name)
	}
	ctx.imports[name] = value
	return nil
}
