package runtime

import (
	"testing"

	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/stretchr/testify/assert"
)

// test allocate module:
// first allocate module should be the main module

var globalValuesI = map[string]Element{
	"A": MockValue{"A"},
	"B": MockValue{"B"},
}

var mockProgram = &syntax.Program{
	Lexer: &syntax.Lexer{
		Lines: []syntax.LineInfo{
			{
				LineText: []rune("LINE1"),
				StartIdx: 0,
			},
			{
				LineText: []rune("LINE2"),
				StartIdx: 5,
			},
		},
	},
}

func TestAllocateModule_ModuleGraph(t *testing.T) {
	vm := InitVM(globalValuesI)
	// allocate main module
	module := vm.AllocateModule("main", mockProgram)
	assert.Equal(t, module.GetID(), 0)
	assert.Equal(t, module.GetName(), "main")

	// allocate another module
	module2 := vm.AllocateModule("MA", mockProgram)
	assert.Equal(t, module2.GetID(), 1)
	assert.Equal(t, vm.FindModuleByName("MA"), module2)

	// assert module graph
	expectedGraph := [][2]int{{0, 1}}
	assert.Equal(t, expectedGraph, vm.moduleGraph.graph)

	// enter module2
	vm.PushCallFrame(NewScriptCallFrame(module2))

	// allocate another module 3
	module3 := vm.AllocateModule("MB", mockProgram)
	assert.Equal(t, module3.GetID(), 2)
	assert.Equal(t, vm.FindModuleByName("MB"), module3)

	expectedGraph2 := [][2]int{{0, 1}, {1, 2}}
	assert.Equal(t, expectedGraph2, vm.moduleGraph.graph)

	// enter function
	vm.PushCallFrame(NewFunctionCallFrame(module3, nil))
	// NOTE: module graph should not change
	expectedGraph3 := [][2]int{{0, 1}, {1, 2}}
	assert.Equal(t, expectedGraph3, vm.moduleGraph.graph)
}

func TestFindModuleByName(t *testing.T) {
	vm := InitVM(globalValuesI)
	//module1
	vm.AllocateModule("__MAIN__", mockProgram)
	assert.Equal(t, vm.FindModuleByName("__MAIN__").GetID(), 0)

	// lookup non-existing module
	assert.Nil(t, vm.FindModuleByName("non-existing"))

	// module2
	module2 := vm.AllocateModule("MODULE2", mockProgram)
	// enter module2
	vm.PushCallFrame(NewScriptCallFrame(module2))

	assert.Equal(t, module2.GetID(), 1)
	assert.Equal(t, vm.FindModuleByName("MODULE2"), module2)

	// re-fetch __MAIN__
	assert.Equal(t, vm.FindModuleByName("__MAIN__").GetID(), 0)
}

func TestPushCallFrame_PopCallFrame(t *testing.T) {
	vm := InitVM(globalValuesI)

	module1 := vm.AllocateModule("main", mockProgram)
	vm.PushCallFrame(NewScriptCallFrame(module1))

	// assert csCount and csModuleID and valueStack
	assert.Equal(t, vm.csCount, 1)
	assert.Equal(t, vm.csModuleID, 0)
	assert.Equal(t, len(vm.valueStack), 1)
	assert.Equal(t, vm.valueStack[module1.GetID()] == nil, false)

	vm.PopCallFrame()
	// assert csCount and csModuleID and valueStack
	assert.Equal(t, vm.csCount, 0)
	assert.Equal(t, vm.csModuleID, -1)
	assert.Equal(t, len(vm.valueStack), 1) // value stack won't be changed

	// again, push script frame and function frame
	module2 := vm.AllocateModule("module2", mockProgram)
	vm.PushCallFrame(NewScriptCallFrame(module2))
	vm.PushCallFrame(NewFunctionCallFrame(module1, nil))
	vm.PushCallFrame(NewExceptionCallFrame(module2, MockValue{"ERROR"}))

	// assert csCount and csModuleID and valueStack
	assert.Equal(t, vm.csCount, 3)
	assert.Equal(t, vm.csModuleID, module2.GetID())
	assert.Equal(t, len(vm.valueStack), 2)
	// assert callframe type
	assert.Equal(t, vm.callStack[0].callType, CALL_TYPE_SCRIPT)
	assert.Equal(t, vm.callStack[1].callType, CALL_TYPE_FUNCTION)
	assert.Equal(t, vm.callStack[2].callType, CALL_TYPE_EXCEPTION_BLOCK)
	// assert callframe thisValue
	assert.Equal(t, vm.callStack[0].thisValue, nil)
	assert.Equal(t, vm.callStack[1].thisValue, nil)
	assert.Equal(t, vm.callStack[2].thisValue, MockValue{"ERROR"})
	assert.Equal(t, vm.GetThisValue(), MockValue{"ERROR"})

	// pop one callframe and csMOdule = 0
	vm.PopCallFrame()
	// assert csCount and csModuleID and valueStack
	assert.Equal(t, vm.csCount, 2)
	assert.Equal(t, vm.csModuleID, module1.GetID())
	assert.Equal(t, len(vm.valueStack), 2)
}

func TestDeclareElement(t *testing.T) {

}
