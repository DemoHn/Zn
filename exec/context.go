package exec

import (
	"fmt"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
	"github.com/DemoHn/Zn/syntax"
)

// Context - code lifecycle management
// TODO: this is a tmp solution. in the future, we will
// gradually obselete this tree-walk based interperter.
type Context struct {
	*SymbolTable
	*ArithInstance
	// lastValue is set during the execution, usually stands for 'the return value' of a function.
	lastValue ZnValue
}

const defaultPrecision = 8

// Result - context execution result structure
// NOTICE: when HasError = true, Value = nil, while execution yields error
//         when HasError = false, Error = nil, Value = <result Value>
//
// Currently only one value is supported as return argument.
type Result struct {
	HasError bool
	Value    ZnValue
	Error    *error.Error
}

// NewContext - create new Zn Context for furthur execution
func NewContext() *Context {
	ctx := new(Context)
	ctx.SymbolTable = NewSymbolTable()
	ctx.ArithInstance = NewArithInstance(defaultPrecision)
	ctx.lastValue = NewZnNull()
	return ctx
}

// ExecuteCode - execute program from input Zn code (whether from file or REPL)
func (ctx *Context) ExecuteCode(in *lex.InputStream) Result {
	l := lex.NewLexer(in)
	p := syntax.NewParser(l)
	// start
	block, err := p.Parse()
	if err != nil {
		return Result{true, nil, err}
	}
	// construct root (program) node
	program := syntax.NewProgramNode(block)

	if err := EvalProgram(ctx, program); err != nil {
		return Result{true, nil, err}
	}
	return Result{false, ctx.lastValue, nil}
}

// ExecuteBlockAST - execute blockStmt AST
// usually for executing function template
func (ctx *Context) ExecuteBlockAST(block *syntax.BlockStmt) Result {
	if err := EvalStmtBlock(ctx, block, true); err != nil {
		// handle returnValue Interrupts
		if err.GetErrorClass() != error.InterruptsClass {
			return Result{true, nil, err}
		}
	}

	lastValue := ctx.lastValue
	if ctx.lastValue == nil {
		lastValue = NewZnNull()
	}
	return Result{false, lastValue, nil}
}

// ResetLastValue - set ctx.lastValue -> nil
func (ctx *Context) ResetLastValue() {
	ctx.lastValue = nil
}

//// Execute (Evaluate) statements

// EvalProgram - evaluate global program (root node)
func EvalProgram(ctx *Context, program *syntax.Program) *error.Error {
	return EvalStmtBlock(ctx, program.Content, true)
}

// EvalStatement - eval statement
func EvalStatement(ctx *Context, stmt syntax.Statement) *error.Error {
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
		fn := NewZnFunction(v)
		return ctx.Bind(v.FuncName.GetLiteral(), fn, false)
	case *syntax.FunctionReturnStmt:
		res, err := EvalExpression(ctx, v.ReturnExpr)
		if err != nil {
			return err
		}
		ctx.lastValue = res
		// send interrupt (NOT AN ACTUAL ERROR)
		return error.ReturnValueInterrupt()
	case syntax.Expression:
		res, err := EvalExpression(ctx, v)
		if err != nil {
			return err
		}
		// set lastValue
		ctx.lastValue = res
		return nil
	default:
		return error.NewErrorSLOT("invalid statement type")
	}
}

func evalVarDeclareStmt(ctx *Context, stmt *syntax.VarDeclareStmt) *error.Error {
	for _, vpair := range stmt.AssignPair {
		obj, err := EvalExpression(ctx, vpair.AssignExpr)
		if err != nil {
			return err
		}
		for _, v := range vpair.Variables {
			vtag := v.GetLiteral()
			// TODO: need copy object!
			if err := ctx.Bind(vtag, obj, false); err != nil {
				return err
			}
		}

	}

	return nil
}

