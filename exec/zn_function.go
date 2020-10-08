package exec

import (
	"fmt"

	"github.com/DemoHn/Zn/syntax"
)

// ZnFunction -
type ZnFunction struct {
	*ZnObject
	*ClosureRef
}

func (zf *ZnFunction) String() string {
	return fmt.Sprintf("方法： %s", zf.ClosureRef.Name)
}

// NewZnFunction - new Zn native function
func NewZnFunction(name string, executor funcExecutor) *ZnFunction {
	closureRef := NewClosureRef(name, executor)
	return &ZnFunction{
		ClosureRef: closureRef,
	}
}

// BuildZnFunctionFromNode -
func BuildZnFunctionFromNode(node *syntax.FunctionDeclareStmt) *ZnFunction {
	funcName := node.FuncName.GetLiteral()
	closureRef := BuildClosureRefFromNode(funcName, node.ParamList, node.ExecBlock)
	return &ZnFunction{
		ClosureRef: closureRef,
	}
}
