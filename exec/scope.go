package exec

import (
	"github.com/DemoHn/Zn/lex"
)

// Scope - tmp Scope solution TODO: will move in the future!
type Scope interface {
	// create new (nested) scope from current scope
	// fails if return scope is nil
	NewScope(ctx *Context, sType string) Scope
	// set current execution line
	SetCurrentLine(line int)
	// GetParent - get parent scope
	GetParent() Scope
	// GetSymbol - get symbol from internal symbol map
	GetSymbol(name string) (SymbolInfo, bool)
	// SetSymbol - set symbol from internal symbol map
	SetSymbol(name string, value ZnValue, isConstant bool)
}

const (
	sTypeRoot = "scopeRoot"
	sTypeFunc = "scopeFunc"
)

//// implementations

// SymbolInfo - symbol info
type SymbolInfo struct {
	Value      ZnValue
	IsConstant bool // if isConstant = true, the value of this symbol is prohibited from any modification.
}

// RootScope - as named, this is the root scope for execution one program.
// usually it contains all active variables, scopes, etc
type RootScope struct {
	//// lexical scope
	// file - current execution file directory
	file string
	// currentLine - current exeuction line
	currentLine int
	// lineStack - lexical info of (parsed) current file
	lineStack *lex.LineStack
	// lastValue - get last valid value even if there's no return statement
	lastValue ZnValue
	// symbolMap - store variables within this scope
	symbolMap map[string]SymbolInfo
}

// NewRootScope - create a rootScope from existing Lexer that
// derives from a program file, a piece of code, etc.
//
// That implies a program has one and only one RootScope.
//
// NOTE: When a program file "requires" another one, another RootScope is created
// for that "required" program file.
func NewRootScope() *RootScope {
	return &RootScope{
		lastValue: NewZnNull(),
		symbolMap: map[string]SymbolInfo{},
	}
}

// Init - init rootScope using new Lexer
func (rs *RootScope) Init(l *lex.Lexer) {
	rs.file = l.InputStream.Scope
	rs.currentLine = 0
	rs.lineStack = l.LineStack
	rs.lastValue = NewZnNull()
}

// SetCurrentLine -
func (rs *RootScope) SetCurrentLine(line int) {
	rs.currentLine = line
}

// NewScope - create new (nested) scope from current scope
// fails if return scope is nil
func (rs *RootScope) NewScope(ctx *Context, sType string) Scope {
	if sType == sTypeFunc {
		return &FuncScope{
			returnValue: NewZnNull(),
			root:        rs,
			parent:      rs,
			symbolMap:   map[string]SymbolInfo{},
		}
	}
	return nil
}

// GetParent -
func (rs *RootScope) GetParent() Scope {
	return nil
}

// GetSymbol - get symbol
func (rs *RootScope) GetSymbol(name string) (SymbolInfo, bool) {
	sym, ok := rs.symbolMap[name]
	return sym, ok
}

// SetSymbol - set symbol
func (rs *RootScope) SetSymbol(name string, value ZnValue, isConstant bool) {
	rs.symbolMap[name] = SymbolInfo{
		value, isConstant,
	}
}

// SetLastValue - set last value
func (rs *RootScope) SetLastValue(value ZnValue) {
	rs.lastValue = value
}

// GetLastValue -
func (rs *RootScope) GetLastValue() ZnValue {
	return rs.lastValue
}

// FuncScope - function scope
type FuncScope struct {
	//// returnValue - the final return value of current scope
	returnValue ZnValue
	root        *RootScope
	parent      Scope
	returnFlag  bool
	symbolMap   map[string]SymbolInfo
}

// NewScope - create new (nested) scope from current scope
// fails if return scope is nil
func (fs *FuncScope) NewScope(ctx *Context, sType string) Scope {
	if sType == sTypeFunc {
		return &FuncScope{
			returnValue: NewZnNull(),
			root:        fs.root,
			parent:      fs,
			symbolMap:   map[string]SymbolInfo{},
		}
	}
	return nil
}

// GetParent -
func (fs *FuncScope) GetParent() Scope {
	return fs.parent
}

// SetCurrentLine - set current execution line
func (fs *FuncScope) SetCurrentLine(line int) {
	fs.root.SetCurrentLine(line)
}

// SetReturnValue - set last value
func (fs *FuncScope) SetReturnValue(value ZnValue) {
	fs.returnValue = value
}

// GetReturnValue -
func (fs *FuncScope) GetReturnValue() ZnValue {
	return fs.returnValue
}

// SetReturnFlag - set last value
func (fs *FuncScope) SetReturnFlag(flag bool) {
	fs.returnFlag = flag
}

// GetReturnFlag -
func (fs *FuncScope) GetReturnFlag() bool {
	return fs.returnFlag
}

// GetSymbol - get symbol
func (fs *FuncScope) GetSymbol(name string) (SymbolInfo, bool) {
	sym, ok := fs.symbolMap[name]
	return sym, ok
}

// SetSymbol - set symbol
func (fs *FuncScope) SetSymbol(name string, value ZnValue, isConstant bool) {
	fs.symbolMap[name] = SymbolInfo{
		value, isConstant,
	}
}

//// helpers
