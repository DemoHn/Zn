package runtime

import (
	"github.com/DemoHn/Zn/pkg/syntax"
)

type Module struct {
	// name - module name
	name        string
	lexer       *syntax.Lexer
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
}

func NewModule(name string, lexer *syntax.Lexer) *Module {
	return &Module{
		name:        name,
		lexer:       lexer,
		currentLine: 0,
		// init root scope to ensure scopeStack NOT empty
		scopeStack: []*Scope{NewScope()},
	}
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

func (m *Module) PopScope() {
	stackLen := len(m.scopeStack)
	if stackLen == 0 {
		return
	}

	// pop last (current) scope
	m.scopeStack = m.scopeStack[:stackLen-1]
}

// SetCurrentLine - set lineIdx to current running scope of the module
func (m *Module) SetCurrentLine(line int) {
	if sp := m.GetCurrentScope(); sp != nil {
		sp.SetExecLineIdx(line)
	}
}
