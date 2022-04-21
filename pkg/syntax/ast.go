package syntax

import "github.com/DemoHn/Zn/lex"

// ast.go defines general AST structure for every language variation (zh, jp, ... etc.)

// AST - abstract syntax tree
type AST struct {
}

// Node - general Node type
type Node interface{
	// record line number of start token
	GetStartLine() int
	SetStartLine(tk Token)
}

// Statement type
type Statement struct {
	StartLine int
}

// GetStartLine -
func (b *Statement) GetStartLine() int { return b.StartLine }

// SetStartLine -
func (b *Statement) SetStartLine(tk Token) {
	//b.StartLine = tk.
}

// Expression - a special type of statement - that yields value after execution
type Expression struct {
	StartLine int
}

// GetStartLine -
func (e *Expression) GetStartLine() int { return e.StartLine }

// SetStartLine -
func (e *Expression) SetStartLine(tk *lex.Token) { e.StartLine = tk.Range.StartLine }


// Assignable - a special type of expression - that is, it could be assigned as
// a value.
//
// Example:
// ID 为 Expr   --> (ID) is an assignable node
// Array # index 为 Expr    --> (Array # index) is an assignable node
type Assignable interface {
	Expression
	assignable()
}

// UnionMapList - HashMap or ArrayList, since they shares similar grammer.
// e.g.  ArrayList  => 【1、2、3、4、5】
//       HashMap    => 【A = 1，B = 2】
type UnionMapList interface {
	Expression
	mapList()
}

//// program (struct)

// Program -
type Program struct {
	Statement
	*lex.Lexer // include lexing info
	Content    *BlockStmt
}

// NodeList - a simple struct that packs several nodes, with custom tag to indicate its feature.
type NodeList struct {
	Tag      int
	Children []Node
}

//// Statements (struct)

// VarDeclareStmt - declare variables as init its values
type VarDeclareStmt struct {
	Statement
	AssignPair []VDAssignPair
}

// EmptyStmt - contains nothing - generated by a semicolon token
type EmptyStmt struct {
	Statement
}

// VDAssignPair - helper type
type VDAssignPair struct {
	Type       vdAssignPairTypeE
	Variables  []*ID
	AssignExpr Expression
	RefMark    bool
	ObjClass   *ID          // 成为 XX： 1，2，3 ... valid only when Type = 2 (VDTypeObjNew)
	ObjParams  []Expression // 成为 XX：P1，P2，P3，... valid only when Type = 2 (VDTypeObjNew)
}

type vdAssignPairTypeE uint8

// declare VD Assign type
const (
	VDTypeAssign      = 1 // 为
	VDTypeObjNew      = 2 // 成为
	VDTypeAssignConst = 3 // 恒为
)

// BranchStmt - conditional (if-else) statement
type BranchStmt struct {
	Statement
	// 'if' block
	IfTrueExpr  Expression
	IfTrueBlock *BlockStmt
	// 'else' block
	IfFalseBlock *BlockStmt
	// 'else if' block
	OtherExprs  []Expression
	OtherBlocks []*BlockStmt
	// if 'else' block exists
	HasElse bool
}

// WhileLoopStmt - loop (while) statement
type WhileLoopStmt struct {
	Statement
	// while this expression satisfies (return TRUE), loop block will be executed.
	TrueExpr Expression
	// execution block
	LoopBlock *BlockStmt
}

// IterateStmt - 以 ... 遍历 ... statement
type IterateStmt struct {
	Statement
	IterateExpr  Expression
	IndexNames   []*ID
	IterateBlock *BlockStmt
}

// declare import libType enum
const (
	// LibTypeStd - standard lib
	LibTypeStd uint8 = 1
)

// ImportStmt - 导入《 ... 》 statement
type ImportStmt struct {
	Statement
	ImportLibType uint8
	ImportName    *String
	ImportItems   []*ID
}

// BlockStmt -
type BlockStmt struct {
	Statement
	Children []Statement
}

// FunctionDeclareStmt - function declaration
type FunctionDeclareStmt struct {
	Statement
	FuncName  *ID
	ParamList []*ParamItem
	ExecBlock *BlockStmt
}

// GetterDeclareStmt - getter declaration (何为)
type GetterDeclareStmt struct {
	Statement
	GetterName *ID
	ExecBlock  *BlockStmt
}

