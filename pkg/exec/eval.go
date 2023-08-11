package exec

import (
	"fmt"

	"github.com/DemoHn/Zn/stdlib"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/syntax/zh"
	"github.com/DemoHn/Zn/pkg/value"
)

// eval.go evaluates program from generated AST tree with specific scopes
// common signature of eval functions:
//
// evalXXXXStmt(c *r.Context, node Node) error
//
// or
//
// evalXXXXExpr(c *r.Context, node Node) (r.Value, error)
//
// NOTICE:
// `evalXXXXStmt` will change the value of its corresponding scope; However, `evalXXXXExpr` will export
// a r.Value object and mostly won't change scopes (but search a variable from scope is frequently used)

func evalProgram(c *r.Context, program *syntax.Program) error {
	otherStmts, err := evalPreStmtBlock(c, program.Content)
	if err != nil {
		return err
	}
	errBlock := evalStmtBlock(c, otherStmts)
	if errBlock != nil {
		val, err := extractSignalValue(errBlock, zerr.SigTypeReturn)
		if err != nil {
			return err
		}
		// set return value
		c.GetCurrentScope().SetReturnValue(val)
		return nil
	}
	return errBlock
}

//// eval statements

// EvalStatement - eval statement
func evalStatement(c *r.Context, stmt syntax.Statement) error {
	var returnValue r.Element
	var sp = c.GetCurrentScope()
	// set current line
	c.SetCurrentLine(stmt.GetCurrentLine())

	// set return value
	defer func() {
		var finalReturnValue r.Element = value.NewNull()
		// set current return value
		if returnValue != nil {
			finalReturnValue = returnValue
		}
		sp.SetReturnValue(finalReturnValue)

		// set parent return value
		parentScope := c.FindParentScope()
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
		fn := compileFunction(c, v.ParamList, v.ExecBlock)
		return c.BindSymbol(v.FuncName.GetLiteral(), fn)
	case *syntax.ClassDeclareStmt:
		className := v.ClassName.GetLiteral()
		if c.FindParentScope() != nil {
			return zerr.ClassNotOnRoot(className)
		}
		// bind classRef
		classRef := BuildClassFromNode(c, className, v)

		return c.BindSymbol(className, classRef)
	case *syntax.IterateStmt:
		return evalIterateStmt(c, v)
	case *syntax.FunctionReturnStmt:
		val, err := evalExpression(c, v.ReturnExpr)
		if err != nil {
			return err
		}
		// send RETURN break
		return zerr.NewReturnSignal(val)
	case *syntax.ThrowExceptionStmt:
		// profoundly return an ERROR to terminate the process
		name := v.ExceptionClass.GetLiteral()
		expClassRef, err := c.FindElement(name)
		if err != nil {
			return err
		}

		if ref, ok := expClassRef.(*value.ClassModel); ok {
			// exec expressions
			var exprs []r.Element
			for _, param := range v.Params {
				exprI, err := evalExpression(c, param)
				if err != nil {
					return err
				}
				exprs = append(exprs, exprI)
			}
			// get exception value!
			val, err := ref.Construct(c, exprs)
			if err != nil {
				return err
			}
			// val MUST BE an Exception Value!
			if expVal, ok := val.(*value.Exception); ok {
				return zerr.NewRuntimeException(expVal.GetMessage())
			}
			return zerr.InvalidExceptionObjectType(name)
		}
		return zerr.InvalidExceptionType(name)
	case *syntax.ContinueStmt:
		// send continue signal
		return zerr.NewContinueSignal()
	case *syntax.BreakStmt:
		return zerr.NewBreakSignal()
	case syntax.Expression:
		expr, err := evalExpression(c, v)
		returnValue = expr
		return err
	default:
		return zerr.UnexpectedCase("语句类型", fmt.Sprintf("%T", v))
	}
}

