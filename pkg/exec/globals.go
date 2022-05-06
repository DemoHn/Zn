package exec

import (
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

type funcExecutor = func(*r.Context, []r.Value) (r.Value, error)

// GlobalValues -
var GlobalValues map[string]r.Value

// init function
func init() {
	var funcNameMap = map[string]funcExecutor{
		"显示": value.DisplayExecutor,
	}

	//// predefined values - those variables (symbols) are defined before
	//// any execution procedure.
	//// NOTICE: those variables are all constants!
	GlobalValues = map[string]r.Value{
		"真": value.NewBool(true),
		"假": value.NewBool(false),
		"空": value.NewNull(),
	}

	// append executor
	for name, executor := range funcNameMap {
		GlobalValues[name] = value.NewFunction(name, executor)
	}
}
