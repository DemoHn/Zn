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
