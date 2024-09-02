package exec

import (
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
	"github.com/DemoHn/Zn/pkg/value/cmodels"
)

type funcExecutor = func(*r.Context, []r.Element) (r.Element, error)

// globalValues -
var globalValues map[string]r.Element

var GlobalValues map[string]r.Element

// init function
func init() {
	var funcNameMap = map[string]funcExecutor{
		"显示": value.DisplayExecutor,
	}

	//// predefined values - those variables (symbols) are defined before
	//// any execution procedure.
	//// NOTICE: those variables are all constants!
	globalValues = map[string]r.Element{
		"真":  value.NewBool(true),
		"假":  value.NewBool(false),
		"空":  value.NewNull(),
		"异常": cmodels.NewExceptionModel(),
	}

	// append executor
	for name, executor := range funcNameMap {
		globalValues[name] = value.NewFunction(nil, executor)
	}

	GlobalValues = globalValues
}