// evalVarDeclareStmt - consists of three branches:
// 1. A，B 为 C
// 2. A，B 成为 X：P1，P2，...
// 3. A，B 恒为 C
func evalVarDeclareStmt(c *r.Context, node *syntax.VarDeclareStmt) error {
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
					obj = value.DuplicateValue(obj)
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

// eval A,B 成为（C：P1，P2，P3，...）
// ensure VDAssignPair.Type MUST BE syntax.VDTypeObjNew
func evalNewObject(c *r.Context, node syntax.VDAssignPair) error {
	vtag := node.ObjClass.GetLiteral()
	// get class definition
	importVal, err := c.FindElement(vtag)
	if err != nil {
		return err
	}
	classRef, ok := importVal.(*value.ClassModel)
	if !ok {
		return zerr.InvalidParamType("classRef")
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

		if err := c.BindSymbol(vtag, finalObj); err != nil {
			return err
		}
	}
	return nil
}

// evalWhileLoopStmt -
func evalWhileLoopStmt(c *r.Context, node *syntax.WhileLoopStmt) error {
	// create new scope
	c.PushScope()
	defer c.PopScope()

	// set context's current scope with new one

	for {
		// #1. first execute expr
		trueExpr, err := evalExpression(c, node.TrueExpr)
		if err != nil {
			return err
		}
		// #2. assert trueExpr to be Bool
		vTrueExpr, ok := trueExpr.(*value.Bool)
		if !ok {
			return zerr.InvalidExprType("bool")
		}
		// break the loop if expr yields not true
		if !vTrueExpr.GetValue() {
			return nil
		}
		// #3. stmt block
		if err := evalStmtBlock(c, node.LoopBlock); err != nil {
			if s, ok := err.(*zerr.Signal); ok {
				if s.SigType == zerr.SigTypeContinue {
					continue
				}
				if s.SigType == zerr.SigTypeBreak {
					return nil
				}
			}
			return err
		}
	}
}

// evalPreStmtBlock - execute classRef, functionDeclare, imports first, then other statements inside the block
func evalPreStmtBlock(c *r.Context, block *syntax.BlockStmt) (*syntax.BlockStmt, error) {
	// current module MUST exists
	module := c.GetCurrentModule()
	otherStmts := &syntax.BlockStmt{
		Children: []syntax.Statement{},
	}

	for _, stmtI := range block.Children {
		// set current execLine
		c.SetCurrentLine(stmtI.GetCurrentLine())

		switch v := stmtI.(type) {
		case *syntax.FunctionDeclareStmt:
			fn := compileFunction(c, v.ParamList, v.ExecBlock)
			vtag := v.FuncName.GetLiteral()

			// add symbol to current scope first
			if err := c.BindSymbol(vtag, fn); err != nil {
				return nil, err
			}

			// then add symbol to export value
			if err := module.AddExportValue(vtag, fn); err != nil {
				return nil, err
			}
		case *syntax.ClassDeclareStmt:
			// bind classRef
			className := v.ClassName.GetLiteral()
			classRef := BuildClassFromNode(c, className, v)
			// add symbol to current scope first
			if err := c.BindSymbol(className, classRef); err != nil {
				return nil, err
			}

			// then add symbol to export value
			if err := module.AddExportValue(className, classRef); err != nil {
				return nil, err
			}
		case *syntax.ImportStmt:
			libName := v.ImportName.GetLiteral()

			var extModule *r.Module
			if v.ImportLibType == syntax.LibTypeStd {
				var err error
				// check if the dependency is valid (i.e. not import itself/no duplicate import)
				if err := c.CheckDepedency(libName, true); err != nil {
					return nil, err
				}
				extModule, err = stdlib.FindModule(libName)
				if err != nil {
					return nil, err
				}
			} else if v.ImportLibType == syntax.LibTypeCustom {
				// check if the dependency is valid (i.e. not import itself/no duplicate import/no circular dependency)
				if err := c.CheckDepedency(libName, false); err != nil {
					return nil, err
				}
				// execute custom module first (in order to get all importable elements)
				if extModule = c.FindModuleCache(libName); extModule == nil {
					newModule, err := execAnotherModule(c, libName)
					if err != nil {
						return nil, err
					}
					extModule = newModule
				}
			}

			if extModule != nil {
				// import all symbols to current module's importRefs
				if len(v.ImportItems) == 0 {
					for name, val := range extModule.GetAllExportValues() {
						if err := c.BindImportSymbol(name, val, extModule); err != nil {
							return nil, err
						}
					}
				} else {
					// import selected symbols
					for _, id := range v.ImportItems {
						name := id.GetLiteral()
						if val, err2 := extModule.GetExportValue(name); err2 == nil {
							if err := c.BindImportSymbol(name, val, extModule); err != nil {
								return nil, err
							}
						}
					}
				}
			}
		default:
			otherStmts.Children = append(otherStmts.Children, stmtI)
		}
	}
	return otherStmts, nil
}

// EvalStmtBlock -
func evalStmtBlock(c *r.Context, block *syntax.BlockStmt) error {
	for _, stmt := range block.Children {
		if err := evalStatement(c, stmt); err != nil {
			return err
		}
	}

	return nil
}

func evalBranchStmt(c *r.Context, node *syntax.BranchStmt) error {
	// create inner scope for if statement
	c.PushScope()
	defer c.PopScope()

	// #1. condition header
	ifExpr, err := evalExpression(c, node.IfTrueExpr)
	if err != nil {
		return err
	}
	vIfExpr, ok := ifExpr.(*value.Bool)
	if !ok {
		return zerr.InvalidExprType("bool")
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
		vOtherExprI, ok := otherExprI.(*value.Bool)
		if !ok {
			return zerr.InvalidExprType("bool")
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

func evalIterateStmt(c *r.Context, node *syntax.IterateStmt) error {
	c.PushScope()
	defer c.PopScope()

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
	execIterationBlockFn := func(key r.Element, v r.Element) error {
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
		if err := c.BindSymbol(valueSlot, value.NewNull()); err != nil {
			return err
		}
	} else if nameLen == 2 {
		keySlot = node.IndexNames[0].Literal
		valueSlot = node.IndexNames[1].Literal
		if err := c.BindSymbol(keySlot, value.NewNull()); err != nil {
			return err
		}
		if err := c.BindSymbol(valueSlot, value.NewNull()); err != nil {
			return err
		}
	} else if nameLen > 2 {
		return zerr.MostParamsError(2)
	}

	// execute iterations
	switch tv := targetExpr.(type) {
	case *value.Array:
		for idx, v := range tv.GetValue() {
			// in iterate statement, index starts from 1 instead of 0
			realIdx := idx + 1
			idxVar := value.NewNumber(float64(realIdx))
			if err := execIterationBlockFn(idxVar, v); err != nil {
				if s, ok := err.(*zerr.Signal); ok {
					if s.SigType == zerr.SigTypeContinue {
						continue
					}
					if s.SigType == zerr.SigTypeBreak {
						return nil
					}
				}
				return err
			}
		}
	case *value.HashMap:
		for _, key := range tv.GetKeyOrder() {
			v := tv.GetValue()[key]
			keyVar := value.NewString(key)
			// handle interrupts
			if err := execIterationBlockFn(keyVar, v); err != nil {
				if s, ok := err.(*zerr.Signal); ok {
					if s.SigType == zerr.SigTypeContinue {
						continue
					}
					if s.SigType == zerr.SigTypeBreak {
						return nil
					}
				}
				return err
			}
		}
	default:
		return zerr.InvalidExprType("array", "hashmap")
	}

	return nil
}

//// execute expressions
func evalExpression(c *r.Context, expr syntax.Expression) (r.Element, error) {
	switch e := expr.(type) {
	case *syntax.VarAssignExpr:
		return evalVarAssignExpr(c, e)
	case *syntax.LogicExpr:
		if e.Type == syntax.LogicAND || e.Type == syntax.LogicOR {
			return evalLogicCombiner(c, e)
		}
		return evalLogicComparator(c, e)
	case *syntax.ArithExpr:
		return evalArithExpr(c, e)
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
	case *syntax.MemberMethodExpr:
		return evalMemberMethodExpr(c, e)
	default:
		return nil, zerr.InvalidExprType()
	}
}

//// checkout eval_function.go for evalFunctionCall()

// 以 A （执行：B、C、D）
func evalMemberMethodExpr(c *r.Context, expr *syntax.MemberMethodExpr) (r.Element, error) {
	// 1. parse root expr
	rootExpr, err := evalExpression(c, expr.Root)
	if err != nil {
		return nil, err
	}

	var vlast r.Element = rootExpr

	for _, methodExpr := range expr.MethodChain {
		// eval method
		funcName := methodExpr.FuncName.GetLiteral()

		// exec params
		params, err := exprsToValues(c, methodExpr.Params)
		if err != nil {
			return nil, err
		}

		v, err := execMethodFunction(c, vlast, funcName, params)
		if err != nil {
			return nil, err
		}
		vlast = v
	}

	// add yield result
	if expr.YieldResult != nil {
		vtag := expr.YieldResult.GetLiteral()
		// bind yield result
		if err := c.BindSymbolDecl(vtag, vlast, false); err != nil {
			return nil, err
		}
	}

	return vlast, nil
}

// evaluate logic combination expressions
// such as A 且 B
// or A 或 B
func evalLogicCombiner(c *r.Context, expr *syntax.LogicExpr) (*value.Bool, error) {
	logicType := expr.Type
	// #1. eval left
	left, err := evalExpression(c, expr.LeftExpr)
	if err != nil {
		return nil, err
	}
	// #2. assert left expr type to be ZnBool
	vleft, ok := left.(*value.Bool)
	if !ok {
		return nil, zerr.InvalidExprType("bool")
	}
	// #3. check if the result could be retrieved earlier
	//
	// 1) for Y = A and B, if A = false, then Y must be false
	// 2) for Y = A or  B, if A = true, then Y must be true
	//
	// for those cases, we can yield result directly
	if logicType == syntax.LogicAND && !vleft.GetValue() {
		return value.NewBool(false), nil
	}
	if logicType == syntax.LogicOR && vleft.GetValue() {
		return value.NewBool(true), nil
	}
	// #4. eval right
	right, err := evalExpression(c, expr.RightExpr)
	if err != nil {
		return nil, err
	}
	vright, ok := right.(*value.Bool)
	if !ok {
		return nil, zerr.InvalidExprType("bool")
	}
	// then evalute data
	switch logicType {
	case syntax.LogicAND:
		return value.NewBool(vleft.GetValue() && vright.GetValue()), nil
	default: // logicOR
		return value.NewBool(vleft.GetValue() || vright.GetValue()), nil
	}
}

// evaluate logic comparator
// ensure both expressions are comparable (i.e. subtype of ZnComparable)
func evalLogicComparator(c *r.Context, expr *syntax.LogicExpr) (*value.Bool, error) {
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
	var cmpErr error
	// #3. do comparison
	switch logicType {
	case syntax.LogicEQ:
		cmpRes, cmpErr = value.CompareValues(left, right, value.CmpEq)
	case syntax.LogicNEQ:
		cmpRes, cmpErr = value.CompareValues(left, right, value.CmpEq)
		cmpRes = !cmpRes // reverse result
	case syntax.LogicGT:
		cmpRes, cmpErr = value.CompareValues(left, right, value.CmpGt)
	case syntax.LogicGTE:
		var cmp1, cmp2 bool
		cmp1, cmpErr = value.CompareValues(left, right, value.CmpGt)
		if cmpErr != nil {
			return nil, cmpErr
		}
		cmp2, cmpErr = value.CompareValues(left, right, value.CmpEq)
		cmpRes = cmp1 || cmp2
	case syntax.LogicLT:
		cmpRes, cmpErr = value.CompareValues(left, right, value.CmpLt)
	case syntax.LogicLTE:
		var cmp1, cmp2 bool
		cmp1, cmpErr = value.CompareValues(left, right, value.CmpLt)
		if cmpErr != nil {
			return nil, cmpErr
		}
		cmp2, cmpErr = value.CompareValues(left, right, value.CmpEq)
		cmpRes = cmp1 || cmp2
	default:
		return nil, zerr.UnexpectedCase("比较类型", fmt.Sprintf("%d", logicType))
	}

	return value.NewBool(cmpRes), cmpErr
}

func evalArithExpr(c *r.Context, expr *syntax.ArithExpr) (*value.Number, error) {
	// exec left Expr
	leftExpr, err := evalExpression(c, expr.LeftExpr)
	if err != nil {
		return nil, err
	}

	leftNum, ok := leftExpr.(*value.Number)
	if !ok {
		return nil, zerr.InvalidExprType("number")
	}
	// exec right expr
	rightExpr, err := evalExpression(c, expr.RightExpr)
	if err != nil {
		return nil, err
	}

	rightNum, ok := rightExpr.(*value.Number)
	if !ok {
		return nil, zerr.InvalidExprType("number")
	}

	// calculate num
	switch expr.Type {
	case syntax.ArithAdd:
		return value.NewNumber(leftNum.GetValue() + rightNum.GetValue()), nil
	case syntax.ArithSub:
		return value.NewNumber(leftNum.GetValue() - rightNum.GetValue()), nil
	case syntax.ArithMul:
		return value.NewNumber(leftNum.GetValue() * rightNum.GetValue()), nil
	case syntax.ArithDiv:
		if rightNum.GetValue() == 0 {
			return nil, zerr.ArithDivZero()
		}
		return value.NewNumber(leftNum.GetValue() / rightNum.GetValue()), nil
	}
	return nil, zerr.UnexpectedCase("运算项", fmt.Sprintf("%d", expr.Type))
}

// eval prime expr
func evalPrimeExpr(c *r.Context, expr syntax.Expression) (r.Element, error) {
	switch e := expr.(type) {
	case *syntax.Number:
		return value.NewNumberFromString(e.GetLiteral())
	case *syntax.String:
		return value.NewString(e.GetLiteral()), nil
	case *syntax.ID:
		vtag := e.GetLiteral()
		return c.FindElement(vtag)
	case *syntax.ArrayExpr:
		var znObjs []r.Element
		for _, item := range e.Items {
			expr, err := evalExpression(c, item)
			if err != nil {
				return nil, err
			}
			znObjs = append(znObjs, expr)
		}

		return value.NewArray(znObjs), nil
	case *syntax.HashMapExpr:
		var znPairs []value.KVPair
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
				return nil, zerr.InvalidExprType("string", "number", "id")
			}

			exprVal, err := evalExpression(c, item.Value)
			if err != nil {
				return nil, err
			}
			znPairs = append(znPairs, value.KVPair{
				Key:   exprKey,
				Value: exprVal,
			})
		}
		return value.NewHashMap(znPairs), nil
	default:
		return nil, zerr.UnexpectedCase("表达式类型", fmt.Sprintf("%T", e))
	}
}

// eval variable assign
func evalVarAssignExpr(c *r.Context, expr *syntax.VarAssignExpr) (r.Element, error) {
	// Right Side
	vr, err := evalExpression(c, expr.AssignExpr)
	if err != nil {
		return nil, err
	}
	// if var assignment is NOT by reference, then duplicate value
	if !expr.RefMark {
		vr = value.DuplicateValue(vr)
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
		return nil, zerr.UnexpectedAssign()
	default:
		return nil, zerr.UnexpectedCase("被赋值", fmt.Sprintf("%T", v))
	}
}

// execAnotherModule - load source code of the module, parse the coe, execute the program, and build depCache!
func execAnotherModule(c *r.Context, name string) (*r.Module, error) {
	if finder := c.GetModuleCodeFinder(); finder != nil {
		source, err := finder(name)
		if err != nil {
			return nil, zerr.ModuleNotFound(name)
		}
		// #1.  create & enter module
		lexer := syntax.NewLexer(source)
		module := r.NewModule(name, lexer)
		c.EnterModule(module)
		defer c.ExitModule()

		// #2. parse program
		p := syntax.NewParser(lexer, zh.NewParserZH())

		program, err := p.Parse()
		if err != nil {
			return nil, WrapSyntaxError(lexer, module, err)
		}

		// #3. eval program
		if err := evalProgram(c, program); err != nil {
			return nil, WrapRuntimeError(c, err)
		}

		return module, nil
	}
	// no finder defined, return nil directly (no throw error)
	return nil, nil
}

func getMemberExprIV(c *r.Context, expr *syntax.MemberExpr) (*value.IV, error) {
	switch expr.RootType {
	case syntax.RootTypeProp: // 其 XX
		thisValue := c.GetThisValue()
		if thisValue == nil {
			return nil, zerr.ThisValueNotFound()
		}
		return value.NewMemberIV(thisValue, expr.MemberID.GetLiteral()), nil
	case syntax.RootTypeExpr: // A 之 B
		valRoot, err := evalExpression(c, expr.Root)
		if err != nil {
			return nil, err
		}
		switch expr.MemberType {
		case syntax.MemberID: // A 之 B
			return value.NewMemberIV(valRoot, expr.MemberID.GetLiteral()), nil
		case syntax.MemberIndex: // A # 0
			idx, err := evalExpression(c, expr.MemberIndex)
			if err != nil {
				return nil, err
			}
			switch v := valRoot.(type) {
			case *value.Array:
				vr, ok := idx.(*value.Number)
				if !ok {
					return nil, zerr.InvalidExprType("integer")
				}
				vri := int(vr.GetValue())
				return value.NewArrayIV(v, vri), nil
			case *value.HashMap:
				var s string
				switch x := idx.(type) {
				// regard decimal value directly as string
				case *value.Number:
					s = x.String()
				case *value.String:
					s = x.String()
				default:
					return nil, zerr.InvalidExprType("integer", "string")
				}
				return value.NewHashMapIV(v, s), nil
			}
			return nil, zerr.InvalidExprType("array", "hashmap")
		}
		return nil, zerr.UnexpectedCase("子项类型", fmt.Sprintf("%d", expr.MemberType))
	}

	return nil, zerr.UnexpectedCase("根元素类型", fmt.Sprintf("%d", expr.RootType))
}

//// helpers
// exprsToValues - []syntax.Expression -> []eval.r.Value
func exprsToValues(c *r.Context, exprs []syntax.Expression) ([]r.Element, error) {
	params := []r.Element{}
	for _, paramExpr := range exprs {
		pval, err := evalExpression(c, paramExpr)
		if err != nil {
			return nil, err
		}
		params = append(params, pval)
	}
	return params, nil
}

// extractSignalValue - signal is a special type of error, so we try to extract signal value from input error if it's really a signal - otherwise output the same error directly.
func extractSignalValue(err error, sigType uint8) (r.Element, error) {
	// if recv breaks
	if sig, ok := err.(*zerr.Signal); ok {
		if sig.SigType == sigType {
			if extra, ok2 := sig.Extra.(r.Element); ok2 {
				return extra, nil
			}
		}
	}
	return nil, err
}

// BuildClassFromNode -
func BuildClassFromNode(upperCtx *r.Context, name string, classNode *syntax.ClassDeclareStmt) *value.ClassModel {
	ref := value.NewClassModel(name, upperCtx.GetCurrentModule())

	// define default constructor
	var constructor = func(c *r.Context, params []r.Element) (r.Element, error) {
		obj := value.NewObject(ref)
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
			return nil, zerr.MismatchParamLengthError(len(params), len(classNode.ConstructorIDList))
		}
		for idx, objParamVal := range params {
			param := classNode.ConstructorIDList[idx]
			// if param is NOT a reference, then we need to copy its value
			if !param.RefMark {
				objParamVal = value.DuplicateValue(objParamVal)
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
		ref.CompPropList[getterTag] = compileFunction(upperCtx, []*syntax.ParamItem{}, gNode.ExecBlock)
	}

	// add methods
	for _, mNode := range classNode.MethodList {
		mTag := mNode.FuncName.GetLiteral()
		ref.MethodList[mTag] = compileFunction(upperCtx, mNode.ParamList, mNode.ExecBlock)
	}

	return ref
}
