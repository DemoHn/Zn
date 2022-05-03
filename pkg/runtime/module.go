package runtime

import "github.com/DemoHn/Zn/pkg/syntax"

type Module struct {
	// name = nil when module is anonymous
	name *string
	// exported symbols on root scope of this module
	symbols map[string]SymbolInfo
	// import modules
	imports []Module
	// rootScope of this module
	rootScope *Scope
	// lexer
	lexer *syntax.Lexer
}

// NewModule - create module with specific name
func NewModule(name string, l *syntax.Lexer) *Module {
	return &Module{
		name: &name,
		symbols: map[string]SymbolInfo{},
		imports: []Module{},
		rootScope: NewScope(),
		lexer: l,
	}
}

// NewAnonymousModule - create module but no module name
func NewAnonymousModule(l *syntax.Lexer) *Module {
	return &Module{
		name: nil,
		symbols: map[string]SymbolInfo{},
		imports: []Module{},
		rootScope: NewScope(),
		lexer: l,
	}
}



