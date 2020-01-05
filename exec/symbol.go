package exec

// SymbolTable - a global symbol table stores all variables
type SymbolTable struct {
	Symbols map[string]ZnObject
}

// Insert - insert one symbol into symbolTable
// if ID has exists, return false;
// if not, return true;
func (st *SymbolTable) Insert(id string, obj ZnObject) bool {
	if _, ok := st.Symbols[id]; !ok {
		st.Symbols[id] = obj
		return true
	}
	return false
}

// Lookup - find the value from symbol table
func (st *SymbolTable) Lookup(id string) (ZnObject, bool) {
	if val, ok := st.Symbols[id]; ok {
		return val, true
	}
	return nil, false
}

// SetData - set symbol table data
func (st *SymbolTable) SetData(id string, obj ZnObject) bool {
	if _, ok := st.Symbols[id]; ok {
		st.Symbols[id] = obj
		return true
	}
	return false
}
