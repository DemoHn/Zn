package exec

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/syntax"
)

// Interpreter - the main interpreter to execute the program and yield results
type Interpreter struct {
	Program *syntax.ProgramNode
	Symbol  *SymbolTable
}

// Execute - execute the program and yield the result
func (it *Interpreter) Execute() string {
	pg := it.Program
	var err *error.Error
	for _, stmt := range pg.Children {
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
			continue
		}
	}

	// yield result
	return it.print(err)
}

// print - print result
func (it *Interpreter) print(err *error.Error) string {
	return "233"
}

func (it *Interpreter) handleVarDeclareStmt(stmt *syntax.VarDeclareStmt) *error.Error {
	obj := execExpression(stmt.AssignExpr)
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
	// TODO
	return nil
}

func execExpression(expr syntax.Expression) ZnObject {
	return nil
}
