package exec

import (
	"github.com/DemoHn/Zn/syntax"
)

// Interpreter - the main interpreter to execute the program and yield results
type Interpreter struct {
	Program *syntax.ProgramNode
}

// Execute - execute the program and yield the result
func (i *Interpreter) Execute() string {
	return ""
}
