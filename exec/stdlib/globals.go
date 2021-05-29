package stdlib

import (
	"fmt"
	"strings"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/exec/val"
)

type funcExecutor = func(*ctx.Context, []ctx.Value) (ctx.Value, *error.Error)

var globalValues map[string]ctx.Value

// （显示） 方法的执行逻辑
var DisplayExecutor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	// display format string
	var items = []string{}

	for _, param := range params {
		// if param is a string, display its value (without 「 」 quotes) directly
		if str, ok := param.(*val.String); ok {
			items = append(items, str.String())
		} else {
			items = append(items, StringifyValue(param))
		}
	}
	fmt.Printf("%s\n", strings.Join(items, " "))
	return val.NewNull(), nil
}

// （递增）方法的执行逻辑
var AddValueExecutor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	var decimals = []*val.Decimal{}

	if err := validateLeastParams(params, "decimal+"); err != nil {
		return nil, err
	}
	// validate types
	for _, param := range params {
		vparam := param.(*val.Decimal)
		decimals = append(decimals, vparam)
	}

	if len(decimals) == 1 {
		return decimals[0], nil
	}

	sum := decimals[0].Add(decimals[1:]...)
	return sum, nil
}

// （递减）方法的执行逻辑
var SubValueExecutor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	var decimals = []*val.Decimal{}

	if err := validateLeastParams(params, "decimal+"); err != nil {
		return nil, err
	}
	// validate types
	for _, param := range params {
		vparam := param.(*val.Decimal)
		decimals = append(decimals, vparam)
	}

	if len(decimals) == 1 {
		return decimals[0], nil
	}

	sum := decimals[0].Sub(decimals[1:]...)
	return sum, nil
}

var MulValueExecutor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	var decimals = []*val.Decimal{}

	if err := validateLeastParams(params, "decimal+"); err != nil {
		return nil, err
	}
	// validate types
	for _, param := range params {
		vparam := param.(*val.Decimal)
		decimals = append(decimals, vparam)
	}

	if len(decimals) == 1 {
		return decimals[0], nil
	}

	sum := decimals[0].Mul(decimals[1:]...)
	return sum, nil
}

var DivValueExecutor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	var decimals = []*val.Decimal{}

	if err := validateLeastParams(params, "decimal+"); err != nil {
		return nil, err
	}
	// validate types
	for _, param := range params {
		vparam := param.(*val.Decimal)
		decimals = append(decimals, vparam)
	}
	if len(decimals) == 1 {
		return decimals[0], nil
	}

	res, err := decimals[0].Div(decimals[1:]...)
	return res, err
}

var ProbeExecutor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	if len(params) != 2 {
		return nil, error.ExactParamsError(2)
	}

	vtag, ok := params[0].(*val.String)
	if !ok {
		return nil, error.InvalidParamType("string")
	}
	// add probe data to log
	c._probe.AddLog(vtag.value, params[1])
	return params[1], nil
}

// init function
func init() {
	var funcNameMap = map[string]funcExecutor{
		"显示":      DisplayExecutor,
		"X+Y":     AddValueExecutor,
		"相加":      AddValueExecutor,
		"X-Y":     SubValueExecutor,
		"相减":      SubValueExecutor,
		"X*Y":     MulValueExecutor,
		"相乘":      MulValueExecutor,
		"X/Y":     DivValueExecutor,
		"相除":      DivValueExecutor,
		"__probe": ProbeExecutor,
	}

	//// predefined values - those variables (symbols) are defined before
	//// any execution procedure.
	//// NOTICE: those variables are all constants!
	globalValues = map[string]ctx.Value{
		"真": val.NewBool(true),
		"假": val.NewBool(false),
		"空": val.NewNull(),
	}

	// append executor
	for name, executor := range funcNameMap {
		globalValues[name] = val.NewFunction(name, executor)
	}
}
