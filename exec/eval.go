package exec

import (
	"fmt"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/exec/stdlib"
	"github.com/DemoHn/Zn/exec/val"
	"github.com/DemoHn/Zn/syntax"
)

const errCodeMethodNotFound = 0x2505

// eval.go evaluates program from generated AST tree with specific scopes
// common signature of eval functions:
//
// evalXXXXStmt(c *ctx.Context, node Node) *error.Error
//
// or
//
// evalXXXXExpr(c *ctx.Context, node Node) (ctx.Value, *error.Error)
//
// NOTICE:
// `evalXXXXStmt` will change the value of its corresponding scope; However, `evalXXXXExpr` will export
// a ctx.Value object and mostly won't change scopes (but search a variable from scope is frequently used)

// val.DuplicateValue - deepcopy values' structure, including bool, string, decimal, array, hashmap
// for function or object or null, pass the original reference instead.
// This is due to the 'copycat by default' policy

func evalProgram(c *ctx.Context, program *syntax.Program) *error.Error {
	otherStmts, err := evalPreStmtBlock(c, program.Content)
	if err != nil {
		return err
	}
	return evalStmtBlock(c, otherStmts)
}

//// eval statements

// EvalStatement - eval statement
func evalStatement(c *ctx.Context, stmt syntax.Statement) *error.Error {
	var returnValue ctx.Value
	var sp = c.GetScope()
	// set currentLine
	c.GetFileInfo().SetCurrentLine(stmt.GetCurrentLine())

	// set return value
	defer func() {
		var finalReturnValue ctx.Value = val.NewNull()
		// set current return value
		if returnValue != nil {
			finalReturnValue = returnValue
		}
		sp.SetReturnValue(finalReturnValue)

		// set parent return value
		parentScope := sp.FindParentScope()
		if parentScope != nil {
			parentScope.SetReturnValue(finalReturnValue)
		}
	}()

	switch v := stmt.(type) {
	case *syntax.VarDeclareStmt:
		return evalVarDeclareStmt(c, v)
	case *syntax.WhileLoopStmt:
		return evalWhileLoopStmt(c, v)
	case *syntax.BranchStmt:
		return evalBranchStmt(c, v)
	case *syntax.EmptyStmt:
		return nil
	case *syntax.FunctionDeclareStmt:
		fn := BuildFunctionFromNode(v)
		return c.BindSymbol(v.FuncName.GetLiteral(), fn)
	case *syntax.ClassDeclareStmt:
		if sp.FindParentScope() != nil {
			return error.NewErrorSLOT("只能在代码主层级定义类")
		}
		// bind classRef
		className := v.ClassName.GetLiteral()
		classRef := BuildClassFromNode(className, v)

		return c.SetImportValue(className, &classRef)
	case *syntax.IterateStmt:
		return evalIterateStmt(c, v)
	case *syntax.FunctionReturnStmt:
		val, err := evalExpression(c, v.ReturnExpr)
		if err != nil {
			return err
		}
		// send RETURN break
		return error.ReturnBreakError(val)
	case syntax.Expression:
		expr, err := evalExpression(c, v)
		returnValue = expr
		return err
	default:
		return error.UnExpectedCase("语句类型", fmt.Sprintf("%T", v))
	}
}

// evalVarDeclareStmt - consists of three branches:
// 1. A，B 为 C
// 2. A，B 成为 X：P1，P2，...
// 3. A，B 恒为 C
func evalVarDeclareStmt(c *ctx.Context, node *syntax.VarDeclareStmt) *error.Error {
	for _, vpair := range node.AssignPair {
		switch vpair.Type {
		case syntax.VDTypeAssign, syntax.VDTypeAssignConst: // 为，恒为
			obj, err := evalExpression(c, vpair.AssignExpr)
			if err != nil {
				return err
			}
			// if assign const
			isConst := false
			if vpair.Type == syntax.VDTypeAssignConst {
				isConst = true
			}

			for _, v := range vpair.Variables {
				vtag := v.GetLiteral()
				if !vpair.RefMark {
					obj = val.DuplicateValue(obj)
				}

				if err := c.BindSymbolDecl(vtag, obj, isConst); err != nil {
					return err
				}
			}
		case syntax.VDTypeObjNew: // 成为
			if err := evalNewObject(c, vpair); err != nil {
				return err
			}
		}
	}
	return nil
}

