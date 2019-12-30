package syntax

import "github.com/DemoHn/Zn/error"

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
	/**firstExpr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	*/
	// TODO: add more cases
	return nil, nil
}
