package exec

import (
	"fmt"
	"math/big"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/exec/val"
	"github.com/DemoHn/Zn/syntax"
)

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

// duplicateValue - deepcopy values' structure, including bool, string, decimal, array, hashmap
// for function or object or null, pass the original reference instead.
// This is due to the 'copycat by default' policy

func evalProgram(c *ctx.Context, program *syntax.Program) *error.Error {
	return evalStmtBlock(c, program.Content)
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
		return bindValue(c, v.FuncName.GetLiteral(), fn)
	case *syntax.ClassDeclareStmt:
		if sp.FindParentScope() != nil {
			return error.NewErrorSLOT("只能在代码主层级定义类")
		}
		return bindClassRef(c, v)
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
					obj = duplicateValue(obj)
				}

				if err := bindValueDecl(c, vtag, obj, isConst); err != nil {
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
	classRef, err := getClassRef(c, vtag)
	if err != nil {
		return err
	}

	cParams, err := exprsToValues(c, node.ObjParams)
	if err != nil {
		return err
	}

	// assign new object to variables
	for _, v := range node.Variables {
		vtag := v.GetLiteral()

		fctx := c.DuplicateNewScope()
		finalObj, err := classRef.Construct(fctx, cParams)
		if err != nil {
			return err
		}

		if bindValue(c, vtag, finalObj); err != nil {
			return err
		}
	}
	return nil
}

