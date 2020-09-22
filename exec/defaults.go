package exec

import (
	"fmt"
	"strings"

	"github.com/DemoHn/Zn/error"
)

var predefinedValues map[string]ZnValue

// （显示） 方法的执行逻辑
var displayExecutor = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
	// display format string
	var items = []string{}

	for _, param := range params {
		if v, ok := param.(*ZnString); ok {
			items = append(items, v.Value)
		} else {
			items = append(items, param.String())
		}
	}
	fmt.Printf("%s\n", strings.Join(items, " "))
	return NewZnNull(), nil
}

// （递增）方法的执行逻辑
var addValueExecutor = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
	var decimals = []*ZnDecimal{}
	if len(params) == 0 {
		return nil, error.LeastParamsError(1)
	}
	// validate types
	for _, param := range params {
		vparam, ok := param.(*ZnDecimal)
		if !ok {
			return nil, error.InvalidParamType("decimal")
		}
		decimals = append(decimals, vparam)
	}

	if len(decimals) == 1 {
		return decimals[0], nil
	}

	sum := ctx.arith.Add(decimals[0], decimals[1:]...)
	return sum, nil
}

// （递减）方法的执行逻辑
var subValueExecutor = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
	var decimals = []*ZnDecimal{}
	if len(params) == 0 {
		return nil, error.LeastParamsError(1)
	}
	// validate types
	for _, param := range params {
		vparam, ok := param.(*ZnDecimal)
		if !ok {
			return nil, error.InvalidParamType("decimal")
		}
		decimals = append(decimals, vparam)
	}

	if len(decimals) == 1 {
		return decimals[0], nil
	}

	sum := ctx.arith.Sub(decimals[0], decimals[1:]...)
	return sum, nil
}

var mulValueExecutor = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
	var decimals = []*ZnDecimal{}
	if len(params) == 0 {
		return nil, error.LeastParamsError(1)
	}
	// validate types
	for _, param := range params {
		vparam, ok := param.(*ZnDecimal)
		if !ok {
			return nil, error.InvalidParamType("decimal")
		}
		decimals = append(decimals, vparam)
	}

	if len(decimals) == 1 {
		return decimals[0], nil
	}

	sum := ctx.arith.Mul(decimals[0], decimals[1:]...)
	return sum, nil
}

var divValueExecutor = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
	var decimals = []*ZnDecimal{}
	if len(params) == 0 {
		return nil, error.LeastParamsError(1)
	}
	// validate types
	for _, param := range params {
		vparam, ok := param.(*ZnDecimal)
		if !ok {
			return nil, error.InvalidParamType("decimal")
		}
		decimals = append(decimals, vparam)
	}
	if len(decimals) == 1 {
		return decimals[0], nil
	}

	return ctx.arith.Div(decimals[0], decimals[1:]...)
}

var probeExecutor = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
	if len(params) != 2 {
		return nil, error.ExactParamsError(2)
	}

	vtag, ok := params[0].(*ZnString)
	if !ok {
		return nil, error.InvalidParamType("string")
	}
	// add probe data to log
	ctx._probe.AddLog(vtag.Value, params[1])
	return params[1], nil
}

// init function
func init() {
	var funcNameMap = map[string]funcExecutor{
		"显示":      displayExecutor,
		"X+Y":     addValueExecutor,
		"相加":      addValueExecutor,
		"X-Y":     subValueExecutor,
		"相减":      subValueExecutor,
		"X*Y":     mulValueExecutor,
		"相乘":      mulValueExecutor,
		"X/Y":     divValueExecutor,
		"相除":      divValueExecutor,
		"__probe": probeExecutor,
	}

	//// predefined values - those variables (symbols) are defined before
	//// any execution procedure.
	//// NOTICE: those variables are all constants!
	predefinedValues = map[string]ZnValue{
		"真": NewZnBool(true),
		"假": NewZnBool(false),
		"空": NewZnNull(),
	}

	// append executor
	for name, executor := range funcNameMap {
		predefinedValues[name] = NewZnNativeFunction(name, executor)
	}
}
