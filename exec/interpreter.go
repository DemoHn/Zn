package exec

import (
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
	var err *error.Error
	for _, stmt := range program.Content.Children {
		switch s := stmt.(type) {
		case *syntax.VarDeclareStmt:
			err = evalVarDeclareStmt(it, s)
			if err != nil {
				break
			}
		default:
			// regard as unknown statement and ignore it
			// TODO: to be continued...
			continue
		}
	}

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
	case *syntax.Number, *syntax.String, *syntax.ID, *syntax.ArrayExpr:
		return evalPrimeExpr(it, e)
	default:
		return nil, error.NewErrorSLOT("unrecognized type")
	}
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
		return nil, error.NewErrorSLOT("参与的表达式须为「二象」类型")
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
		return nil, error.NewErrorSLOT("参与的表达式须为「二象」类型")
	}
	// then evalute data
	switch logicType {
	case syntax.LogicAND:
		return NewZnBool(vleft.Value && vright.Value), nil
	case syntax.LogicOR:
		return NewZnBool(vleft.Value || vright.Value), nil
	default:
		return nil, error.NewErrorSLOT("不合法的类型") // 这个一般走不到
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
	// #2. assert left to be comparable
	vleft, ok := left.(ZnComparable)
	if !ok {
		return nil, error.NewErrorSLOT("must be comparable")
	}
	// #3. eval right
	right, err := EvalExpression(it, expr.RightExpr)
	if err != nil {
		return nil, err
	}
	// #4. assert right to be comparable
	vright, ok := right.(ZnComparable)
	if !ok {
		return nil, error.NewErrorSLOT("must be comparable")
	}

	// #5. evaluate
	switch logicType {
	case syntax.LogicEQ:
		return vleft.Equals(vright)
	case syntax.LogicNEQ:
		zb, err := vleft.Equals(vright)
		return zb.Rev(), err
	case syntax.LogicIS:
		return vleft.Is(vright)
	case syntax.LogicISN:
		zb, err := vleft.Is(vright)
		return zb.Rev(), err
	case syntax.LogicGT:
		return vleft.GreaterThan(vright)
	case syntax.LogicGTE:
		zb1, err := vleft.GreaterThan(vright)
		if err != nil {
			return nil, err
		}
		zb2, err := vleft.Equals(vright)
		if err != nil {
			return nil, err
		}

		return NewZnBool(zb1.Value || zb2.Value), nil
	case syntax.LogicLT:
		return vleft.LessThan(vright)
	case syntax.LogicLTE:
		zb1, err := vleft.LessThan(vright)
		if err != nil {
			return nil, err
		}
		zb2, err := vleft.Equals(vright)
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
