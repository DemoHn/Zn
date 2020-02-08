package syntax

import (
	"fmt"
	"strings"
)

// StringifyAST - stringify an abstract ProgramNode into a readable string
func StringifyAST(node Node) string {
	switch v := node.(type) {
	case *Program:
		return fmt.Sprintf("$PG(%s)", StringifyAST(v.Content))
	// expressions
	case *ArrayExpr:
		var exprs = []string{}
		for _, expr := range v.Items {
			exprs = append(exprs, StringifyAST(expr))
		}
		return fmt.Sprintf("$ARR(%s)", strings.Join(exprs, " "))
	case *Number:
		return fmt.Sprintf("$NUM(%s)", v.Literal)
	case *String:
		return fmt.Sprintf("$STR(%s)", v.Literal)
	case *ID:
		return fmt.Sprintf("$ID(%s)", v.Literal)
	case *LogicExpr:
		var typeStrMap = map[LogicType]string{
			LogicIS:  "$IS",
			LogicEQ:  "$EQ",
			LogicNEQ: "$NEQ",
			LogicAND: "$AND",
			LogicOR:  "$OR",
			LogicGT:  "$GT",
			LogicGTE: "$GTE",
			LogicLT:  "$LT",
			LogicLTE: "$LTE",
		}

		lstr := StringifyAST(v.LeftExpr)
		rstr := StringifyAST(v.RightExpr)
		return fmt.Sprintf("%s(%s %s)", typeStrMap[v.Type], lstr, rstr)
	// var assign expressions
	case *EmptyStmt:
		return "$"
	case *VarDeclareStmt:
		var items = []string{}
		// parse vars
		for _, vpair := range v.AssignPair {
			var vars = []string{}
			var expr string
			for _, vd := range vpair.Variables {
				vars = append(vars, StringifyAST(vd))
			}
			expr = StringifyAST(vpair.AssignExpr)
			items = append(items, fmt.Sprintf("vars[]=(%s) expr[]=(%s)", strings.Join(vars, " "), expr))
		}
		// parse exprs
		return fmt.Sprintf("$VD(%s)", strings.Join(items, " "))
	case *VarAssignExpr:
		var target, assign string
		target = StringifyAST(v.TargetVar)
		assign = StringifyAST(v.AssignExpr)

		return fmt.Sprintf("$VA(target=(%s) assign=(%s))", target, assign)
	case *FuncCallExpr:
		var params = []string{}
		var name = StringifyAST(v.FuncName)

		for _, p := range v.Params {
			params = append(params, StringifyAST(p))
		}

		return fmt.Sprintf("$FN(name=(%s) params=(%s))", name, strings.Join(params, " "))
	case *BranchStmt:
		var conds = []string{}
		// add if-branch
		ifExpr := fmt.Sprintf("ifExpr=(%s)", StringifyAST(v.IfTrueExpr))
		ifBlock := fmt.Sprintf("ifBlock=(%s)", StringifyAST(v.IfTrueBlock))
		conds = append(conds, ifExpr, ifBlock)

		// add else-branch
		if v.HasElse {
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
	case *WhileLoopStmt:
		return fmt.Sprintf("$WL(expr=(%s) block=(%s))", StringifyAST(v.TrueExpr), StringifyAST(v.LoopBlock))
	case *BlockStmt:
		var statements = []string{}
		for _, stmt := range v.Children {
			statements = append(statements, StringifyAST(stmt))
		}
		return fmt.Sprintf("$BK(%s)", strings.Join(statements, " "))
	default:
		return ""
	}
}
