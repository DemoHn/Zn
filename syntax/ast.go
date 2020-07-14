package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

//////// Node types

//// interfaces

// Node -
type Node interface{}

type consumerFunc func(idx int, nodes []Node) Node

// Statement -
type Statement interface {
	Node
	GetCurrentLine() int
	SetCurrentLine(tk *lex.Token)
}

// StmtBase - Statement Base
type StmtBase struct {
	currentLine int
}

func (b *StmtBase) stmtNode() {}

// GetCurrentLine -
func (b *StmtBase) GetCurrentLine() int { return b.currentLine }

// SetCurrentLine -
func (b *StmtBase) SetCurrentLine(tk *lex.Token) {
	b.currentLine = tk.Range.StartLine
}

// Expression - a speical type of statement - that yields value after execution
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
func (e *ExprBase) SetCurrentLine(tk *lex.Token) { e.currentLine = tk.Range.StartLine }
func (e *ExprBase) stmtNode()                    {}
func (e *ExprBase) exprNode()                    {}

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
//       HashMap    => 【A == 1，B == 2】
type UnionMapList interface {
	Expression
	mapList()
}

//// program (struct)

// Program -
type Program struct {
	StmtBase
	Content *BlockStmt
}

// NodeList - a node that packs several nodes
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
	Variables  []*ID
	AssignExpr Expression
}

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

// BlockStmt -
type BlockStmt struct {
	StmtBase
	Children []Statement
}

// FunctionDeclareStmt - function declaration
type FunctionDeclareStmt struct {
	StmtBase
	FuncName  *ID
	ParamList []*ID
	ExecBlock *BlockStmt
}

// FunctionReturnStmt - return (expr)
type FunctionReturnStmt struct {
	StmtBase
	ReturnExpr Expression
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
	ExprBase
	Items []Expression
}

// HashMapExpr - hashMap expression
type HashMapExpr struct {
	ExprBase
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
	ExprBase
	TargetVar  Assignable
	AssignExpr Expression
}

// FuncCallExpr - function call
type FuncCallExpr struct {
	ExprBase
	FuncName *ID
	Params   []Expression
}

// MemberExpr - declare a member (dot) relation
// Example:
//    此之 代码
//    【1，2】 之 和
//    此之 （结束）
//    100 之 （+1）
type MemberExpr struct {
	ExprBase
	Root         Expression
	IsSelfRoot   bool
	MemberType   memberTypeE // 1 - memberID, 2 - memberMethod, 3 - memberIndex
	MemberID     *ID
	MemberMethod *FuncCallExpr
	MemberIndex  Expression
}

// memberTypeE - member type enumeration
type memberTypeE uint8

// declare member types
const (
	MemberID     memberTypeE = 1 // T 之 prop
	MemberMethod memberTypeE = 2 // T 之 （method）
	MemberIndex  memberTypeE = 3 // T # num
)

// LogicTypeE - enumerates several logic type (OR, AND, EQ, etc)
type LogicTypeE uint8

// declare some logic types
const (
	LogicOR  LogicTypeE = 1  // 或
	LogicAND LogicTypeE = 2  // 且
	LogicIS  LogicTypeE = 3  // 此 ... 为 ... TODO: deprecate IT!
	LogicEQ  LogicTypeE = 4  // 等于
	LogicNEQ LogicTypeE = 5  // 不等于
	LogicGT  LogicTypeE = 6  // 大于
	LogicGTE LogicTypeE = 7  // 不小于
	LogicLT  LogicTypeE = 8  // 小于
	LogicLTE LogicTypeE = 9  // 不大于
	LogicISN LogicTypeE = 10 // 此 ... 不为 ... TODO: deprecate IT!
)

