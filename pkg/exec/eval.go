package exec

import (
	"fmt"
	"math"
	"strings"

	"github.com/DemoHn/Zn/stdlib"

	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/runtime"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/syntax/zh"
	"github.com/DemoHn/Zn/pkg/value"
)

const (
	EVConstExceptionClassName       = "异常"
	EVConstExceptionContentProperty = "内容"
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
// a r.Value object and mostly won't change scope values (but searching a variable from scope is frequently used)

// evalProgram - eval the statements of the program with the following order:
//
// 1. IMPORT statement(s) - `导入《文件》`
// 2. INPUTVAR statement - `输入长、宽、高`
// 3. CLASSDEF statement(s) - `定义货件`
// 4. FUNCDEF statement(s) - `如何执行？`
// 5. other statements
// 6. CATCH statement(s) - `拦截异常`
//
// If the program doesn't follow the order (e.g. the func declare block at the end of program),
// it doesn't matter - we will order the statements in the program automatically before execution.
// (This will not affect line numbers)
func evalProgram(vm *r.VM, program *syntax.Program, varInputs r.ElementMap) (r.Element, error) {
	vm.PushCallFrame(runtime.NewScriptCallFrame(-1, program, nil))
	defer vm.PopCallFrame()

	// 1. import libs
	for _, importStmt := range program.ImportBlock {
		if err := evalImportStmt(vm, importStmt); err != nil {
			return nil, err
		}
	}

	if program.ExecBlock != nil {
		// 2. do exec block -> load values from context.varInputs -> current scope
		paramList := []r.Element{}
		for _, inputV := range program.ExecBlock.InputBlock {
			inputName, err := MatchIDName(inputV)
			if err != nil {
				return nil, err
			}
			inputNameStr := inputName.GetLiteral()

			// match name from idList and append the value from varInputMap
			if elem, ok := varInputs[inputNameStr]; ok {
				paramList = append(paramList, elem)
			} else {
				return nil, zerr.InputValueNotFound(inputNameStr)
			}
		}
		return evalExecBlock(vm, program.ExecBlock, paramList)
	}

	return value.NewNull(), nil
}

func evalExecBlock(vm *r.VM, execBlock *syntax.ExecBlock, params []r.Element) (r.Element, error) {
	// 1.1 check param length
	inputParamNum := len(execBlock.InputBlock)
	if len(params) != inputParamNum {
		return nil, zerr.MismatchParamLengthError(inputParamNum, len(params))
	}

	for idx, param := range execBlock.InputBlock {
		idTag, err := MatchIDName(param)
		if err != nil {
			return nil, err
		}

		// set inputValue to current scope
		if err := vm.DeclareElement(idTag, params[idx]); err != nil {
			return nil, err
		}
	}

	var errBlock error
	var rtnValue r.Element
	rtnValue, errBlock = evalStmtBlock(vm, execBlock.StmtBlock)

	// extract & handle exception value (抛出)
	if errBlock != nil {
		exception, realErr := extractSignalValue(errBlock, zerr.SigTypeException)
		if realErr == nil {
			errBlock = handleExceptions(vm, execBlock.CatchBlock, exception)
		}
	}
	return rtnValue, errBlock
}

func evalStmtBlock(vm *r.VM, stmtBlock *syntax.StmtBlock) (r.Element, error) {
	// create inner scope for if statement
	vm.BeginScope()
	defer vm.EndScope()

	classDefStmts := make([]syntax.Statement, 0)
	funcDefStmts := make([]syntax.Statement, 0)
	otherStmts := make([]syntax.Statement, 0)

	allStmts := make([]syntax.Statement, 0)

	for _, stmtX := range stmtBlock.Children {
		switch v := stmtX.(type) {
		case *syntax.ClassDeclareStmt:
			classDefStmts = append(classDefStmts, v)
		case *syntax.FunctionDeclareStmt:
			funcDefStmts = append(funcDefStmts, v)
		default:
			otherStmts = append(otherStmts, v)
		}
	}

	// reorder the statements
	allStmts = append(allStmts, classDefStmts...)
	allStmts = append(allStmts, funcDefStmts...)
	allStmts = append(allStmts, otherStmts...)

	var rtnValue r.Element
	var err error
	for _, stmt := range allStmts {
		if rtnValue, err = evalStatement(vm, stmt); err != nil {
			return nil, err
		}
	}

	return rtnValue, nil
}

func handleExceptions(c *r.Context, catchBlock []*syntax.CatchBlockPair, exception r.Element) error {
	// by default, we use "异常" to match *value.Exception type exceptions
	var objClassName = ""
	switch v := exception.(type) {
	case *value.Exception:
		objClassName = EVConstExceptionClassName
	case *value.Object:
		objClassName = v.GetObjectName()
	}

	// iterate catchBlocks to match
	for _, catchBlockItem := range catchBlock {
		classID, err := MatchIDName(catchBlockItem.ExceptionClass)
		if err != nil {
			return err
		}

		// if exception block matches exception className
		if objClassName != "" && classID.GetLiteral() == objClassName {
			newScope := c.PushScope()
			defer c.PopScope()

			newScope.SetThisValue(exception)

			// do execution (with "this" value = exception value)
			_, err := evalStmtBlock(c, catchBlockItem.StmtBlock)
			return err
		}
	}

	// no handle block catches this error, then throw it anyway
	// for default *value.Exception class
	if objE, ok := exception.(*value.Exception); ok {
		return objE
	} else {
		// other custom exception class
		finalStr := ""
		s, _ := exception.GetProperty(c, EVConstExceptionContentProperty)
		if s != nil {
			if ss, ok := s.(*value.String); ok {
				finalStr = ss.GetValue()
			}
		}
		return value.NewException(finalStr)
	}
}

//// eval statements

// EvalStatement - eval statement
func evalStatement(vm *r.VM, stmt syntax.Statement) (r.Element, error) {
	// set current line
	vm.SetCurrentLine(stmt.GetCurrentLine())

	switch v := stmt.(type) {
	case *syntax.VarDeclareStmt:
		return value.NewNull(), evalVarDeclareStmt(vm, v)
	case *syntax.WhileLoopStmt:
		return value.NewNull(), evalWhileLoopStmt(vm, v)
	case *syntax.BranchStmt:
		return value.NewNull(), evalBranchStmt(vm, v)
	case *syntax.EmptyStmt:
		return value.NewNull(), nil
	case *syntax.FunctionDeclareStmt:
		if v.DeclareType == syntax.DeclareTypeConstructor {
			return value.NewNull(), evalConstructorDeclareStmt(vm, v)
		} else {
			return value.NewNull(), evalFunctionDeclareStmt(vm, v)
		}
	case *syntax.ClassDeclareStmt:
		return value.NewNull(), evalClassDeclareStmt(vm, v)
	case *syntax.IterateStmt:
		return value.NewNull(), evalIterateStmt(vm, v)
	case *syntax.FunctionReturnStmt:
		return evalExpression(vm, v.ReturnExpr)
	case *syntax.ThrowExceptionStmt:
		return value.NewNull(), evalThrowExceptionStmt(vm, v)
	case *syntax.ContinueStmt:
		// send continue signal
		return value.NewNull(), zerr.NewContinueSignal()
	case *syntax.BreakStmt:
		return value.NewNull(), zerr.NewBreakSignal()
	case syntax.Expression:
		return evalExpression(vm, v)
	default:
		return nil, zerr.UnexpectedCase("语句类型", fmt.Sprintf("%T", v))
	}
}

// evalVarDeclareStmt - consists of three branches:
// 1. A，B 设为 C
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
				vtag, err := MatchIDName(v)
				if err != nil {
					return err
				}

				obj = value.DuplicateValue(obj)
				if err := c.BindSymbolDecl(vtag, obj, isConst); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// 定义XX
func evalClassDeclareStmt(c *r.Context, node *syntax.ClassDeclareStmt) error {
	module := c.GetCurrentModule()

	className, err := MatchIDName(node.ClassName)
	if err != nil {
		return err
	}

	if c.FindParentScope() != nil {
		return zerr.ClassNotOnRoot(className.GetLiteral())
	}
	// bind classRef
	classRef, err := compileClass(c, className, node)
	if err != nil {
		return err
	}

	// add symbol to current scope first
	if err := c.BindSymbolConst(className, classRef); err != nil {
		return err
	}

	// then add symbol to export value
	if err := module.AddExportValue(className.GetLiteral(), classRef); err != nil {
		return err
	}
	return nil
}

// 如何XX？
func evalFunctionDeclareStmt(c *r.Context, node *syntax.FunctionDeclareStmt) error {
	module := c.GetCurrentModule()

	// declare as normal function
	vtag, err := MatchIDName(node.Name)
	if err != nil {
		return err
	}
	fn := compileFunction(c, node)

	// add symbol to current scope first
	if err := c.BindSymbol(vtag, fn); err != nil {
		return err
	}

	// then add symbol to export value
	if err := module.AddExportValue(vtag.GetLiteral(), fn); err != nil {
		return err
	}
	return nil
}

// 如何新建XX？
func evalConstructorDeclareStmt(upperCtx *r.Context, node *syntax.FunctionDeclareStmt) error {
	// 1. check if class type is valid
	className, err := MatchIDName(node.Name)
	if err != nil {
		return err
	}
	// 2. find class model
	classModel, err := upperCtx.FindElement(className)
	if err != nil {
		return err
	}

	cmodel, ok := classModel.(*value.ClassModel)
	if !ok {
		return zerr.InvalidClassType(className.GetLiteral())
	}

	//// there are some different Factors from normal method function:
	// 1. no outerScope (clousure scope)
	// 2. no 此 const variable inside the fn scope
	constructorLogic := func(c *r.Context, elems []r.Element) (r.Element, error) {
		fnScope := c.PushScope()
		defer c.PopScope()

		newObject := value.NewObject(cmodel, map[string]r.Element{})
		// set "this" value
		fnScope.SetThisValue(newObject)

		if _, err := evalExecBlock(c, node.ExecBlock, elems); err != nil {
			return nil, err
		}
		return newObject, nil
	}
	cmodel.SetConstructor(constructorLogic)

	return nil
}

// eval 创建XX：P1，P2，P3，...！
// ensure VDAssignPair.Type MUST BE syntax.VDTypeObjNew
func evalNewObject(vm *r.VM, node *syntax.ObjNewExpr) (r.Element, error) {
	classID, err := MatchIDName(node.ClassName)
	if err != nil {
		return nil, err
	}
	// get class definition
	importVal, err := vm.FindElement(classID)
	if err != nil {
		return nil, err
	}
	classRef, ok := importVal.(*value.ClassModel)
	if !ok {
		return nil, zerr.InvalidParamType("classRef")
	}

	cParams, err := exprsToValues(vm, node.Params)
	if err != nil {
		return nil, err
	}

	return classRef.Construct(vm, cParams)
}

// eval 导入《模块A》
func evalImportStmt(c *r.Context, node *syntax.ImportStmt) error {
	libName := node.ImportName.GetLiteral()

	var extModule *r.Module
	if node.ImportLibType == syntax.LibTypeStd {
		var err error
		// check if the dependency is valid (i.e. not import itself/no duplicate import)
		if err := c.CheckDepedency(libName, true); err != nil {
			return err
		}
		extModule, err = stdlib.FindModule(libName)
		if err != nil {
			return err
		}
	} else if node.ImportLibType == syntax.LibTypeCustom {
		// check if the dependency is valid (i.e. not import itself/no duplicate import/no circular dependency)
		if err := c.CheckDepedency(libName, false); err != nil {
			return err
		}
		// execute custom module first (in order to get all importable elements)
		if extModule = c.FindModuleCache(libName); extModule == nil {
			newModule, err := execAnotherModule(c, libName)
			if err != nil {
				return err
			}
			extModule = newModule
		}
	}

	if extModule != nil {
		// import all symbols to current module's importRefs
		if len(node.ImportItems) == 0 {
			for name, val := range extModule.GetAllExportValues() {
				if err := c.BindImportSymbol(name, val, extModule); err != nil {
					return err
				}
			}
		} else {
			// import selected symbols
			for _, id := range node.ImportItems {
				name := id.GetLiteral()
				if val, err2 := extModule.GetExportValue(name); err2 == nil {
					if err := c.BindImportSymbol(name, val, extModule); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// evalWhileLoopStmt -
func evalWhileLoopStmt(vm *r.VM, node *syntax.WhileLoopStmt) error {
	// set context's current scope with new one

	for {
		// #1. first execute expr
		trueExpr, err := evalExpression(vm, node.TrueExpr)
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
		if _, err := evalStmtBlock(vm, node.LoopBlock); err != nil {
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

func evalBranchStmt(vm *r.VM, node *syntax.BranchStmt) error {
	// #1. condition header
	ifExpr, err := evalExpression(vm, node.IfTrueExpr)
	if err != nil {
		return err
	}
	vIfExpr, ok := ifExpr.(*value.Bool)
	if !ok {
		return zerr.InvalidExprType("bool")
	}

	// exec if-branch
	if vIfExpr.GetValue() {
		// create inner scope for if statement
		_, err := evalStmtBlock(vm, node.IfTrueBlock)
		return err
	}
	// exec else-if branches
	for idx, otherExpr := range node.OtherExprs {
		otherExprI, err := evalExpression(vm, otherExpr)
		if err != nil {
			return err
		}
		vOtherExprI, ok := otherExprI.(*value.Bool)
		if !ok {
			return zerr.InvalidExprType("bool")
		}
		// exec else-if branch
		if vOtherExprI.GetValue() {
			// create inner scope for if statement
			_, err := evalStmtBlock(vm, node.OtherBlocks[idx])
			return err
		}
	}
	// exec else branch if possible
	if node.HasElse {
		_, err := evalStmtBlock(vm, node.IfFalseBlock)
		return err
	}
	return nil
}

func evalIterateStmt(vm *r.VM, node *syntax.IterateStmt) error {
	// pre-defined key, value variable name
	var keySlot, valueSlot *r.IDName
	var matchErr error
	var nameLen = len(node.IndexNames)

	// 以A，B遍历C： D
	// execute expr: C
	targetExpr, err := evalExpression(vm, node.IterateExpr)
	if err != nil {
		return err
	}

	// execIterationBlock, including set "currentKey" and "currentValue" to scope,
	// and preDefined indication variables
	execIterationBlockFn := func(key r.Element, v r.Element) error {
		// set pre-defined value
		if nameLen == 1 {
			if err := vm.DeclareElement(valueSlot, v); err != nil {
				return err
			}
		} else if nameLen == 2 {
			if err := vm.DeclareElement(keySlot, key); err != nil {
				return err
			}
			if err := vm.DeclareElement(valueSlot, v); err != nil {
				return err
			}
		}
		_, err := evalStmtBlock(vm, node.IterateBlock)
		return err
	}

	// define indication variables as "currentKey" and "currentValue" under new iterScope
	// of course since there's no any iteration is executed yet, the initial values are all "Null"
	switch nameLen {
	case 0:
		// do nothing
	case 1:
		// Accept IDName ONLY
		valueSlot, matchErr = MatchIDName(node.IndexNames[0])
		if matchErr != nil {
			return matchErr
		}

		// init valueSlot as Null
		if err := vm.SetElement(valueSlot, value.NewNull()); err != nil {
			return err
		}
	case 2:
		keySlot, matchErr = MatchIDName(node.IndexNames[0])
		if matchErr != nil {
			return matchErr
		}
		valueSlot, matchErr = MatchIDName(node.IndexNames[1])
		if matchErr != nil {
			return matchErr
		}

		// init symbol value as Null
		if err := vm.SetElement(keySlot, value.NewNull()); err != nil {
			return err
		}
		if err := vm.SetElement(valueSlot, value.NewNull()); err != nil {
			return err
		}
	default:
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

// // execute expressions
func evalExpression(vm *r.VM, expr syntax.Expression) (r.Element, error) {
	switch e := expr.(type) {
	case *syntax.VarAssignExpr:
		return evalVarAssignExpr(vm, e)
	case *syntax.LogicExpr:
		if e.Type == syntax.LogicAND || e.Type == syntax.LogicOR {
			return evalLogicCombiner(vm, e)
		}
		return evalLogicComparator(vm, e)
	case *syntax.ArithExpr:
		if e.Type == syntax.ArithModulo {
			return evalArithTypeModuloExpr(vm, e)
		}
		return evalArithExpr(vm, e)
	case *syntax.MemberExpr:
		iv, err := getMemberExprIV(vm, e)
		if err != nil {
			return nil, err
		}
		return iv.ReduceRHS(vm)
	case *syntax.String, *syntax.ID, *syntax.ArrayExpr, *syntax.HashMapExpr:
		return evalPrimeExpr(vm, e)
	case *syntax.FuncCallExpr:
		return evalFunctionCall(vm, e)
	case *syntax.MemberMethodExpr:
		return evalMemberMethodExpr(vm, e)
	case *syntax.ObjNewExpr:
		return evalNewObject(vm, e)
	default:
		return nil, zerr.InvalidExprType()
	}
}

//// checkout eval_function.go for evalFunctionCall()

// 以 A （执行：B、C、D）
func evalMemberMethodExpr(vm *r.VM, expr *syntax.MemberMethodExpr) (r.Element, error) {
	// 1. parse root expr
	rootExpr, err := evalExpression(vm, expr.Root)
	if err != nil {
		return nil, err
	}

	var vlast r.Element = rootExpr

	for _, methodExpr := range expr.MethodChain {
		// eval method
		funcName, err := MatchIDName(methodExpr.FuncName)
		if err != nil {
			return nil, err
		}

		// exec params
		params, err := exprsToValues(vm, methodExpr.Params)
		if err != nil {
			return nil, err
		}

		v, err := execMethodFunction(vm, vlast, funcName, params)
		if err != nil {
			return nil, err
		}
		vlast = v
	}

	// add yield result
	if expr.YieldResult != nil {
		vtag, err := MatchIDName(expr.YieldResult)
		if err != nil {
			return nil, err
		}

		// bind yield result
		if err := vm.BindSymbolDecl(vtag, vlast, false); err != nil {
			return nil, err
		}
	}

	return vlast, nil
}

// evaluate logic combination expressions
// such as A 且 B
// or A 或 B
func evalLogicCombiner(vm *r.VM, expr *syntax.LogicExpr) (*value.Bool, error) {
	logicType := expr.Type
	// #1. eval left
	left, err := evalExpression(vm, expr.LeftExpr)
	if err != nil {
		return nil, err
	}
	// #2. assert left expr type to be ZnBool
	vleft, ok := left.(*value.Bool)
	if !ok {
		return nil, zerr.InvalidExprType("bool")
	}
	// #3. check if the result could be retrieved earlier (short-circuit)
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
	right, err := evalExpression(vm, expr.RightExpr)
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

// 抛出XX异常：“xxx”！
func evalThrowExceptionStmt(vm *r.VM, node *syntax.ThrowExceptionStmt) error {
	// profoundly return an ERROR to terminate the execution flow
	expClassID, err := MatchIDName(node.ExceptionClass)
	if err != nil {
		return err
	}
	expClassModel, err := vm.FindElement(expClassID)
	if err != nil {
		return err
	}

	cmodel, ok := expClassModel.(*value.ClassModel)
	if !ok {
		return zerr.InvalidExceptionType(expClassID.GetLiteral())
	}
	// exec expressions, similiar to "新建XX" statement
	var exprs []r.Element
	for _, param := range node.Params {
		exprI, err := evalExpression(vm, param)
		if err != nil {
			return err
		}
		exprs = append(exprs, exprI)
	}

	// build exception value!
	exceptionObj, err := cmodel.Construct(vm, exprs)
	if err != nil {
		return err
	}
	return zerr.NewExceptionSignal(exceptionObj)
}

// evaluate logic comparator
// ensure both expressions are comparable
func evalLogicComparator(vm *r.VM, expr *syntax.LogicExpr) (*value.Bool, error) {
	logicType := expr.Type
	// #1. eval left
	left, err := evalExpression(vm, expr.LeftExpr)
	if err != nil {
		return nil, err
	}
	// #2. eval right
	right, err := evalExpression(vm, expr.RightExpr)
	if err != nil {
		return nil, err
	}

	var cmpRes bool
	var cmpErr error
	// #3. do comparison
	switch logicType {
	case syntax.LogicXEQ:
		cmpRes, cmpErr = compareLogicXEQ(left, right)
	case syntax.LogicXNEQ:
		cmpRes, cmpErr = compareLogicXEQ(left, right)
		cmpRes = !cmpRes // reverse result
	case syntax.LogicEQ: // logicEQ, only used in Number
		cmpRes, cmpErr = compareLogicEQ(left, right)
	case syntax.LogicNEQ:
		cmpRes, cmpErr = compareLogicEQ(left, right)
		cmpRes = !cmpRes // reverse result
	case syntax.LogicGT:
		cmpRes, cmpErr = compareLogicGT(left, right)
	case syntax.LogicGTE:
		cmpRes, cmpErr = compareLogicGTE(left, right)
	case syntax.LogicLT:
		cmpRes, cmpErr = compareLogicLT(left, right)
	case syntax.LogicLTE:
		cmpRes, cmpErr = compareLogicLTE(left, right)
	default:
		return nil, zerr.UnexpectedCase("比较类型", fmt.Sprintf("%d", logicType))
	}

	return value.NewBool(cmpRes), cmpErr
}

// [elem] 为 [elem] -> [bool]
func compareLogicXEQ(left r.Element, right r.Element) (bool, error) {
	switch vl := left.(type) {
	case *value.Null:
		if _, ok := right.(*value.Null); ok {
			return true, nil
		}
		return false, nil
	case *value.Number:
		// compare right value - number only
		if vr, ok := right.(*value.Number); ok {
			return vl.GetValue() == vr.GetValue(), nil
		}
		return false, nil
	case *value.String:
		// compare right value - string only
		if vr, ok := right.(*value.String); ok {
			cmpResult := strings.Compare(vl.GetValue(), vr.GetValue()) == 0
			return cmpResult, nil
		}
		return false, nil
	case *value.Bool:
		// compare right value - bool only
		if vr, ok := right.(*value.Bool); ok {
			cmpResult := vl.GetValue() == vr.GetValue()
			return cmpResult, nil
		}
		return false, nil
	case *value.Array:
		if vr, ok := right.(*value.Array); ok {
			vla := vl.GetValue()
			vra := vr.GetValue()
			if len(vla) != len(vra) {
				return false, nil
			}
			// compare each item
			for idx := range vla {
				cmpVal, err := compareLogicXEQ(vla[idx], vra[idx])
				if err != nil {
					return false, err
				}
				// break the loop only when cmpVal = false
				if !cmpVal {
					return false, nil
				}
			}
			return true, nil
		}
		return false, nil
	case *value.HashMap:
		if vr, ok := right.(*value.HashMap); ok {
			vla := vl.GetValue()
			vra := vr.GetValue()

			if len(vla) != len(vra) {
				return false, nil
			}
			// cmp each item
			for idx := range vla {
				// ensure the key exists on vr
				vrr, ok := vra[idx]
				if !ok {
					return false, nil
				}
				cmpVal, err := compareLogicXEQ(vla[idx], vrr)
				if err != nil {
					return false, err
				}
				return cmpVal, nil
			}
			return true, nil
		}
		return false, nil
	}
	return false, zerr.InvalidCompareLType("number", "string", "bool", "array", "hashmap")
}

// [number] == [number] -> [bool]
func compareLogicEQ(left r.Element, right r.Element) (bool, error) {
	if vl, ok := left.(*value.Number); ok {
		if vr, ok2 := right.(*value.Number); ok2 {
			return vl.GetValue() == vr.GetValue(), nil
		}
		return false, zerr.InvalidCompareRType("number")
	}
	return false, zerr.InvalidCompareLType("number")
}

// [number] < [number] -> [bool]
func compareLogicLT(left r.Element, right r.Element) (bool, error) {
	if vl, ok := left.(*value.Number); ok {
		if vr, ok2 := right.(*value.Number); ok2 {
			return vl.GetValue() < vr.GetValue(), nil
		}
		return false, zerr.InvalidCompareRType("number")
	}
	return false, zerr.InvalidCompareLType("number")
}

// [number] <= [number] -> [bool]
func compareLogicLTE(left r.Element, right r.Element) (bool, error) {
	if vl, ok := left.(*value.Number); ok {
		if vr, ok2 := right.(*value.Number); ok2 {
			return vl.GetValue() <= vr.GetValue(), nil
		}
		return false, zerr.InvalidCompareRType("number")
	}
	return false, zerr.InvalidCompareLType("number")
}

// [number] > [number] -> [bool]
func compareLogicGT(left r.Element, right r.Element) (bool, error) {
	if vl, ok := left.(*value.Number); ok {
		if vr, ok2 := right.(*value.Number); ok2 {
			return vl.GetValue() > vr.GetValue(), nil
		}
		return false, zerr.InvalidCompareRType("number")
	}
	return false, zerr.InvalidCompareLType("number")
}

// [number] >= [number] -> [bool]
func compareLogicGTE(left r.Element, right r.Element) (bool, error) {
	if vl, ok := left.(*value.Number); ok {
		if vr, ok2 := right.(*value.Number); ok2 {
			return vl.GetValue() >= vr.GetValue(), nil
		}
		return false, zerr.InvalidCompareRType("number")
	}
	return false, zerr.InvalidCompareLType("number")
}

func evalArithExpr(vm *r.VM, expr *syntax.ArithExpr) (*value.Number, error) {
	// exec left Expr
	leftExpr, err := evalExpression(vm, expr.LeftExpr)
	if err != nil {
		return nil, err
	}

	leftNum, ok := leftExpr.(*value.Number)
	if !ok {
		return nil, zerr.InvalidExprType("number")
	}
	// exec right expr
	rightExpr, err := evalExpression(vm, expr.RightExpr)
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
	case syntax.ArithIntDiv:
		// python style intDiv, where result close to the closet lower integer e.g. -15 // 2 = -8 (instead of -7)
		if rightNum.GetValue() == 0 {
			return nil, zerr.ArithDivZero()
		}
		return value.NewNumber(
			math.Floor(leftNum.GetValue() / rightNum.GetValue()),
		), nil
	}
	return nil, zerr.UnexpectedCase("运算项", fmt.Sprintf("%d", expr.Type))
}

// evalArithTypeModuloExpr - handle special case of ArithExpr where Type = ArithModulo (%)
// A % B has two types:
//  1. ArithModulo: [Number] % [Number] -> [Number] (e.g  5 % 2 = 1)
//  2. String format: [String] % [Array] -> [String] (e.g. “{}-{}” % 【1、2】= “1-2”)
func evalArithTypeModuloExpr(vm *r.VM, expr *syntax.ArithExpr) (r.Element, error) {
	leftExpr, err := evalExpression(vm, expr.LeftExpr)
	if err != nil {
		return nil, err
	}

	rightExpr, err := evalExpression(vm, expr.RightExpr)
	if err != nil {
		return nil, err
	}

	// handle CASE 1
	if leftNum, okL := leftExpr.(*value.Number); okL {
		if rightNum, okR := rightExpr.(*value.Number); okR {
			// a/b = q with remainder r, where b*q + r = a and 0 <= abs(r) < b
			// so q = a 'intdiv' b, r = a - q * b
			a := leftNum.GetValue()
			b := rightNum.GetValue()
			if b == 0 {
				return nil, zerr.ArithDivZero()
			}
			q := math.Floor(a / b)

			return value.NewNumber(a - q*b), nil
		}
	}

	// handle CASE 2
	if leftStr, okL := leftExpr.(*value.String); okL {
		if rightArr, okR := rightExpr.(*value.Array); okR {
			formattedStr, err := formatString(leftStr, rightArr)
			if err != nil {
				return nil, err
			}

			return formattedStr, nil
		}
	}

	return nil, zerr.InvalidExprType("")
}

// eval prime expr
func evalPrimeExpr(vm *r.VM, expr syntax.Expression) (r.Element, error) {
	switch e := expr.(type) {
	case *syntax.String:
		return value.NewString(e.GetLiteral()), nil
	case *syntax.ID:
		idValue, err := MatchIDType(e)
		if err != nil {
			return nil, err
		}
		switch t := idValue.(type) {
		case *r.IDName:
			return vm.FindElement(t)
		case *r.IDNumber:
			return value.NewNumber(t.GetValue()), nil
		default:
			// currently idValue only have IDName or IDNumber
			return nil, zerr.UnexpectedCase("ID格式", fmt.Sprintf("%T", t))
		}
	case *syntax.ArrayExpr:
		var znObjs []r.Element
		for _, item := range e.Items {
			expr, err := evalExpression(vm, item)
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
				if _, err := MatchIDType(k); err != nil {
					return nil, err
				} else {
					exprKey = k.GetLiteral()
				}
			default:
				return nil, zerr.InvalidExprType("string", "number", "id")
			}

			exprVal, err := evalExpression(vm, item.Value)
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
func evalVarAssignExpr(vm *r.VM, expr *syntax.VarAssignExpr) (r.Element, error) {
	// Right Side
	vr, err := evalExpression(vm, expr.AssignExpr)
	if err != nil {
		return nil, err
	}
	vr = value.DuplicateValue(vr)

	// Left Side
	switch v := expr.TargetVar.(type) {
	case *syntax.ID:
		// set ID
		vtag, err := MatchIDName(v)
		if err != nil {
			return nil, err
		}
		err2 := vm.SetElement(vtag, vr)
		return vr, err2
	case *syntax.MemberExpr:
		if v.MemberType == syntax.MemberID || v.MemberType == syntax.MemberIndex {
			iv, err := getMemberExprIV(vm, v)
			if err != nil {
				return nil, err
			}
			return vr, iv.ReduceLHS(vm, vr)
		}
		return nil, zerr.UnexpectedAssign()
	default:
		return nil, zerr.UnexpectedCase("被赋值", fmt.Sprintf("%T", v))
	}
}

// execAnotherModule - load source code of the module, parse the coe, execute the program, and build depCache!
func execAnotherModule(c *r.Context, name string) (*r.Module, error) {
	if finder := c.GetModuleCodeFinder(); finder != nil {
		source, err := finder(false, name)
		if err != nil {
			return nil, zerr.ModuleNotFound(name)
		}

		// #1. parse program
		p := syntax.NewParser(source, zh.NewParserZH())

		program, err := p.Compile()
		if err != nil {
			// moduleName
			return nil, WrapSyntaxError(p, name, err)
		}

		// #2. create module & enter module
		newModule := r.NewModule(name, program.Lines)
		c.EnterModule(newModule)
		defer c.ExitModule()

		// #3. eval program
		if _, err := evalProgram(c, program, nil); err != nil {
			return nil, WrapRuntimeError(c, err)
		}

		return newModule, nil
	}
	// no finder defined, return nil directly (no throw error)
	return nil, nil
}

func getMemberExprIV(vm *r.VM, expr *syntax.MemberExpr) (*value.IV, error) {
	switch expr.RootType {
	case syntax.RootTypeProp: // 其 XX
		thisValue := vm.GetThisValue()
		if thisValue == nil {
			return nil, zerr.ThisValueNotFound()
		}
		return value.NewMemberIV(thisValue, expr.MemberID.GetLiteral()), nil
	case syntax.RootTypeExpr: // A 之 B
		valRoot, err := evalExpression(vm, expr.Root)
		if err != nil {
			return nil, err
		}
		switch expr.MemberType {
		case syntax.MemberID: // A 之 B
			return value.NewMemberIV(valRoot, expr.MemberID.GetLiteral()), nil
		case syntax.MemberIndex: // A # 0
			idx, err := evalExpression(vm, expr.MemberIndex)
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

// // helpers
// exprsToValues - []syntax.Expression -> []eval.r.Value
func exprsToValues(vm *r.VM, exprs []syntax.Expression) ([]r.Element, error) {
	params := []r.Element{}
	for _, paramExpr := range exprs {
		pval, err := evalExpression(vm, paramExpr)
		if err != nil {
			return nil, err
		}
		params = append(params, pval)
	}
	return params, nil
}

// extractSignalValue - signal is a special type of error, so we try to extract signal value from input error if it's really a signal - otherwise output the REAL error directly.
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

var EvaluateProgram = evalProgram
