package exec

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/exec/val"
)

type funcExecutor = func(*ctx.Context, []ctx.Value) (ctx.Value, *error.Error)

// GlobalValues -
var GlobalValues map[string]ctx.Value

// init function
func init() {
	var funcNameMap = map[string]funcExecutor{
		"显示":      val.DisplayExecutor,
		"X+Y":     val.AddValueExecutor,
		"相加":      val.AddValueExecutor,
		"X-Y":     val.SubValueExecutor,
		"相减":      val.SubValueExecutor,
		"X*Y":     val.MulValueExecutor,
		"相乘":      val.MulValueExecutor,
		"X/Y":     val.DivValueExecutor,
		"相除":      val.DivValueExecutor,
		"__probe": val.ProbeExecutor,
	}

	//// predefined values - those variables (symbols) are defined before
	//// any execution procedure.
	//// NOTICE: those variables are all constants!
	GlobalValues = map[string]ctx.Value{
		"真": val.NewBool(true),
		"假": val.NewBool(false),
		"空": val.NewNull(),
	}

	// append executor
	for name, executor := range funcNameMap {
		GlobalValues[name] = val.NewFunction(name, executor)
	}
}
