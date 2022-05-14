package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/syntax"
)

type Module struct {
	// name = nil when module is anonymous
	name string
	// exported symbols on root scope of this module
	symbols map[string]SymbolInfo
	// import modules
	// <module name> -> [symbolNames]
	importSymbolMap map[string][]string
	// lexer
	lexer *syntax.Lexer
}

// NewModule - create module with specific name
func NewModule(name string, l *syntax.Lexer) *Module {
	return &Module{
		name: name,
		symbols: map[string]SymbolInfo{},
		importSymbolMap: map[string][]string{},
		lexer: l,
	}
}

// AddSymbol -
func (m *Module) AddSymbol(symbol string, value Value, isConst bool) error {
	if _, ok := m.symbols[symbol]; ok {
		return zerr.NameRedeclared(symbol)
	}
	m.symbols[symbol] = SymbolInfo{
		value:   value,
		isConst: isConst,
	}
	return nil
}

func (m *Module) GetSymbol(symbol string) (Value, error) {
	if v, ok := m.symbols[symbol]; ok {
		return v.value, nil
	}

	return nil, zerr.NameNotDefined(symbol)
}




