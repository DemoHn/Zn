package syntax

import (
	"fmt"
	"strings"
)

// Parser - parse source file into syntax tree for further execution.
type Parser struct {
	*Lexer
	ASTBuilder
}

// ASTBuilder - build AST from tokens. Its logic varies from different languages.
// Currently, only Chinese ASTBuilder is supported
type ASTBuilder interface {
	ParseAST(lexer *Lexer) (*Program, error)
}

// NewParser - create a new parser from source
func NewParser(lexer *Lexer, astBuilder ASTBuilder) *Parser {
	return &Parser{
		Lexer:        lexer,
		ASTBuilder:   astBuilder,
	}
}

// Parser - parse all tokens into syntax tree
// TODO: in the future we'll parse it into bytecodes directly, instead.
func (p *Parser) Parse() (ast *Program, err error) {
	// handle panics
	defer func() {
		var ok bool
		if r := recover(); r != nil {
			err, ok = r.(error)
			if !ok {
				panic(r)
			}
		}
	}()

	ast, err = p.ParseAST(p.Lexer)
	return
}

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
	case *HashMapExpr:
		var exprs = []string{}
		for _, expr := range v.KVPair {
			exprs = append(exprs, fmt.Sprintf("key[]=(%s) value[]=(%s)", StringifyAST(expr.Key), StringifyAST(expr.Value)))
		}
		return fmt.Sprintf("$HM(%s)", strings.Join(exprs, " "))
	case *Number:
		return fmt.Sprintf("$NUM(%s)", v.Literal)
	case *String:
		return fmt.Sprintf("$STR(%s)", v.Literal)
	case *ID:
		return fmt.Sprintf("$ID(%s)", v.Literal)
	case *LogicExpr:
		var typeStrMap = map[uint8]string{
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
		return fmt.Sprintf("%s(L=(%s) R=(%s))", typeStrMap[v.Type], lstr, rstr)
	case *ArithExpr:
		t := ""
		switch  v.Type {
		case ArithAdd:
			t = "ADD"
		case ArithSub:
			t = "SUB"
		case ArithMul:
			t = "MUL"
		case ArithDiv:
			t = "DIV"
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
			case VDTypeObjNew:
				paramsStr := []string{}
				for _, vp := range vpair.ObjParams {
					paramsStr = append(paramsStr, StringifyAST(vp))
				}
				items = append(items, fmt.Sprintf(
					"$VP(object vars[]=(%s) class=(%s) params[]=(%s))",
					strings.Join(vars, " "),
					StringifyAST(vpair.ObjClass),
					strings.Join(paramsStr, " "),
				))
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
		paramsStr := []string{}
		for _, p := range v.ParamList {
			paramsStr = append(paramsStr, StringifyAST(p))
		}

		return fmt.Sprintf("$FN(name=(%s) params=(%s) blockTokens=(%s))",
			StringifyAST(v.FuncName),
			strings.Join(paramsStr, " "),
			StringifyAST(v.ExecBlock))
	case *GetterDeclareStmt:
		return fmt.Sprintf("$GT(name=(%s) blockTokens=(%s))",
			StringifyAST(v.GetterName),
			StringifyAST(v.ExecBlock))
	case *BlockStmt:
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
		constructorStr := []string{}
		propertyStr := []string{}
		methodStr := []string{}
		getterStr := []string{}

		for _, c := range v.ConstructorIDList {
			constructorStr = append(constructorStr, StringifyAST(c))
		}
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
			"$CLS(name=(%s) properties=(%s) constructor=(%s) methods=(%s) getters=(%s))",
			StringifyAST(v.ClassName),
			strings.Join(propertyStr, " "),
			strings.Join(constructorStr, " "),
			strings.Join(methodStr, " "),
			strings.Join(getterStr, " "),
		)
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
