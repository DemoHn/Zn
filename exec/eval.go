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
	scope.SetCurrentLine(stmt.GetCurrentLine())
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
	case *syntax.FunctionReturnStmt:
		val, err := evalExpression(ctx, scope, v.ReturnExpr)
		if err != nil {
			return err
		}
		if fs, ok := scope.(*FuncScope); ok {
			fs.SetReturnFlag(true)
			fs.SetReturnValue(val)
		}
		return nil
	case syntax.Expression:
		val, err := evalExpression(ctx, scope, v)
		if err != nil {
			return err
		}
		if rs, ok := scope.(*RootScope); ok {
			rs.SetLastValue(val)
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
	loopScope := scope.NewScope(ctx, "while")
	// TODO: more handler on scope
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

//// execute expressions

func evalExpression(ctx *Context, scope Scope, expr syntax.Expression) (ZnValue, *error.Error) {
	scope.SetCurrentLine(expr.GetCurrentLine())
	switch e := expr.(type) {
	case *syntax.VarAssignExpr:
		return evalVarAssignExpr(ctx, scope, e)
	case *syntax.LogicExpr:
		if e.Type == syntax.LogicAND || e.Type == syntax.LogicOR {
			return evalLogicCombiner(ctx, scope, e)
		}
		return evalLogicComparator(ctx, scope, e)
	case *syntax.ArrayListIndexExpr:
		iv, err := getArrayListIV(ctx, scope, e)
		if err != nil {
			return nil, err
		}
		// regard iv as a RHS value
		return iv.Reduce(nil, false)
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
	fScope, _ := scope.NewScope(ctx, sTypeFunc).(*FuncScope)
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
	params := []ZnValue{}
	for _, paramExpr := range expr.Params {
		pval, err := evalExpression(ctx, scope, paramExpr)
		if err != nil {
			return nil, err
		}
		params = append(params, pval)
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

		var res ZnValue = NewZnNull()
		execBlock := zf.Node.ExecBlock
		// iterate block
		for _, stmt := range execBlock.Children {
			if err := evalStatement(ctx, scope, stmt); err != nil {
				return nil, err
			}
			// if returnFlag = true (after executing the statement)
			// return result immediately
			if scope.GetReturnFlag() {
				goto result
			}
		}

	result:
		res = scope.GetReturnValue()
		return res, nil
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
	case *syntax.ArrayListIndexExpr:
		iv, err := getArrayListIV(ctx, scope, v)
		if err != nil {
			return nil, err
		}
		return iv.Reduce(val, true)
	default:
		return nil, error.InvalidCaseType()
	}
}

func getArrayListIV(ctx *Context, scope Scope, expr *syntax.ArrayListIndexExpr) (ZnIV, *error.Error) {
	// val # index  --> 【1，２，３】#2
	val, err := evalExpression(ctx, scope, expr.Root)
	if err != nil {
		return nil, err
	}
	idx, err := evalExpression(ctx, scope, expr.Index)
	if err != nil {
		return nil, err
	}
	switch v := val.(type) {
	case *ZnArray:
		vr, ok := idx.(*ZnDecimal)
		if !ok {
			return nil, error.InvalidExprType("integer")
		}
		return &ZnArrayIV{v, vr}, nil
	case *ZnHashMap:
		var s *ZnString
		switch x := idx.(type) {
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
	scope.SetSymbol(name, value, isConstatnt)
	return nil
}
