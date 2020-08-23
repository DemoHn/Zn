package exec

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
	"github.com/DemoHn/Zn/syntax"
)

// Scope -
type Scope interface {
	// GetParent - get parent scope
	GetParent() Scope
	// GetRoot - get its root scope
	GetRoot() *RootScope
	// GetSymbol - get symbol from internal symbol map
	GetSymbol(name string) (SymbolInfo, bool)
	// SetSymbol - set symbol from internal symbol map
	SetSymbol(name string, value ZnValue, isConstant bool)
	// HasSymbol - if the scope has stand-alone valueMap
	HasSymbol() bool
}

const (
	sTypeRoot    = "scopeRoot"
	sTypeFunc    = "scopeFunc"
	sTypeWhile   = "scopeWhile"
	sTypeIterate = "scopeIterate"
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
	// classRefMap - class definition template (reference)
	// this item only exists on RootScope since class defition block IS allowed
	// ONLY in root block
	classRefMap map[string]*syntax.ClassDeclareStmt
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
		lastValue:   NewZnNull(),
		symbolMap:   map[string]SymbolInfo{},
		classRefMap: map[string]*syntax.ClassDeclareStmt{},
	}
}

// Init - init rootScope using new Lexer
func (rs *RootScope) Init(l *lex.Lexer) {
	rs.file = l.InputStream.GetFile()
	rs.currentLine = 0
	rs.lineStack = l.LineStack
	rs.lastValue = NewZnNull()
}

// SetCurrentLine -
func (rs *RootScope) SetCurrentLine(line int) {
	rs.currentLine = line
}

// GetParent -
func (rs *RootScope) GetParent() Scope {
	return nil
}

// GetRoot -
func (rs *RootScope) GetRoot() *RootScope {
	return rs
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

// HasSymbol -
func (rs *RootScope) HasSymbol() bool {
	return true
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
	returnValue ZnValue
	root        *RootScope
	parent      Scope
	symbolMap   map[string]SymbolInfo
}

// GetParent -
func (fs *FuncScope) GetParent() Scope {
	return fs.parent
}

// GetRoot -
func (fs *FuncScope) GetRoot() *RootScope {
	return fs.root
}

// SetCurrentLine - set current execution line
func (fs *FuncScope) SetCurrentLine(line int) {
	fs.root.SetCurrentLine(line)
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

// HasSymbol -
func (fs *FuncScope) HasSymbol() bool {
	return true
}

// GetReturnValue -
func (fs *FuncScope) GetReturnValue() ZnValue {
	return fs.returnValue
}

// SetReturnValue -
func (fs *FuncScope) SetReturnValue(value ZnValue) {
	fs.returnValue = value
}

// WhileScope - a scope within *while* statement
// NOTICE: there's no standalone symbolMap inside this scope,
// instead, use it's parent for get/set symbols
type WhileScope struct {
	root   *RootScope
	parent Scope
}

// HasSymbol - while scope has NO standalone symbol system
func (ws *WhileScope) HasSymbol() bool {
	return false
}

// SetSymbol - set symbol
func (ws *WhileScope) SetSymbol(name string, value ZnValue, isConstant bool) {
	return
}

// GetSymbol - get symbol
func (ws *WhileScope) GetSymbol(name string) (SymbolInfo, bool) {
	return SymbolInfo{}, false
}

// GetParent - get parent
func (ws *WhileScope) GetParent() Scope {
	return ws.parent
}

// GetRoot -
func (ws *WhileScope) GetRoot() *RootScope {
	return ws.root
}

// execSpecialMethods - a weird way to execute internal "scope"-bound functions
// example:
// 每当 Cond：
//     此之（结束）
//     此之（继续）
//
// where `此之（结束）` means under this whileScope, execute the (结束) method to break the loop (same as "break" keyword)
// where `此之（继续）` means under this whileScope, execute the (继续) method to continue the loop (same as "continue" keyword)
func (ws *WhileScope) execSpecialMethods(name string, params []ZnValue) (ZnValue, *error.Error) {
	if name == "结束" {
		return NewZnNull(), error.BreakBreakError()
	}
	if name == "继续" {
		return NewZnNull(), error.ContinueBreakError()
	}
	// for other keywords, return error directly
	return nil, error.NewErrorSLOT("no appropriate method name for while loop to execute")
}

// createScope - create new (nested) scope from current scope
// fails if return scope is nil
func createScope(ctx *Context, scope Scope, sType string) Scope {
	switch sType {
	case sTypeFunc:
		return &FuncScope{
			returnValue: NewZnNull(),
			root:        scope.GetRoot(),
			parent:      scope,
			symbolMap:   map[string]SymbolInfo{},
		}
	case sTypeWhile:
		return &WhileScope{
			root:   scope.GetRoot(),
			parent: scope,
		}
	case sTypeIterate:
		return &IterateScope{
			root:      scope.GetRoot(),
			parent:    scope,
			symbolMap: map[string]SymbolInfo{},
		}
	}

	return nil
}

// IterateScope - iterate stmt scope
type IterateScope struct {
	root      *RootScope
	parent    Scope
	symbolMap map[string]SymbolInfo
	// current iteration: keys & values
	currentIndex ZnValue
	currentValue ZnValue
}

// GetParent -
func (its *IterateScope) GetParent() Scope {
	return its.parent
}

// GetRoot -
func (its *IterateScope) GetRoot() *RootScope {
	return its.root
}

// GetSymbol -
func (its *IterateScope) GetSymbol(name string) (SymbolInfo, bool) {
	sym, ok := its.symbolMap[name]
	return sym, ok
}

// SetSymbol -
func (its *IterateScope) SetSymbol(name string, value ZnValue, isConstant bool) {
	its.symbolMap[name] = SymbolInfo{
		value, isConstant,
	}
}

// HasSymbol -
func (its *IterateScope) HasSymbol() bool {
	return true
}

func (its *IterateScope) setCurrentKV(index ZnValue, value ZnValue) {
	its.currentIndex = index
	its.currentValue = value
}

// get props: 此之值，此之索引
func (its *IterateScope) getSpecialProps(name string) ZnValue {
	if name == "值" {
		return its.currentValue
	}
	if name == "索引" {
		return its.currentIndex
	}
	panic(error.NewErrorSLOT("no appropriate prop name to get"))
}

func (its *IterateScope) execSpecialMethods(name string, params []ZnValue) (ZnValue, *error.Error) {
	if name == "结束" {
		return NewZnNull(), error.BreakBreakError()
	}
	if name == "继续" {
		return NewZnNull(), error.ContinueBreakError()
	}
	// for other keywords, return error directly
	return nil, error.NewErrorSLOT("no appropriate method name for while loop to execute")
}
