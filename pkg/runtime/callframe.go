package runtime

import "github.com/DemoHn/Zn/pkg/syntax"

const (
	CALL_TYPE_SCRIPT          uint8 = 1
	CALL_TYPE_FUNCTION        uint8 = 2
	CALL_TYPE_EXCEPTION_BLOCK uint8 = 3
)

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

func NewScriptCallFrame(moduleID int, programAST *syntax.Program) *CallFrame {
	return &CallFrame{
		moduleID:    moduleID,
		callType:    CALL_TYPE_SCRIPT,
		currentLine: 0,
		programAST:  programAST,
		// thisValue is valid only for CALL_TYPE_FUNCTION
		thisValue: nil,
	}
}

func NewFunctionCallFrame(moduleID int, programAST *syntax.Program, thisValue Element) *CallFrame {
	return &CallFrame{
		moduleID:    moduleID,
		callType:    CALL_TYPE_FUNCTION,
		currentLine: 0,
		programAST:  programAST,
		thisValue:   thisValue,
	}
}

func (cf *CallFrame) GetCurrentLine() int {
	return cf.currentLine
}

func (cf *CallFrame) SetCurrentLine(line int) {
	cf.currentLine = line
}

func (cf *CallFrame) IsFunctionCallFrame() bool {
	return cf.callType == CALL_TYPE_FUNCTION
}

func (cf *CallFrame) IsScriptCallFrame() bool {
	return cf.callType == CALL_TYPE_SCRIPT
}
