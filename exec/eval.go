package exec

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/syntax"
)

// eval.go evaluates program from generated AST tree with specific scopes
// common signature of eval functions:
//
// evalXXXXStmt(ctx *Context, scope Scope, node Node) *error.Error
//
// or
//
// evalXXXXExpr(ctx *Context, scope Scope, node Node) (ZnValue, *error.Error)
//
// NOTICE:
// `evalXXXXStmt` will change the value of its corresponding scope; However, `evalXXXXExpr` will export
// a ZnValue object and mostly won't change scopes (but search a variable from scope is frequently used)

// TODO: find a better way to handle this
func duplicateValue(in ZnValue) ZnValue {
	return in
}

//// eval program
func evalProgram(ctx *Context, scope *RootScope, program *syntax.Program) *error.Error {
	return evalStmtBlock(ctx, scope, program.Content)
}

//// eval statements

// EvalStatement - eval statement
func evalStatement(ctx *Context, scope Scope, stmt syntax.Statement) *error.Error {
	scope.GetRoot().SetCurrentLine(stmt.GetCurrentLine())
	switch v := stmt.(type) {
	case *syntax.VarDeclareStmt:
		return evalVarDeclareStmt(ctx, scope, v)
	case *syntax.WhileLoopStmt:
		return evalWhileLoopStmt(ctx, scope, v)
	case *syntax.BranchStmt:
		return evalBranchStmt(ctx, scope, v)
	case *syntax.EmptyStmt:
		return nil
	case *syntax.FunctionDeclareStmt:
		fn := NewZnFunction(v)
		return bindValue(ctx, scope, v.FuncName.GetLiteral(), fn, false)
	case *syntax.IterateStmt:
		return evalIterateStmt(ctx, scope, v)
	case *syntax.FunctionReturnStmt:
		val, err := evalExpression(ctx, scope, v.ReturnExpr)
		if err != nil {
			return err
		}
		// send RETURN break
		return error.ReturnBreakError(val)
	case syntax.Expression:
		val, err := evalExpression(ctx, scope, v)
		if err != nil {
			return err
		}
		// set last value (of rootScope or funcScope)
		sp := scope
		for sp != nil {
			switch v := sp.(type) {
			case *RootScope:
				v.SetLastValue(val)
				return nil
			case *FuncScope:
				v.SetReturnValue(val)
				return nil
			}
			sp = sp.GetParent()
		}
		return nil
	default:
		return error.InvalidCaseType()
	}
}