// eval A,B 成为 C：P1，P2，P3，...
// ensure VDAssignPair.Type MUST BE syntax.VDTypeObjNew
func evalNewObject(c *ctx.Context, node syntax.VDAssignPair) *error.Error {
	vtag := node.ObjClass.GetLiteral()
	// get class definition
	importVal, err := c.GetImportValue(vtag)
	if err != nil {
		return err
	}
	classRef, ok := importVal.(*val.ClassRef)
	if !ok {
		return error.InvalidParamType("classRef")
	}

	cParams, err := exprsToValues(c, node.ObjParams)
	if err != nil {
		return err
	}

	// assign new object to variables
	for _, v := range node.Variables {
		vtag := v.GetLiteral()

		finalObj, err := classRef.Construct(c, cParams)
		if err != nil {
			return err
		}

		if c.BindSymbol(vtag, finalObj); err != nil {
			return err
		}
	}
	return nil
}

// evalWhileLoopStmt -
func evalWhileLoopStmt(c *ctx.Context, node *syntax.WhileLoopStmt) *error.Error {
	currentScope := c.GetScope()
	// create new scope
	newScope := currentScope.CreateChildScope()
	newScope.SetThisValue(val.NewLoopCtl())
	// set context's current scope with new one
	c.SetScope(newScope)
	// after finish executing this block, revert scope to old one
	defer c.SetScope(currentScope)

	for {
		// #1. first execute expr
		trueExpr, err := evalExpression(c, node.TrueExpr)
		if err != nil {
			return err
		}
		// #2. assert trueExpr to be Bool
		vTrueExpr, ok := trueExpr.(*val.Bool)
		if !ok {
			return error.InvalidExprType("bool")
		}
		// break the loop if expr yields not true
		if !vTrueExpr.GetValue() {
			return nil
		}
		// #3. stmt block
		if err := evalStmtBlock(c, node.LoopBlock); err != nil {
			if err.GetCode() == error.ContinueBreakSignal {
				// continue next turn
				continue
			}
			if err.GetCode() == error.BreakBreakSignal {
				// break directly
				return nil
			}
			return err
		}
	}
}

// evalPreStmtBlock - execute classRef, functionDeclare, imports first, then other statements inside the block
func evalPreStmtBlock(c *ctx.Context, block *syntax.BlockStmt) (*syntax.BlockStmt, *error.Error) {
	otherStmts := &syntax.BlockStmt{
		Children: []syntax.Statement{},
	}

	for _, stmtI := range block.Children {
		switch v := stmtI.(type) {
		case *syntax.FunctionDeclareStmt:
			fn := BuildFunctionFromNode(v)
			if err := c.BindSymbol(v.FuncName.GetLiteral(), fn); err != nil {
				return nil, err
			}
		case *syntax.ClassDeclareStmt:
			// bind classRef
			className := v.ClassName.GetLiteral()
			classRef := BuildClassFromNode(className, v)
			if err := c.SetImportValue(className, &classRef); err != nil {
				return nil, err
			}
		case *syntax.ImportStmt:
			// TODO: support non-stdlib imports
			libName := v.ImportName.GetLiteral()
			libData, ok := stdlib.PackageList[libName]
			if !ok {
				return nil, error.NewErrorSLOT("对应的类库不存在")
			}

			itemsList := []string{}
			// if itemList is [] (e.g. 导入《 ... 》)
			// then we import all itmes in this libaray
			if len(v.ImportItems) == 0 {
				for k := range libData {
					itemsList = append(itemsList, k)
				}
			} else {
				for _, id := range v.ImportItems {
					itemsList = append(itemsList, id.GetLiteral())
				}
			}
			// import name globally
			for _, itemName := range itemsList {
				libItem, ok2 := libData[itemName]
				if ok2 {
					c.SetImportValue(itemName, libItem)
				}
			}
		default:
			otherStmts.Children = append(otherStmts.Children, stmtI)
		}
	}
	return otherStmts, nil
}

// EvalStmtBlock -
func evalStmtBlock(c *ctx.Context, block *syntax.BlockStmt) *error.Error {
	for _, stmt := range block.Children {
		if err := evalStatement(c, stmt); err != nil {
			return err
		}
	}

	return nil
}

