package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// VarDeclareStmt - declare variables as init its values
type VarDeclareStmt struct {
	Variables  []*ID
	AssignExpr Expression
}

func (vn *VarDeclareStmt) statementNode() {}

// parsing process

// ParseVarDeclare - yield VarDeclare node
// CFG:
// VarDeclare -> 令 IdfList 为 Expr
//    IdfList -> I I'
//         I' -> ，I I'
//            ->
//
func (p *Parser) ParseVarDeclare() (*VarDeclareStmt, *error.Error) {
	// #0. consume LING keyword
	if err := p.consume(lex.TypeDeclareW); err != nil {
		return nil, err
	}

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
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	vNode.AssignExpr = expr
	return vNode, nil
}

func parseIdentifierList(p *Parser, vNode *VarDeclareStmt) *error.Error {
	// #0. consume Identifier
	if err := p.consumeFunc(cbIdentifier(vNode), lex.TypeVarQuote, lex.TypeIdentifier); err != nil {
		return err
	}
	// #1. parse identifier tail
	return parseIdentifierTail(p, vNode)
}

func parseIdentifierTail(p *Parser, vNode *VarDeclareStmt) *error.Error {
	// skip all parsing
	if p.current().Type != lex.TypeCommaSep {
		return nil
	}
	// #0. consume comma
	if err := p.consume(lex.TypeCommaSep); err != nil {
		return err
	}
	// #1. consume Identifier
	if err := p.consumeFunc(cbIdentifier(vNode), lex.TypeVarQuote, lex.TypeIdentifier); err != nil {
		return err
	}
	// #2. parse tail nested again
	return parseIdentifierTail(p, vNode)
}

// callback -
func cbIdentifier(vNode *VarDeclareStmt) func(tk *lex.Token) {
	return func(tk *lex.Token) {
		tid := new(ID)
		tid.SetLiteral(string(tk.Literal))
		// append variables
		vNode.Variables = append(vNode.Variables, tid)
	}
}
