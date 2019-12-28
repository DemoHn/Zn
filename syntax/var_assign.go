package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// VarAssignNode - var assignment statement
type VarAssignNode struct {
	Variables  []*IdentifierNode
	AssignExpr Expression
}

func (vn *VarAssignNode) getType() nodeType {
	return TypeVarAssign
}

// parsing process

// ParseVarAssign - yield VarAssign node
// CFG:
// VarAssign -> 令 IdfList 为 Expr
//   IdfList -> I I'
//        I' -> ，I I'
//           ->
//
func (p *Parser) ParseVarAssign() (*VarAssignNode, *error.Error) {
	// #0. consume LING keyword
	if err := p.consume(lex.TypeDeclareW); err != nil {
		return nil, err
	}

	vNode := &VarAssignNode{
		Variables:  make([]*IdentifierNode, 0),
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
	// #3. parse expression
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	vNode.AssignExpr = expr
	return vNode, nil
}

func parseIdentifierList(p *Parser, vNode *VarAssignNode) *error.Error {
	// #0. consume Identifier
	if err := p.consumeFunc(cbIdentifier(vNode), lex.TypeVarQuote, lex.TypeIdentifier); err != nil {
		return err
	}
	// #1. parse identifier tail
	return parseIdentifierTail(p, vNode)
}

func parseIdentifierTail(p *Parser, vNode *VarAssignNode) *error.Error {
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
func cbIdentifier(vNode *VarAssignNode) func(tk *lex.Token) {
	return func(tk *lex.Token) {
		// append variables
		vNode.Variables = append(vNode.Variables, &IdentifierNode{
			literal: string(tk.Literal),
		})
	}
}