func evalBranchStmt(c *ctx.Context, node *syntax.BranchStmt) *error.Error {
	// set scope
	currentScope := c.GetScope()
	c.SetScope(currentScope.CreateChildScope())
	defer c.SetScope(currentScope)

	// #1. condition header
	ifExpr, err := evalExpression(c, node.IfTrueExpr)
	if err != nil {
		return err
	}
	vIfExpr, ok := ifExpr.(*val.Bool)
	if !ok {
		return error.InvalidExprType("bool")
	}

	// exec if-branch
	if vIfExpr.GetValue() {
		return evalStmtBlock(c, node.IfTrueBlock)
	}
	// exec else-if branches
	for idx, otherExpr := range node.OtherExprs {
		otherExprI, err := evalExpression(c, otherExpr)
		if err != nil {
			return err
		}
		vOtherExprI, ok := otherExprI.(*val.Bool)
		if !ok {
			return error.InvalidExprType("bool")
		}
		// exec else-if branch
		if vOtherExprI.GetValue() {
			return evalStmtBlock(c, node.OtherBlocks[idx])
		}
	}
	// exec else branch if possible
	if node.HasElse {
		return evalStmtBlock(c, node.IfFalseBlock)
	}
	return nil
}

func evalIterateStmt(c *ctx.Context, node *syntax.IterateStmt) *error.Error {
	currentScope := c.GetScope()
	newScope := currentScope.CreateChildScope()
	newScope.SetThisValue(val.NewLoopCtl())
	// set new scope (and revert to old when done)
	c.SetScope(newScope)
	defer c.SetScope(currentScope)

	// pre-defined key, value variable name
	var keySlot, valueSlot string
	var nameLen = len(node.IndexNames)

	// 以A，B遍历C： D
	// execute expr: C
	targetExpr, err := evalExpression(c, node.IterateExpr)
	if err != nil {
		return err
	}

	// execIterationBlock, including set "currentKey" and "currentValue" to scope,
	// and preDefined indication variables
	execIterationBlockFn := func(key ctx.Value, v ctx.Value) *error.Error {
		// set pre-defined value
		if nameLen == 1 {
			if err := c.SetSymbol(valueSlot, v); err != nil {
				return err
			}
		} else if nameLen == 2 {
			if err := c.SetSymbol(keySlot, key); err != nil {
				return err
			}
			if err := c.SetSymbol(valueSlot, v); err != nil {
				return err
			}
		}
		return evalStmtBlock(c, node.IterateBlock)
	}

	// define indication variables as "currentKey" and "currentValue" under new iterScope
	// of course since there's no any iteration is executed yet, the initial values are all "Null"
	if nameLen == 1 {
		valueSlot = node.IndexNames[0].Literal
		if err := c.BindSymbol(valueSlot, val.NewNull()); err != nil {
			return err
		}
	} else if nameLen == 2 {
		keySlot = node.IndexNames[0].Literal
		valueSlot = node.IndexNames[1].Literal
		if err := c.BindSymbol(keySlot, val.NewNull()); err != nil {
			return err
		}
		if err := c.BindSymbol(valueSlot, val.NewNull()); err != nil {
			return err
		}
	} else if nameLen > 2 {
		return error.MostParamsError(2)
	}

	// execute iterations
	switch tv := targetExpr.(type) {
	case *val.Array:
		for idx, v := range tv.GetValue() {
			idxVar := val.NewDecimalFromInt(idx, 0)
			if err := execIterationBlockFn(idxVar, v); err != nil {
				if err.GetCode() == error.ContinueBreakSignal {
					// continue next turn
					continue
				}
				if err.GetCode() == error.BreakBreakSignal {
					// break directly
					return nil
				}
				return err
			}
		}
	case *val.HashMap:
		for _, key := range tv.GetKeyOrder() {
			v := tv.GetValue()[key]
			keyVar := val.NewString(key)
			// handle interrupts
			if err := execIterationBlockFn(keyVar, v); err != nil {
				if err.GetCode() == error.ContinueBreakSignal {
					// continue next turn
					continue
				}
				if err.GetCode() == error.BreakBreakSignal {
					// break directly
					return nil
				}
				return err
			}
		}
	default:
		return error.InvalidExprType("array", "hashmap")
	}
	return nil
}