func evalVarDeclareStmt(ctx *Context, scope Scope, node *syntax.VarDeclareStmt) *error.Error {
	for _, vpair := range node.AssignPair {
		obj, err := evalExpression(ctx, scope, vpair.AssignExpr)
		if err != nil {
			return err
		}
		for _, v := range vpair.Variables {
			vtag := v.GetLiteral()
			finalObj := duplicateValue(obj)
			if bindValue(ctx, scope, vtag, finalObj, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func evalWhileLoopStmt(ctx *Context, scope Scope, node *syntax.WhileLoopStmt) *error.Error {
	loopScope := createScope(ctx, scope, sTypeWhile)
	for {
		// #1. first execute expr
		trueExpr, err := evalExpression(ctx, loopScope, node.TrueExpr)
		if err != nil {
			return err
		}
		// #2. assert trueExpr to be ZnBool
		vTrueExpr, ok := trueExpr.(*ZnBool)
		if !ok {
			return error.InvalidExprType("bool")
		}
		// break the loop if expr yields not true
		if vTrueExpr.Value == false {
			return nil
		}
		// #3. stmt block
		if err := evalStmtBlock(ctx, loopScope, node.LoopBlock); err != nil {
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

// EvalStmtBlock -
func evalStmtBlock(ctx *Context, scope Scope, block *syntax.BlockStmt) *error.Error {
	for _, stmt := range block.Children {
		err := evalStatement(ctx, scope, stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func evalBranchStmt(ctx *Context, scope Scope, node *syntax.BranchStmt) *error.Error {
	// #1. condition header
	ifExpr, err := evalExpression(ctx, scope, node.IfTrueExpr)
	if err != nil {
		return err
	}
	vIfExpr, ok := ifExpr.(*ZnBool)
	if !ok {
		return error.InvalidExprType("bool")
	}
	// exec if-branch
	if vIfExpr.Value == true {
		return evalStmtBlock(ctx, scope, node.IfTrueBlock)
	}
	// exec else-if branches
	for idx, otherExpr := range node.OtherExprs {
		otherExprI, err := evalExpression(ctx, scope, otherExpr)
		if err != nil {
			return err
		}
		vOtherExprI, ok := otherExprI.(*ZnBool)
		if !ok {
			return error.InvalidExprType("bool")
		}
		// exec else-if branch
		if vOtherExprI.Value == true {
			return evalStmtBlock(ctx, scope, node.OtherBlocks[idx])
		}
	}
	// exec else branch if possible
	if node.HasElse == true {
		return evalStmtBlock(ctx, scope, node.IfFalseBlock)
	}
	return nil
}

func evalIterateStmt(ctx *Context, scope Scope, node *syntax.IterateStmt) *error.Error {
	// pre-defined key, value variable name
	var keySlot, valueSlot string
	var nameLen = len(node.IndexNames)

	iterScope := createIterateScope(ctx, scope)
	// 以A，B遍历C： D
	// execute expr: C
	targetExpr, err := evalExpression(ctx, scope, node.IterateExpr)
	if err != nil {
		return err
	}

	// execIterationBlock, including set "currentKey" and "currentValue" to scope,
	// and preDefined indication variables
	execIterationBlockFn := func(key ZnValue, val ZnValue) *error.Error {
		// set values of 此之值 and 此之
		iterScope.setCurrentKV(key, val)

		// set pre-defined value
		if nameLen == 1 {
			if err := setValue(ctx, iterScope, valueSlot, val); err != nil {
				return err
			}
		} else if nameLen == 2 {
			if err := setValue(ctx, iterScope, keySlot, key); err != nil {
				return err
			}
			if err := setValue(ctx, iterScope, valueSlot, val); err != nil {
				return err
			}
		}
		return evalStmtBlock(ctx, iterScope, node.IterateBlock)
	}

	// define indication variables as "currentKey" and "currentValue" under new iterScope
	// of course since there's no any iteration is executed yet, the initial values are all "Null"
	if nameLen == 1 {
		valueSlot = node.IndexNames[0].Literal
		if err := bindValue(ctx, iterScope, valueSlot, NewZnNull(), false); err != nil {
			return err
		}
	} else if nameLen == 2 {
		keySlot = node.IndexNames[0].Literal
		valueSlot = node.IndexNames[1].Literal
		if err := bindValue(ctx, iterScope, keySlot, NewZnNull(), false); err != nil {
			return err
		}
		if err := bindValue(ctx, iterScope, valueSlot, NewZnNull(), false); err != nil {
			return err
		}
	} else if nameLen > 2 {
		return error.NewErrorSLOT("过多的前置变量个数")
	}

	// execute iterations
	switch tv := targetExpr.(type) {
	case *ZnArray:
		for idx, val := range tv.Value {
			idxVar := NewZnDecimalFromInt(idx)
			if err := execIterationBlockFn(idxVar, val); err != nil {
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
	case *ZnHashMap:
		for key, val := range tv.Value {
			keyVar := NewZnString(key)
			// handle interrupts
			if err := execIterationBlockFn(keyVar, val); err != nil {
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

func evalExpression(ctx *Context, scope Scope, expr syntax.Expression) (ZnValue, *error.Error) {
	scope.GetRoot().SetCurrentLine(expr.GetCurrentLine())
	switch e := expr.(type) {
	case *syntax.VarAssignExpr:
		return evalVarAssignExpr(ctx, scope, e)
	case *syntax.LogicExpr:
		if e.Type == syntax.LogicAND || e.Type == syntax.LogicOR {
			return evalLogicCombiner(ctx, scope, e)
		}
		return evalLogicComparator(ctx, scope, e)
	case *syntax.MemberExpr:
		iv, err := getMemberExprIV(ctx, scope, e)
		if err != nil {
			return nil, err
		}
		return iv.Reduce(ctx, nil, false)
	case *syntax.Number, *syntax.String, *syntax.ID, *syntax.ArrayExpr, *syntax.HashMapExpr:
		return evalPrimeExpr(ctx, scope, e)
	case *syntax.FuncCallExpr:
		return evalFunctionCall(ctx, scope, e)
	default:
		return nil, error.InvalidExprType()
	}
}

// （显示：A，B，C）
func evalFunctionCall(ctx *Context, scope Scope, expr *syntax.FuncCallExpr) (ZnValue, *error.Error) {
	fScope, _ := createScope(ctx, scope, sTypeFunc).(*FuncScope)
	vtag := expr.FuncName.GetLiteral()
	// find function definctxion
	val, err := getValue(ctx, scope, vtag)
	if err != nil {
		return nil, err
	}
	// assert value
	zf, ok := val.(*ZnFunction)
	if !ok {
		return nil, error.InvalidFuncVariable(vtag)
	}
	// exec params
	params, err := exprsToValues(ctx, scope, expr.Params)
	if err != nil {
		return nil, err
	}
	return evalFunctionValue(ctx, fScope, params, zf)
}

func evalFunctionValue(ctx *Context, scope *FuncScope, params []ZnValue, zf *ZnFunction) (ZnValue, *error.Error) {
	// if executor = nil, then use default function executor
	if zf.Executor == nil {
		// check param length
		if len(params) != len(zf.Node.ParamList) {
			return nil, error.MismatchParamLengthError(len(zf.Node.ParamList), len(params))
		}

		// set id
		for idx, param := range params {
			paramID := zf.Node.ParamList[idx]
			if err := bindValue(ctx, scope, paramID.GetLiteral(), param, false); err != nil {
				return nil, err
			}
		}

		execBlock := zf.Node.ExecBlock
		// iterate block
		for _, stmt := range execBlock.Children {
			if err := evalStatement(ctx, scope, stmt); err != nil {
				// if recv breaks
				if err.GetCode() == error.ReturnBreakSignal {
					if extra, ok := err.GetExtra().(ZnValue); ok {
						return extra, nil
					}
				}
				return nil, err
			}
		}
		return scope.GetReturnValue(), nil
	}
	// use pre-defined execution logic
	return zf.Executor(ctx, scope, params)

}

// evaluate logic combination expressions
// such as A 且 B
// or A 或 B
func evalLogicCombiner(ctx *Context, scope Scope, expr *syntax.LogicExpr) (*ZnBool, *error.Error) {
	logicType := expr.Type
	// #1. eval left
	left, err := evalExpression(ctx, scope, expr.LeftExpr)
	if err != nil {
		return nil, err
	}
	// #2. assert left expr type to be ZnBool
	vleft, ok := left.(*ZnBool)
	if !ok {
		return nil, error.InvalidExprType("bool")
	}
	// #3. check if the result could be retrieved earlier
	//
	// 1) for Y = A and B, if A = false, then Y must be false
	// 2) for Y = A or  B, if A = true, then Y must be true
	//
	// for those cases, we can yield result directly
	if logicType == syntax.LogicAND && vleft.Value == false {
		return NewZnBool(false), nil
	}
	if logicType == syntax.LogicOR && vleft.Value == true {
		return NewZnBool(true), nil
	}
	// #4. eval right
	right, err := evalExpression(ctx, scope, expr.RightExpr)
	if err != nil {
		return nil, err
	}
	vright, ok := right.(*ZnBool)
	if !ok {
		return nil, error.InvalidExprType("bool")
	}
	// then evalute data
	switch logicType {
	case syntax.LogicAND:
		return NewZnBool(vleft.Value && vright.Value), nil
	default: // logicOR
		return NewZnBool(vleft.Value || vright.Value), nil
	}
}

// evaluate logic comparator
// ensure both expressions are comparable (i.e. subtype of ZnComparable)
func evalLogicComparator(ctx *Context, scope Scope, expr *syntax.LogicExpr) (*ZnBool, *error.Error) {
	logicType := expr.Type
	// #1. eval left
	left, err := evalExpression(ctx, scope, expr.LeftExpr)
	if err != nil {
		return nil, err
	}
	// #3. eval right
	right, err := evalExpression(ctx, scope, expr.RightExpr)
	if err != nil {
		return nil, err
	}

	// #5. evaluate
	switch logicType {
	case syntax.LogicEQ:
		return left.Compare(right, compareTypeEq)
	case syntax.LogicNEQ:
		zb, err := left.Compare(right, compareTypeEq)
		return zb.Rev(), err
	case syntax.LogicIS:
		return left.Compare(right, compareTypeIs)
	case syntax.LogicISN:
		zb, err := left.Compare(right, compareTypeIs)
		return zb.Rev(), err
	case syntax.LogicGT:
		return left.Compare(right, compareTypeGt)
	case syntax.LogicGTE:
		zb1, err := left.Compare(right, compareTypeGt)
		if err != nil {
			return nil, err
		}
		zb2, err := left.Compare(right, compareTypeEq)
		if err != nil {
			return nil, err
		}

		return NewZnBool(zb1.Value || zb2.Value), nil
	case syntax.LogicLT:
		return left.Compare(right, compareTypeLt)
	case syntax.LogicLTE:
		zb1, err := left.Compare(right, compareTypeLt)
		if err != nil {
			return nil, err
		}
		zb2, err := left.Compare(right, compareTypeEq)
		if err != nil {
			return nil, err
		}

		return NewZnBool(zb1.Value || zb2.Value), nil
	default:
		return nil, error.InvalidCaseType()
	}
}

// eval prime expr
func evalPrimeExpr(ctx *Context, scope Scope, expr syntax.Expression) (ZnValue, *error.Error) {
	switch e := expr.(type) {
	case *syntax.Number:
		return NewZnDecimal(e.GetLiteral())
	case *syntax.String:
		return NewZnString(e.GetLiteral()), nil
	case *syntax.ID:
		vtag := e.GetLiteral()
		return getValue(ctx, scope, vtag)
	case *syntax.ArrayExpr:
		znObjs := []ZnValue{}
		for _, item := range e.Items {
			expr, err := evalExpression(ctx, scope, item)
			if err != nil {
				return nil, err
			}
			znObjs = append(znObjs, expr)
		}

		return NewZnArray(znObjs), nil
	case *syntax.HashMapExpr:
		znPairs := []KVPair{}
		for _, item := range e.KVPair {
			expr, err := evalExpression(ctx, scope, item.Key)
			if err != nil {
				return nil, err
			}
			exprKey, ok := expr.(*ZnString)
			if !ok {
				return nil, error.InvalidExprType("string", "integer")
			}
			exprVal, err := evalExpression(ctx, scope, item.Value)
			if err != nil {
				return nil, err
			}
			znPairs = append(znPairs, KVPair{
				Key:   exprKey.Value,
				Value: exprVal,
			})
		}
		return NewZnHashMap(znPairs), nil
	default:
		return nil, error.InvalidCaseType()
	}
}

// eval var assign
func evalVarAssignExpr(ctx *Context, scope Scope, expr *syntax.VarAssignExpr) (ZnValue, *error.Error) {
	// Right Side
	val, err := evalExpression(ctx, scope, expr.AssignExpr)
	if err != nil {
		return nil, err
	}

	// Left Side
	switch v := expr.TargetVar.(type) {
	case *syntax.ID:
		// set ID
		vtag := v.GetLiteral()
		err2 := setValue(ctx, scope, vtag, val)
		return val, err2
	case *syntax.MemberExpr:
		iv, err := getMemberExprIV(ctx, scope, v)
		if err != nil {
			return nil, err
		}
		return iv.Reduce(ctx, val, true)
	default:
		return nil, error.InvalidCaseType()
	}
}

func getMemberExprIV(ctx *Context, scope Scope, expr *syntax.MemberExpr) (ZnIV, *error.Error) {
	if expr.IsSelfRoot { // 此之 XX
		switch expr.MemberType {
		case syntax.MemberID:
			tag := expr.MemberID.Literal
			return &ZnScopeMemberIV{scope, tag}, nil
		case syntax.MemberMethod:
			m := expr.MemberMethod
			funcName := m.FuncName.Literal
			paramVals, err := exprsToValues(ctx, scope, m.Params)
			if err != nil {
				return nil, err
			}
			return &ZnScopeMethodIV{scope, funcName, paramVals}, nil
		}
		return nil, error.NewErrorSLOT("unsupport memberType (should not throw)")
	}

	// IsSelfRoot = false (with root)
	valRoot, err := evalExpression(ctx, scope, expr.Root)
	if err != nil {
		return nil, err
	}
	switch expr.MemberType {
	case syntax.MemberID: // A 之 B
		tag := expr.MemberID.Literal
		return &ZnMemberIV{valRoot, tag}, nil
	case syntax.MemberMethod:
		m := expr.MemberMethod
		funcName := m.FuncName.Literal
		paramVals, err := exprsToValues(ctx, scope, m.Params)
		if err != nil {
			return nil, err
		}
		return &ZnMethodIV{valRoot, funcName, paramVals}, nil
	case syntax.MemberIndex:
		idx, err := evalExpression(ctx, scope, expr.MemberIndex)
		if err != nil {
			return nil, err
		}
		switch v := valRoot.(type) {
		case *ZnArray:
			vr, ok := idx.(*ZnDecimal)
			if !ok {
				return nil, error.InvalidExprType("integer")
			}
			return &ZnArrayIV{v, vr}, nil
		case *ZnHashMap:
			var s *ZnString
			switch x := idx.(type) {
			// regard decimal value directly as string
			case *ZnDecimal:
				// transform decimal value to string
				// x.exp < 0 express that its a decimal value with point mark, not an integer
				if x.exp < 0 {
					return nil, error.InvalidExprType("integer", "string")
				}
				s = NewZnString(x.String())
			case *ZnString:
				s = x
			default:
				return nil, error.InvalidExprType("integer", "string")
			}
			return &ZnHashMapIV{v, s}, nil
		default:
			return nil, error.InvalidExprType("array", "hashmap")
		}
	}

	return nil, error.NewErrorSLOT("unsupport memberType (should not throw)")
}

//// scope value setters/getters
func getValue(ctx *Context, scope Scope, name string) (ZnValue, *error.Error) {
	// find on globals first
	if symVal, inGlobals := ctx.globals[name]; inGlobals {
		return symVal, nil
	}
	// ...then in symbols
	sp := scope
	for sp != nil {
		sym, ok := sp.GetSymbol(name)
		if ok {
			return sym.Value, nil
		}
		// if not found, search its parent
		sp = sp.GetParent()
	}
	return nil, error.NameNotDefined(name)
}

func setValue(ctx *Context, scope Scope, name string, value ZnValue) *error.Error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	// ...then in symbols
	sp := scope
	for sp != nil {
		sym, ok := sp.GetSymbol(name)
		if ok {
			if sym.IsConstant {
				return error.AssignToConstant()
			}
			sp.SetSymbol(name, value, false)
			return nil
		}
		// if not found, search its parent
		sp = sp.GetParent()
	}
	return error.NameNotDefined(name)
}

func bindValue(ctx *Context, scope Scope, name string, value ZnValue, isConstatnt bool) *error.Error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	// bind directly
	if _, ok := scope.GetSymbol(name); ok {
		return error.NameRedeclared(name)
	}
	sp := scope
	for sp != nil {
		if sp.HasSymbol() {
			sp.SetSymbol(name, value, isConstatnt)
			return nil
		}
		sp = sp.GetParent()
	}

	return nil
}

//// helpers

// exprsToValues - []syntax.Expression -> []eval.ZnValue
func exprsToValues(ctx *Context, scope Scope, exprs []syntax.Expression) ([]ZnValue, *error.Error) {
	params := []ZnValue{}
	for _, paramExpr := range exprs {
		pval, err := evalExpression(ctx, scope, paramExpr)
		if err != nil {
			return nil, err
		}
		params = append(params, pval)
	}
	return params, nil
}
