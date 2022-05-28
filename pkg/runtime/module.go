package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/value"
)

type funcExecutor = func(*Context, []Value) (Value, error)

type Module struct {
	// name = nil when module is anonymous
	name string
	// exported symbols on root scope of this module
	symbols map[string]SymbolInfo
}

// NewModule - create module with specific name
func NewModule(name string) *Module {
	return &Module{
		name: name,
		symbols: map[string]SymbolInfo{},
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

func (m *Module) GetSymbols() map[string]SymbolInfo {
	return m.symbols
}

// RegisterFunction - add function into module
func (m *Module) RegisterFunction(name string, fn funcExecutor) {
	m.symbols[name] = SymbolInfo{
		value:   value.NewFunction(name, fn),
		isConst: true,
	}
}

// RegisterClass - add class info module
func (m *Module) RegisterClass(name string, ref value.ClassRef) {
	m.symbols[name] = SymbolInfo{
		value:   &ref,
		isConst: true,
	}
}
