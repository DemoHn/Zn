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