//// execute expressions
func evalExpression(c *ctx.Context, expr syntax.Expression) (ctx.Value, *error.Error) {
	c.GetFileInfo().SetCurrentLine(expr.GetCurrentLine())

	switch e := expr.(type) {
	case *syntax.VarAssignExpr:
		return evalVarAssignExpr(c, e)
	case *syntax.LogicExpr:
		if e.Type == syntax.LogicAND || e.Type == syntax.LogicOR {
			return evalLogicCombiner(c, e)
		}
		return evalLogicComparator(c, e)
	case *syntax.MemberExpr:
		iv, err := getMemberExprIV(c, e)
		if err != nil {
			return nil, err
		}
		return iv.ReduceRHS(c)
	case *syntax.Number, *syntax.String, *syntax.ID, *syntax.ArrayExpr, *syntax.HashMapExpr:
		return evalPrimeExpr(c, e)
	case *syntax.FuncCallExpr:
		return evalFunctionCall(c, e)
	case *syntax.ObjDFuncCallExpr:
		return evalObjDFuncCallExpr(c, e)
	default:
		return nil, error.InvalidExprType()
	}
}

// （显示：A、B、C），得到D
func evalFunctionCall(c *ctx.Context, expr *syntax.FuncCallExpr) (ctx.Value, *error.Error) {
	var zf *val.ClosureRef
	vtag := expr.FuncName.GetLiteral()

	// for a function call, if thisValue NOT FOUND, that means the target closure is a FUNCTION
	// instead of a METHOD (which is defined on class definition statement)
	//
	// If thisValue != nil, we will attempt to find clsoure from its method list;
	// then look up from scope's values.
	//
	// If thisValue == nil, we will look up target closure from scope's values directly.
	thisValue, _ := c.FindThisValue()

	// exec params first
	params, err := exprsToValues(c, expr.Params)
	if err != nil {
		return nil, err
	}

	// if thisValue exists, find ID from its method list
	if thisValue != nil {
		v, err := thisValue.ExecMethod(c, vtag, params)
		if err == nil {
			if expr.YieldResult != nil {
				// add yield result
				vtag := expr.YieldResult.GetLiteral()
				// bind yield result
				c.BindScopeSymbolDecl(c.GetScope(), vtag, v)
			}
			// return result
			return v, nil
		}
		if err.GetCode() != errCodeMethodNotFound {
			return nil, err
		}
	}

	// if function value not found from object scope, look up from local scope
	if zf == nil {
		// find function definction
		v, err := c.FindSymbol(vtag)
		if err != nil {
			return nil, err
		}
		// assert value
		zval, ok := v.(*val.Function)
		if !ok {
			return nil, error.InvalidFuncVariable(vtag)
		}
		zf = zval.GetValue()
	}

	// exec function call via its ClosureRef
	v2, err := zf.Exec(c, thisValue, params)
	if err != nil {
		return nil, err
	}

	if expr.YieldResult != nil {
		// add yield result
		ytag := expr.YieldResult.GetLiteral()
		// bind yield result
		c.BindScopeSymbolDecl(c.GetScope(), ytag, v2)
	}

	// return result
	return v2, nil
}

// 对于A （执行：1、2、3）
func evalObjDFuncCallExpr(c *ctx.Context, expr *syntax.ObjDFuncCallExpr) (ctx.Value, *error.Error) {
	currentScope := c.GetScope()
	newScope := currentScope.CreateChildScope()
	// set scope
	c.SetScope(newScope)
	defer c.SetScope(currentScope)

	// 1. parse root expr
	rootExpr, err := evalExpression(c, expr.RootObject)
	if err != nil {
		return nil, err
	}
	// set this value
	c.GetScope().SetThisValue(rootExpr)

	v, err := evalFunctionCall(c, expr.FuncExpr)
	if err != nil {
		return nil, err
	}
	// add yield result
	if expr.YieldResult != nil {
		vtag := expr.YieldResult.GetLiteral()
		// bind yield result
		c.BindScopeSymbolDecl(currentScope, vtag, v)
	}

	return v, nil
}

