package syntax

import (
	"fmt"
	"strings"
)

// StringifyAST - stringify an abstract ProgramNode into a readable string, mainly for UNITTEST purpose
func StringifyAST(node Node) string {
	switch v := node.(type) {
	case *Program: // replace program in the future
		importBlockStr := []string{}
		for _, importBlock := range v.ImportBlock {
			importBlockStr = append(importBlockStr, StringifyAST(importBlock))
		}

		var importSS = ""
		if len(importBlockStr) > 0 {
			importSS = strings.Join(importBlockStr, " ") + " "
		}
		return fmt.Sprintf("$PG(%s%s)",
			importSS,
			StringifyAST(v.ExecBlock),
		)
	case *ExecBlock:
		// input block
		var idList = make([]string, 0)
		for _, id := range v.InputBlock {
			idList = append(idList, StringifyAST(id))
		}
		// catch block
		var catchStrList = make([]string, 0)
		for _, c := range v.CatchBlock {
			catchStrList = append(catchStrList, fmt.Sprintf("cls[]=(%s) stmt[]=(%s)",
				StringifyAST(c.ExceptionClass), StringifyAST(c.StmtBlock)))
		}

		return fmt.Sprintf("$X(I=(%s) S=(%s) C=(%s))",
			strings.Join(idList, " "),
			StringifyAST(v.StmtBlock),
			strings.Join(catchStrList, " "),
		)
	// expressions
	case *ArrayExpr:
		var exprs = []string{}
		for _, expr := range v.Items {
			exprs = append(exprs, StringifyAST(expr))
		}
		return fmt.Sprintf("$ARR(%s)", strings.Join(exprs, " "))
	case *HashMapExpr:
		var exprs = []string{}
		for _, expr := range v.KVPair {
			exprs = append(exprs, fmt.Sprintf("key[]=(%s) value[]=(%s)", StringifyAST(expr.Key), StringifyAST(expr.Value)))
		}
		return fmt.Sprintf("$HM(%s)", strings.Join(exprs, " "))
	case *String:
		return fmt.Sprintf("$STR(%s)", v.GetLiteral())
	case *ID:
		return fmt.Sprintf("$ID(%s)", v.GetLiteral())
	case *LogicExpr:
		var typeStrMap = map[uint8]string{
			LogicEQ:   "$EQ",
			LogicNEQ:  "$NEQ",
			LogicAND:  "$AND",
			LogicOR:   "$OR",
			LogicGT:   "$GT",
			LogicGTE:  "$GTE",
			LogicLT:   "$LT",
			LogicLTE:  "$LTE",
			LogicXEQ:  "$XEQ",
			LogicXNEQ: "$XNEQ",
		}

		lstr := StringifyAST(v.LeftExpr)
		rstr := StringifyAST(v.RightExpr)
		return fmt.Sprintf("%s(L=(%s) R=(%s))", typeStrMap[v.Type], lstr, rstr)
	case *ArithExpr:
		t := ""
		switch v.Type {
		case ArithAdd:
			t = "ADD"
		case ArithSub:
			t = "SUB"
		case ArithMul:
			t = "MUL"
		case ArithDiv:
			t = "DIV"
		case ArithIntDiv:
			t = "INTDIV"
		case ArithModulo:
			t = "MODULO"
		}

		return fmt.Sprintf("$AR(type=(%s) left=(%s) right=(%s))", t, StringifyAST(v.LeftExpr), StringifyAST(v.RightExpr))
	case *MemberExpr:
		var str = ""
		var sType = ""
		switch v.MemberType {
		case MemberID:
			str = StringifyAST(v.MemberID)
			sType = "mID"
		case MemberIndex:
			str = StringifyAST(v.MemberIndex)
			sType = "mIndex"
		}
		rootTypeStr := "rootScope"
		if v.RootType == RootTypeProp {
			if v.RootType == RootTypeProp {
				rootTypeStr = "rootProp"
			}
			return fmt.Sprintf("$MB(%s type=(%s) object=(%s))", rootTypeStr, sType, str)
		}
		rootStr := StringifyAST(v.Root)
		return fmt.Sprintf("$MB(root=(%s) type=(%s) object=(%s))", rootStr, sType, str)
	case *MemberMethodExpr:
		rootStr := StringifyAST(v.Root)
		// display chain
		chain := []string{}
		for _, method := range v.MethodChain {
			chain = append(chain, StringifyAST(method))
		}
		yieldRes := ""
		if v.YieldResult != nil {
			yieldRes = fmt.Sprintf(" yield=(%s)", StringifyAST(v.YieldResult))
		}
		return fmt.Sprintf("$MMF(root=(%s) chain=(%s)%s)", rootStr, strings.Join(chain, " "), yieldRes)
	case *ObjNewExpr:
		// chain
		var params = []string{}
		for _, p := range v.Params {
			params = append(params, StringifyAST(p))
		}
		return fmt.Sprintf("$NEW(class=(%s) params=(%s))", StringifyAST(v.ClassName), strings.Join(params, " "))
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

			switch vpair.Type {
			case VDTypeAssign, VDTypeAssignConst:
				expr = StringifyAST(vpair.AssignExpr)
				if vpair.Type == VDTypeAssignConst {
					items = append(items, fmt.Sprintf("$VP(const vars[]=(%s) expr[]=(%s))", strings.Join(vars, " "), expr))
				} else {
					items = append(items, fmt.Sprintf("$VP(vars[]=(%s) expr[]=(%s))", strings.Join(vars, " "), expr))
				}
			}
		}
		// parse exprs
		return fmt.Sprintf("$VD(%s)", strings.Join(items, " "))
	// var assign expressions
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

		yieldResult := ""
		if v.YieldResult != nil {
			yieldResult = fmt.Sprintf(" yield=(%s)", StringifyAST(v.YieldResult))
		}
		return fmt.Sprintf("$FN(name=(%s) params=(%s)%s)", name, strings.Join(params, " "), yieldResult)
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

		// add other branches
		for idx, expr := range v.OtherExprs {
			otherExpr := fmt.Sprintf("otherExpr[]=(%s)", StringifyAST(expr))
			otherBlock := fmt.Sprintf("otherBlock[]=(%s)", StringifyAST(v.OtherBlocks[idx]))

			conds = append(conds, otherExpr, otherBlock)
		}

		return fmt.Sprintf("$IF(%s)", strings.Join(conds, " "))
	case *WhileLoopStmt:
		return fmt.Sprintf("$WL(expr=(%s) block=(%s))", StringifyAST(v.TrueExpr), StringifyAST(v.LoopBlock))
	case *ImportStmt:
		itemsStr := []string{}
		for _, vi := range v.ImportItems {
			itemsStr = append(itemsStr, StringifyAST(vi))
		}
		return fmt.Sprintf("$IM(name=(%s) items=(%s))", StringifyAST(v.ImportName), strings.Join(itemsStr, " "))
	case *FunctionReturnStmt:
		return fmt.Sprintf("$RT(%s)", StringifyAST(v.ReturnExpr))
	case *FunctionDeclareStmt:
		fnTypeStr := "FN"
		if v.DeclareType == DeclareTypeConstructor {
			fnTypeStr = "COS"
		} else if v.DeclareType == DeclareTypeGetter {
			fnTypeStr = "GET"
		}
		return fmt.Sprintf("$FN(type=%s block=(%s))",
			fnTypeStr,
			StringifyAST(v.ExecBlock))
	case *StmtBlock:
		var statements = []string{}
		for _, stmt := range v.Children {
			statements = append(statements, StringifyAST(stmt))
		}
		return fmt.Sprintf("$BK(%s)", strings.Join(statements, " "))
	case *IterateStmt:
		paramsStr := []string{}
		for _, p := range v.IndexNames {
			paramsStr = append(paramsStr, StringifyAST(p))
		}
		return fmt.Sprintf("$IT(target=(%s) idxList=(%s) block=(%s))",
			StringifyAST(v.IterateExpr),
			strings.Join(paramsStr, " "),
			StringifyAST(v.IterateBlock))
	case *ClassDeclareStmt:
		propertyStr := []string{}
		methodStr := []string{}
		getterStr := []string{}
		for _, p := range v.PropertyList {
			propertyStr = append(propertyStr, StringifyAST(p))
		}
		for _, m := range v.MethodList {
			methodStr = append(methodStr, StringifyAST(m))
		}
		for _, g := range v.GetterList {
			getterStr = append(getterStr, StringifyAST(g))
		}

		return fmt.Sprintf(
			"$CLS(name=(%s) properties=(%s) methods=(%s) getters=(%s))",
			StringifyAST(v.ClassName),
			strings.Join(propertyStr, " "),
			strings.Join(methodStr, " "),
			strings.Join(getterStr, " "),
		)
	case *BreakStmt:
		return "$BREAK"
	case *ContinueStmt:
		return "$CONTINUE"
	case *PropertyDeclareStmt:
		return fmt.Sprintf(
			"$PD(id=(%s) expr=(%s))",
			StringifyAST(v.PropertyID),
			StringifyAST(v.InitValue),
		)
	case *ParamItem:
		refMark := "false"
		if v.RefMark {
			refMark = "true"
		}
		return fmt.Sprintf(
			"$PM(id=(%s) ref=(%s))",
			StringifyAST(v.ID),
			refMark,
		)
	default:
		return ""
	}
}
