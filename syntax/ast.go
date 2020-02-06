package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

//////// Node types

// Node -
type Node interface{}

type consumerFunc func(idx int, nodes []Node) (Node, *error.Error)

// Statement -
type Statement interface {
	Node
	stmtNode()
}

// Expression - a speical type of statement - that yields value after execution
type Expression interface {
	Statement
	exprNode()
	IsPrimitive() bool
}

//// program (struct)

// Program -
type Program struct {
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
	AssignPair []VDAssignPair
}

// VDAssignPair - helper type
type VDAssignPair struct {
	Variables  []*ID
	AssignExpr Expression
}

// VarAssignStmt - variable assignment statement
// assign <TargetExpr> from <AssignExpr>
//
// i.e.
// TargetExpr := AssignExpr
type VarAssignStmt struct {
	TargetVar  *ID
	AssignExpr Expression
}

// BranchStmt - conditional (if-else) statement
type BranchStmt struct {
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

// BlockStmt -
type BlockStmt struct {
	Children []Statement
}

// implement statement inteface
func (va *VarAssignStmt) stmtNode()  {}
func (vn *VarDeclareStmt) stmtNode() {}
func (bk *BlockStmt) stmtNode()      {}
func (bs *BranchStmt) stmtNode()     {}

//// Expressions (struct)

// PrimeExpr - primitive expression
type PrimeExpr struct {
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
	PrimeExpr
	Items []Expression
}

// LogicType - define several logic type (OR, AND, EQ, etc)
type LogicType uint8

// declare some logic types
const (
	LogicOR  LogicType = 1 // 或
	LogicAND LogicType = 2 // 且
	LogicIS  LogicType = 3 // 此 ... 为 ...
	LogicEQ  LogicType = 4 // 等于
	LogicNEQ LogicType = 5 // 不等于
	LogicGT  LogicType = 6 // 大于
	LogicGTE LogicType = 7 // 不小于
	LogicLT  LogicType = 8 // 小于
	LogicLTE LogicType = 9 // 不大于
)

// LogicExpr - logical expression return TRUE (真) or FALSE (假) only
type LogicExpr struct {
	Type      LogicType
	LeftExpr  Expression
	RightExpr Expression
}

// IsPrimitive -
func (le *LogicExpr) IsPrimitive() bool { return false }

// implement expression interface

// IsPrimitive - a primeExpr must be primitive, that is, no longer additional
// calculation required.
func (pe *PrimeExpr) IsPrimitive() bool { return true }

// SetLiteral - set literal for primeExpr
func (pe *PrimeExpr) SetLiteral(literal string) { pe.Literal = literal }

// GetLiteral -
func (pe *PrimeExpr) GetLiteral() string { return pe.Literal }

func (pe *PrimeExpr) exprNode() {}
func (pe *PrimeExpr) stmtNode() {}

func (le *LogicExpr) exprNode() {}
func (le *LogicExpr) stmtNode() {}

//////// Parse Methods

// ParseStatement - a program consists of statements
//
// CFG:
// Statement -> VarDeclareStmt
//           -> VarAssignStmt
//           -> ；
func ParseStatement(p *Parser) (Statement, *error.Error) {
	var validTypes = []lex.TokenType{
		lex.TypeStmtSep,
		lex.TypeDeclareW,
		lex.TypeCondW,
	}
	match, tk := p.tryConsume(validTypes)
	if match {
		switch tk.Type {
		case lex.TypeStmtSep:
			// skip
			return nil, nil
		case lex.TypeDeclareW:
			return ParseVarDeclare(p)
		case lex.TypeCondW:
			mainIndent := p.getPeekIndent()
			return ParseBranchStmt(p, mainIndent)
		}
	}
	return ParseExpression(p)
	//return ParseVarAssignStmt(p)
}

// ParseExpression - parse an expression, see the following CFG for details
//
// CFG:
// E  -> AndE E'
// E' -> 或 AndE E'
//    ->
//
// AndE  -> EqE AndE'
// AndE' -> 且 EqE AndE'
//       ->
//
// EqE   -> VaE
// EqE'  -> 等于 VaE
//       -> 不等于 VaE
//       -> 小于 VaE
//       -> 不小于 VaE
//       -> 大于 VaE
//       -> 不大于 VaE
//       ->
//
// VaE   -> BsE VaE'
// VaE'  -> 为 BsE
//       ->
//
// BsE   -> { E }
//       -> （ ID ： E，E，...）
//       -> 此 E 为 E
//       -> ID
//       -> Number
//       -> String
//       -> 【 E，E，...】
//
func ParseExpression(p *Parser) (Expression, *error.Error) {
	var logicItemParser func(int) (Expression, *error.Error)
	var logicItemTailParser func(int, Expression) (Expression, *error.Error)
	// logicKeywords, ordered by precedence asc
	// that means, the very begin logicKeyword ([]lex.TokenType) has lowest precedence
	var logicKeywords = [4][]lex.TokenType{
		[]lex.TokenType{lex.TypeLogicOrW},
		[]lex.TokenType{lex.TypeLogicAndW},
		[]lex.TokenType{
			lex.TypeLogicEqualW,
			lex.TypeLogicNotEqW,
			lex.TypeLogicGtW,
			lex.TypeLogicGteW,
			lex.TypeLogicLtW,
			lex.TypeLogicLteW,
		},
		[]lex.TokenType{lex.TypeLogicYesW}, // notice: this represents for Variable Assignment!
	}
	var logicTypeMap = map[lex.TokenType]LogicType{
		lex.TypeLogicOrW:    LogicOR,
		lex.TypeLogicAndW:   LogicAND,
		lex.TypeLogicEqualW: LogicEQ,
		lex.TypeLogicNotEqW: LogicNEQ,
		lex.TypeLogicGtW:    LogicGT,
		lex.TypeLogicGteW:   LogicGTE,
		lex.TypeLogicLtW:    LogicLT,
		lex.TypeLogicLteW:   LogicLTE,
		lex.TypeLogicYesW:   LogicIS,
	}
	var logicAllowTails = [4]bool{true, true, false, false}

	//// anynomous function definition
	logicItemParser = func(idx int) (Expression, *error.Error) {
		if idx >= len(logicKeywords) {
			return ParseBasicExpr(p)
		}
		// #1. match item
		expr1, err := logicItemParser(idx + 1)
		if err != nil {
			return nil, err
		}
		return logicItemTailParser(idx, expr1)
	}

	//// anynomous function definition
	logicItemTailParser = func(idx int, leftExpr Expression) (Expression, *error.Error) {
		// #1. consume keyword
		match, tk := p.tryConsume(logicKeywords[idx])
		if !match {
			return leftExpr, nil
		}
		// #2. consume Y
		rightExpr, err := logicItemParser(idx + 1)
		if err != nil {
			return nil, err
		}
		// compose logic expr
		logicExpr := &LogicExpr{
			Type:      logicTypeMap[tk.Type],
			LeftExpr:  leftExpr,
			RightExpr: rightExpr,
		}
		// #3. consume X' (X-tail)
		if !logicAllowTails[idx] {
			return logicItemTailParser(idx, logicExpr)
		} else {
			return logicExpr, nil
		}
	}

	return logicItemParser(0)
}

// ParseBasicExpr - parse general basic expression
//
// currently, basic expression only contains
// ID
// Number
// String
// ArrayExpr
func ParseBasicExpr(p *Parser) (Expression, *error.Error) {
	var validTypes = []lex.TokenType{
		lex.TypeIdentifier,
		lex.TypeVarQuote,
		lex.TypeNumber,
		lex.TypeString,
		lex.TypeArrayQuoteL,
	}

	match, tk := p.tryConsume(validTypes)
	if match {
		switch tk.Type {
		case lex.TypeIdentifier, lex.TypeVarQuote:
			expr := new(ID)
			expr.SetLiteral(string(tk.Literal))
			return expr, nil
		case lex.TypeNumber:
			expr := new(Number)
			expr.SetLiteral(string(tk.Literal))
			return expr, nil
		case lex.TypeString:
			expr := new(String)
			expr.SetLiteral(string(tk.Literal))
			return expr, nil
		case lex.TypeArrayQuoteL:
			arrExpr, err := ParseArrayExpr(p)
			if err != nil {
				return nil, err
			}
			return arrExpr, nil
		}
	}
	return nil, error.InvalidSyntax()
}

// ParseArrayExpr - yield ArrayExpr node
// CFG:
// ArrayExpr -> 【 ItemList 】
// ItemList  -> Expr ExprTail
//           ->
// ExprTail  -> ， Expr ExprTail
//           ->
//
// Expr      -> PrimaryExpr
//
// PrimaryExpr -> Number
//             -> String
//             -> ID
//             -> ArrayExpr
func ParseArrayExpr(p *Parser) (*ArrayExpr, *error.Error) {
	ar := &ArrayExpr{
		Items: []Expression{},
	}
	// #1. consume item list (comma list)
	exprs, err := parseCommaList(p, func(idx int, nodes []Node) (Node, *error.Error) {
		return ParseExpression(p)
	})
	if err != nil {
		return nil, err
	}

	// type cast (because there's no GENERIC TYPE in golang!!!)
	for _, expr := range exprs {
		exprT, _ := expr.(Expression)
		ar.Items = append(ar.Items, exprT)
	}

	// #2. consume right brancket
	if err := p.consume(lex.TypeArrayQuoteR); err != nil {
		return nil, err
	}
	return ar, nil
}

// ParseVarAssignStmt - parse general variable assign statement
//
// CFG:
// VarAssignStmt -> TargetV 为 ExprA           (1)
//               -> ExprA ， 得到 TargetV       (2)
//
func ParseVarAssignStmt(p *Parser) (*VarAssignStmt, *error.Error) {
	p.setLineMask(modeInline)
	defer p.unsetLineMask(modeInline)

	var stmt = new(VarAssignStmt)
	var isTargetFirst = true
	// #0. parse first expression
	// either ExprT (case 1) or ExprA (case 2)
	firstExpr, err := ParseExpression(p)
	if err != nil {
		return nil, err
	}

	// #1. parse the middle one
	switch p.peek().Type {
	case lex.TypeLogicYesW:
		// NOTICE: currently we only support ID as the first expr, so we
		// check the type here
		if id, ok := firstExpr.(*ID); ok {
			stmt.TargetVar = id
			p.next()
		} else {
			return nil, error.InvalidSyntax()
		}
	case lex.TypeCommaSep:
		if p.peek2().Type == lex.TypeFuncYieldW {
			stmt.AssignExpr = firstExpr
			p.next()
			p.next()
			isTargetFirst = false
		} else {
			return nil, error.InvalidSyntax()
		}
	default:
		return nil, error.InvalidSyntax()
	}

	// #2. parse the second expression
	secondExpr, err := ParseExpression(p)
	if err != nil {
		return nil, err
	}
	if isTargetFirst {
		stmt.AssignExpr = secondExpr
	} else {
		// assert the second expr as ID
		if id, ok := secondExpr.(*ID); ok {
			stmt.TargetVar = id
		} else {
			return nil, error.InvalidSyntax()
		}
	}

	return stmt, nil
}

// ParseVarDeclare - yield VarDeclare node
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
func ParseVarDeclare(p *Parser) (*VarDeclareStmt, *error.Error) {
	vNode := &VarDeclareStmt{
		AssignPair: []VDAssignPair{},
	}

	idTypes := []lex.TokenType{
		lex.TypeVarQuote,
		lex.TypeIdentifier,
	}

	const (
		tagWithAssignExpr = 10
	)
	var nodes = []Node{}
	var err *error.Error

	var consumer = func(idx int, nodes []Node) (Node, *error.Error) {
		// subExpr -> ID
		//         -> ID 为 expr
		var idExpr *ID
		// #1. consume ID first
		if match, tk := p.tryConsume(idTypes); match {
			idExpr = new(ID)
			idExpr.SetLiteral(string(tk.Literal))
		} else {
			return nil, error.InvalidSyntax()
		}

		// #2. consume LogicYes - if not, return ID directly
		if match2, _ := p.tryConsume([]lex.TokenType{lex.TypeLogicYesW}); !match2 {
			return idExpr, nil
		}

		// #3. consume expr
		assignExpr, err2 := ParseExpression(p)
		if err2 != nil {
			return nil, err2
		}
		return &NodeList{
			Tag:      tagWithAssignExpr,
			Children: []Node{idExpr, assignExpr},
		}, nil
	}
	// #01. try to read colon
	// if colon exists -> parse comma list by block
	// if colon not exists -> parse comma list inline
	if match, _ := p.tryConsume([]lex.TokenType{lex.TypeFuncCall}); match {
		expected, blockIndent := p.expectBlockIndent()
		if !expected {
			return nil, error.InvalidSyntax()
		}
		if nodes, err = parseCommaListBlock(p, blockIndent, consumer); err != nil {
			return nil, err
		}
	} else {
		// #02. consume identifier declare list (comma list) inline
		if nodes, err = parseCommaList(p, consumer); err != nil {
			return nil, err
		}
	}

	var idPtrList = []*ID{}
	// #03. translate & append nodes to pair
	for _, node := range nodes {
		switch v := node.(type) {
		case *ID:
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
			}
		}
	}

	return vNode, nil
}

