package exec

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

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
	var returnValue Value
	// set currentLine
	ctx.scope.fileInfo.currentLine = stmt.GetCurrentLine()

	// set return value
	defer func() {
		var finalReturnValue Value = NewNull()
		// set current return value
		if returnValue != nil {
			finalReturnValue = returnValue
		}
		ctx.scope.returnValue = finalReturnValue

		// set parent return value
		if ctx.scope.parent != nil {
			ctx.scope.parent.returnValue = finalReturnValue
		}
	}()

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
		expr, err := evalExpression(ctx, v)
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
			if err := evalNewObject(ctx, vpair); err != nil {
				return err
			}
		}
	}
	return nil
}

// eval A,B 成为 C：P1，P2，P3，...
// ensure VDAssignPair.Type MUST BE syntax.VDTypeObjNew
func evalNewObject(ctx *Context, node syntax.VDAssignPair) *error.Error {
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

		fctx := ctx.DuplicateNewScope()
		finalObj, err := classRef.Construct(fctx, cParams)
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
	lctx := ctx.DuplicateNewScope()
	lctx.scope.sgValue = NewLoopCtl()
	for {
		// #1. first execute expr
		trueExpr, err := evalExpression(lctx, node.TrueExpr)
		if err != nil {
			return err
		}
		// #2. assert trueExpr to be Bool
		vTrueExpr, ok := trueExpr.(*Bool)
		if !ok {
			return error.InvalidExprType("bool")
		}
		// break the loop if expr yields not true
		if vTrueExpr.value == false {
			return nil
		}
		// #3. stmt block
		if err := evalStmtBlock(lctx, node.LoopBlock); err != nil {
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
	// only rootScope could enable hoist
	if ctx.scope.parent == nil {
		enableHoist = true
	}

	if enableHoist {
		// ROUND I: declare function stmt FIRST
		for _, stmtI := range block.Children {
			switch v := stmtI.(type) {
			case *syntax.FunctionDeclareStmt:
				fn := BuildFunctionFromNode(v)
				if err := bindValue(ctx, v.FuncName.GetLiteral(), fn); err != nil {
					return err
				}
			case *syntax.ClassDeclareStmt:
				if err := bindClassRef(ctx, v); err != nil {
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

	bctx := ctx.DuplicateNewScope()
	// exec if-branch
	if vIfExpr.value == true {
		return evalStmtBlock(bctx, node.IfTrueBlock)
	}
	// exec else-if branches
	for idx, otherExpr := range node.OtherExprs {
		otherExprI, err := evalExpression(bctx, otherExpr)
		if err != nil {
			return err
		}
		vOtherExprI, ok := otherExprI.(*Bool)
		if !ok {
			return error.InvalidExprType("bool")
		}
		// exec else-if branch
		if vOtherExprI.value == true {
			return evalStmtBlock(bctx, node.OtherBlocks[idx])
		}
	}
	// exec else branch if possible
	if node.HasElse == true {
		return evalStmtBlock(bctx, node.IfFalseBlock)
	}
	return nil
}

func evalIterateStmt(ctx *Context, node *syntax.IterateStmt) *error.Error {
	// pre-defined key, value variable name
	var keySlot, valueSlot string
	var nameLen = len(node.IndexNames)

	// 以A，B遍历C： D
	// execute expr: C
	targetExpr, err := evalExpression(ctx, node.IterateExpr)
	if err != nil {
		return err
	}

	// create new child scope
	ictx := ctx.DuplicateNewScope()
	ictx.scope.sgValue = NewLoopCtl()
	// execIterationBlock, including set "currentKey" and "currentValue" to scope,
	// and preDefined indication variables
	execIterationBlockFn := func(key Value, val Value) *error.Error {
		// set values of 此之值 and 此之
		sgValue, _ := ictx.scope.sgValue.(*LoopCtl)
		sgValue.SetCurrentKeyValue(key, val)

		// set pre-defined value
		if nameLen == 1 {
			if err := setValue(ictx, valueSlot, val); err != nil {
				return err
			}
		} else if nameLen == 2 {
			if err := setValue(ictx, keySlot, key); err != nil {
				return err
			}
			if err := setValue(ictx, valueSlot, val); err != nil {
				return err
			}
		}
		return evalStmtBlock(ictx, node.IterateBlock)
	}

	// define indication variables as "currentKey" and "currentValue" under new iterScope
	// of course since there's no any iteration is executed yet, the initial values are all "Null"
	if nameLen == 1 {
		valueSlot = node.IndexNames[0].Literal
		if err := bindValue(ictx, valueSlot, NewNull()); err != nil {
			return err
		}
	} else if nameLen == 2 {
		keySlot = node.IndexNames[0].Literal
		valueSlot = node.IndexNames[1].Literal
		if err := bindValue(ictx, keySlot, NewNull()); err != nil {
			return err
		}
		if err := bindValue(ictx, valueSlot, NewNull()); err != nil {
			return err
		}
	} else if nameLen > 2 {
		return error.MostParamsError(2)
	}

	// execute iterations
	switch tv := targetExpr.(type) {
	case *Array:
		for idx, val := range tv.value {
			idxVar := NewDecimalFromInt(idx, 0)
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
	case *HashMap:
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

func evalExpression(ctx *Context, expr syntax.Expression) (Value, *error.Error) {
	ctx.scope.fileInfo.currentLine = expr.GetCurrentLine()

	switch e := expr.(type) {
	case *syntax.VarAssignExpr:
		return evalVarAssignExpr(ctx, e)
	case *syntax.LogicExpr:
		if e.Type == syntax.LogicAND || e.Type == syntax.LogicOR {
			return evalLogicCombiner(ctx, e)
		}
		return evalLogicComparator(ctx, e)
	case *syntax.MemberExpr:
		// when MemberType = memberID or memberIndex, it's a typical "getter" expression;
		// when MemberType = memberMethod, it's a method call, so we could not use IV logic
		// to handle it
		switch e.MemberType {
		case syntax.MemberID, syntax.MemberIndex:
			iv, err := getMemberExprIV(ctx, e)
			if err != nil {
				return nil, err
			}
			return iv.ReduceRHS(ctx)
		case syntax.MemberMethod:
			// get root expr
			var rootValue Value
			switch e.RootType {
			case syntax.RootTypeExpr:
				root, err := evalExpression(ctx, e.Root)
				if err != nil {
					return nil, err
				}
				rootValue = root
			case syntax.RootTypeScope:
				sgValue, err := findSgValue(ctx)
				if err != nil {
					return nil, err
				}
				rootValue = sgValue
			default: // 其他 rootType 不支持
				return nil, error.UnExpectedCase("根元素类型", fmt.Sprintf("%d", e.MemberType))
			}

			// execute method
			methodName := e.MemberMethod.FuncName.GetLiteral()
			paramValues, err := exprsToValues(ctx, e.MemberMethod.Params)
			if err != nil {
				return nil, err
			}
			fctx := ctx.DuplicateNewScope()
			fctx.scope.thisValue = rootValue
			return rootValue.ExecMethod(fctx, methodName, paramValues)
		}
		return nil, error.UnExpectedCase("成员类型", fmt.Sprintf("%d", e.MemberType))
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
	var zf *ClosureRef
	vtag := expr.FuncName.GetLiteral()

	// for a function call, if thisValue NOT FOUND, that means the target closure is a FUNCTION
	// instead of a METHOD (which is defined on class definition statement)
	//
	// If thisValue != nil, we will attempt to find clsoure from its method list;
	// then look up from scope's values.
	//
	// If thisValue == nil, we will look up target closure from scope's values directly.
	thisValue, _ := findThisValue(ctx)

	// if thisValue exists, find ID from its method list
	if thisValue != nil {
		if obj, ok := thisValue.(*Object); ok {
			// find value
			if method, ok2 := obj.ref.MethodList[vtag]; ok2 {
				zf = &method
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
		zf = &zval.value
	}

	// exec params
	params, err := exprsToValues(ctx, expr.Params)
	if err != nil {
		return nil, err
	}

	// exec function call via its ClosureRef
	fctx := ctx.DuplicateNewScope()
	return zf.Exec(fctx, params)
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
		return nil, error.UnExpectedCase("比较类型", fmt.Sprintf("%d", logicType))
	}

	return NewBool(cmpRes), cmpErr
}

// eval prime expr
func evalPrimeExpr(ctx *Context, expr syntax.Expression) (Value, *error.Error) {
	switch e := expr.(type) {
	case *syntax.Number:
		return NewDecimal(e.GetLiteral())
	case *syntax.String:
		return NewString(e.GetLiteral()), nil
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

		return NewArray(znObjs), nil
	case *syntax.HashMapExpr:
		znPairs := []KVPair{}
		for _, item := range e.KVPair {
			expr, err := evalExpression(ctx, item.Key)
			if err != nil {
				return nil, err
			}
			exprKey, ok := expr.(*String)
			if !ok {
				return nil, error.InvalidExprType("string")
			}
			exprVal, err := evalExpression(ctx, item.Value)
			if err != nil {
				return nil, err
			}
			znPairs = append(znPairs, KVPair{
				Key:   exprKey.value,
				Value: exprVal,
			})
		}
		return NewHashMap(znPairs), nil
	default:
		return nil, error.UnExpectedCase("表达式类型", fmt.Sprintf("%T", e))
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
		if v.MemberType == syntax.MemberID || v.MemberType == syntax.MemberIndex {
			iv, err := getMemberExprIV(ctx, v)
			if err != nil {
				return nil, err
			}
			return val, iv.ReduceLHS(ctx, val)
		}
		return nil, error.NewErrorSLOT("方法不能被赋值")
	default:
		return nil, error.UnExpectedCase("被赋值", fmt.Sprintf("%T", v))
	}
}

func getMemberExprIV(ctx *Context, expr *syntax.MemberExpr) (*IV, *error.Error) {
	switch expr.RootType {
	case syntax.RootTypeScope: // 此之 XX
		sgValue, err := findSgValue(ctx)
		if err != nil {
			return nil, err
		}

		return &IV{
			reduceType: IVTypeMember,
			root:       sgValue,
			member:     expr.MemberID.GetLiteral(),
		}, nil

	case syntax.RootTypeProp: // 其 XX
		thisValue, err := findThisValue(ctx)
		if err != nil {
			return nil, err
		}
		return &IV{
			reduceType: IVTypeMember,
			root:       thisValue,
			member:     expr.MemberID.GetLiteral(),
		}, nil
	case syntax.RootTypeExpr: // A 之 B
		valRoot, err := evalExpression(ctx, expr.Root)
		if err != nil {
			return nil, err
		}
		switch expr.MemberType {
		case syntax.MemberID: // A 之 B
			return &IV{
				reduceType: IVTypeMember,
				root:       valRoot,
				member:     expr.MemberID.GetLiteral(),
			}, nil
		case syntax.MemberIndex: // A # 0
			idx, err := evalExpression(ctx, expr.MemberIndex)
			if err != nil {
				return nil, err
			}
			switch v := valRoot.(type) {
			case *Array:
				vr, ok := idx.(*Decimal)
				if !ok {
					return nil, error.InvalidExprType("integer")
				}
				vri, e := vr.asInteger()
				if e != nil {
					return nil, error.InvalidExprType("integer")
				}
				return &IV{
					reduceType: IVTypeArray,
					root:       v,
					index:      vri,
				}, nil
			case *HashMap:
				var s string
				switch x := idx.(type) {
				// regard decimal value directly as string
				case *Decimal:
					// transform decimal value to string
					// x.exp < 0 express that its a decimal value with point mark, not an integer
					if x.exp < 0 {
						return nil, error.InvalidExprType("integer", "string")
					}
					s = x.String()
				case *String:
					s = x.String()
				default:
					return nil, error.InvalidExprType("integer", "string")
				}
				return &IV{
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

func getClassRef(ctx *Context, name string) (*ClassRef, *error.Error) {
	ref, ok := ctx.scope.classRefMap[name]
	if ok {
		return &ref, nil
	}
	return nil, error.NameNotDefined(name)
}

func bindClassRef(ctx *Context, classStmt *syntax.ClassDeclareStmt) *error.Error {
	name := classStmt.ClassName.GetLiteral()
	_, ok := ctx.scope.classRefMap[name]
	if ok {
		return error.NameRedeclared(name)
	}
	ctx.scope.classRefMap[name] = BuildClassFromNode(name, classStmt)
	return nil
}

// bind non-const value with re-declaration check on same scope
func bindValue(ctx *Context, name string, value Value) *error.Error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	// bind directly
	if ctx.scope != nil {
		if _, ok := ctx.scope.symbolMap[name]; ok {
			return error.NameRedeclared(name)
		}
		// set value
		ctx.scope.symbolMap[name] = SymbolInfo{value, false}
	}
	return nil
}

// bind value for declaration statement - that variables could be re-bind.
func bindValueDecl(ctx *Context, name string, value Value, isConst bool) *error.Error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return error.NameRedeclared(name)
	}
	if ctx.scope != nil {
		ctx.scope.symbolMap[name] = SymbolInfo{value, isConst}
	}
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

/// findSgValue - find the suitable sgValue in current context
// Rules:
//
// - if sgValue in current scope (ctx.scope.sgValue) != nil, returns the current one;
// - if sgValue in current scope == nil, then look up its parent util to the root;
func findSgValue(ctx *Context) (Value, *error.Error) {
	sp := ctx.scope
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

// findThisValue - similar with findSgValue(ctx), it looks up for nearest valid
// thisValue value.
func findThisValue(ctx *Context) (Value, *error.Error) {
	sp := ctx.scope
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
func duplicateValue(in Value) Value {
	switch v := in.(type) {
	case *Bool:
		return NewBool(v.value)
	case *String:
		return NewString(v.value)
	case *Decimal:
		x := new(big.Int)
		return &Decimal{
			co:  x.Set(v.co),
			exp: v.exp,
		}
	case *Null:
		return in // no need to copy since all "NULL" values are same
	case *Array:
		newArr := []Value{}
		for _, val := range v.value {
			newArr = append(newArr, duplicateValue(val))
		}
		return NewArray(newArr)
	case *HashMap:
		kvPairs := []KVPair{}
		for _, key := range v.keyOrder {
			dupVal := duplicateValue(v.value[key])
			kvPairs = append(kvPairs, KVPair{key, dupVal})
		}
		return NewHashMap(kvPairs)
	case *Function: // function itself is immutable, so return directly
		return in
	case *Object: // we don't copy object value at all
		return in
	}
	return in
}

// compareValues - some ZnValues are comparable from specific types of right value
// otherwise it will throw error.
//
// There are three types of compare verbs (actions): Eq, Lt and Gt.
//
// Eq - compare if two values are "equal". Usually there are two rules:
// 1. types of left and right value are same. A number MUST BE equals to a number, that means
// (string) “2” won't be equals to (number) 2;
// 2. each items SHOULD BE identical, even for composited types (i.e. array, hashmap)
//
// Lt - for two decimals ONLY. If leftValue < rightValue.
//
// Gt - for two decimals ONLY. If leftValue > rightValue.
//
func compareValues(left Value, right Value, verb compareVerb) (bool, *error.Error) {
	switch vl := left.(type) {
	case *Null:
		if _, ok := right.(*Null); ok {
			return true, nil
		}
		return false, nil
	case *Decimal:
		// compare right value - decimal only
		if vr, ok := right.(*Decimal); ok {
			r1, r2 := rescalePair(*vl, *vr)
			cmpResult := false
			switch verb {
			case CmpEq:
				cmpResult = (r1.co.Cmp(r2.co) == 0)
			case CmpLt:
				cmpResult = (r1.co.Cmp(r2.co) < 0)
			case CmpGt:
				cmpResult = (r1.co.Cmp(r2.co) > 0)
			default:
				return false, error.UnExpectedCase("比较原语", strconv.Itoa(int(verb)))
			}
			return cmpResult, nil
		}
		// if vert == CmbEq and rightValue is not decimal type
		// then return `false` directly
		if verb == CmpEq {
			return false, nil
		}
		return false, error.InvalidCompareRType("decimal")
	case *String:
		// Only CmpEq is valid for comparison
		if verb != CmpEq {
			return false, error.InvalidCompareLType("decimal", "string", "bool", "array", "hashmap")
		}
		// compare right value - string only
		if vr, ok := right.(*String); ok {
			cmpResult := (strings.Compare(vl.value, vr.value) == 0)
			return cmpResult, nil
		}
		return false, nil
	case *Bool:
		if verb != CmpEq {
			return false, error.InvalidCompareLType("decimal", "string", "bool", "array", "hashmap")
		}
		// compare right value - bool only
		if vr, ok := right.(*Bool); ok {
			cmpResult := vl.value == vr.value
			return cmpResult, nil
		}
		return false, nil
	case *Array:
		if verb != CmpEq {
			return false, error.InvalidCompareLType("decimal", "string", "bool", "array", "hashmap")
		}

		if vr, ok := right.(*Array); ok {
			if len(vl.value) != len(vr.value) {
				return false, nil
			}
			// cmp each item
			for idx := range vl.value {
				cmpVal, err := compareValues(vl.value[idx], vr.value[idx], CmpEq)
				if err != nil {
					return false, err
				}
				return cmpVal, nil
			}
			return true, nil
		}
		return false, nil
	case *HashMap:
		if verb != CmpEq {
			return false, error.InvalidCompareLType("decimal", "string", "bool", "array", "hashmap")
		}

		if vr, ok := right.(*HashMap); ok {
			if len(vl.value) != len(vr.value) {
				return false, nil
			}
			// cmp each item
			for idx := range vl.value {
				// ensure the key exists on vr
				vrr, ok := vr.value[idx]
				if !ok {
					return false, nil
				}
				cmpVal, err := compareValues(vl.value[idx], vrr, CmpEq)
				if err != nil {
					return false, err
				}
				return cmpVal, nil
			}
			return true, nil
		}
		return false, nil
	}
	return false, error.InvalidCompareLType("decimal", "string", "bool", "array", "hashmap")
}

// StringifyValue - yield a string from Value
func StringifyValue(value Value) string {
	switch v := value.(type) {
	case *String:
		return fmt.Sprintf("「%s」", v.value)
	case *Decimal:
		return v.String()
	case *Array:
		strs := []string{}
		for _, item := range v.value {
			strs = append(strs, StringifyValue(item))
		}

		return fmt.Sprintf("【%s】", strings.Join(strs, "，"))
	case *Bool:
		data := "真"
		if v.value == false {
			data = "假"
		}
		return data
	case *Function:
		return fmt.Sprintf("[方法]")
	case *Null:
		return "空"
	case *Object:
		return fmt.Sprintf("[对象]")
	case *HashMap:
		strs := []string{}
		for _, key := range v.keyOrder {
			value := v.value[key]
			strs = append(strs, fmt.Sprintf("%s == %s", key, StringifyValue(value)))
		}
		return fmt.Sprintf("【%s】", strings.Join(strs, "，"))
	}
	return ""
}