// evalWhileLoopStmt -
func evalWhileLoopStmt(c *ctx.Context, node *syntax.WhileLoopStmt) *error.Error {
	newScope := c.ShiftChildScope()
	newScope.SetSgValue(val.NewLoopCtl())
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
		if vTrueExpr.GetValue() == false {
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

// EvalStmtBlock -
func evalStmtBlock(c *ctx.Context, block *syntax.BlockStmt) *error.Error {
	enableHoist := false
	// only rootScope could enable hoist
	if c.GetScope().FindChildScope() == nil {
		enableHoist = true
	}

	if enableHoist {
		// ROUND I: declare function stmt FIRST
		for _, stmtI := range block.Children {
			switch v := stmtI.(type) {
			case *syntax.FunctionDeclareStmt:
				fn := BuildFunctionFromNode(v)
				if err := bindValue(c, v.FuncName.GetLiteral(), fn); err != nil {
					return err
				}
			case *syntax.ClassDeclareStmt:
				if err := bindClassRef(c, v); err != nil {
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
				if err := evalStatement(c, stmtII); err != nil {
					return err
				}
			}
		}
	} else {
		for _, stmt := range block.Children {
			if err := evalStatement(c, stmt); err != nil {
				return err
			}
		}
	}
	return nil
}

func evalBranchStmt(c *ctx.Context, node *syntax.BranchStmt) *error.Error {
	// #1. condition header
	ifExpr, err := evalExpression(c, node.IfTrueExpr)
	if err != nil {
		return err
	}
	vIfExpr, ok := ifExpr.(*val.Bool)
	if !ok {
		return error.InvalidExprType("bool")
	}

	c.ShiftChildScope()
	// exec if-branch
	if vIfExpr.GetValue() == true {
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
		if vOtherExprI.GetValue() == true {
			return evalStmtBlock(c, node.OtherBlocks[idx])
		}
	}
	// exec else branch if possible
	if node.HasElse == true {
		return evalStmtBlock(c, node.IfFalseBlock)
	}
	return nil
}

func evalIterateStmt(c *ctx.Context, node *syntax.IterateStmt) *error.Error {
	// pre-defined key, value variable name
	var keySlot, valueSlot string
	var nameLen = len(node.IndexNames)

	// 以A，B遍历C： D
	// execute expr: C
	targetExpr, err := evalExpression(c, node.IterateExpr)
	if err != nil {
		return err
	}

	// shift child scope
	newScope := c.ShiftChildScope()
	newScope.SetSgValue(val.NewLoopCtl())
	// execIterationBlock, including set "currentKey" and "currentValue" to scope,
	// and preDefined indication variables
	execIterationBlockFn := func(key ctx.Value, v ctx.Value) *error.Error {
		// set values of 此之值 and 此之
		sgValueT := newScope.GetSgValue()
		sgValue, _ := sgValueT.(*val.LoopCtl)
		sgValue.SetCurrentKeyValue(key, v)

		// set pre-defined value
		if nameLen == 1 {
			if err := setValue(c, valueSlot, v); err != nil {
				return err
			}
		} else if nameLen == 2 {
			if err := setValue(c, keySlot, key); err != nil {
				return err
			}
			if err := setValue(c, valueSlot, v); err != nil {
				return err
			}
		}
		return evalStmtBlock(c, node.IterateBlock)
	}

	// define indication variables as "currentKey" and "currentValue" under new iterScope
	// of course since there's no any iteration is executed yet, the initial values are all "Null"
	if nameLen == 1 {
		valueSlot = node.IndexNames[0].Literal
		if err := bindValue(c, valueSlot, val.NewNull()); err != nil {
			return err
		}
	} else if nameLen == 2 {
		keySlot = node.IndexNames[0].Literal
		valueSlot = node.IndexNames[1].Literal
		if err := bindValue(c, keySlot, val.NewNull()); err != nil {
			return err
		}
		if err := bindValue(c, valueSlot, val.NewNull()); err != nil {
			return err
		}
	} else if nameLen > 2 {
		return error.MostParamsError(2)
	}

	// execute iterations
	switch tv := targetExpr.(type) {
	case *val.Array:
		for idx, v := range tv.value {
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
		for _, key := range tv.keyOrder {
			val := tv.value[key]
			keyVar := NewString(key)
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
		// when MemberType = memberID or memberIndex, it's a typical "getter" expression;
		// when MemberType = memberMethod, it's a method call, so we could not use IV logic
		// to handle it
		switch e.MemberType {
		case syntax.MemberID, syntax.MemberIndex:
			iv, err := getMemberExprIV(c, e)
			if err != nil {
				return nil, err
			}
			return iv.ReduceRHS(c)
		case syntax.MemberMethod:
			// get root expr
			var rootValue ctx.Value
			switch e.RootType {
			case syntax.RootTypeExpr:
				root, err := evalExpression(c, e.Root)
				if err != nil {
					return nil, err
				}
				rootValue = root
			case syntax.RootTypeScope:
				sgValue, err := findSgValue(c)
				if err != nil {
					return nil, err
				}
				rootValue = sgValue
			default: // 其他 rootType 不支持
				return nil, error.UnExpectedCase("根元素类型", fmt.Sprintf("%d", e.MemberType))
			}

			// execute method
			methodName := e.MemberMethod.FuncName.GetLiteral()
			paramValues, err := exprsToValues(c, e.MemberMethod.Params)
			if err != nil {
				return nil, err
			}
			fctx := c.DuplicateNewScope()
			fctx.scope.thisValue = rootValue
			return rootValue.ExecMethod(fctx, methodName, paramValues)
		}
		return nil, error.UnExpectedCase("成员类型", fmt.Sprintf("%d", e.MemberType))
	case *syntax.Number, *syntax.String, *syntax.ID, *syntax.ArrayExpr, *syntax.HashMapExpr:
		return evalPrimeExpr(c, e)
	case *syntax.FuncCallExpr:
		return evalFunctionCall(c, e)
	default:
		return nil, error.InvalidExprType()
	}
}

// （显示：A，B，C）
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
	thisValue, _ := findThisValue(c)

	// if thisValue exists, find ID from its method list
	if thisValue != nil {
		if obj, ok := thisValue.(*val.Object); ok {
			// find value
			if method, ok2 := obj.ref.MethodList[vtag]; ok2 {
				zf = &method
			}
		}
	}

	// if function value not found from object scope, look up from local scope
	if zf == nil {
		// find function definction
		v, err := getValue(c, vtag)
		if err != nil {
			return nil, err
		}
		// assert value
		zval, ok := v.(*val.Function)
		if !ok {
			return nil, error.InvalidFuncVariable(vtag)
		}
		zf = &zval.value
	}

	// exec params
	params, err := exprsToValues(c, expr.Params)
	if err != nil {
		return nil, err
	}

	// exec function call via its ClosureRef
	fctx := c.DuplicateNewScope()
	return zf.Exec(fctx, params)
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
	if logicType == syntax.LogicAND && vleft.GetValue() == false {
		return val.NewBool(false), nil
	}
	if logicType == syntax.LogicOR && vleft.GetValue() == true {
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
		return getValue(c, vtag)
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
	val, err := evalExpression(c, expr.AssignExpr)
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
		err2 := setValue(c, vtag, val)
		return val, err2
	case *syntax.MemberExpr:
		if v.MemberType == syntax.MemberID || v.MemberType == syntax.MemberIndex {
			iv, err := getMemberExprIV(c, v)
			if err != nil {
				return nil, err
			}
			return val, iv.ReduceLHS(c, val)
		}
		return nil, error.NewErrorSLOT("方法不能被赋值")
	default:
		return nil, error.UnExpectedCase("被赋值", fmt.Sprintf("%T", v))
	}
}

func getMemberExprIV(c *ctx.Context, expr *syntax.MemberExpr) (*val.IV, *error.Error) {
	switch expr.RootType {
	case syntax.RootTypeScope: // 此之 XX
		sgValue, err := findSgValue(c)
		if err != nil {
			return nil, err
		}

		return &val.IV{
			reduceType: val.IVTypeMember,
			root:       sgValue,
			member:     expr.MemberID.GetLiteral(),
		}, nil

	case syntax.RootTypeProp: // 其 XX
		thisValue, err := findThisValue(c)
		if err != nil {
			return nil, err
		}
		return &val.IV{
			reduceType: val.IVTypeMember,
			root:       thisValue,
			member:     expr.MemberID.GetLiteral(),
		}, nil
	case syntax.RootTypeExpr: // A 之 B
		valRoot, err := evalExpression(c, expr.Root)
		if err != nil {
			return nil, err
		}
		switch expr.MemberType {
		case syntax.MemberID: // A 之 B
			return &val.IV{
				reduceType: val.IVTypeMember,
				root:       valRoot,
				member:     expr.MemberID.GetLiteral(),
			}, nil
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
				return &val.IV{
					reduceType: val.IVTypeArray,
					root:       v,
					index:      vri,
				}, nil
			case *val.HashMap:
				var s string
				switch x := idx.(type) {
				// regard decimal value directly as string
				case *val.Decimal:
					// transform decimal value to string
					// x.exp < 0 express that its a decimal value with point mark, not an integer
					if x.exp < 0 {
						return nil, error.InvalidExprType("integer", "string")
					}
					s = x.String()
				case *val.String:
					s = x.String()
				default:
					return nil, error.InvalidExprType("integer", "string")
				}
				return &val.IV{
					reduceType: IVTypeHashMap,
					root:       v,
					member:     s,
				}, nil
			}
			return nil, error.InvalidExprType("array", "hashmap")
		}
		return nil, error.UnExpectedCase("子项类型", fmt.Sprintf("%d", expr.MemberType))
	}

	return nil, error.UnExpectedCase("根元素类型", fmt.Sprintf("%d", expr.RootType))
}

//// scope value setters/getters
func getValue(c *ctx.Context, name string) (ctx.Value, *error.Error) {
	// find on globals first
	if symVal, inGlobals := c.globals[name]; inGlobals {
		return symVal, nil
	}
	// ...then in symbols
	sp := c.scope
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

func setValue(c *ctx.Context, name string, value ctx.Value) *error.Error {
	if _, inGlobals := c.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	// ...then in symbols
	sp := c.scope
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

func getClassRef(c *ctx.Context, name string) (*ClassRef, *error.Error) {
	ref, ok := c.scope.classRefMap[name]
	if ok {
		return &ref, nil
	}
	return nil, error.NameNotDefined(name)
}

func bindClassRef(c *ctx.Context, classStmt *syntax.ClassDeclareStmt) *error.Error {
	name := classStmt.ClassName.GetLiteral()
	_, ok := c.scope.classRefMap[name]
	if ok {
		return error.NameRedeclared(name)
	}
	c.scope.classRefMap[name] = BuildClassFromNode(name, classStmt)
	return nil
}

// bind non-const value with re-declaration check on same scope
func bindValue(c *ctx.Context, name string, value ctx.Value) *error.Error {
	if _, inGlobals := c.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	// bind directly
	if c.scope != nil {
		if _, ok := c.scope.symbolMap[name]; ok {
			return error.NameRedeclared(name)
		}
		// set value
		c.scope.symbolMap[name] = SymbolInfo{value, false}
	}
	return nil
}

// bind value for declaration statement - that variables could be re-bind.
func bindValueDecl(c *ctx.Context, name string, value ctx.Value, isConst bool) *error.Error {
	if _, inGlobals := c.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	if c.scope != nil {
		c.scope.symbolMap[name] = SymbolInfo{value, isConst}
	}
	return nil
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

/// findSgValue - find the suitable sgValue in current context
// Rules:
//
// - if sgValue in current scope (c.scope.sgValue) != nil, returns the current one;
// - if sgValue in current scope == nil, then look up its parent util to the root;
func findSgValue(c *ctx.Context) (ctx.Value, *error.Error) {
	sp := c.scope
	for sp != nil {
		sgValue := sp.sgValue
		if sgValue != nil {
			return sgValue, nil
		}

		// otherwise, find sgValue from parent scope
		sp = sp.parent
	}

	return nil, error.PropertyNotFound("sgValue")
}

// findThisValue - similar with findSgValue(c), it looks up for nearest valid
// thisValue value.
func findThisValue(c *ctx.Context) (ctx.Value, *error.Error) {
	sp := c.scope
	for sp != nil {
		thisValue := sp.thisValue
		if thisValue != nil {
			return thisValue, nil
		}

		// otherwise, find thisValue from parent scope
		sp = sp.parent
	}

	return nil, error.PropertyNotFound("thisValue")
}

// duplicateValue - deepcopy values' structure, including bool, string, decimal, array, hashmap
// for function or object or null, pass the original reference instead.
// This is due to the 'copycat by default' policy
func duplicateValue(in ctx.Value) ctx.Value {
	switch v := in.(type) {
	case *val.Bool:
		return val.NewBool(v.value)
	case *val.String:
		return val.NewString(v.value)
	case *val.Decimal:
		x := new(big.Int)
		return &val.Decimal{
			co:  x.Set(v.co),
			exp: v.exp,
		}
	case *val.Null:
		return in // no need to copy since all "NULL" values are same
	case *val.Array:
		newArr := []ctx.Value{}
		for _, val := range v.value {
			newArr = append(newArr, duplicateValue(val))
		}
		return val.NewArray(newArr)
	case *val.HashMap:
		kvPairs := []val.KVPair{}
		for _, key := range v.keyOrder {
			dupVal := duplicateValue(v.value[key])
			kvPairs = append(kvPairs, KVPair{key, dupVal})
		}
		return val.NewHashMap(kvPairs)
	case *val.Function: // function itself is immutable, so return directly
		return in
	case *val.Object: // we don't copy object value at all
		return in
	}
	return in
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
				if err := bindValue(c, v.FuncName.GetLiteral(), fn); err != nil {
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
		return ctx.scope.returnValue, nil
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
				paramVal = duplicateValue(paramVal)
			}
			paramName := param.ID.GetLiteral()
			if err := bindValue(ctx, paramName, paramVal); err != nil {
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
		// init prop list
		for _, propPair := range classNode.PropertyList {
			propID := propPair.PropertyID.GetLiteral()
			expr, err := evalExpression(c, propPair.InitValue)
			if err != nil {
				return nil, err
			}
			obj.propList[propID] = expr
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
				objParamVal = duplicateValue(objParamVal)
			}
			paramName := param.ID.GetLiteral()
			obj.propList[paramName] = objParamVal
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
	return val.NewFunction("", closureRef.Exec)
}