// EvalStmtBlock -
func EvalStmtBlock(ctx *Context, block *syntax.BlockStmt, sameScope bool) *error.Error {
	if !sameScope {
		ctx.EnterScope()
		defer ctx.ExitScope()
	}
	for _, stmt := range block.Children {
		err := EvalStatement(ctx, stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func evalWhileLoopStmt(ctx *Context, loopStmt *syntax.WhileLoopStmt) *error.Error {
	for {
		// #1. first execute expr
		trueExpr, err := EvalExpression(ctx, loopStmt.TrueExpr)
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
		if err := EvalStmtBlock(ctx, loopStmt.LoopBlock, false); err != nil {
			return nil
		}
	}
}

func evalBranchStmt(ctx *Context, branchStmt *syntax.BranchStmt) *error.Error {
	// #1. if branch
	ifExpr, err := EvalExpression(ctx, branchStmt.IfTrueExpr)
	if err != nil {
		return err
	}
	vIfExpr, ok := ifExpr.(*ZnBool)
	if !ok {
		return error.InvalidExprType("bool")
	}
	// exec if-branch
	if vIfExpr.Value == true {
		return EvalStmtBlock(ctx, branchStmt.IfTrueBlock, false)
	}
	// exec else-if branches
	for idx, otherExpr := range branchStmt.OtherExprs {
		otherExprI, err := EvalExpression(ctx, otherExpr)
		if err != nil {
			return err
		}
		vOtherExprI, ok := otherExprI.(*ZnBool)
		if !ok {
			return error.InvalidExprType("bool")
		}
		// exec else-if branch
		if vOtherExprI.Value == true {
			return EvalStmtBlock(ctx, branchStmt.OtherBlocks[idx], false)
		}
	}
	// exec else branch if possible
	if branchStmt.HasElse == true {
		return EvalStmtBlock(ctx, branchStmt.IfFalseBlock, false)
	}
	return nil
}

//// Execute (Evaluate) expressions

// EvalExpression - execute expression
func EvalExpression(ctx *Context, expr syntax.Expression) (ZnValue, *error.Error) {
	switch e := expr.(type) {
	case *syntax.VarAssignExpr:
		return evalVarAssignExpr(ctx, e)
	case *syntax.LogicExpr:
		if e.Type == syntax.LogicAND || e.Type == syntax.LogicOR {
			return evalLogicCombiner(ctx, e)
		}
		return evalLogicComparator(ctx, e)
	case *syntax.ArrayListIndexExpr:
		iv, err := getArrayListIV(ctx, e)
		if err != nil {
			return nil, err
		}
		// regard iv as a RHS value
		return iv.Reduce(nil, false)
	case *syntax.Number, *syntax.String, *syntax.ID, *syntax.ArrayExpr, *syntax.HashMapExpr:
		// TODO: add HashMapExpr
		return evalPrimeExpr(ctx, e)
	case *syntax.FuncCallExpr:
		return evalFunctionCall(ctx, e)
	default:
		return nil, error.NewErrorSLOT("unrecognized type")
	}
}

// （显示：A，B，C）
func evalFunctionCall(ctx *Context, expr *syntax.FuncCallExpr) (ZnValue, *error.Error) {
	vtag := expr.FuncName.GetLiteral()
	// find function definctxion
	val, err := ctx.Lookup(vtag)
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
		pval, err := EvalExpression(ctx, paramExpr)
		if err != nil {
			return nil, err
		}
		params = append(params, pval)
	}
	// exec function
	return vval.Exec(params, ctx)
}

// evaluate logic combination expressions
// such as A 且 B
// or A 或 B
func evalLogicCombiner(ctx *Context, expr *syntax.LogicExpr) (*ZnBool, *error.Error) {
	logicType := expr.Type
	// #1. eval left
	left, err := EvalExpression(ctx, expr.LeftExpr)
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
	right, err := EvalExpression(ctx, expr.RightExpr)
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
func evalLogicComparator(ctx *Context, expr *syntax.LogicExpr) (*ZnBool, *error.Error) {
	logicType := expr.Type
	// #1. eval left
	left, err := EvalExpression(ctx, expr.LeftExpr)
	if err != nil {
		return nil, err
	}
	// #3. eval right
	right, err := EvalExpression(ctx, expr.RightExpr)
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
func evalPrimeExpr(ctx *Context, expr syntax.Expression) (ZnValue, *error.Error) {
	switch e := expr.(type) {
	case *syntax.Number:
		return NewZnDecimal(e.GetLiteral())
	case *syntax.String:
		return NewZnString(e.GetLiteral()), nil
	case *syntax.ID:
		vtag := e.GetLiteral()
		return ctx.Lookup(vtag)
	case *syntax.ArrayExpr:
		znObjs := []ZnValue{}
		for _, ctxem := range e.Items {
			expr, err := EvalExpression(ctx, ctxem)
			if err != nil {
				return nil, err
			}
			znObjs = append(znObjs, expr)
		}

		return NewZnArray(znObjs), nil
	case *syntax.HashMapExpr:
		znPairs := []KVPair{}
		for _, ctxem := range e.KVPair {
			expr, err := EvalExpression(ctx, ctxem.Key)
			if err != nil {
				return nil, err
			}
			exprKey, ok := expr.(*ZnString)
			if !ok {
				return nil, error.NewErrorSLOT("key should be string")
			}
			exprVal, err := EvalExpression(ctx, ctxem.Value)
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
func evalVarAssignExpr(ctx *Context, expr *syntax.VarAssignExpr) (ZnValue, *error.Error) {
	// Right Side
	val, err := EvalExpression(ctx, expr.AssignExpr)
	if err != nil {
		return nil, err
	}

	// Left Side
	switch v := expr.TargetVar.(type) {
	case *syntax.ID:
		// set ID
		vtag := v.GetLiteral()
		err2 := ctx.SetData(vtag, val)
		return val, err2
	case *syntax.ArrayListIndexExpr:
		iv, err := getArrayListIV(ctx, v)
		if err != nil {
			return nil, err
		}
		return iv.Reduce(val, true)
	default:
		return nil, error.NewErrorSLOT("invalid type")
	}
}

// eval A#n A#{ e }, etc.
// NOTE: RHV stands for Right Hand Value, which means the expression will yield values directly
// like what a RHV does.
/**
func evalArrayListIndexExpr(ctx *Context, expr *syntax.ArrayListIndexExpr) (ZnValue, *error.Error) {
	// #1. eval root expr
	val, err := EvalExpression(ctx, expr.Root)
	if err != nil {
		return nil, err
	}
	valIdx, err := EvalExpression(ctx, expr.Index)
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
		return &ZnArrayIV{vl, vr}, nil
	case *ZnHashMap:
		vr, ok := valIdx.(*ZnString)
		if !ok {
			return nil, error.InvalidExprType("string")
		}
		// retrieve value by key
		_, ok = vl.Value[vr.Value]
		if !ok {
			return nil, error.IndexKeyNotFound(vr.Value)
		}
		return &ZnHashMapIV{vl, vr}, nil
	default:
		return nil, error.InvalidExprType("array", "hashmap")
	}
}
*/

func getArrayListIV(ctx *Context, expr *syntax.ArrayListIndexExpr) (ZnIV, *error.Error) {
	val, err := EvalExpression(ctx, expr.Root)
	if err != nil {
		return nil, err
	}
	idx, err := EvalExpression(ctx, expr.Index)
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
		var s *ZnString = nil
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
