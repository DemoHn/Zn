package exec

import (
	"reflect"
	"strconv"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/syntax"
)

type compareVerb uint8

// Define compareVerbs, for details of each verb, check the following comments
// on compareValues() function.
const (
	CmpEq compareVerb = 1
	CmpLt compareVerb = 2
	CmpGt compareVerb = 3
)

// eval.go evaluates program from generated AST tree with specific scopes
// common signature of eval functions:
//
// evalXXXXStmt(ctx *Context, node Node) *error.Error
//
// or
//
// evalXXXXExpr(ctx *Context, node Node) (Value, *error.Error)
//
// NOTICE:
// `evalXXXXStmt` will change the value of its corresponding scope; However, `evalXXXXExpr` will export
// a Value object and mostly won't change scopes (but search a variable from scope is frequently used)

// duplicateValue - deepcopy values' structure, including bool, string, decimal, array, hashmap
// for function or object or null, pass the original reference instead.
// This is due to the 'copycat by default' policy

func evalProgram(ctx *Context, program *syntax.Program) *error.Error {
	return evalStmtBlock(ctx, program.Content)
}

//// eval statements

// EvalStatement - eval statement
func evalStatement(ctx *Context, stmt syntax.Statement) *error.Error {
	// TODO: set current line
	switch v := stmt.(type) {
	case *syntax.VarDeclareStmt:
		return evalVarDeclareStmt(ctx, v)
	case *syntax.WhileLoopStmt:
		return evalWhileLoopStmt(ctx, v)
	case *syntax.BranchStmt:
		return evalBranchStmt(ctx, v)
	case *syntax.EmptyStmt:
		return nil
	case *syntax.FunctionDeclareStmt:
		fn := BuildFunctionFromNode(v)
		return bindValue(ctx, v.FuncName.GetLiteral(), fn)
	case *syntax.ClassDeclareStmt:
		if ctx.scope.parent != nil {
			return error.NewErrorSLOT("只能在代码主层级定义类")
		}
		return bindClassRef(ctx, v)
	case *syntax.IterateStmt:
		return evalIterateStmt(ctx, v)
	case *syntax.FunctionReturnStmt:
		val, err := evalExpression(ctx, v.ReturnExpr)
		if err != nil {
			return err
		}
		// send RETURN break
		return error.ReturnBreakError(val)
	case syntax.Expression:
		resetLastValue = false
		val, err := evalExpression(ctx, v)
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
		return error.UnExpectedCase("语句类型", reflect.TypeOf(v).Name())
	}
}

// evalVarDeclareStmt - consists of three branches:
// 1. A，B 为 C
// 2. A，B 成为 X：P1，P2，...
// 3. A，B 恒为 C
func evalVarDeclareStmt(ctx *Context, node *syntax.VarDeclareStmt) *error.Error {
	for _, vpair := range node.AssignPair {
		switch vpair.Type {
		case syntax.VDTypeAssign, syntax.VDTypeAssignConst: // 为，恒为
			obj, err := evalExpression(ctx, vpair.AssignExpr)
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
					obj = duplicateValue(obj)
				}

				if err := bindValueDecl(ctx, vtag, obj, isConst); err != nil {
					return err
				}
			}
		case syntax.VDTypeObjNew: // 成为
			if err := evalNewObjectPart(ctx, vpair); err != nil {
				return err
			}
		}
	}
	return nil
}

// eval A,B 成为 C：P1，P2，P3，...
// ensure VDAssignPair.Type MUST BE syntax.VDTypeObjNew
func evalNewObjectPart(ctx *Context, node syntax.VDAssignPair) *error.Error {
	vtag := node.ObjClass.GetLiteral()
	// get class definition
	classRef, err := getClassRef(ctx, vtag)
	if err != nil {
		return err
	}

	cParams, err := exprsToValues(ctx, node.ObjParams)
	if err != nil {
		return err
	}

	// assign new object to variables
	for _, v := range node.Variables {
		vtag := v.GetLiteral()
		// compose a new object instance
		fScope := NewFuncScope(scope, nil)
		finalObj, err := classRef.Construct(ctx, cParams)
		if err != nil {
			return err
		}

		if bindValue(ctx, vtag, finalObj); err != nil {
			return err
		}
	}
	return nil
}

