package syntax

//////// Node types

//// interfaces

// Node -
type Node interface{}

// Statement -
type Statement interface {
	Node
	GetCurrentLine() int
	SetCurrentLine(line int)
}

// StmtBase - Statement Base
type StmtBase struct {
	currentLine int
}

func (b *StmtBase) stmtNode() {}

// GetCurrentLine -
func (b *StmtBase) GetCurrentLine() int { return b.currentLine }

// SetCurrentLine -
func (b *StmtBase) SetCurrentLine(line int) {
	b.currentLine = line
}

// Expression - a special type of statement - that yields value after execution
type Expression interface {
	Statement
	exprNode()
}

// ExprBase -
type ExprBase struct {
	currentLine int
}

// GetCurrentLine -
func (e *ExprBase) GetCurrentLine() int { return e.currentLine }

// SetCurrentLine -
func (e *ExprBase) SetCurrentLine(line int) { e.currentLine = line }
func (e *ExprBase) stmtNode()               {}
func (e *ExprBase) exprNode()               {}

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
// e.g.  ArrayList  => 【1，2，3，4，5】
//       HashMap    => 【A = 1，B = 2】
type UnionMapList interface {
	Expression
	mapList()
}

//// program (struct)

// Program -
type Program struct {
	StmtBase
	*Lexer  // include lexing info
	Content *BlockStmt
}

// NodeList - a simple struct that packs several nodes, with custom tag to indicate its feature.
type NodeList struct {
	Tag      int
	Children []Node
}

//// Statements (struct)

// VarDeclareStmt - declare variables as init its values
type VarDeclareStmt struct {
	StmtBase
	AssignPair []VDAssignPair
}

// EmptyStmt - contains nothing - generated by a semicolon token
type EmptyStmt struct {
	StmtBase
}

// VDAssignPair - helper type
type VDAssignPair struct {
	Type       vdAssignPairTypeE
	Variables  []*ID
	AssignExpr Expression
	RefMark    bool
}

type vdAssignPairTypeE uint8

// declare VD Assign type
const (
	VDTypeAssign      = 1 // 设为
	VDTypeAssignConst = 3 // 恒为
)

// BranchStmt - conditional (if-else) statement
type BranchStmt struct {
	StmtBase
	// if
	IfTrueExpr  Expression
	IfTrueBlock *BlockStmt
	// else
	IfFalseBlock *BlockStmt
	// else if
	OtherExprs  []Expression
	OtherBlocks []*BlockStmt
	// else-branch exists or not
	HasElse bool
}

// WhileLoopStmt - (while) statement
type WhileLoopStmt struct {
	StmtBase
	// while this expression satisfies (return TRUE), the following block executes.
	TrueExpr Expression
	// execution block
	LoopBlock *BlockStmt
}

// IterateStmt - 以 ... 遍历 ... statement
type IterateStmt struct {
	StmtBase
	IterateExpr  Expression
	IndexNames   []*ID
	IterateBlock *BlockStmt
}

// ImportStmt - 导入《 ... 》 statement
type ImportStmt struct {
	StmtBase
	ImportLibType uint8
	ImportName    *String
	ImportItems   []*ID
}

type BreakStmt struct {
	StmtBase
}

type ContinueStmt struct {
	StmtBase
}

// declare import libType enum
const (
	// LibTypeStd - standard lib
	LibTypeStd    uint8 = 1
	LibTypeCustom uint8 = 2
)

// BlockStmt -
type BlockStmt struct {
	StmtBase
	Children []Statement
}

// FunctionDeclareStmt - function declaration
type FunctionDeclareStmt struct {
	StmtBase
	FuncName    *ID
	ParamList   []*ParamItem
	ExecBlock   *BlockStmt
	CatchBlocks []*CatchBlockPair
}

type CatchBlockPair struct {
	ExceptionClass *ID
	ExecBlock      *BlockStmt
}

// ConstructorDeclareStmt - (如何新建) constructor is a special function
// to help create a new Object with some pre-defined logic
type ConstructorDeclareStmt struct {
	StmtBase
	DelcareClassName *ID
	ParamList        []*ParamItem
	ExecBlock        *BlockStmt
	CatchBlocks      []*CatchBlockPair
}

// GetterDeclareStmt - getter declaration (何为)
type GetterDeclareStmt struct {
	StmtBase
	GetterName *ID
	ExecBlock  *BlockStmt
	// NOTE: Intentionally NO CatchBlocks!!
	// since getter is a property access, it should not have any exception
}

// FunctionReturnStmt - return (expr)
type FunctionReturnStmt struct {
	StmtBase
	ReturnExpr Expression
}

// 输入XX、XX、XX...
type VarInputStmt struct {
	StmtBase
	IDList []*ID
}

// ClassDeclareStmt - class definition (定义XX：)
type ClassDeclareStmt struct {
	StmtBase
	ClassName *ID
	// 其XX为XX
	PropertyList []*PropertyDeclareStmt
	// 如何XXX？
	MethodList []*FunctionDeclareStmt
	// 何为XXX？
	GetterList []*GetterDeclareStmt
}

