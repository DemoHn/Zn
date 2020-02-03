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
	TargetVar  *ID
	AssignExpr Expression
}

func (va *VarAssignStmt) statementNode() {}

// ParseVarAssignStmt - parse general variable assign statement
//
// CFG:
// VarAssignStmt -> TargetV 为 ExprA           (1)
// VarAssignStmt -> TargetV 是 ExprA           (1A)
//               -> ExprA ， 得到 TargetV       (2)
//
func (p *Parser) ParseVarAssignStmt() (*VarAssignStmt, *error.Error) {
	p.setLineMask(modeInline)
	defer p.unsetLineMask(modeInline)

	var stmt = new(VarAssignStmt)
	var isTargetFirst = true
	// #0. parse first expression
	// either ExprT (case 1) or ExprA (case 2)
	firstExpr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}

	// #1. parse the middle one
	switch p.current().Type {
	case lex.TypeLogicYesW, lex.TypeLogicYesIIW:
		// NOTICE: currently we only support ID as the first expr, so we
		// check the type here
		if id, ok := firstExpr.(*ID); ok {
			stmt.TargetVar = id
			p.next()
		} else {
			return nil, error.InvalidSyntax()
		}
	case lex.TypeCommaSep:
		if p.peek().Type == lex.TypeFuncYieldW {
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
	secondExpr, err := p.ParseExpression()
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
