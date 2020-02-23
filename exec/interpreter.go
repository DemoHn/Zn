package exec

import (
	"fmt"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/syntax"
)

// Interpreter - the main interpreter to execute the program and yield results
type Interpreter struct {
	*SymbolTable
}

// NewInterpreter -
func NewInterpreter() *Interpreter {
	return &Interpreter{
		SymbolTable: NewSymbolTable(),
	}
}

// Execute - execute the program and yield the result
func (it *Interpreter) Execute(program *syntax.Program) string {
	if program.Content == nil {
		return ""
	}
	err := evalBlockStatement(it, program.Content, true)

	// yield result
	return it.print(err)
}

// print - print result
func (it *Interpreter) print(err *error.Error) string {
	if err != nil {
		return err.Error()
	}

	return it.printSymbols()
}

//// Execute (Evaluate) statements

// EvalStatement - eval statement
func EvalStatement(it *Interpreter, stmt syntax.Statement) *error.Error {
	switch v := stmt.(type) {
	case *syntax.VarDeclareStmt:
		return evalVarDeclareStmt(it, v)
	case *syntax.WhileLoopStmt:
		return evalWhileLoopStmt(it, v)
	case *syntax.BranchStmt:
		return evalBranchStmt(it, v)
	case *syntax.EmptyStmt:
		return nil
	case syntax.Expression:
		_, err := EvalExpression(it, v)
		return err
	default:
		return error.NewErrorSLOT("invalid statement type")
	}
}

func evalVarDeclareStmt(it *Interpreter, stmt *syntax.VarDeclareStmt) *error.Error {
	for _, vpair := range stmt.AssignPair {
		obj, err := EvalExpression(it, vpair.AssignExpr)
		if err != nil {
			return err
		}
		for _, v := range vpair.Variables {
			vtag := v.GetLiteral()
			// TODO: need copy object!
			if err := it.Bind(vtag, obj, false); err != nil {
				return err
			}
		}

	}

	return nil
}

