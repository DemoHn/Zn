package runtime

import (
	"testing"

	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/stretchr/testify/assert"
)

// assert new script callframe
func TestNewScriptCallFrame(t *testing.T) {
	module := &Module{
		id:       1,
		fullName: "test",
		program:  &syntax.Program{},
	}
	cf := NewScriptCallFrame(module)

	// assert cf
	assert.Equal(t, cf.module, module)
	assert.Equal(t, cf.callType, CALL_TYPE_SCRIPT)
	assert.Equal(t, cf.currentLine, 0)
	assert.Equal(t, cf.programAST, module.program)
	assert.Equal(t, cf.thisValue, nil)
	assert.Equal(t, cf.returnValue, nil)
	assert.Equal(t, cf.IsExceptionCallFrame(), false)
	assert.Equal(t, cf.IsFunctionCallFrame(), false)
	assert.Equal(t, cf.IsScriptCallFrame(), true)
}

// TestNewFunctionCallFrame
func TestNewFunctionCallFrame(t *testing.T) {
	module := &Module{
		id:       1,
		fullName: "test",
		program:  &syntax.Program{},
	}
	cf := NewFunctionCallFrame(module, nil)
	// assert cf
	assert.Equal(t, cf.module, module)
	assert.Equal(t, cf.callType, CALL_TYPE_FUNCTION)
	assert.Equal(t, cf.currentLine, 0)
	assert.Equal(t, cf.programAST, module.program)
	assert.Equal(t, cf.thisValue, nil)
	assert.Equal(t, cf.returnValue, nil)
	assert.Equal(t, cf.IsExceptionCallFrame(), false)
	assert.Equal(t, cf.IsFunctionCallFrame(), true)
	assert.Equal(t, cf.IsScriptCallFrame(), false)
}

// TestNewExceptionCallFrame
func TestNewExceptionCallFrame(t *testing.T) {
	module := &Module{
		id:       1,
		fullName: "test",
		program:  &syntax.Program{},
	}
	cf := NewExceptionCallFrame(module, nil)

	// assert cf
	assert.Equal(t, cf.module, module)
	assert.Equal(t, cf.callType, CALL_TYPE_EXCEPTION_BLOCK)
	assert.Equal(t, cf.currentLine, 0)
	assert.Equal(t, cf.programAST, module.program)
	assert.Equal(t, cf.thisValue, nil)
	assert.Equal(t, cf.returnValue, nil)
	assert.Equal(t, cf.IsExceptionCallFrame(), true)
	assert.Equal(t, cf.IsFunctionCallFrame(), false)
	assert.Equal(t, cf.IsScriptCallFrame(), false)
}