// PropertyDeclareStmt - valid inside Class
type PropertyDeclareStmt struct {
	StmtBase
	PropertyID *ID
	InitValue  Expression
}

// ThrowExceptionStmt - throw error
type ThrowExceptionStmt struct {
	StmtBase
	ExceptionClass *ID
	Params         []Expression
}

// ParamItem - parameter item
type ParamItem struct {
	ID      *ID
	RefMark bool
}

//// Expressions (struct)

// PrimeExpr - primitive expression
type PrimeExpr struct {
	ExprBase
	Literal string
}

// ID - Identifier type
type ID struct {
	PrimeExpr
}

// String -
type String struct {
	PrimeExpr
}

// ArrayExpr - array expression
type ArrayExpr struct {
	ExprBase
	Items []Expression
}

// HashMapExpr - hashMap expression
type HashMapExpr struct {
	ExprBase
	KVPair []HashMapKeyValuePair
}

// HashMapKeyValuePair -
type HashMapKeyValuePair struct {
	Key   Expression
	Value Expression
}

// VarAssignExpr - variable assignment statement
// assign <TargetExpr> from <AssignExpr>
//
// i.e.
// TargetExpr := AssignExpr
type VarAssignExpr struct {
	ExprBase
	TargetVar  Assignable
	RefMark    bool
	AssignExpr Expression
}

type ObjNewExpr struct {
	ExprBase
	ClassName *ID
	Params    []Expression
}

// FuncCallExpr - function call
type FuncCallExpr struct {
	ExprBase
	FuncName    *ID
	Params      []Expression
	YieldResult *ID
}

// MemberExpr - declare a member (dot) relation
// Example:
//    【1，2】 之 和
//    数组#2
type MemberExpr struct {
	ExprBase
	Root       Expression // root Expr (maybe null when rootType is 2 or 3)
	RootType   uint8      // 1 - RootTypeExpr, 2 - RootTypeProp (aka. 其)
	MemberType uint8      // 1 - memberID, 3 - memberIndex
	// union: memberItem
	MemberID    *ID
	MemberIndex Expression
}

// MemberMethodExpr - declare a member method
// Example:
// 以 X （执行方法：YYY、ZZZ）
type MemberMethodExpr struct {
	ExprBase
	Root        Expression
	MethodChain []*FuncCallExpr
	YieldResult *ID
}

// declare root types
const (
	RootTypeExpr uint8 = 1 // T 之 X
	RootTypeProp uint8 = 2 // 其 X
)

// declare member types
const (
	MemberID    uint8 = 1 // T 之 prop
	MemberIndex uint8 = 2 // T # num
)

// declare some logic types
const (
	LogicOR  uint8 = 1 // 或
	LogicAND uint8 = 2 // 且
	// LogicEQ ~ LogicLTE only valid when both left & right values are Number
	LogicEQ  uint8 = 4 // 等于
	LogicNEQ uint8 = 5 // 不等于
	LogicGT  uint8 = 6 // 大于
	LogicGTE uint8 = 7 // 不小于
	LogicLT  uint8 = 8 // 小于
	LogicLTE uint8 = 9 // 不大于
	// LogicXEQ & LogicXNEQ are similar to LogicEQ, but has wider usage
	// Number, String, Array, HashMap... even Objects are valid!
	LogicXEQ  uint8 = 10 // 为
	LogicXNEQ uint8 = 11 // 不为

	// arith types
	ArithAdd    uint8 = 12 // +
	ArithSub    uint8 = 13 // -
	ArithMul    uint8 = 14 // *
	ArithDiv    uint8 = 15 // /
	ArithIntDiv uint8 = 16 // |
	ArithModulo uint8 = 17 // %
)

// LogicExpr - logical expression return TRUE (真) or FALSE (假) only
type LogicExpr struct {
	ExprBase
	Type      uint8
	LeftExpr  Expression
	RightExpr Expression
}

// ArithExpr - arithmetic expression like (+ - * /)
type ArithExpr struct {
	ExprBase
	Type      uint8
	LeftExpr  Expression
	RightExpr Expression
}

// EqMarkConfig - configure the representation of equal mark
type EqMarkConfig struct {
	AsMapSign   bool // '=' represents for '对应', e.g. 【A = 1】
	AsVarAssign bool // '=' represents for '设为'
}

// implement expression interface

// SetLiteral - set literal for primeExpr
func (pe *PrimeExpr) SetLiteral(literal []rune) { pe.Literal = string(literal) }

// GetLiteral -
func (pe *PrimeExpr) GetLiteral() string { return pe.Literal }
func (ar *ArrayExpr) mapList()           {} // belongs to unionMapList
func (ar *HashMapExpr) mapList()         {} // belongs to unionMapList
func (id *ID) assignable()               {}
func (me *MemberExpr) assignable()       {}
