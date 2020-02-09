package exec

import (
	"fmt"
	"strings"

	"github.com/DemoHn/Zn/error"
)

// SymbolInfo - symbol info
type SymbolInfo struct {
	nestLevel int
	value     ZnValue
}

// SymbolTable - a global hash-table manages all symbols
// including variable, func name, etcs.
// (support nesting)
type SymbolTable struct {
	// key   -> symbol name
	// value -> symbol info
	symbolMap map[string][]SymbolInfo
	nestLevel int
}

// NewSymbolTable -
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbolMap: map[string][]SymbolInfo{},
		nestLevel: 0,
	}
}

// Bind - add value to symbol table
func (st *SymbolTable) Bind(id string, obj ZnValue) *error.Error {
	newInfo := SymbolInfo{
		nestLevel: st.nestLevel,
		value:     obj,
	}

	symArr, ok := st.symbolMap[id]
	if !ok {
		// init symbolInfo array
		st.symbolMap[id] = []SymbolInfo{newInfo}
		return nil
	}

	// check if there's variable re-declaration
	if len(symArr) > 0 && symArr[0].nestLevel == st.nestLevel {
		return error.NewErrorSLOT("variable re-declaration")
	}

	// prepend data
	st.symbolMap[id] = append([]SymbolInfo{newInfo}, st.symbolMap[id]...)
	return nil
}

// Lookup - find the corresponded value from ID,
// if nothing found, return error
func (st *SymbolTable) Lookup(id string) (ZnValue, *error.Error) {
	symArr, ok := st.symbolMap[id]
	if !ok {
		return nil, error.NewErrorSLOT("no valid variable found")
	}

	// find the nearest level of value
	if symArr == nil || len(symArr) == 0 {
		return nil, error.NewErrorSLOT("no valid variable found")
	}
	return symArr[0].value, nil
}

// EnterScope - enter a nested scope
func (st *SymbolTable) EnterScope() {
	st.nestLevel++
}

// ExitScope - exit from a nested scope
func (st *SymbolTable) ExitScope() {
	// find all variable
	for idx, symArr := range st.symbolMap {
		if symArr != nil && len(symArr) > 0 {
			if symArr[0].nestLevel == st.nestLevel {
				// remove first item since it's outdated
				st.symbolMap[idx] = symArr[1:]
			}
		}
	}

	st.nestLevel--
}

// SetData -
func (st *SymbolTable) SetData(id string, obj ZnValue) *error.Error {
	symArr, ok := st.symbolMap[id]
	if !ok {
		return error.NewErrorSLOT("variable not defined!")
	}

	if symArr != nil && len(symArr) > 0 {
		symArr[0].value = obj
		return nil
	}

	return error.NewErrorSLOT("variable not defined!")
}

func (st *SymbolTable) printSymbols() string {
	strs := []string{}
	for k, symArr := range st.symbolMap {
		if symArr != nil {
			for _, symItem := range symArr {
				symStr := "ε"
				if symItem.value != nil {
					symStr = symItem.value.String()
				}
				strs = append(strs, fmt.Sprintf("‹%s, %d› => %s", k, symItem.nestLevel, symStr))
			}
		}

	}

	return strings.Join(strs, "\n")
}