// ParseBlockStmt - parse all statements inside a block
func ParseBlockStmt(p *Parser, blockIndent int) (*BlockStmt, *error.Error) {
	bStmt := &BlockStmt{
		Children: []Statement{},
	}

	for (p.peek().Type != lex.TypeEOF) && p.getPeekIndent() == blockIndent {
		stmt, err := ParseStatement(p)
		if err != nil {
			return nil, err
		}
		bStmt.Children = append(bStmt.Children, stmt)
	}

	return bStmt, nil
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
func ParseBranchStmt(p *Parser, mainIndent int) (*BranchStmt, *error.Error) {
	var condExpr Expression
	var condBlock *BlockStmt
	var err *error.Error

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
				return stmt, nil
			}
			// parse related keywords (如果 expr： , 再如 expr：, 否则：)
			if match, tk := p.tryConsume(condKeywords); match {
				if tk.Type == lex.TypeCondOtherW {
					hState = stateOtherBranch
				} else {
					hState = stateElseBranch
				}
			} else {
				return stmt, nil
			}
		case stateElseBranch:
			if p.getPeekIndent() != mainIndent {
				return stmt, nil
			}
			if match, _ := p.tryConsume([]lex.TokenType{lex.TypeCondElseW}); !match {
				return stmt, nil
			}
		}

		p.setLineMask(modeInline)
		// #1. parse expr
		if hState != stateElseBranch {
			if condExpr, err = ParseExpression(p); err != nil {
				return nil, err
			}
		}

		// #2. parse colon
		if err = p.consume(lex.TypeFuncCall); err != nil {
			return nil, err
		}
		p.unsetLineMask(modeInline)

		// #3. parse block statements
		ok, blockIndent := p.expectBlockIndent()
		if !ok {
			return nil, error.NewErrorSLOT("unexpected indent")
		}
		if condBlock, err = ParseBlockStmt(p, blockIndent); err != nil {
			return nil, err
		}

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
			return stmt, nil
		}
	}
	return stmt, nil
}

