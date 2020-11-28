package exec

import (
	"fmt"
	"strings"

	"github.com/DemoHn/Zn/error"
)

var globalValues map[string]Value

// （显示） 方法的执行逻辑
var displayExecutor = func(ctx *Context, params []Value) (Value, *error.Error) {
	// display format string
	var items = []string{}

	for _, param := range params {
		// if param is a string, display its value (without 「 」 quotes) directly
		if str, ok := param.(*String); ok {
			items = append(items, str.value)
		} else {
			items = append(items, StringifyValue(param))
		}
	}
	fmt.Printf("%s\n", strings.Join(items, " "))
	return NewNull(), nil
}

// （递增）方法的执行逻辑
var addValueExecutor = func(ctx *Context, params []Value) (Value, *error.Error) {
	var decimals = []Decimal{}
	if len(params) == 0 {
		return nil, error.LeastParamsError(1)
	}
	// validate types
	for _, param := range params {
		vparam, ok := param.(*Decimal)
		if !ok {
			return nil, error.InvalidParamType("decimal")
		}
		decimals = append(decimals, *vparam)
	}

	if len(decimals) == 1 {
		return &decimals[0], nil
	}

	sum := ctx.arith.Add(decimals[0], decimals[1:]...)
	return &sum, nil
}

// （递减）方法的执行逻辑
var subValueExecutor = func(ctx *Context, params []Value) (Value, *error.Error) {
	var decimals = []Decimal{}
	if len(params) == 0 {
		return nil, error.LeastParamsError(1)
	}
	// validate types
	for _, param := range params {
		vparam, ok := param.(*Decimal)
		if !ok {
			return nil, error.InvalidParamType("decimal")
		}
		decimals = append(decimals, *vparam)
	}

	if len(decimals) == 1 {
		return &decimals[0], nil
	}

	sum := ctx.arith.Sub(decimals[0], decimals[1:]...)
	return &sum, nil
}

var mulValueExecutor = func(ctx *Context, params []Value) (Value, *error.Error) {
	var decimals = []Decimal{}
	if len(params) == 0 {
		return nil, error.LeastParamsError(1)
	}
	// validate types
	for _, param := range params {
		vparam, ok := param.(*Decimal)
		if !ok {
			return nil, error.InvalidParamType("decimal")
		}
		decimals = append(decimals, *vparam)
	}

	if len(decimals) == 1 {
		return &decimals[0], nil
	}

	sum := ctx.arith.Mul(decimals[0], decimals[1:]...)
	return &sum, nil
}

var divValueExecutor = func(ctx *Context, params []Value) (Value, *error.Error) {
	var decimals = []Decimal{}
	if len(params) == 0 {
		return nil, error.LeastParamsError(1)
	}
	// validate types
	for _, param := range params {
		vparam, ok := param.(*Decimal)
		if !ok {
			return nil, error.InvalidParamType("decimal")
		}
		decimals = append(decimals, *vparam)
	}
	if len(decimals) == 1 {
		return &decimals[0], nil
	}

	res, err := ctx.arith.Div(decimals[0], decimals[1:]...)
	return &res, err
}

var probeExecutor = func(ctx *Context, params []Value) (Value, *error.Error) {
	if len(params) != 2 {
		return nil, error.ExactParamsError(2)
	}

	vtag, ok := params[0].(*String)
	if !ok {
		return nil, error.InvalidParamType("string")
	}
	// add probe data to log
	ctx._probe.AddLog(vtag.value, params[1])
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
	globalValues = map[string]Value{
		"真": NewBool(true),
		"假": NewBool(false),
		"空": NewNull(),
	}

	// append executor
	for name, executor := range funcNameMap {
		globalValues[name] = NewFunction(name, executor)
	}
}
