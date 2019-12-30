package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// VarAssignStmt - variable assignment statement
// assign <TargetExpr> from <AssignExpr>
//
// i.e.
// TargetExpr := AssignExpr
type VarAssignStmt struct {
	TargetExpr Expression
	AssignExpr Expression
}

func (va *VarAssignStmt) getType() nodeType {
	return TypeVarAssign
}

func (va *VarAssignStmt) statementNode() {}

// ParseVarAssignStmt - parse general variable assign statement
//
// CFG:
// VarAssignStmt -> ExprT 设为 ExprA           (1)
//               -> ExprA ， 得 ExprT          (2)
//
// TODO:
// we need special handling for
//
// FuncName ： A，B，C，得ExprT
func (p *Parser) ParseVarAssignStmt() (*VarAssignStmt, *error.Error) {
	// #0. parse first expression
	// either ExprT (case 1) or ExprA (case 2)
	firstExpr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}

	// #1. parse the middle one
	var order int // type (1)
	t1 := p.current().Type
	if t1 == lex.TypeAssignW {
		order = 1
		p.next()
	} else if t1 == lex.TypeCommaSep {
		if p.peek().Type == lex.TypeFuncYieldW {
			order = 2
			p.next()
			p.next()
		} else {
			return nil, error.NewErrorSLOT("parsing comma wrongly (should not exist here)")
		}
	}

	// #2. parse the second expression
	secondExpr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}

	if order == 1 {
		return &VarAssignStmt{
			TargetExpr: firstExpr,
			AssignExpr: secondExpr,
		}, nil
	} else if order == 2 {
		return &VarAssignStmt{
			TargetExpr: secondExpr,
			AssignExpr: firstExpr,
		}, nil
	}
	return nil, error.NewErrorSLOT("unknown error (should not exists here)")
}
