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
		target = StringifyAST(v.TargetExpr)
		assign = StringifyAST(v.AssignExpr)

		return fmt.Sprintf("$VA(target=(%s) assign=(%s))", target, assign)
	default:
		return ""
	}
}
