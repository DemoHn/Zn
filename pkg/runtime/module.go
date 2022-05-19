package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
)

const (
	ModuleTypeStd uint8 = 1
	ModuleTypeCustom uint8 = 2
)

// ImportSymbol -
type ImportSymbol struct {
	name string
	moduleType uint8
}

type Module struct {
	// name = nil when module is anonymous
	name string
	// exported symbols on root scope of this module
	symbols map[string]SymbolInfo
	// import modules
	// [symbolName] -> <module>
	importSymbolMap map[string]ImportSymbol
}

// NewModule - create module with specific name
func NewModule(name string) *Module {
	return &Module{
		name: name,
		symbols: map[string]SymbolInfo{},
		importSymbolMap: map[string]ImportSymbol{},
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

func (m *Module) GetSymbols() []string {
	var res []string
	for sym := range m.symbols {
		res = append(res, sym)
	}
	return res
}

func (m *Module) AddImportSymbols(moduleType uint8, name string, items []string) error {
	for _, item := range items {
		if _, ok := m.importSymbolMap[item]; ok {
			return zerr.NameRedeclared(item)
		}

		m.importSymbolMap[item] = ImportSymbol{
			name:       name,
			moduleType: moduleType,
		}
	}
	return nil
}