// evaluate logic combination expressions
// such as A 且 B
// or A 或 B
func evalLogicCombiner(c *ctx.Context, expr *syntax.LogicExpr) (*val.Bool, *error.Error) {
	logicType := expr.Type
	// #1. eval left
	left, err := evalExpression(c, expr.LeftExpr)
	if err != nil {
		return nil, err
	}
	// #2. assert left expr type to be ZnBool
	vleft, ok := left.(*val.Bool)
	if !ok {
		return nil, error.InvalidExprType("bool")
	}
	// #3. check if the result could be retrieved earlier
	//
	// 1) for Y = A and B, if A = false, then Y must be false
	// 2) for Y = A or  B, if A = true, then Y must be true
	//
	// for those cases, we can yield result directly
	if logicType == syntax.LogicAND && !vleft.GetValue() {
		return val.NewBool(false), nil
	}
	if logicType == syntax.LogicOR && vleft.GetValue() {
		return val.NewBool(true), nil
	}
	// #4. eval right
	right, err := evalExpression(c, expr.RightExpr)
	if err != nil {
		return nil, err
	}
	vright, ok := right.(*val.Bool)
	if !ok {
		return nil, error.InvalidExprType("bool")
	}
	// then evalute data
	switch logicType {
	case syntax.LogicAND:
		return val.NewBool(vleft.GetValue() && vright.GetValue()), nil
	default: // logicOR
		return val.NewBool(vleft.GetValue() || vright.GetValue()), nil
	}
}

// evaluate logic comparator
// ensure both expressions are comparable (i.e. subtype of ZnComparable)
func evalLogicComparator(c *ctx.Context, expr *syntax.LogicExpr) (*val.Bool, *error.Error) {
	logicType := expr.Type
	// #1. eval left
	left, err := evalExpression(c, expr.LeftExpr)
	if err != nil {
		return nil, err
	}
	// #2. eval right
	right, err := evalExpression(c, expr.RightExpr)
	if err != nil {
		return nil, err
	}

	var cmpRes bool
	var cmpErr *error.Error
	// #3. do comparison
	switch logicType {
	case syntax.LogicEQ:
		cmpRes, cmpErr = val.CompareValues(left, right, val.CmpEq)
	case syntax.LogicNEQ:
		cmpRes, cmpErr = val.CompareValues(left, right, val.CmpEq)
		cmpRes = !cmpRes // reverse result
	case syntax.LogicGT:
		cmpRes, cmpErr = val.CompareValues(left, right, val.CmpGt)
	case syntax.LogicGTE:
		var cmp1, cmp2 bool
		cmp1, cmpErr = val.CompareValues(left, right, val.CmpGt)
		if cmpErr != nil {
			return nil, cmpErr
		}
		cmp2, cmpErr = val.CompareValues(left, right, val.CmpEq)
		cmpRes = cmp1 || cmp2
	case syntax.LogicLT:
		cmpRes, cmpErr = val.CompareValues(left, right, val.CmpLt)
	case syntax.LogicLTE:
		var cmp1, cmp2 bool
		cmp1, cmpErr = val.CompareValues(left, right, val.CmpLt)
		if cmpErr != nil {
			return nil, cmpErr
		}
		cmp2, cmpErr = val.CompareValues(left, right, val.CmpEq)
		cmpRes = cmp1 || cmp2
	default:
		return nil, error.UnExpectedCase("比较类型", fmt.Sprintf("%d", logicType))
	}

	return val.NewBool(cmpRes), cmpErr
}

// eval prime expr
func evalPrimeExpr(c *ctx.Context, expr syntax.Expression) (ctx.Value, *error.Error) {
	switch e := expr.(type) {
	case *syntax.Number:
		return val.NewDecimal(e.GetLiteral())
	case *syntax.String:
		return val.NewString(e.GetLiteral()), nil
	case *syntax.ID:
		vtag := e.GetLiteral()
		return c.FindSymbol(vtag)
	case *syntax.ArrayExpr:
		znObjs := []ctx.Value{}
		for _, item := range e.Items {
			expr, err := evalExpression(c, item)
			if err != nil {
				return nil, err
			}
			znObjs = append(znObjs, expr)
		}

		return val.NewArray(znObjs), nil
	case *syntax.HashMapExpr:
		znPairs := []val.KVPair{}
		for _, item := range e.KVPair {
			// the key of KVPair MUST BE one of the following types
			//   - Number (its literal as key)
			//   - ID (its literal as key)
			//   - String
			// Other types are NOT ACCEPTED.
			// Specially, we regard all those exprs as literals, i.e.
			// Number 123    <==> “123”
			// Number 1.5*10^8 <==> “1.5*10^8”
			// ID 标识符      <==> “标识符”
			// String “世界”  <==> “世界”
			var exprKey string
			switch k := item.Key.(type) {
			case *syntax.String:
				exprKey = k.GetLiteral()
			case *syntax.ID:
				exprKey = k.GetLiteral()
			case *syntax.Number:
				exprKey = k.GetLiteral()
			default:
				return nil, error.InvalidExprType("string", "decimal", "id")
			}

			exprVal, err := evalExpression(c, item.Value)
			if err != nil {
				return nil, err
			}
			znPairs = append(znPairs, val.KVPair{
				Key:   exprKey,
				Value: exprVal,
			})
		}
		return val.NewHashMap(znPairs), nil
	default:
		return nil, error.UnExpectedCase("表达式类型", fmt.Sprintf("%T", e))
	}
}

