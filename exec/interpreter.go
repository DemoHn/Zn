package exec

import (
	"fmt"
	"strings"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/syntax"
)

// Interpreter - the main interpreter to execute the program and yield results
type Interpreter struct {
	Symbol *SymbolTable
}

// NewInterpreter -
func NewInterpreter() *Interpreter {
	return &Interpreter{
		Symbol: &SymbolTable{
			Symbols: map[string]ZnObject{},
		},
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
		case *syntax.VarAssignStmt:
			err = it.handleVarAssignStmt(s)
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

	strs := []string{}
	for k, symbol := range it.Symbol.Symbols {
		symStr := "ε"
		if symbol != nil {
			symStr = symbol.String()
		}
		strs = append(strs, fmt.Sprintf("‹%s› => %s", k, symStr))
	}

	return strings.Join(strs, "\n")
}

func (it *Interpreter) handleVarDeclareStmt(stmt *syntax.VarDeclareStmt) *error.Error {
	obj := execExpression(it, stmt.AssignExpr)
	for _, v := range stmt.Variables {
		vtag := v.GetLiteral()
		// TODO: need copy object!
		if !it.Symbol.Insert(vtag, obj) {
			return error.NewErrorSLOT("variable redeclaration!")
		}
	}

	return nil
}

func (it *Interpreter) handleVarAssignStmt(stmt *syntax.VarAssignStmt) *error.Error {
	obj := execExpression(it, stmt.AssignExpr)
	vtag := stmt.TargetVar.GetLiteral()

	if _, ok := it.Symbol.Lookup(vtag); ok {
		it.Symbol.SetData(vtag, obj)
		return nil
	}
	return error.NewErrorSLOT("variable not defined!")
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
		if obj, ok := it.Symbol.Lookup(vtag); ok {
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
