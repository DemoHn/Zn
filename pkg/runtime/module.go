package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/syntax"
)

type Module struct {
	// name - module name (when anonymous = false, it should be non-empty)
	name string
	// anonymous - when anonymous = true, the module has no name
	// but one context only allow ONE anonymous module
	anonymous bool
	// internal = true for internal modules (e.g. standard library, plugins)  these are imported via 《 》mark instead of “ ” mark. usually there's no source code in zinc (logics are written in Golang la...) so `module.lexer` = nil
	internal bool

	// lexer - the lexer of module's source code for mapping objects to source code.
	// if the module is a standard library, lexer = nil
	lexer *syntax.Lexer

	/* scopeStack - the call stack of execution scope
	   the stack looks like the following diagram:

	   +----------+
	N  | current  |
	   +----------+
	   | parent1  |
	   +----------+
	   | parent2  |
	   +----------+
	   |   ...    |
	   +----------+
	0  |  root    |
	   +----------+

	- push child scope when executing child block (e.g. 如何、每当)
	- pop the top scope when exiting child block.

	- The ROOT scope of module is `scopeStack[0]`
	- The parent scope of `scopeStack[N]` is `scopeStack[N-1]`
	- The child scope of `scopeStack[N]` is `scopeStack[N+1]`
	*/
	scopeStack []*Scope
	// exportValues - all classes and functions are exported for external
	// imports - so here we insert all exportable values to this map after first scan
	// note: all export values are constants.
	exportValues map[string]Element
}

func NewModule(name string, lexer *syntax.Lexer) *Module {
	return &Module{
		name:      name,
		anonymous: false,
		internal:  false,
		lexer:     lexer,
		// init root scope to ensure scopeStack NOT empty
		scopeStack:   []*Scope{NewScope(nil)},
		exportValues: map[string]Element{},
	}
}

func NewAnonymousModule(lexer *syntax.Lexer) *Module {
	return &Module{
		name:      "",
		anonymous: true,
		lexer:     lexer,
		// init root scope to ensure scopeStack NOT empty
		scopeStack:   []*Scope{NewScope(nil)},
		exportValues: map[string]Element{},
	}
}

// called
func NewInternalModule(name string) *Module {
	return &Module{
		name:      name,
		anonymous: false,
		internal:  true,
		lexer:     nil,
		// init root scope to ensure scopeStack NOT empty
		scopeStack:   []*Scope{NewScope(nil)},
		exportValues: map[string]Element{},
	}
}

func (m *Module) SetLexer(l *syntax.Lexer) {
	m.lexer = l
}

func (m *Module) GetName() string {
	return m.name
}

func (m *Module) GetLexer() *syntax.Lexer {
	return m.lexer
}

func (m *Module) IsAnonymous() bool {
	return m.anonymous
}

//// scopeStack operation
////
func (m *Module) GetCurrentScope() *Scope {
	stackLen := len(m.scopeStack)
	if stackLen == 0 {
		return nil
	}

	return m.scopeStack[stackLen-1]
}

func (m *Module) PushScope() *Scope {
	sp := m.GetCurrentScope()
	if sp == nil {
		return nil
	}

	childScope := NewScope(sp.thisValue)
	// push scope into ScopeStack
	m.scopeStack = append(m.scopeStack, childScope)

	return m.GetCurrentScope()
}

func (m *Module) AddScope(scope *Scope) {
	m.scopeStack = append(m.scopeStack, scope)
}

func (m *Module) PopScope() {
	stackLen := len(m.scopeStack)
	if stackLen == 0 {
		return
	}

	// pop last (current) scope
	m.scopeStack = m.scopeStack[:stackLen-1]
}

// FindScopeValue - find symbol in the context from the latest scope
// up to its first one
func (m *Module) FindScopeValue(name string) (Element, error) {
	// iterate from last to very first
	for cursor := len(m.scopeStack) - 1; cursor >= 0; cursor-- {
		sp := m.scopeStack[cursor]
		if ok, val := sp.GetSymbolValue(name); ok {
			return val, nil
		}
	}

	return nil, zerr.NameNotDefined(name)
}

// FindScopeValue - find symbol in the context from the latest scope
// up to its first one
func (m *Module) FindScopeSymbol(name string) (SymbolInfo, error) {
	// iterate from last to very first
	for cursor := len(m.scopeStack) - 1; cursor >= 0; cursor-- {
		sp := m.scopeStack[cursor]
		if ok, sym := sp.GetSymbol(name); ok {
			return sym, nil
		}
	}

	return SymbolInfo{}, zerr.NameNotDefined(name)
}

// SetScopeValue - set value of an existing symbol (whatever in current scope or root scope la..)
// there, the process includes 3 steps:
// 1. find the symbol in scope stack
// 2. set new value of the symbol
// 3. if no symbol found, throw error directly
func (m *Module) SetScopeValue(name string, value Element) error {
	// iterate from last to very first
	for cursor := len(m.scopeStack) - 1; cursor >= 0; cursor-- {
		sp := m.scopeStack[cursor]
		if ok, sym := sp.GetSymbol(name); ok {
			if sym.isConst {
				return zerr.AssignToConstant()
			}
			sp.SetSymbolValue(name, value, false, sym.module)
			return nil
		}
	}
	return zerr.NameNotDefined(name)
}

// BindSymbol - bind a non-const value on current scope - however, if the same symbol has bound, then an error occurs.
func (m *Module) BindSymbol(name string, value Element, isConst bool, rebindCheck bool) error {
	if sp := m.GetCurrentScope(); sp != nil {
		// bind value on current scope
		if rebindCheck {
			if ok, _ := sp.GetSymbol(name); ok {
				return zerr.NameRedeclared(name)
			}
		}

		// set value
		sp.SetSymbolValue(name, value, isConst, m)
	}
	return nil
}

// BindImportSymbol - bind a non-const value on current scope from another module. by default, if the same symbol has bound, then an error occurs.
func (m *Module) BindImportSymbol(name string, value Element, refModule *Module) error {
	if sp := m.GetCurrentScope(); sp != nil {
		// flushing is NOT allowed
		if ok, _ := sp.GetSymbol(name); ok {
			return zerr.NameRedeclared(name)
		}

		// set value
		sp.SetSymbolValue(name, value, true, refModule)
	}
	return nil
}

//// imports & exports
func (m *Module) AddExportValue(name string, value Element) error {
	if _, ok := m.exportValues[name]; ok {
		return zerr.NameRedeclared(name)
	}

	m.exportValues[name] = value
	return nil
}

func (m *Module) GetAllExportValues() map[string]Element {
	return m.exportValues
}

func (m *Module) GetExportValue(name string) (Element, error) {
	if v, ok := m.exportValues[name]; ok {
		return v, nil
	}

	return nil, zerr.NameNotDefined(name)
}