// eval variable assign
func evalVarAssignExpr(c *ctx.Context, expr *syntax.VarAssignExpr) (ctx.Value, *error.Error) {
	// Right Side
	vr, err := evalExpression(c, expr.AssignExpr)
	if err != nil {
		return nil, err
	}
	// if var assignment is NOT by reference, then duplicate value
	if !expr.RefMark {
		vr = val.DuplicateValue(vr)
	}

	// Left Side
	switch v := expr.TargetVar.(type) {
	case *syntax.ID:
		// set ID
		vtag := v.GetLiteral()
		err2 := c.SetSymbol(vtag, vr)
		return vr, err2
	case *syntax.MemberExpr:
		if v.MemberType == syntax.MemberID || v.MemberType == syntax.MemberIndex {
			iv, err := getMemberExprIV(c, v)
			if err != nil {
				return nil, err
			}
			return vr, iv.ReduceLHS(c, vr)
		}
		return nil, error.NewErrorSLOT("方法不能被赋值")
	default:
		return nil, error.UnExpectedCase("被赋值", fmt.Sprintf("%T", v))
	}
}

func getMemberExprIV(c *ctx.Context, expr *syntax.MemberExpr) (*val.IV, *error.Error) {
	switch expr.RootType {
	case syntax.RootTypeProp: // 其 XX
		thisValue, err := c.FindThisValue()
		if err != nil {
			return nil, err
		}
		return val.NewMemberIV(thisValue, expr.MemberID.GetLiteral()), nil
	case syntax.RootTypeExpr: // A 之 B
		valRoot, err := evalExpression(c, expr.Root)
		if err != nil {
			return nil, err
		}
		switch expr.MemberType {
		case syntax.MemberID: // A 之 B
			return val.NewMemberIV(valRoot, expr.MemberID.GetLiteral()), nil
		case syntax.MemberIndex: // A # 0
			idx, err := evalExpression(c, expr.MemberIndex)
			if err != nil {
				return nil, err
			}
			switch v := valRoot.(type) {
			case *val.Array:
				vr, ok := idx.(*val.Decimal)
				if !ok {
					return nil, error.InvalidExprType("integer")
				}
				vri, e := vr.AsInteger()
				if e != nil {
					return nil, error.InvalidExprType("integer")
				}
				return val.NewArrayIV(v, vri), nil
			case *val.HashMap:
				var s string
				switch x := idx.(type) {
				// regard decimal value directly as string
				case *val.Decimal:
					// transform decimal value to string
					// x.exp < 0 express that its a decimal value with point mark, not an integer
					if x.GetExp() < 0 {
						return nil, error.InvalidExprType("integer", "string")
					}
					s = x.String()
				case *val.String:
					s = x.String()
				default:
					return nil, error.InvalidExprType("integer", "string")
				}
				return val.NewHashMapIV(v, s), nil
			}
			return nil, error.InvalidExprType("array", "hashmap")
		}
		return nil, error.UnExpectedCase("子项类型", fmt.Sprintf("%d", expr.MemberType))
	}

	return nil, error.UnExpectedCase("根元素类型", fmt.Sprintf("%d", expr.RootType))
}

//// helpers
// exprsToValues - []syntax.Expression -> []eval.ctx.Value
func exprsToValues(c *ctx.Context, exprs []syntax.Expression) ([]ctx.Value, *error.Error) {
	params := []ctx.Value{}
	for _, paramExpr := range exprs {
		pval, err := evalExpression(c, paramExpr)
		if err != nil {
			return nil, err
		}
		params = append(params, pval)
	}
	return params, nil
}

