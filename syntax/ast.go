package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

//////// Node types

// Node -
type Node interface{}

// Statement -
type Statement interface {
	Node
	stmtNode()
}

// Expression - a speical type of statement - that yields value after execution
type Expression interface {
	Node
	exprNode()
	IsPrimitive() bool
}

//// program (struct)

// Program -
type Program struct {
	Content *BlockStmt
}

//// Statements (struct)

// VarDeclareStmt - declare variables as init its values
type VarDeclareStmt struct {
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

// implement expression interface

// IsPrimitive - a primeExpr must be primitive, that is, no longer additional
// calculation required.
func (pe *PrimeExpr) IsPrimitive() bool { return true }

// SetLiteral - set literal for primeExpr
func (pe *PrimeExpr) SetLiteral(literal string) { pe.Literal = literal }

// GetLiteral -
func (pe *PrimeExpr) GetLiteral() string { return pe.Literal }

func (pe *PrimeExpr) exprNode() {}

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

	return ParseVarAssignStmt(p)
}

// ParseExpression - parse general expression (abstract expression type)
//
// currently, expression only contains
// ID
// Number
// String
// ArrayExpr
func ParseExpression(p *Parser) (Expression, *error.Error) {
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
		Items: make([]Expression, 0),
	}
	// #1. consume item list
	if err := parseItemList(p, ar); err != nil {
		return nil, err
	}
	// #2. consume right brancket
	if err := p.consume(lex.TypeArrayQuoteR); err != nil {
		return nil, err
	}
	return ar, nil
}

func parseItemList(p *Parser, ar *ArrayExpr) *error.Error {
	// #0. parse expression
	expr, err := ParseExpression(p)
	if err != nil {
		return err
	}
	ar.Items = append(ar.Items, expr)

	// #1. parse list tail
	return parseItemListTail(p, ar)
}

func parseItemListTail(p *Parser, ar *ArrayExpr) *error.Error {
	// #0. consume comma
	if match, _ := p.tryConsume([]lex.TokenType{lex.TypeCommaSep}); !match {
		return nil
	}
	// #1. parse expression
	expr, err := ParseExpression(p)
	if err != nil {
		return err
	}
	ar.Items = append(ar.Items, expr)

	// #2. parse tail nested again
	return parseItemListTail(p, ar)
}

// ParseVarAssignStmt - parse general variable assign statement
//
// CFG:
// VarAssignStmt -> TargetV 为 ExprA           (1)
// VarAssignStmt -> TargetV 是 ExprA           (1A)
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
func ParseVarDeclare(p *Parser) (*VarDeclareStmt, *error.Error) {
	p.setLineMask(modeInline)
	defer p.unsetLineMask(modeInline)

	vNode := &VarDeclareStmt{
		Variables:  make([]*ID, 0),
		AssignExpr: nil,
	}
	// #1. consume identifier list
	if err := parseIdentifierList(p, vNode); err != nil {
		return nil, err
	}
	// #2. consume logicYes
	if err := p.consume(lex.TypeLogicYesW); err != nil {
		return nil, err
	}

	// parse expression
	expr, err := ParseExpression(p)
	if err != nil {
		return nil, err
	}

	vNode.AssignExpr = expr
	return vNode, nil
}

func parseIdentifierList(p *Parser, vNode *VarDeclareStmt) *error.Error {
	validTypes := []lex.TokenType{
		lex.TypeVarQuote,
		lex.TypeIdentifier,
	}
	// #1. consume Identifier
	if match, tk := p.tryConsume(validTypes); match {
		tid := new(ID)
		tid.SetLiteral(string(tk.Literal))
		// append variables
		vNode.Variables = append(vNode.Variables, tid)
	} else {
		return error.InvalidSyntax()
	}

	// #2. parse identifier tail
	return parseIdentifierTail(p, vNode)
}

func parseIdentifierTail(p *Parser, vNode *VarDeclareStmt) *error.Error {
	commaTypes := []lex.TokenType{
		lex.TypeCommaSep,
	}
	idTypes := []lex.TokenType{
		lex.TypeVarQuote,
		lex.TypeIdentifier,
	}
	// #1. consume Comma
	if match, _ := p.tryConsume(commaTypes); !match {
		return nil
	}
	// #2. consume Identifier
	if match, tk := p.tryConsume(idTypes); match {
		tid := new(ID)
		tid.SetLiteral(string(tk.Literal))
		// append variables
		vNode.Variables = append(vNode.Variables, tid)
	} else {
		return error.InvalidSyntax()
	}

	// #3. parse tail nested again
	return parseIdentifierTail(p, vNode)
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
	var firstBranch = true
	var branchTokenType = lex.TypeCondW

	for p.peek().Type != lex.TypeEOF {
		// when firstBranch = true, we are parsing (if-branch) by default
		if !firstBranch {
			// check indent
			if p.getPeekIndent() != mainIndent {
				return stmt, nil
			}
			// consume other condKeywords (否则，再如)
			if match, tk := p.tryConsume(condKeywords); match {
				branchTokenType = tk.Type
			} else {
				// for other keywords, this parsing task has
				return stmt, nil
			}
		}
		p.setLineMask(modeInline)
		// #1 parse expression
		if branchTokenType != lex.TypeCondElseW {
			condExpr, err = ParseExpression(p)
			if err != nil {
				return nil, err
			}
		}
		// #2. parse colon
		if err := p.consume(lex.TypeFuncCall); err != nil {
			return nil, err
		}
		// #3 parse block stmts
		p.unsetLineMask(modeInline)
		ok, blockIndent := p.expectBlockIndent()
		if !ok {
			return nil, error.NewErrorSLOT("unexpected indent")
		}
		blockStmt, err := ParseBlockStmt(p, blockIndent)
		if err != nil {
			return nil, err
		}

		// #4. merge data
		switch branchTokenType {
		case lex.TypeCondW:
			stmt.IfTrueExpr = condExpr
			stmt.IfTrueBlock = blockStmt
			firstBranch = false
		case lex.TypeCondElseW:
			// set false-branch flag
			stmt.HasElse = true
			stmt.IfFalseBlock = blockStmt
		case lex.TypeCondOtherW:
			stmt.OtherExprs = append(stmt.OtherExprs, condExpr)
			stmt.OtherBlocks = append(stmt.OtherBlocks, blockStmt)
		}
	}
	return stmt, nil
}