// evalWhileLoopStmt -
func evalWhileLoopStmt(ctx *Context, node *syntax.WhileLoopStmt) *error.Error {
	loopScope := NewWhileScope(scope)
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
func evalStmtBlock(ctx *Context, block *syntax.BlockStmt) *error.Error {
	enableHoist := false
	rootScope, ok := scope.(*RootScope)
	if ok {
		enableHoist = true
	}

	if enableHoist {
		// ROUND I: declare function stmt FIRST
		for _, stmtI := range block.Children {
			switch v := stmtI.(type) {
			case *syntax.FunctionDeclareStmt:
				fn := BuildZnFunctionFromNode(v)
				if err := bindValue(ctx, v.FuncName.GetLiteral(), fn); err != nil {
					return err
				}
			case *syntax.ClassDeclareStmt:
				if err := bindClassRef(ctx, rootScope, v); err != nil {
					return err
				}
			}
		}
		// ROUND II: exec statement except functionDecl stmt
		for _, stmtII := range block.Children {
			switch stmtII.(type) {
			case *syntax.FunctionDeclareStmt, *syntax.ClassDeclareStmt:
				continue
			default:
				if err := evalStatement(ctx, stmtII); err != nil {
					return err
				}
			}
		}
	} else {
		for _, stmt := range block.Children {
			if err := evalStatement(ctx, stmt); err != nil {
				return err
			}
		}
	}
	return nil
}

func evalBranchStmt(ctx *Context, node *syntax.BranchStmt) *error.Error {
	// #1. condition header
	ifExpr, err := evalExpression(ctx, node.IfTrueExpr)
	if err != nil {
		return err
	}
	vIfExpr, ok := ifExpr.(*Bool)
	if !ok {
		return error.InvalidExprType("bool")
	}
	// exec if-branch
	if vIfExpr.Value == true {
		return evalStmtBlock(ctx, node.IfTrueBlock)
	}
	// exec else-if branches
	for idx, otherExpr := range node.OtherExprs {
		otherExprI, err := evalExpression(ctx, otherExpr)
		if err != nil {
			return err
		}
		vOtherExprI, ok := otherExprI.(*Bool)
		if !ok {
			return error.InvalidExprType("bool")
		}
		// exec else-if branch
		if vOtherExprI.Value == true {
			return evalStmtBlock(ctx, node.OtherBlocks[idx])
		}
	}
	// exec else branch if possible
	if node.HasElse == true {
		return evalStmtBlock(ctx, node.IfFalseBlock)
	}
	return nil
}

func evalIterateStmt(ctx *Context, node *syntax.IterateStmt) *error.Error {
	// pre-defined key, value variable name
	var keySlot, valueSlot string
	var nameLen = len(node.IndexNames)

	iterScope := NewIterateScope(scope)
	// 以A，B遍历C： D
	// execute expr: C
	targetExpr, err := evalExpression(ctx, node.IterateExpr)
	if err != nil {
		return err
	}

	// execIterationBlock, including set "currentKey" and "currentValue" to scope,
	// and preDefined indication variables
	execIterationBlockFn := func(key Value, val Value) *error.Error {
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
		if err := bindValue(ctx, iterScope, valueSlot, NewZnNull()); err != nil {
			return err
		}
	} else if nameLen == 2 {
		keySlot = node.IndexNames[0].Literal
		valueSlot = node.IndexNames[1].Literal
		if err := bindValue(ctx, iterScope, keySlot, NewZnNull()); err != nil {
			return err
		}
		if err := bindValue(ctx, iterScope, valueSlot, NewZnNull()); err != nil {
			return err
		}
	} else if nameLen > 2 {
		return error.MostParamsError(2)
	}

	// execute iterations
	switch tv := targetExpr.(type) {
	case *ZnArray:
		for idx, val := range tv.Value {
			idxVar := NewZnDecimalFromInt(idx, 0)
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
		for _, key := range tv.KeyOrder {
			val := tv.Value[key]
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

func evalExpression(ctx *Context, expr syntax.Expression) (Value, *error.Error) {
	scope.GetRoot().SetCurrentLine(expr.GetCurrentLine())
	switch e := expr.(type) {
	case *syntax.VarAssignExpr:
		return evalVarAssignExpr(ctx, e)
	case *syntax.LogicExpr:
		if e.Type == syntax.LogicAND || e.Type == syntax.LogicOR {
			return evalLogicCombiner(ctx, e)
		}
		return evalLogicComparator(ctx, e)
	case *syntax.MemberExpr:
		iv, err := getMemberExprIV(ctx, e)
		if err != nil {
			return nil, err
		}
		return iv.Reduce(ctx, nil, false)
	case *syntax.Number, *syntax.String, *syntax.ID, *syntax.ArrayExpr, *syntax.HashMapExpr:
		return evalPrimeExpr(ctx, e)
	case *syntax.FuncCallExpr:
		return evalFunctionCall(ctx, e)
	default:
		return nil, error.InvalidExprType()
	}
}

// （显示：A，B，C）
func evalFunctionCall(ctx *Context, expr *syntax.FuncCallExpr) (Value, *error.Error) {
	vtag := expr.FuncName.GetLiteral()
	var zf *ClosureRef

	// if current scope is FuncScope, find ID from funcScope's "targetThis" method list
	if sp, ok := scope.(*FuncScope); ok {
		targetThis := sp.GetTargetThis()
		if targetThis != nil {
			if val, err := targetThis.GetMethod(vtag); err == nil {
				zf = val
			}
		}
	}

	// if function value not found from object scope, look up from local scope
	if zf == nil {
		// find function definction
		val, err := getValue(ctx, vtag)
		if err != nil {
			return nil, err
		}
		// assert value
		zval, ok := val.(*Function)
		if !ok {
			return nil, error.InvalidFuncVariable(vtag)
		}
		zf = zval.ClosureRef
	}

	// exec params
	params, err := exprsToValues(ctx, expr.Params)
	if err != nil {
		return nil, err
	}

	fScope := NewFuncScope(scope, nil)
	// exec function call via its ClosureRef
	return zf.Exec(ctx, fScope, params)
}

// evaluate logic combination expressions
// such as A 且 B
// or A 或 B
func evalLogicCombiner(ctx *Context, expr *syntax.LogicExpr) (*Bool, *error.Error) {
	logicType := expr.Type
	// #1. eval left
	left, err := evalExpression(ctx, expr.LeftExpr)
	if err != nil {
		return nil, err
	}
	// #2. assert left expr type to be ZnBool
	vleft, ok := left.(*Bool)
	if !ok {
		return nil, error.InvalidExprType("bool")
	}
	// #3. check if the result could be retrieved earlier
	//
	// 1) for Y = A and B, if A = false, then Y must be false
	// 2) for Y = A or  B, if A = true, then Y must be true
	//
	// for those cases, we can yield result directly
	if logicType == syntax.LogicAND && vleft.value == false {
		return NewBool(false), nil
	}
	if logicType == syntax.LogicOR && vleft.value == true {
		return NewBool(true), nil
	}
	// #4. eval right
	right, err := evalExpression(ctx, expr.RightExpr)
	if err != nil {
		return nil, err
	}
	vright, ok := right.(*Bool)
	if !ok {
		return nil, error.InvalidExprType("bool")
	}
	// then evalute data
	switch logicType {
	case syntax.LogicAND:
		return NewBool(vleft.value && vright.value), nil
	default: // logicOR
		return NewBool(vleft.value || vright.value), nil
	}
}

// evaluate logic comparator
// ensure both expressions are comparable (i.e. subtype of ZnComparable)
func evalLogicComparator(ctx *Context, expr *syntax.LogicExpr) (*Bool, *error.Error) {
	logicType := expr.Type
	// #1. eval left
	left, err := evalExpression(ctx, expr.LeftExpr)
	if err != nil {
		return nil, err
	}
	// #2. eval right
	right, err := evalExpression(ctx, expr.RightExpr)
	if err != nil {
		return nil, err
	}

	var cmpRes bool
	var cmpErr *error.Error
	// #3. do comparison
	switch logicType {
	case syntax.LogicEQ:
		cmpRes, cmpErr = compareValues(left, right, CmpEq)
	case syntax.LogicNEQ:
		cmpRes, cmpErr = compareValues(left, right, CmpEq)
		cmpRes = !cmpRes // reverse result
	case syntax.LogicGT:
		cmpRes, cmpErr = compareValues(left, right, CmpGt)
	case syntax.LogicGTE:
		var cmp1, cmp2 bool
		cmp1, cmpErr = compareValues(left, right, CmpGt)
		if cmpErr != nil {
			return nil, cmpErr
		}
		cmp2, cmpErr = compareValues(left, right, CmpEq)
		cmpRes = cmp1 || cmp2
	case syntax.LogicLT:
		cmpRes, cmpErr = compareValues(left, right, CmpLt)
	case syntax.LogicLTE:
		var cmp1, cmp2 bool
		cmp1, cmpErr = compareValues(left, right, CmpLt)
		if cmpErr != nil {
			return nil, cmpErr
		}
		cmp2, cmpErr = compareValues(left, right, CmpEq)
		cmpRes = cmp1 || cmp2
	default:
		return nil, error.UnExpectedCase("比较类型", strconv.Itoa(int(logicType)))
	}

	return NewZnBool(cmpRes), cmpErr
}

// eval prime expr
func evalPrimeExpr(ctx *Context, expr syntax.Expression) (Value, *error.Error) {
	switch e := expr.(type) {
	case *syntax.Number:
		return NewZnDecimal(e.GetLiteral())
	case *syntax.String:
		return NewZnString(e.GetLiteral()), nil
	case *syntax.ID:
		vtag := e.GetLiteral()
		return getValue(ctx, vtag)
	case *syntax.ArrayExpr:
		znObjs := []Value{}
		for _, item := range e.Items {
			expr, err := evalExpression(ctx, item)
			if err != nil {
				return nil, err
			}
			znObjs = append(znObjs, expr)
		}

		return NewZnArray(znObjs), nil
	case *syntax.HashMapExpr:
		znPairs := []KVPair{}
		for _, item := range e.KVPair {
			expr, err := evalExpression(ctx, item.Key)
			if err != nil {
				return nil, err
			}
			exprKey, ok := expr.(*ZnString)
			if !ok {
				return nil, error.InvalidExprType("string")
			}
			exprVal, err := evalExpression(ctx, item.Value)
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
		return nil, error.UnExpectedCase("表达式类型", reflect.TypeOf(e).Name())
	}
}

// eval variable assign
func evalVarAssignExpr(ctx *Context, expr *syntax.VarAssignExpr) (Value, *error.Error) {
	// Right Side
	val, err := evalExpression(ctx, expr.AssignExpr)
	if err != nil {
		return nil, err
	}
	// if var assignment is NOT by reference, then duplicate value
	if !expr.RefMark {
		val = duplicateValue(val)
	}

	// Left Side
	switch v := expr.TargetVar.(type) {
	case *syntax.ID:
		// set ID
		vtag := v.GetLiteral()
		err2 := setValue(ctx, vtag, val)
		return val, err2
	case *syntax.MemberExpr:
		iv, err := getMemberExprIV(ctx, v)
		if err != nil {
			return nil, err
		}
		return iv.Reduce(ctx, val, true)
	default:
		return nil, error.UnExpectedCase("被赋值", reflect.TypeOf(v).Name())
	}
}

func getMemberExprIV(ctx *Context, expr *syntax.MemberExpr) (ZnIV, *error.Error) {
	if expr.RootType == syntax.RootTypeScope { // 此之 XX
		switch expr.MemberType {
		case syntax.MemberID:
			tag := expr.MemberID.Literal
			return &ZnScopeMemberIV{tag}, nil
		case syntax.MemberMethod:
			m := expr.MemberMethod
			funcName := m.FuncName.Literal
			paramVals, err := exprsToValues(ctx, m.Params)
			if err != nil {
				return nil, err
			}
			return &ZnScopeMethodIV{funcName, paramVals}, nil
		}
		return nil, error.UnExpectedCase("子项类型", strconv.Itoa(int(expr.MemberType)))
	}

	if expr.RootType == syntax.RootTypeProp { // 其 XX
		if expr.MemberType == syntax.MemberID {
			tag := expr.MemberID.Literal
			return &ZnPropIV{tag}, nil
		}
		return nil, error.UnExpectedCase("子项类型", strconv.Itoa(int(expr.MemberType)))
	}

	// RootType = RootTypeExpr
	valRoot, err := evalExpression(ctx, expr.Root)
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
		paramVals, err := exprsToValues(ctx, m.Params)
		if err != nil {
			return nil, err
		}

		return &ZnMethodIV{valRoot, funcName, paramVals}, nil
	case syntax.MemberIndex:
		idx, err := evalExpression(ctx, expr.MemberIndex)
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
	return nil, error.UnExpectedCase("子项类型", reflect.TypeOf(expr.MemberType).Name())
}

//// scope value setters/getters
func getValue(ctx *Context, name string) (Value, *error.Error) {
	// find on globals first
	if symVal, inGlobals := ctx.globals[name]; inGlobals {
		return symVal, nil
	}
	// ...then in symbols
	sp := ctx.scope
	for sp != nil {
		sym, ok := sp.symbolMap[name]
		if ok {
			return sym.value, nil
		}
		// if not found, search its parent
		sp = sp.parent
	}
	return nil, error.NameNotDefined(name)
}

func setValue(ctx *Context, name string, value Value) *error.Error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	// ...then in symbols
	sp := ctx.scope
	for sp != nil {
		sym, ok := sp.symbolMap[name]
		if ok {
			if sym.isConst {
				return error.AssignToConstant()
			}
			sp.symbolMap[name] = SymbolInfo{value, false}
			return nil
		}
		// if not found, search its parent
		sp = sp.parent
	}
	return error.NameNotDefined(name)
}

func getClassRef(ctx *Context, scope *RootScope, name string) (*ClassRef, *error.Error) {
	ref, ok := scope.classRefMap[name]
	if ok {
		return ref, nil
	}
	return nil, error.NameNotDefined(name)
}

func bindClassRef(ctx *Context, scope *RootScope, classStmt *syntax.ClassDeclareStmt) *error.Error {
	name := classStmt.ClassName.GetLiteral()
	_, ok := scope.classRefMap[name]
	if ok {
		return error.NameRedeclared(name)
	}
	scope.classRefMap[name] = BuildClassRefFromNode(name, classStmt)
	return nil
}

// bind non-const value with re-declaration check on same scope
func bindValue(ctx *Context, name string, value Value) *error.Error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	// bind directly
	if _, ok := ctx.GetSymbol(name); ok {
		return error.NameRedeclared(name)
	}
	ctx.SetSymbol(name, value, false)
	return nil
}

// bind value for declaration statement - that variables could be re-bind.
func bindValueDecl(ctx *Context, name string, value Value, isConst bool) *error.Error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	ctx.SetSymbol(name, value, isConst)
	return nil
}

//// helpers

// exprsToValues - []syntax.Expression -> []eval.Value
func exprsToValues(ctx *Context, exprs []syntax.Expression) ([]Value, *error.Error) {
	params := []Value{}
	for _, paramExpr := range exprs {
		pval, err := evalExpression(ctx, paramExpr)
		if err != nil {
			return nil, err
		}
		params = append(params, pval)
	}
	return params, nil
}
