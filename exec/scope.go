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
}

//// implementations

// ScopeBase - basic scope structure
type ScopeBase struct {
	root      *RootScope
	parent    Scope
	symbolMap map[string]SymbolInfo
}

// GetRoot - get root scope
func (sb *ScopeBase) GetRoot() *RootScope {
	return sb.root
}

// GetParent - get parent scope
func (sb *ScopeBase) GetParent() Scope {
	return sb.parent
}

// GetSymbol -
func (sb *ScopeBase) GetSymbol(name string) (SymbolInfo, bool) {
	sym, ok := sb.symbolMap[name]
	return sym, ok
}

// SetSymbol -
func (sb *ScopeBase) SetSymbol(name string, value ZnValue, isConstant bool) {
	sb.symbolMap[name] = SymbolInfo{
		value, isConstant,
	}
}

// SymbolInfo - symbol info
type SymbolInfo struct {
	Value      ZnValue
	IsConstant bool // if isConstant = true, the value of this symbol is prohibited from any modification.
}

// RootScope - as named, this is the root scope for execution one program.
// usually it contains all active variables, scopes, etc
type RootScope struct {
	*ScopeBase
	//// lexical scope
	// file - current execution file directory
	file string
	// currentLine - current exeuction line
	currentLine int
	// lineStack - lexical info of (parsed) current file
	lineStack *lex.LineStack
	// lastValue - get last valid value even if there's no return statement
	lastValue ZnValue
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
	rs := &RootScope{
		lastValue:   NewZnNull(),
		classRefMap: map[string]*syntax.ClassDeclareStmt{},
	}
	rs.ScopeBase = &ScopeBase{
		root:      rs,
		parent:    nil,
		symbolMap: map[string]SymbolInfo{},
	}

	return rs
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
	*ScopeBase
	returnValue ZnValue
}

// NewFuncScope -
func NewFuncScope(parent Scope) *FuncScope {
	return &FuncScope{
		returnValue: NewZnNull(),
		ScopeBase: &ScopeBase{
			root:      parent.GetRoot(),
			parent:    parent,
			symbolMap: map[string]SymbolInfo{},
		},
	}
}

// SetCurrentLine - set current execution line
func (fs *FuncScope) SetCurrentLine(line int) {
	fs.root.SetCurrentLine(line)
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
type WhileScope struct {
	*ScopeBase
}

// NewWhileScope -
func NewWhileScope(parent Scope) *WhileScope {
	return &WhileScope{
		ScopeBase: &ScopeBase{
			root:      parent.GetRoot(),
			parent:    parent,
			symbolMap: map[string]SymbolInfo{},
		},
	}
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
	switch name {
	case "结束":
		return NewZnNull(), error.BreakBreakError()
	case "继续":
		return NewZnNull(), error.ContinueBreakError()
	default:
		// for other keywords, return error directly
		return nil, error.NewErrorSLOT("no appropriate method name for while loop to execute")
	}
}

// IterateScope - iterate stmt scope
type IterateScope struct {
	*ScopeBase
	// current iteration: keys & values
	currentIndex ZnValue
	currentValue ZnValue
}

// NewIterateScope -
func NewIterateScope(parent Scope) *IterateScope {
	return &IterateScope{
		ScopeBase: &ScopeBase{
			root:      parent.GetRoot(),
			parent:    parent,
			symbolMap: map[string]SymbolInfo{},
		},
	}
}

func (its *IterateScope) setCurrentKV(index ZnValue, value ZnValue) {
	its.currentIndex = index
	its.currentValue = value
}

// get props: 此之值，此之索引
func (its *IterateScope) getSpecialProps(name string) (ZnValue, *error.Error) {
	switch name {
	case "值":
		return its.currentValue, nil
	case "索引":
		return its.currentIndex, nil
	default:
		return nil, error.NewErrorSLOT("no appropriate prop name to get")
	}
}

func (its *IterateScope) execSpecialMethods(name string, params []ZnValue) (ZnValue, *error.Error) {
	switch name {
	case "结束":
		return NewZnNull(), error.BreakBreakError()
	case "继续":
		return NewZnNull(), error.ContinueBreakError()
	default:
		// for other keywords, return error directly
		return nil, error.NewErrorSLOT("no appropriate method name for while loop to execute")
	}
}

// ObjectScope -
type ObjectScope struct {
	*ScopeBase
	rootObject ZnValue
}

// NewObjectScope -
func NewObjectScope(parent Scope, rootObject ZnValue) *ObjectScope {
	return &ObjectScope{
		ScopeBase: &ScopeBase{
			root:      parent.GetRoot(),
			parent:    parent,
			symbolMap: map[string]SymbolInfo{},
		},
		rootObject: rootObject,
	}
}

//// helpers

// findObjectScope
// returns: (found bool, objScope *ObjectScope)
func findObjectScope(scope Scope) (bool, *ObjectScope) {
	var sp = scope
	var objectScope *ObjectScope
	// find valid scope
	for sp != nil {
		if osp, ok := sp.(*ObjectScope); ok {
			objectScope = osp
			break
		}
		sp = sp.GetParent()
	}

	if sp == nil {
		return false, nil
	}
	return true, objectScope
}
