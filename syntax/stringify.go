package syntax

import (
	"fmt"
	"strings"
)

// StringifyAST - stringify an abstract ProgramNode into a readable string
func StringifyAST(node Node) string {
	switch v := node.(type) {
	case *ProgramNode:
		var statements = []string{}
		for _, stmt := range v.Children {
			statements = append(statements, StringifyAST(stmt))
		}
		return fmt.Sprintf("$PG(%s)", strings.Join(statements, " "))
	// expressions
	case *ArrayExpr:
		var exprs = []string{}
		for _, expr := range v.Items {
			exprs = append(exprs, StringifyAST(expr))
		}
		return fmt.Sprintf("$ARR(%s)", strings.Join(exprs, " "))
	case *Number:
		return fmt.Sprintf("$NUM(%s)", v.literal)
	case *String:
		return fmt.Sprintf("$STR(%s)", v.literal)
	case *ID:
		return fmt.Sprintf("$ID(%s)", v.literal)
	// var assign expressions
	case *VarDeclareStmt:
		var vars = []string{}
		var expr string

		// parse vars
		for _, vd := range v.Variables {
			vars = append(vars, StringifyAST(vd))
		}
		// parse exprs
		expr = StringifyAST(v.AssignExpr)

		return fmt.Sprintf("$VD(vars=(%s) expr=(%s))", strings.Join(vars, " "), expr)
	case *VarAssignStmt:
		var target, assign string
		target = StringifyAST(v.TargetVar)
		assign = StringifyAST(v.AssignExpr)

		return fmt.Sprintf("$VA(target=(%s) assign=(%s))", target, assign)
	case *ConditionStmt:
		var conds = []string{}
		// add if-branch
		ifExpr := fmt.Sprintf("ifExpr=(%s)", StringifyAST(v.IfTrueExpr))
		ifBlock := fmt.Sprintf("ifBlock=(%s)", StringifyAST(v.IfTrueBlock))
		conds = append(conds, ifExpr, ifBlock)

		// add else-branch
		if v.HasIfFalse {
			elseBlock := fmt.Sprintf("elseBlock=(%s)", StringifyAST(v.IfFalseBlock))

			conds = append(conds, elseBlock)
		}

		// add other branchs
		for idx, expr := range v.OtherExprs {
			otherExpr := fmt.Sprintf("otherExpr[]=(%s)", StringifyAST(expr))
			otherBlock := fmt.Sprintf("otherBlock[]=(%s)", StringifyAST(v.OtherBlocks[idx]))

			conds = append(conds, otherExpr, otherBlock)
		}

		return fmt.Sprintf("$IF(%s)", strings.Join(conds, " "))
	case *BlockStmt:
		var statements = []string{}
		for _, stmt := range v.Content {
			statements = append(statements, StringifyAST(stmt))
		}
		return fmt.Sprintf("$BK(%s)", strings.Join(statements, " "))
	default:
		return ""
	}
}
