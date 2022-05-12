package runtime

import "github.com/DemoHn/Zn/pkg/syntax"

type Module struct {
	// name = nil when module is anonymous
	name *string
	// exported symbols on root scope of this module
	symbols map[string]SymbolInfo
	// import modules
	imports []Module
	// lexer
	lexer *syntax.Lexer
}

// NewModule - create module with specific name
func NewModule(name string, l *syntax.Lexer) *Module {
	return &Module{
		name: &name,
		symbols: map[string]SymbolInfo{},
		imports: []Module{},
		lexer: l,
	}
}




