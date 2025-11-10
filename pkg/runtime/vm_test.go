package runtime

import (
	"testing"

	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/stretchr/testify/assert"
)

// test allocate module:
// first allocate module should be the main module

var globalValues = map[string]Element{
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

func TestAllocateModule(t *testing.T) {
	vm := InitVM(globalValues)
	// allocate main module
	module := vm.AllocateModule("main", mockProgram)
	assert.Equal(t, module.GetID(), 0)
	assert.Equal(t, module.GetName(), "main")

	// allocate another module
	module2 := vm.AllocateModule("MA", mockProgram)
	assert.Equal(t, module2.GetID(), 1)
	assert.Equal(t, vm.FindModuleByName("MA"), module2)

	// assert module graph
	expectedGraph := [][2]int{{-1, 0}, {0, 1}}
	assert.Equal(t, expectedGraph, vm.moduleGraph.graph)
}