func evalBlockStatement(it *Interpreter, block *syntax.BlockStmt, globalScope bool) *error.Error {
	if !globalScope {
		it.EnterScope()
		defer it.ExitScope()
	}

	for _, stmt := range block.Children {
		err := EvalStatement(it, stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func evalWhileLoopStmt(it *Interpreter, loopStmt *syntax.WhileLoopStmt) *error.Error {
	for {
		// #1. first execute expr
		trueExpr, err := EvalExpression(it, loopStmt.TrueExpr)
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
		if err := evalBlockStatement(it, loopStmt.LoopBlock, false); err != nil {
			return nil
		}
	}
}

func evalBranchStmt(it *Interpreter, branchStmt *syntax.BranchStmt) *error.Error {
	// #1. if branch
	ifExpr, err := EvalExpression(it, branchStmt.IfTrueExpr)
	if err != nil {
		return err
	}
	vIfExpr, ok := ifExpr.(*ZnBool)
	if !ok {
		return error.InvalidExprType("bool")
	}
	// exec if-branch
	if vIfExpr.Value == true {
		return evalBlockStatement(it, branchStmt.IfTrueBlock, false)
	}
	// exec else-if branches
	for idx, otherExpr := range branchStmt.OtherExprs {
		otherExprI, err := EvalExpression(it, otherExpr)
		if err != nil {
			return err
		}
		vOtherExprI, ok := otherExprI.(*ZnBool)
		if !ok {
			return error.InvalidExprType("bool")
		}
		// exec else-if branch
		if vOtherExprI.Value == true {
			return evalBlockStatement(it, branchStmt.OtherBlocks[idx], false)
		}
	}
	// exec else branch if possible
	if branchStmt.HasElse == true {
		return evalBlockStatement(it, branchStmt.IfFalseBlock, false)
	}
	return nil
}

//// Execute (Evaluate) expressions

// EvalExpression - execute expression
func EvalExpression(it *Interpreter, expr syntax.Expression) (ZnValue, *error.Error) {
	switch e := expr.(type) {
	case *syntax.VarAssignExpr:
		return evalVarAssignExpr(it, e)
	case *syntax.LogicExpr:
		if e.Type == syntax.LogicAND || e.Type == syntax.LogicOR {
			return evalLogicCombiner(it, e)
		}
		return evalLogicComparator(it, e)
	case *syntax.ArrayListIndexExpr:
		// TODO: differ LHV & RHV
		return evalArrayListIndexExprRHV(it, e)
	case *syntax.Number, *syntax.String, *syntax.ID, *syntax.ArrayExpr, *syntax.HashMapExpr:
		// TODO: add HashMapExpr
		return evalPrimeExpr(it, e)
	case *syntax.FuncCallExpr:
		return evalFunctionCall(it, e)
	default:
		return nil, error.NewErrorSLOT("unrecognized type")
	}
}

// （显示：A，B，C）
func evalFunctionCall(it *Interpreter, expr *syntax.FuncCallExpr) (ZnValue, *error.Error) {
	vtag := expr.FuncName.GetLiteral()
	// find function definition
	val, err := it.Lookup(vtag)
	if err != nil {
		return nil, err
	}
	// assert value
	vval, ok := val.(*ZnFunction)
	if !ok {
		return nil, error.NewErrorSLOT(fmt.Sprintf("「%s」须为一个方法", vtag))
	}
	// exec params
	params := []ZnValue{}
	for _, paramExpr := range expr.Params {
		pval, err := EvalExpression(it, paramExpr)
		if err != nil {
			return nil, err
		}
		params = append(params, pval)
	}
	// exec function
	return vval.Exec(params, it.SymbolTable)
}

// evaluate logic combination expressions
// such as A 且 B
// or A 或 B
func evalLogicCombiner(it *Interpreter, expr *syntax.LogicExpr) (*ZnBool, *error.Error) {
	logicType := expr.Type
	// #1. eval left
	left, err := EvalExpression(it, expr.LeftExpr)
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
	right, err := EvalExpression(it, expr.RightExpr)
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
func evalLogicComparator(it *Interpreter, expr *syntax.LogicExpr) (*ZnBool, *error.Error) {
	logicType := expr.Type
	// #1. eval left
	left, err := EvalExpression(it, expr.LeftExpr)
	if err != nil {
		return nil, err
	}
	// #3. eval right
	right, err := EvalExpression(it, expr.RightExpr)
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
		return nil, error.NewErrorSLOT("invalid logic type")
	}
}

// eval prime expr
func evalPrimeExpr(it *Interpreter, expr syntax.Expression) (ZnValue, *error.Error) {
	switch e := expr.(type) {
	case *syntax.Number:
		return NewZnDecimal(e.GetLiteral())
	case *syntax.String:
		return NewZnString(e.GetLiteral()), nil
	case *syntax.ID:
		vtag := e.GetLiteral()
		return it.Lookup(vtag)
	case *syntax.ArrayExpr:
		znObjs := []ZnValue{}
		for _, item := range e.Items {
			expr, err := EvalExpression(it, item)
			if err != nil {
				return nil, err
			}
			znObjs = append(znObjs, expr)
		}

		return NewZnArray(znObjs), nil
	case *syntax.HashMapExpr:
		znPairs := []KVPair{}
		for _, item := range e.KVPair {
			expr, err := EvalExpression(it, item.Key)
			if err != nil {
				return nil, err
			}
			exprKey, ok := expr.(*ZnString)
			if !ok {
				return nil, error.NewErrorSLOT("key should be string")
			}
			exprVal, err := EvalExpression(it, item.Value)
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
		// to be continued...
		return nil, error.NewErrorSLOT("invalid type")
	}
}

// eval var assign
func evalVarAssignExpr(it *Interpreter, expr *syntax.VarAssignExpr) (ZnValue, *error.Error) {
	val, err := EvalExpression(it, expr.AssignExpr)
	if err != nil {
		return nil, err
	}
	vtag := expr.TargetVar.GetLiteral()

	err2 := it.SetData(vtag, val)
	return val, err2
}

// eval A#n A#{ e }, etc.
// NOTE: RHV stands for Right Hand Value, which means the expression will yield values directly
// like what a RHV does.
func evalArrayListIndexExprRHV(it *Interpreter, expr *syntax.ArrayListIndexExpr) (ZnValue, *error.Error) {
	// #1. eval root expr
	val, err := EvalExpression(it, expr.Root)
	if err != nil {
		return nil, err
	}
	valIdx, err := EvalExpression(it, expr.Index)
	if err != nil {
		return nil, err
	}
	// #2. assert types
	switch vl := val.(type) {
	case *ZnArray:
		// assert valIdx data
		vr, ok := valIdx.(*ZnDecimal)
		if !ok {
			return nil, error.InvalidExprType("integer")
		}
		idx, err := vr.asInteger()
		if err != nil {
			return nil, error.InvalidExprType("integer")
		}
		if idx < 0 || idx >= len(vl.Value) {
			return nil, error.IndexOutOfRange()
		}
		return vl.Value[idx], nil
	case *ZnHashMap:
		vr, ok := valIdx.(*ZnString)
		if !ok {
			return nil, error.InvalidExprType("string")
		}
		// retrieve value by key
		data, ok := vl.Value[vr.Value]
		if !ok {
			return nil, error.IndexKeyNotFound(vr.Value)
		}
		return data, nil
	default:
		return nil, error.InvalidExprType("array", "hashmap")
	}

	return nil, nil
}