// FunctionReturnStmt - return (expr)
type FunctionReturnStmt struct {
	Statement
	ReturnExpr Expression
}

// ClassDeclareStmt - class definition (定义XX：)
type ClassDeclareStmt struct {
	Statement
	ClassName *ID
	// 其XX为XX
	PropertyList []*PropertyDeclareStmt
	// 是为XX，YY，ZZ
	ConstructorIDList []*ParamItem
	// 如何XXX？
	MethodList []*FunctionDeclareStmt
	// 何为XXX？
	GetterList []*GetterDeclareStmt
}

// PropertyDeclareStmt - valid inside Class
type PropertyDeclareStmt struct {
	Statement
	PropertyID *ID
	InitValue  Expression
}

// ParamItem - parameter item
type ParamItem struct {
	ID      *ID
	RefMark bool
}

//// Expressions (struct)

// PrimeExpr - primitive expression
type PrimeExpr struct {
	Expression
	Literal string
}

// ID - Identifier type
type ID struct {
	PrimeExpr
}

// Number -
type Number struct {
	PrimeExpr
}

// String -
type String struct {
	PrimeExpr
}

// ArrayExpr - array expression
type ArrayExpr struct {
	Expression
	Items []Expression
}

// HashMapExpr - hashMap expression
type HashMapExpr struct {
	Expression
	KVPair []hashMapKeyValuePair
}

// hashMapKeyValuePair -
type hashMapKeyValuePair struct {
	Key   Expression
	Value Expression
}

// VarAssignExpr - variable assignment statement
// assign <TargetExpr> from <AssignExpr>
//
// i.e.
// TargetExpr := AssignExpr
type VarAssignExpr struct {
	Expression
	TargetVar  Assignable
	RefMark    bool
	AssignExpr Expression
}

// FuncCallExpr - function call
type FuncCallExpr struct {
	Expression
	FuncName    *ID
	Params      []Expression
	YieldResult *ID
}

// MemberExpr - declare a member (dot) relation
// Example:
//    【1，2】 之 和
//    数组#2
type MemberExpr struct {
	Expression
	Root       Expression  // root Expr (maybe null when rootType is 2 or 3)
	RootType   rootTypeE   // 1 - RootTypeExpr, 2 - RootTypeProp (aka. 其)
	MemberType memberTypeE // 1 - memberID, 3 - memberIndex
	// union: memberItem
	MemberID    *ID
	MemberIndex Expression
}

// MemberMethodExpr - declare a member method
// Example:
// 以 X （执行方法：YYY、ZZZ）
type MemberMethodExpr struct {
	Expression
	Root        Expression
	MethodChain []*FuncCallExpr
	YieldResult *ID
}

// rootTypeE - root type enumeration
type rootTypeE uint8

// declare root types
const (
	RootTypeExpr rootTypeE = 1 // T 之 X
	RootTypeProp rootTypeE = 2 // 其 X
)

// memberTypeE - member type enumeration
type memberTypeE uint8

// declare member types
const (
	MemberID    memberTypeE = 1 // T 之 prop
	MemberIndex memberTypeE = 2 // T # num
)

// LogicTypeE - enumerates several logic type (OR, AND, EQ, etc)
type LogicTypeE uint8

// declare some logic types
const (
	LogicOR  LogicTypeE = 1 // 或
	LogicAND LogicTypeE = 2 // 且
	LogicEQ  LogicTypeE = 4 // 等于
	LogicNEQ LogicTypeE = 5 // 不等于
	LogicGT  LogicTypeE = 6 // 大于
	LogicGTE LogicTypeE = 7 // 不小于
	LogicLT  LogicTypeE = 8 // 小于
	LogicLTE LogicTypeE = 9 // 不大于
)

// LogicExpr - logical expression return TRUE (真) or FALSE (假) only
type LogicExpr struct {
	Expression
	Type      LogicTypeE
	LeftExpr  Expression
	RightExpr Expression
}

// EqMarkConfig - configure the representation of equal mark
type EqMarkConfig struct {
	AsMapSign   bool // '=' represents for '是', e.g. 【A = 1】
	AsVarAssign bool // '=' represents for '为'
	AsEqual     bool // '=' represents for '等于'
}