// BuildClosureFromNode - create a closure (with default param handler logic)
// from Zn code (*syntax.BlockStmt). It's the constructor of 如何XX or (anoymous function in the future)
func BuildClosureFromNode(paramTags []*syntax.ParamItem, stmtBlock *syntax.BlockStmt) val.ClosureRef {
	var executor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
		// iterate block round I - function hoisting
		// NOTE: function hoisting means bind function definitions at the beginning
		// of execution so that even if "function execution" statement is before
		// "function definition" statement.
		for _, stmtI := range stmtBlock.Children {
			if v, ok := stmtI.(*syntax.FunctionDeclareStmt); ok {
				fn := BuildFunctionFromNode(v)
				if err := c.BindSymbol(v.FuncName.GetLiteral(), fn); err != nil {
					return nil, err
				}
			}
		}
		// iterate block round II - execution of rest code blocks
		for _, stmtII := range stmtBlock.Children {
			if _, ok := stmtII.(*syntax.FunctionDeclareStmt); !ok {
				if err := evalStatement(c, stmtII); err != nil {
					// if recv breaks
					if err.GetCode() == error.ReturnBreakSignal {
						if extra, ok := err.GetExtra().(ctx.Value); ok {
							return extra, nil
						}
					}
					return nil, err
				}
			}
		}
		return c.GetScope().GetReturnValue(), nil
	}

	var paramHandler = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
		// check param length
		if len(params) != len(paramTags) {
			return nil, error.MismatchParamLengthError(len(paramTags), len(params))
		}

		// bind params (as variable) to function scope
		for idx, paramVal := range params {
			param := paramTags[idx]
			// if param is NOT a reference type, then we need additionally
			// copy its value
			if !param.RefMark {
				paramVal = val.DuplicateValue(paramVal)
			}
			paramName := param.ID.GetLiteral()
			if err := c.BindSymbol(paramName, paramVal); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	return val.NewClosure(paramHandler, executor)
}

// BuildClassFromNode -
func BuildClassFromNode(name string, classNode *syntax.ClassDeclareStmt) val.ClassRef {
	ref := val.ClassRef{
		Name:         name,
		Constructor:  nil,
		PropList:     []string{},
		CompPropList: map[string]val.ClosureRef{},
		MethodList:   map[string]val.ClosureRef{},
	}

	// define default constrcutor
	var constructor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
		obj := val.NewObject(ref)
		propMap := obj.GetPropList()
		// init prop list
		for _, propPair := range classNode.PropertyList {
			propID := propPair.PropertyID.GetLiteral()
			expr, err := evalExpression(c, propPair.InitValue)
			if err != nil {
				return nil, err
			}

			propMap[propID] = expr
			ref.PropList = append(ref.PropList, propID)
		}
		// constructor: set some properties' value
		if len(params) != len(classNode.ConstructorIDList) {
			return nil, error.MismatchParamLengthError(len(params), len(classNode.ConstructorIDList))
		}
		for idx, objParamVal := range params {
			param := classNode.ConstructorIDList[idx]
			// if param is NOT a reference, then we need to copy its value
			if !param.RefMark {
				objParamVal = val.DuplicateValue(objParamVal)
			}
			paramName := param.ID.GetLiteral()
			propMap[paramName] = objParamVal
		}

		return obj, nil
	}
	// set constructor
	ref.Constructor = constructor

	// add getters
	for _, gNode := range classNode.GetterList {
		getterTag := gNode.GetterName.GetLiteral()
		ref.CompPropList[getterTag] = BuildClosureFromNode([]*syntax.ParamItem{}, gNode.ExecBlock)
	}

	// add methods
	for _, mNode := range classNode.MethodList {
		mTag := mNode.FuncName.GetLiteral()
		ref.MethodList[mTag] = BuildClosureFromNode(mNode.ParamList, mNode.ExecBlock)
	}

	return ref
}

// BuildFunctionFromNode -
func BuildFunctionFromNode(node *syntax.FunctionDeclareStmt) *val.Function {
	closureRef := BuildClosureFromNode(node.ParamList, node.ExecBlock)
	return val.NewFunctionFromClosure(closureRef)
}
