package runtime

import (
	"github.com/DemoHn/Zn/pkg/syntax"
)

type VM struct {
	globals map[string]Element
	// module-level & scope-level local value stacks
	// including global & local variables
	// KEY: moduleID
	// VALUE: the ScopeStack of corresponding module
	valueStack map[int]ScopeStack

	// callStack - store all call frames
	callStack []CallFrame
	// csCursor - index of current callStack
	csCursor int

	// modules - allocates modules by ID & stores export values
	modules *ModuleGraph
}

type CallFrame struct {
	moduleID    int
	callType    uint8
	currentLine int // current exec line in the module's source code
	programAST  *syntax.Program
	// for SCRIPT callFrame, thisValue = nil
	// for FUCTION callFrame, thisValues depends on the function
	// 	 - for method function, thisValue = [Object Instance]
	//   - for direct function, thisValue = nil
	thisValue Element
}

const (
	CALL_TYPE_SCRIPT   uint8 = 1
	CALL_TYPE_FUNCTION uint8 = 2
)

type ScopeStack []Scope

type ElementMap = map[string]Element

func InitVM(globals map[string]Element) *VM {
	return &VM{
		globals:     globals,
		valueStack:  map[int]ScopeStack{},
		callStack:   []CallFrame{},
		csCursor:    -1,
		moduleGraph: nil,
	}
}
