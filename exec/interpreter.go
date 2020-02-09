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
			err = it.handleVarDeclareStmt(s)
			if err != nil {
				break
			}
		case *syntax.VarAssignExpr:
			err = it.handleVarAssignExpr(s)
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

func (it *Interpreter) handleVarDeclareStmt(stmt *syntax.VarDeclareStmt) *error.Error {
	for _, vpair := range stmt.AssignPair {
		obj := execExpression(it, vpair.AssignExpr)
		for _, v := range vpair.Variables {
			vtag := v.GetLiteral()
			// TODO: need copy object!
			if err := it.Bind(vtag, obj); err != nil {
				return err
			}
		}

	}

	return nil
}

func (it *Interpreter) handleVarAssignExpr(stmt *syntax.VarAssignExpr) *error.Error {
	obj := execExpression(it, stmt.AssignExpr)
	vtag := stmt.TargetVar.GetLiteral()

	return it.SetData(vtag, obj)
}

func execExpression(it *Interpreter, expr syntax.Expression) ZnObject {
	if expr.IsPrimitive() {
		return execPrimitiveExpr(it, expr)
	}
	// TODO: to be continued...
	return nil
}

func execPrimitiveExpr(it *Interpreter, expr syntax.Expression) ZnObject {
	switch e := expr.(type) {
	case *syntax.Number:
		zd := new(ZnDecimal)
		zd.SetValue(e.GetLiteral())

		return zd
	case *syntax.String:
		zstr := new(ZnString)
		zstr.SetValue(e.GetLiteral())

		return zstr
	case *syntax.ID:
		vtag := e.GetLiteral()
		if obj, err := it.Lookup(vtag); err == nil {
			return obj
		}
		return nil
	case *syntax.ArrayExpr:
		znObjs := []ZnObject{}
		znArr := new(ZnArray)
		for _, item := range e.Items {
			znObjs = append(znObjs, execPrimitiveExpr(it, item))
		}

		znArr.Init(znObjs)
		return znArr
	default:
		// TODO: to be continued...
		return nil
	}
}
