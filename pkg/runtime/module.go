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

	// lexer - the lexer of module's source code for mapping objects to source code.
	// if the module is a standard library, lexer = nil
	lexer *syntax.Lexer

	// currentLine - current execution lineNum (index)
	currentLine int
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

	// importRefs - after importing symbols from other external modules,
	// we need to save the relationship between symbol and its original module, so that we could shift to original module's scope
	importRefs map[string]*Module
	// exportValues - all classes and functions are exported for external
	// imports - so here we insert all exportable values to this map after first scan
	// note: all export values are constants.
	exportValues map[string]Value
}

func NewModule(name string, lexer *syntax.Lexer) *Module {
	return &Module{
		name:        name,
		anonymous:   false,
		lexer:       lexer,
		currentLine: 0,
		// init root scope to ensure scopeStack NOT empty
		scopeStack:   []*Scope{NewScope()},
		importRefs:   map[string]*Module{},
		exportValues: map[string]Value{},
	}
}

func NewAnonymousModule(lexer *syntax.Lexer) *Module {
	return &Module{
		name:        "",
		anonymous:   true,
		lexer:       lexer,
		currentLine: 0,
		// init root scope to ensure scopeStack NOT empty
		scopeStack:   []*Scope{NewScope()},
		exportValues: map[string]Value{},
	}
}

// SetCurrentLine - set lineIdx to current running scope of the module
func (m *Module) SetCurrentLine(line int) {
	m.currentLine = line
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
	if sp := m.GetCurrentScope(); sp == nil {
		return nil
	}
	childScope := NewScope()
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
func (m *Module) FindScopeValue(name string) (Value, error) {
	// iterate from last to very first
	for cursor := len(m.scopeStack) - 1; cursor >= 0; cursor-- {
		sp := m.scopeStack[cursor]
		if ok, val := sp.GetSymbolValue(name); ok {
			return val, nil
		}
	}

	return nil, zerr.NameNotDefined(name)
}

// SetScopeValue - set value of an existing symbol (whatever in current scope or root scope la..)
// there, the process includes 3 steps:
// 1. find the symbol in scope stack
// 2. set new value of the symbol
// 3. if no symbol found, throw error directly
func (m *Module) SetScopeValue(name string, value Value) error {
	// iterate from last to very first
	for cursor := len(m.scopeStack) - 1; cursor >= 0; cursor-- {
		sp := m.scopeStack[cursor]
		if ok, sym := sp.GetSymbol(name); ok {
			if sym.isConst {
				return zerr.AssignToConstant()
			}
			sp.SetSymbolValue(name, value, false)
			return nil
		}
	}
	return zerr.NameNotDefined(name)
}

// BindValue - bind a non-const value on current scope - however, if the same symbol has bound, then an error occurs.
func (m *Module) BindSymbol(name string, value Value, isConst bool, rebindCheck bool) error {
	if sp := m.GetCurrentScope(); sp != nil {
		// bind value on current scope
		if rebindCheck {
			if ok, _ := sp.GetSymbol(name); ok {
				return zerr.NameRedeclared(name)
			}
		}

		// set value
		sp.SetSymbolValue(name, value, isConst)
	}
	return nil
}

//// imports & exports
func (m *Module) AddExportValue(name string, value Value) error {
	if _, ok := m.exportValues[name]; ok {
		return zerr.NameRedeclared(name)
	}

	m.exportValues[name] = value
	return nil
}

func (m *Module) GetAllExportValues() map[string]Value {
	return m.exportValues
}

func (m *Module) GetExportValue(name string) (Value, error) {
	if v, ok := m.exportValues[name]; ok {
		return v, nil
	}

	return nil, zerr.NameNotDefined(name)
}

func (m *Module) SetImportRef(symbol string, module *Module) error {
	if _, ok := m.importRefs[symbol]; ok {
		return zerr.NameRedeclared(symbol)
	}

	m.importRefs[symbol] = module
	return nil
}

func (m *Module) GetImportRef(symbol string) *Module {
	return m.importRefs[symbol]
}