// parse helpers
func parseCommaList(p *Parser, consumer consumerFunc) ([]Node, *error.Error) {
	var node Node
	var err *error.Error
	//
	list := []Node{}

	var sepTypes = []lex.TokenType{
		lex.TypeCommaSep,
	}
	// first item MUST be consumed!
	if node, err = consumer(0, list); err != nil {
		return nil, err
	}
	list = append(list, node)

	// iterate to get value
	for {
		// consume comma
		if match, _ := p.tryConsume(sepTypes); !match {
			// stop parsing immediately
			return list, nil
		}
		if node, err = consumer(len(list), list); err != nil {
			return nil, err
		}
		list = append(list, node)
	}
}

func parseCommaListBlock(p *Parser, blockIndent int, consumer consumerFunc) ([]Node, *error.Error) {
	var node Node
	var err *error.Error
	//
	list := []Node{}

	var sepTypes = []lex.TokenType{
		lex.TypeCommaSep,
	}
	// first token MUST be exactly on the indent
	if p.getPeekIndent() != blockIndent {
		return nil, error.NewErrorSLOT("unexpected indent")
	}
	// first item MUST be consumed!
	if node, err = consumer(0, list); err != nil {
		return nil, err
	}
	list = append(list, node)

	// iterate to get value
	for {
		if p.getPeekIndent() != blockIndent {
			return list, nil
		}
		// consume comma
		if match, _ := p.tryConsume(sepTypes); !match {
			// stop parsing immediately
			return list, nil
		}
		if node, err = consumer(len(list), list); err != nil {
			return nil, err
		}
		list = append(list, node)
	}
}