// LogicExpr - logical expression return TRUE (真) or FALSE (假) only
type LogicExpr struct {
	ExprBase
	Type      LogicTypeE
	LeftExpr  Expression
	RightExpr Expression
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

//////// Parse Methods

//// NOTE: the following methods are all using panic -> recover for error management.
//// This is to expect elimilating `err != nil` statements.

// ParseStatement - a program consists of statements
//
// CFG:
// Statement -> VarDeclareStmt
//           -> BranchStmt
//           -> Expr
//           -> ；
func ParseStatement(p *Parser) Statement {
	defer p.resetLineTermFlag()

	var validTypes = []lex.TokenType{
		lex.TypeStmtSep,
		lex.TypeComment,
		lex.TypeDeclareW,
		lex.TypeCondW,
		lex.TypeFuncW,
		lex.TypeReturnW,
		lex.TypeWhileLoopW,
		lex.TypeVarOneW,
		lex.TypeIteratorW,
	}
	match, tk := p.tryConsume(validTypes...)
	if match {
		var s Statement
		switch tk.Type {
		case lex.TypeStmtSep, lex.TypeComment:
			// skip them because it's meaningless for syntax parsing
			s = new(EmptyStmt)
		case lex.TypeDeclareW:
			s = ParseVarDeclareStmt(p)
		case lex.TypeCondW:
			mainIndent := p.getPeekIndent()
			s = ParseBranchStmt(p, mainIndent)
		case lex.TypeFuncW:
			s = ParseFunctionDeclareStmt(p)
		case lex.TypeReturnW:
			s = ParseFunctionReturnStmt(p)
		case lex.TypeWhileLoopW:
			s = ParseWhileLoopStmt(p)
		case lex.TypeVarOneW:
			s = ParseVarOneLeadStmt(p) // parse any statements leads with 「以」
		case lex.TypeIteratorW:
			s = ParseIteratorStmt(p)
		}
		s.SetCurrentLine(tk)
		return s
	}
	// other case, parse expression
	return ParseExpression(p, true)
}

// ParseExpression - parse an expression, see the following CFG for details
//
// CFG:
// Expr  -> AndE Expr'
// Expr' -> 或 AndE Expr'
//       ->
//
// AndE  -> EqE AndE'
// AndE' -> 且 EqE AndE'
//       ->
//
// EqE   -> VaE EqE'
// EqE'  -> 等于 VaE
//       -> 不等于 VaE
//       -> 小于 VaE
//       -> 不小于 VaE
//       -> 大于 VaE
//       -> 不大于 VaE
//       ->
//
// VaE   -> IdxE VaE'
// VaE'  -> 为 IdxE
//       ->
//
// IdxE  -> BsE IdxE'
// IdxE' -> #  Number   IdxE'
// IdxE' -> #  String   IdxE'
//       -> #{  Expr  }  IdxE'
//
// precedences:
//
// # #{}  >  为  >  等于，大于，etc.  >  且  >  或
func ParseExpression(p *Parser, asVarAssign bool) Expression {
	var logicItemParser func(int) Expression
	var logicItemTailParser func(int, Expression) Expression
	// logicKeywords, ordered by precedence asc
	// that means, the very begin logicKeyword ([]lex.TokenType) has lowest precedence
	var logicKeywords = [4][]lex.TokenType{
		{lex.TypeLogicOrW},
		{lex.TypeLogicAndW},
		{
			lex.TypeLogicEqualW,
			lex.TypeLogicNotEqW,
			lex.TypeLogicGtW,
			lex.TypeLogicGteW,
			lex.TypeLogicLtW,
			lex.TypeLogicLteW,
		},
		{lex.TypeLogicYesW, lex.TypeLogicNotW},
	}
	var logicTypeMap = map[lex.TokenType]LogicTypeE{
		lex.TypeLogicOrW:    LogicOR,
		lex.TypeLogicAndW:   LogicAND,
		lex.TypeLogicEqualW: LogicEQ,
		lex.TypeLogicNotEqW: LogicNEQ,
		lex.TypeLogicGtW:    LogicGT,
		lex.TypeLogicGteW:   LogicGTE,
		lex.TypeLogicLtW:    LogicLT,
		lex.TypeLogicLteW:   LogicLTE,
		lex.TypeLogicNotW:   LogicISN, // TODO: merge LogicISN -> logicNEQ
		lex.TypeLogicYesW:   LogicIS,  // TODO: merge LogicIS -> logicEQ
	}
	var logicAllowTails = [4]bool{true, true, false, false}

	//// anynomous function definition
	logicItemParser = func(idx int) Expression {
		if idx >= len(logicKeywords) {
			return ParseMemberExpr(p)
		}
		// #1. match item
		expr1 := logicItemParser(idx + 1)

		return logicItemTailParser(idx, expr1)
	}

	//// anynomous function definition
	logicItemTailParser = func(idx int, leftExpr Expression) Expression {
		var finalExpr Expression
		// #1. consume keyword
		match, tk := p.tryConsume(logicKeywords[idx]...)
		if !match {
			return leftExpr
		}
		// #2. consume Y
		rightExpr := logicItemParser(idx + 1)

		// compose logic expr
		if tk.Type == lex.TypeLogicYesW && asVarAssign {
			// if 为 (LogicYes) is interpreted as varAssign
			// usually for normal expressions (except 如果，每当 expr)
			vid, ok := leftExpr.(Assignable)
			if !ok {
				panic(error.ExprMustTypeID())
			}
			finalExpr = &VarAssignExpr{
				TargetVar:  vid,
				AssignExpr: rightExpr,
			}
		} else {
			finalExpr = &LogicExpr{
				Type:      logicTypeMap[tk.Type],
				LeftExpr:  leftExpr,
				RightExpr: rightExpr,
			}
		}
		// set current line (after finalExpr has been initialized)
		finalExpr.SetCurrentLine(tk)

		// #3. consume X' (X-tail)
		if logicAllowTails[idx] {
			return logicItemTailParser(idx, finalExpr)
		}
		return finalExpr
	}

	return logicItemParser(0)
}

// ParseMemberExpr -
//
// CFG:
//
// MemE  -> 此之 CallE' IdxE'
//       -> BsE IdxE'
//
// IdxE' -> #  Number   IdxE'
//       -> #  String   IdxE'
//       -> #{  Expr  }  IdxE'
//       -> 之  CallE' IdxE'
//       ->
//
// CallE' -> ID
//        -> （ID：E，E，...）
func ParseMemberExpr(p *Parser) Expression {
	// internal functions
	var calleeTailParser func(bool, Expression) *MemberExpr
	var memberTailParser func(Expression) Expression

	// specially parsing items after 之 or 此之
	calleeTailParser = func(hasRoot bool, expr Expression) *MemberExpr {
		var validTypes = []lex.TokenType{
			lex.TypeIdentifier,
			lex.TypeVarQuote,
			lex.TypeFuncQuoteL,
			lex.TypeVarOneW,
		}
		memberExpr := &MemberExpr{
			Root:       nil,
			IsSelfRoot: !hasRoot,
		}
		if hasRoot {
			memberExpr.Root = expr
		}
		match, tk := p.tryConsume(validTypes...)
		if match {
			switch tk.Type {
			case lex.TypeIdentifier, lex.TypeVarQuote:
				id := newID(tk)
				id.SetCurrentLine(tk)
				memberExpr.MemberType = MemberID
				memberExpr.MemberID = id
			case lex.TypeFuncQuoteL:
				e := ParseFuncCallExpr(p)
				e.SetCurrentLine(tk)
				memberExpr.MemberType = MemberMethod
				memberExpr.MemberMethod = e
			case lex.TypeVarOneW:
				e := ParseVarOneLeadExpr(p)
				e.SetCurrentLine(tk)
				memberExpr.MemberType = MemberMethod
				memberExpr.MemberMethod = e
			}

			return memberExpr
		}
		panic(error.InvalidSyntax())
	}

	memberTailParser = func(expr Expression) Expression {
		mExpr := &MemberExpr{}

		match, tk := p.tryConsume(lex.TypeMapHash, lex.TypeMapQHash, lex.TypeObjDotW)
		if !match {
			return expr
		}
		mExpr.SetCurrentLine(tk)
		mExpr.Root = expr
		mExpr.IsSelfRoot = false

		switch tk.Type {
		case lex.TypeMapHash:
			match2, tk2 := p.tryConsume(lex.TypeNumber, lex.TypeString)
			if match2 {
				// set memberType
				mExpr.MemberType = MemberIndex
				switch tk2.Type {
				case lex.TypeNumber:
					mExpr.MemberIndex = newNumber(tk2)
				case lex.TypeString:
					mExpr.MemberIndex = newString(tk2)
				}
				return memberTailParser(mExpr)
			}
			panic(error.InvalidSyntax())
		case lex.TypeMapQHash: // lex.TypeMapQHash
			// #1. parse Expr
			mExpr.MemberType = MemberIndex
			mExpr.MemberIndex = ParseExpression(p, true)

			// #2. parse tail brace
			p.consume(lex.TypeStmtQuoteR)

			return memberTailParser(mExpr)
		case lex.TypeObjDotW:
			newExpr := calleeTailParser(true, expr)
			// replace current memberExpr as newExpr
			return memberTailParser(newExpr)
		}

		panic(error.InvalidSyntax())
	}

	// #1. parse 此之 expr
	match, _ := p.tryConsume(lex.TypeStaticSelfW) // 此之
	if match {
		newExpr := calleeTailParser(false, nil)
		return memberTailParser(newExpr)
	}
	// #1. parse basic expr
	rootExpr := ParseBasicExpr(p)
	return memberTailParser(rootExpr)
}

// ParseBasicExpr - parse general basic expression
//
// CFG:
// BsE   -> { E }
//       -> （ ID ： E，E，...）
//       -> 以 E （ ID ： E，E，...）
//       -> ID
//       -> Number
//       -> String
//       -> ArrayList
func ParseBasicExpr(p *Parser) Expression {
	var validTypes = []lex.TokenType{
		lex.TypeIdentifier,
		lex.TypeVarQuote,
		lex.TypeNumber,
		lex.TypeString,
		lex.TypeArrayQuoteL,
		lex.TypeStmtQuoteL,
		lex.TypeFuncQuoteL,
		lex.TypeLogicNotW,
		lex.TypeVarOneW,
	}

	match, tk := p.tryConsume(validTypes...)
	if match {
		var e Expression
		switch tk.Type {
		case lex.TypeIdentifier, lex.TypeVarQuote:
			e = newID(tk)
		case lex.TypeNumber:
			e = newNumber(tk)
		case lex.TypeString:
			e = newString(tk)
		case lex.TypeArrayQuoteL:
			e = ParseArrayExpr(p)
		case lex.TypeStmtQuoteL:
			e = ParseExpression(p, true)
			p.consume(lex.TypeStmtQuoteR)
		case lex.TypeFuncQuoteL:
			e = ParseFuncCallExpr(p)
		case lex.TypeVarOneW:
			e = ParseVarOneLeadExpr(p)
		}
		e.SetCurrentLine(tk)
		return e
	}
	panic(error.InvalidSyntax())
}

// ParseArrayExpr - yield ArrayExpr node (support both hashMap and arrayList)
// CFG:
// ArrayExpr -> 【 ItemList 】
//           -> 【】
//           -> 【 HashMapList 】
//           -> 【 == 】
// ItemList  -> Expr ExprTail
//           ->
// ExprTail  -> ， Expr ExprTail
//           ->
// HashMapList -> Expr == Expr， Expr2 == Expr2， ...
func ParseArrayExpr(p *Parser) UnionMapList {
	// #0. try to match if empty
	if match, emptyExpr := tryParseEmptyMapList(p); match {
		return emptyExpr
	}

	const (
		tagHashMap     = 11
		subtypeUnknown = 0
		subtypeArray   = 1
		subtypeHashMap = 2
	)
	// #1. consume item list (comma list)
	exprs := parseCommaList(p, func(idx int, nodes []Node) Node {
		expr := ParseExpression(p, true)

		// parse if there's double equals, then cont'd parsing right expr for hashmap
		if match, _ := p.tryConsume(lex.TypeMapData); match {
			exprR := ParseExpression(p, true)

			return &NodeList{
				Tag:      tagHashMap,
				Children: []Node{expr, exprR},
			}
		}
		return expr
	})

	// type cast (because there's no GENERIC TYPE in golang!!!)
	var ar = &ArrayExpr{
		Items: []Expression{},
	}
	var hm = &HashMapExpr{
		KVPair: []hashMapKeyValuePair{},
	}
	var subtype = subtypeUnknown
	for _, expr := range exprs {
		switch v := expr.(type) {
		case Expression:
			if subtype == subtypeUnknown {
				subtype = subtypeArray
			}
			if subtype != subtypeArray {
				panic(error.MixArrayHashMap())
			}
			// add value
			ar.Items = append(ar.Items, v)
		case *NodeList: // tagHashMap
			if subtype == subtypeUnknown {
				subtype = subtypeHashMap
			}
			if subtype != subtypeHashMap {
				panic(error.MixArrayHashMap())
			}
			n0, _ := v.Children[0].(Expression)
			n1, _ := v.Children[1].(Expression)
			hm.KVPair = append(hm.KVPair, hashMapKeyValuePair{
				Key:   n0,
				Value: n1,
			})
		}
	}

	// #2. consume right brancket
	p.consume(lex.TypeArrayQuoteR)

	// #3. return value
	if subtype == subtypeArray {
		return ar
	}
	return hm
}

func tryParseEmptyMapList(p *Parser) (bool, UnionMapList) {
	emptyTrialTypes := []lex.TokenType{
		lex.TypeArrayQuoteR, // for empty array
		lex.TypeMapData,     // for empty hashmap
	}

	if match, tk := p.tryConsume(emptyTrialTypes...); match {
		switch tk.Type {
		case lex.TypeArrayQuoteR:
			e := &ArrayExpr{Items: []Expression{}}
			e.SetCurrentLine(tk)
			return true, e
		case lex.TypeMapData:
			p.consume(lex.TypeArrayQuoteR)
			e := &HashMapExpr{KVPair: []hashMapKeyValuePair{}}
			e.SetCurrentLine(tk)
			return true, e
		}
	}
	return false, nil
}

// ParseFuncCallExpr - yield FuncCallExpr node
//
// CFG:
// FuncCallExpr  -> （ ID ： commaList ）
// commaList     -> E commaListTail
// commaListTail -> ， E commaListTail
//               ->
func ParseFuncCallExpr(p *Parser) *FuncCallExpr {
	var callExpr = &FuncCallExpr{
		Params: []Expression{},
	}
	// #1. parse ID
	match, tk := p.tryConsume(lex.TypeVarQuote, lex.TypeIdentifier)
	if !match {
		panic(error.InvalidSyntaxCurr())
	}
	callExpr.FuncName = newID(tk)
	// #2. parse colon (maybe there's no params)
	match2, _ := p.tryConsume(lex.TypeFuncCall)
	if match2 {
		// #2.1 parse comma list
		nodes := parseCommaList(p, func(idx int, nodes []Node) Node {
			return ParseExpression(p, true)
		})
		// #2.2 translate nodes into params
		for _, node := range nodes {
			v, _ := node.(Expression)
			callExpr.Params = append(callExpr.Params, v)
		}
	}

	// #3. parse right quote
	p.consume(lex.TypeFuncQuoteR)

	return callExpr
}

// ParseVarOneLeadExpr - 以 ... （‹方法名›）
// CFG:
//
// FuncExpr -> 以 Expr，Expr， ... RawFuncExpr
// RawFuncExpr -> （ ID ： commaList ）
func ParseVarOneLeadExpr(p *Parser) *FuncCallExpr {
	// #1. parse exprs
	exprList := []Expression{}
	nodes := parseCommaList(p, func(idx int, nodes []Node) Node {
		return ParseExpression(p, true)
	})
	for _, node := range nodes {
		expr, _ := node.(Expression) // don't worry, it must be an expression
		exprList = append(exprList, expr)
	}
	// #2. parse FuncExpr (maybe)
	match2, tk := p.tryConsume(lex.TypeFuncQuoteL)
	if !match2 {
		panic(error.InvalidSyntaxCurr())
	}

	// then suppose it's a funcCall expr
	funcCallExpr := ParseFuncCallExpr(p)
	// insert first ID into funcCall list
	funcCallExpr.Params = append(exprList, (funcCallExpr.Params)...)
	funcCallExpr.SetCurrentLine(tk)
	return funcCallExpr
}

// ParseVarDeclareStmt - yield VarDeclare node
// CFG:
// VarDeclare -> 令 IdfList 为 Expr
//    IdfList -> I I'
//         I' -> ，I I'
//            ->
//
// or block declaration:
//
// VarDeclare -> 令 ：
//           ...     I1 ， I2，
//           ...     I3 ， I4， I5 ...
func ParseVarDeclareStmt(p *Parser) *VarDeclareStmt {
	vNode := &VarDeclareStmt{
		AssignPair: []VDAssignPair{},
	}

	const (
		tagWithAssignExpr = 10
	)
	var nodes = []Node{}

	var consumer = func(idx int, nodes []Node) Node {
		// subExpr -> ID
		//         -> ID 为 expr
		var idExpr *ID
		// #1. consume ID first
		if match, tk := p.tryConsume(lex.TypeVarQuote, lex.TypeIdentifier); match {
			idExpr = newID(tk)
		} else {
			panic(error.InvalidSyntaxCurr())
		}

		// #2. consume LogicYes - if not, return ID directly
		if match2, _ := p.tryConsume(lex.TypeLogicYesW); !match2 {
			return idExpr
		}

		// #3. consume expr
		assignExpr := ParseExpression(p, true)

		return &NodeList{
			Tag:      tagWithAssignExpr,
			Children: []Node{idExpr, assignExpr},
		}
	}
	// #01. try to read colon
	// if colon exists -> parse comma list by block
	// if colon not exists -> parse comma list inline
	if match, _ := p.tryConsume(lex.TypeFuncCall); match {
		expected, blockIndent := p.expectBlockIndent()
		if !expected {
			panic(error.InvalidSyntaxCurr())
		}
		nodes = parseCommaListBlock(p, blockIndent, consumer)
	} else {
		// #02. consume identifier declare list (comma list) inline
		nodes = parseCommaList(p, consumer)
	}

	var idPtrList = []*ID{}
	// to avoid incomplete AssignPair (i.e. only IDs, but no VA part, like 令A，B，C)
	var completeAssignPair = false
	// #03. translate & append nodes to pair
	for _, node := range nodes {
		switch v := node.(type) {
		case *ID:
			completeAssignPair = false
			idPtrList = append(idPtrList, v)
		case *NodeList:
			if v.Tag == tagWithAssignExpr {
				newPair := VDAssignPair{
					Variables: []*ID{},
				}

				firstID, _ := v.Children[0].(*ID)
				idPtrList = append(idPtrList, firstID)
				//
				secondExpr, _ := v.Children[1].(Expression)
				// copy newPair
				newPair.Variables = append(newPair.Variables, idPtrList...)
				newPair.AssignExpr = secondExpr

				// append newPair
				vNode.AssignPair = append(vNode.AssignPair, newPair)
				// clear idPtrList
				idPtrList = []*ID{}
				completeAssignPair = true
			}
		}
	}

	if !completeAssignPair {
		panic(error.IncompleteStmtCurr())
	}

	return vNode
}

// ParseWhileLoopStmt - yield while loop node
// CFG:
// WhileLoopStmt -> 每当 Expr ：
//               ..     Block
func ParseWhileLoopStmt(p *Parser) *WhileLoopStmt {
	// #1. consume expr
	// 为  as logicYES here
	trueExpr := ParseExpression(p, false)

	// #2. parse colon
	p.consume(lex.TypeFuncCall)
	// #3. parse block
	expected, blockIndent := p.expectBlockIndent()
	if !expected {
		panic(error.InvalidSyntax())
	}
	block := ParseBlockStmt(p, blockIndent)
	return &WhileLoopStmt{
		TrueExpr:  trueExpr,
		LoopBlock: block,
	}
}

// ParseBlockStmt - parse all statements inside a block
func ParseBlockStmt(p *Parser, blockIndent int) *BlockStmt {
	bStmt := &BlockStmt{
		Children: []Statement{},
	}

	for (p.peek().Type != lex.TypeEOF) && p.getPeekIndent() == blockIndent {
		stmt := ParseStatement(p)
		bStmt.Children = append(bStmt.Children, stmt)
	}

	return bStmt
}

// ParseBranchStmt - yield BranchStmt node
// CFG:
// CondStmt -> 如果 IfTrueExpr ：
//         ...     IfTrueBlock
//
//          -> 如果 IfTrueExpr ：
//         ...     IfTrueBlock
//         ... 否则 ：
//         ...     IfFalseBlock
//
//          -> 如果 IfTrueExpr ：
//         ...     IfTrueBlock
//         ... 再如 OtherExpr1 ：
//         ...     OtherBlock1
//         ... 再如 OtherExpr2 ：
//         ...     OtherBlock2
//         ... ....
//             否则 ：
//         ...     IfFalseBlock
func ParseBranchStmt(p *Parser, mainIndent int) *BranchStmt {
	var condExpr Expression
	var condBlock *BlockStmt

	var stmt = new(BranchStmt)

	var condKeywords = []lex.TokenType{
		lex.TypeCondElseW,
		lex.TypeCondOtherW,
	}
	// by definition, the first Branch (if-branch) is required,
	// and the 如果 (if) keyword has been consumed before this function call.
	//
	// thus for other branches (like else-branch and elseif-branch),
	// we should consume the corresponding keyword explicitly. (否则，再如)
	const (
		stateInit        = 0
		stateIfBranch    = 1
		stateElseBranch  = 2
		stateOtherBranch = 3
	)
	var hState = stateInit

	for p.peek().Type != lex.TypeEOF {
		// parse header
		switch hState {
		case stateInit:
			hState = stateIfBranch
		case stateIfBranch, stateOtherBranch:
			if p.getPeekIndent() != mainIndent {
				return stmt
			}
			// parse related keywords (如果 expr： , 再如 expr：, 否则：)
			if match, tk := p.tryConsume(condKeywords...); match {
				if tk.Type == lex.TypeCondOtherW {
					hState = stateOtherBranch
				} else {
					hState = stateElseBranch
				}
			} else {
				return stmt
			}
		case stateElseBranch:
			if p.getPeekIndent() != mainIndent {
				return stmt
			}
			if match, _ := p.tryConsume(lex.TypeCondElseW); !match {
				return stmt
			}
		}

		// #1. parse condition expr
		if hState != stateElseBranch {
			condExpr = ParseExpression(p, false)
		}

		// #2. parse colon
		p.consume(lex.TypeFuncCall)

		// #3. parse block statements
		ok, blockIndent := p.expectBlockIndent()
		if !ok {
			panic(error.UnexpectedIndent())
		}
		condBlock = ParseBlockStmt(p, blockIndent)

		// #4. fill data
		switch hState {
		case stateIfBranch:
			stmt.IfTrueExpr = condExpr
			stmt.IfTrueBlock = condBlock
		case stateOtherBranch:
			stmt.OtherExprs = append(stmt.OtherExprs, condExpr)
			stmt.OtherBlocks = append(stmt.OtherBlocks, condBlock)
		case stateElseBranch:
			stmt.HasElse = true
			stmt.IfFalseBlock = condBlock
			// only one else-branch is accepted
			return stmt
		}
	}
	return stmt
}

// ParseFunctionDeclareStmt - yield FunctionDeclareStmt node
// CFG:
// FDStmt -> 如何 FuncName ？
//       ...     已知 ID1， ID2， ...
//       ...     ExecBlock
//       ...     ....
//
// FDStmt -> 如何 FuncName ？
//       ...     ExecBlock
//       ...     ....
//
func ParseFunctionDeclareStmt(p *Parser) *FunctionDeclareStmt {
	var fdStmt = &FunctionDeclareStmt{
		ParamList: []*ID{},
	}
	// by definition, when 已知 statement exists, it should be at first line
	// of function block
	const (
		stateParamList = 0
		stateFuncBlock = 2
	)
	var hState = stateParamList

	// #1. try to parse ID
	match, tk := p.tryConsume(lex.TypeIdentifier, lex.TypeVarQuote)
	if !match {
		panic(error.InvalidSyntaxCurr())
	}
	fdStmt.FuncName = newID(tk)
	// #2. try to parse question mark
	p.consume(lex.TypeFuncDeclare)

	// #3. parse block manually
	ok, blockIndent := p.expectBlockIndent()
	if !ok {
		panic(error.UnexpectedIndent())
	}
	// #3.1 parse param def list
	for p.peek().Type != lex.TypeEOF && p.getPeekIndent() == blockIndent {
		switch hState {
		case stateParamList:
			// parse 已知 expr
			if match, _ := p.tryConsume(lex.TypeParamAssignW); match {
				fdStmt.ParamList = parseParamDefList(p, true)
				// then change state
				hState = stateFuncBlock
			} else {
				hState = stateFuncBlock
			}
		case stateFuncBlock:
			fdStmt.ExecBlock = ParseBlockStmt(p, blockIndent)
		}
	}
	return fdStmt
}

// ParseVarOneLeadStmt -
// There're 2 possible statements
//
// 1. 以 K，V 遍历...
// 2. 以 A，B，C （执行方法）
//
// CFG:
//
// VOStmt -> 以 ID ， ID ... 遍历 IStmtT'
//        -> 以 Expr ， Expr ... FuncExprT'
func ParseVarOneLeadStmt(p *Parser) Statement {
	validTypes := []lex.TokenType{
		lex.TypeIteratorW,
		lex.TypeFuncQuoteL,
	}
	exprList := parseCommaList(p, func(idx int, nodes []Node) Node {
		return ParseExpression(p, true)
	})

	match, tk := p.tryConsume(validTypes...)
	if match {
		switch tk.Type {
		case lex.TypeIteratorW:
			// validate if each node in exprList is an ID type
			// otherwise an error will be thrown
			idList := []*ID{}
			for _, pExpr := range exprList {
				if id, ok := pExpr.(*ID); ok {
					idList = append(idList, id)
				} else {
					panic(error.InvalidExprType("id"))
				}
			}
			return parseIteratorStmtRest(p, idList)
		case lex.TypeFuncQuoteL:
			targetExpr := ParseFuncCallExpr(p)
			// transform []Node -> []Expression
			pListExprs := []Expression{}
			for _, pExpr := range exprList {
				exprItem, _ := pExpr.(Expression)
				pListExprs = append(pListExprs, exprItem)
			}
			// prepend exprs
			targetExpr.Params = append(pListExprs, targetExpr.Params...)
			return targetExpr
		}
	}
	panic(error.InvalidSyntax())
}

// ParseIteratorStmt - parse iterate stmt that starts with 遍历 keyword only
// CFG:
//
// IStmt -> 遍历 TargetExpr ：  StmtBlock
func ParseIteratorStmt(p *Parser) *IterateStmt {
	return parseIteratorStmtRest(p, []*ID{})
}

// parseIteratorStmtRest - parse after 以 ... and meet 遍历
// IStmtT'  -> [遍历] TargetExpr ：  StmtBlock
func parseIteratorStmtRest(p *Parser, idList []*ID) *IterateStmt {
	// 1. parse target expr
	targetExpr := ParseExpression(p, true)

	// 2. parse colon
	p.consume(lex.TypeFuncCall)

	// 3. parse iterate block
	expected, blockIndent := p.expectBlockIndent()
	if !expected {
		panic(error.InvalidSyntax())
	}
	block := ParseBlockStmt(p, blockIndent)

	return &IterateStmt{
		IterateExpr:  targetExpr,
		IndexNames:   idList,
		IterateBlock: block,
	}
}

// ParseFunctionReturnStmt - yield FuncParamList node (without head token: 返回)
//
// CFG:
// FRStmt -> 返回 Expression
func ParseFunctionReturnStmt(p *Parser) *FunctionReturnStmt {
	expr := ParseExpression(p, true)
	return &FunctionReturnStmt{
		ReturnExpr: expr,
	}
}

// tryParseParamList - try to parse line
func tryParseParamList(p *Parser, fdStmt *FunctionDeclareStmt) bool {
	// #1. parse 已知
	if match, _ := p.tryConsume(lex.TypeParamAssignW); !match {
		return false
	}
	// #2. parse param list (ID or VarQuote)
	nodes := parseCommaList(p, func(i int, n []Node) Node {
		// parse ID or VarQuote
		match, tk := p.tryConsume(lex.TypeVarQuote, lex.TypeIdentifier)
		if !match {
			panic(error.InvalidSyntaxCurr())
		}

		return newID(tk)
	})

	// transform nodes to actual identifiers
	for _, node := range nodes {
		// make sure node must be *ID type
		v, _ := node.(*ID)
		fdStmt.ParamList = append(fdStmt.ParamList, v)
	}
	return true
}

//// parse helpers
func parseCommaList(p *Parser, consumer consumerFunc) []Node {
	var node Node
	list := []Node{}

	// first item MUST be consumed!
	node = consumer(0, list)
	list = append(list, node)

	// iterate to get value
	for {
		// consume comma
		if match, _ := p.tryConsume(lex.TypeCommaSep); !match {
			// stop parsing immediately
			return list
		}
		node = consumer(len(list), list)
		list = append(list, node)
	}
}

func parseCommaListBlock(p *Parser, blockIndent int, consumer consumerFunc) []Node {
	var node Node
	//
	list := []Node{}

	// first token MUST be exactly on the indent
	if p.getPeekIndent() != blockIndent {
		panic(error.UnexpectedIndent())
	}
	// first item MUST be consumed!
	node = consumer(0, list)
	list = append(list, node)

	// iterate to get value
	for {
		if p.getPeekIndent() != blockIndent {
			return list
		}
		// consume comma
		if match, _ := p.tryConsume(lex.TypeCommaSep); !match {
			// stop parsing immediately
			return list
		}
		node = consumer(len(list), list)
		list = append(list, node)
	}
}

func parseParamDefList(p *Parser, allowBreak bool) []*ID {
	defer func() {
		if allowBreak {
			p.resetLineTermFlag()
		}
	}()
	var idList = []*ID{}

	// parse param lists
	nodes := parseCommaList(p, func(idx int, nodes []Node) Node {
		match2, tk := p.tryConsume(lex.TypeVarQuote, lex.TypeIdentifier)
		if !match2 {
			panic(error.InvalidSyntaxCurr())
		}

		return newID(tk)
	})
	// append IDs
	for _, node := range nodes {
		v, _ := node.(*ID)
		idList = append(idList, v)
	}

	return idList
}

func newID(tk *lex.Token) *ID {
	id := new(ID)
	id.SetLiteral(tk.Literal)
	id.SetCurrentLine(tk)
	return id
}

func newNumber(tk *lex.Token) *Number {
	num := new(Number)
	num.SetLiteral(tk.Literal)
	num.SetCurrentLine(tk)
	return num
}

func newString(tk *lex.Token) *String {
	str := new(String)
	// remove first char and last char (that are left & right quotes)
	str.SetLiteral(tk.Literal[1 : len(tk.Literal)-1])
	str.SetCurrentLine(tk)
	return str
}

// public helpers

// NewProgramNode -
func NewProgramNode(block *BlockStmt) *Program {
	return &Program{
		Content: block,
	}
}

// NewIDNode -
func NewIDNode(tk *lex.Token) *ID {
	return newID(tk)
}

// NewNumberNode -
func NewNumberNode(tk *lex.Token) *Number {
	return newNumber(tk)
}

// NewStringNode -
func NewStringNode(tk *lex.Token) *String {
	return newString(tk)
}
