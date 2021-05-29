package ctx

import (
	"github.com/DemoHn/Zn/debug"
	"github.com/DemoHn/Zn/error"
)

// Context is a global variable that stores current execution
// states, global configurations
type Context struct {
	// globals - stores all global variables
	globals map[string]Value
	// import - imported value from stdlib or elsewhere
	imports map[string]Value
	// fileInfo -
	fileInfo *FileInfo
	// a seperate map to store inner debug data
	// usage: call （__probe：「tagName」，variable）
	// it will record all logs (including variable value, curernt scope, etc.)
	// the value is deep-copied so don't worry - the value logged won't be changed
	_probe *debug.Probe
	// Scope -
	scope *Scope
}

// NewContext - create new Zn Context. Notice through the life-cycle
// of one code execution, there's only one running context to store all states.
func NewContext(globalsMap map[string]Value) *Context {
	return &Context{
		globals: globalsMap,
		imports: map[string]Value{},
		_probe:  debug.NewProbe(),
		scope:   NewScope(),
	}
}

// GetScope -
func (ctx *Context) GetScope() *Scope {
	return ctx.scope
}

// SetScope -
func (ctx *Context) SetScope(sp *Scope) {
	ctx.scope = sp
}

// SetFileInfo
func (ctx *Context) SetFileInfo(fileInfo *FileInfo) {
	ctx.fileInfo = fileInfo
}

// GetFileInfo -
func (ctx *Context) GetFileInfo() *FileInfo {
	return ctx.fileInfo
}

//// scope symbols getters / setters

// FindSymbol - find symbol in the context from current scope
// up to its root scope
func (ctx *Context) FindSymbol(name string) (Value, *error.Error) {
	// find on globals first
	if symVal, inGlobals := ctx.globals[name]; inGlobals {
		return symVal, nil
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
	return nil, error.NameNotDefined(name)
}

// SetSymbol -
func (ctx *Context) SetSymbol(name string, value Value) *error.Error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	// ...then in symbols
	sp := ctx.scope
	for sp != nil {
		sym, ok := sp.symbolMap[name]
		if ok {
			if sym.isConst {
				return error.AssignToConstant()
			}
			sp.symbolMap[name] = SymbolInfo{value, false}
			return nil
		}
		// if not found, search its parent
		sp = sp.parent
	}
	return error.NameNotDefined(name)
}

// BindSymbol - bind non-const value with re-declaration check on same scope
func (ctx *Context) BindSymbol(name string, value Value) *error.Error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	// bind directly
	if ctx.scope != nil {
		if _, ok := ctx.scope.symbolMap[name]; ok {
			return error.NameRedeclared(name)
		}
		// set value
		ctx.scope.symbolMap[name] = SymbolInfo{value, false}
	}
	return nil
}

// BindSymbolDecl - bind value for declaration statement - that variables could be re-bind.
func (ctx *Context) BindSymbolDecl(name string, value Value, isConst bool) *error.Error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	if ctx.scope != nil {
		ctx.scope.symbolMap[name] = SymbolInfo{value, isConst}
	}
	return nil

}
