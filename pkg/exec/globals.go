package exec

import (
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

type funcExecutor = func(*r.Context, []r.Element) (r.Element, error)

// globalValues -
var globalValues map[string]r.Element

// init function
func init() {
	var funcNameMap = map[string]funcExecutor{
		"显示": value.DisplayExecutor,
	}

	// construct 异常 class
	expClassRef := value.NewClassModel("异常")
	expClassRef.Constructor = func(c *r.Context, values []r.Element) (r.Element, error) {
		if err := value.ValidateExactParams(values, "string"); err != nil {
			return nil, err
		}

		message := values[0].(*value.String)
		return value.NewException(message.String()), nil
	}

	//// predefined values - those variables (symbols) are defined before
	//// any execution procedure.
	//// NOTICE: those variables are all constants!
	globalValues = map[string]r.Element{
		"真":  value.NewBool(true),
		"假":  value.NewBool(false),
		"空":  value.NewNull(),
		"异常": expClassRef,
	}

	// append executor
	for name, executor := range funcNameMap {
		globalValues[name] = value.NewFunction(nil, executor)
	}
}
