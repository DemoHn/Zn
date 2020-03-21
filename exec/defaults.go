package exec

import (
	"fmt"
	"strings"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/syntax"
)

var predefinedValues map[string]ZnValue

// （显示） 方法的执行逻辑
var displayExecutor = func(params []ZnValue, template *syntax.FunctionDeclareStmt, ctx *Context) (ZnValue, *error.Error) {
	// display format string
	var items = []string{}

	for _, param := range params {
		items = append(items, param.String())
	}
	fmt.Printf("%s\n", strings.Join(items, " "))
	return NewZnNull(), nil
}

// （递增）方法的执行逻辑
var addValueExecutor = func(params []ZnValue, template *syntax.FunctionDeclareStmt, ctx *Context) (ZnValue, *error.Error) {
	var decimals = []*ZnDecimal{}
	// validate types
	for _, param := range params {
		vparam, ok := param.(*ZnDecimal)
		if !ok {
			return nil, error.NewErrorSLOT("入参皆须为「数值」类型")
		}
		decimals = append(decimals, vparam)
	}

	sum := ctx.ArithInstance.Add(decimals...)
	return sum, nil
}

// （递减）方法的执行逻辑
var subValueExecutor = func(params []ZnValue, template *syntax.FunctionDeclareStmt, ctx *Context) (ZnValue, *error.Error) {
	var decimals = []*ZnDecimal{}
	// validate types
	for _, param := range params {
		vparam, ok := param.(*ZnDecimal)
		if !ok {
			return nil, error.NewErrorSLOT("入参皆须为「数值」类型")
		}
		decimals = append(decimals, vparam)
	}

	sum := ctx.ArithInstance.Sub(decimals...)
	return sum, nil
}

// init function
func init() {
	//// predefined values - those variables (symbols) are defined before
	//// any execution procedure.
	//// NOTICE: those variables are all constants!
	predefinedValues = map[string]ZnValue{
		"真":   NewZnBool(true),
		"假":   NewZnBool(false),
		"空":   NewZnNull(),
		"显示":  NewZnNativeFunction("显示", displayExecutor),
		"X+Y": NewZnNativeFunction("X+Y", addValueExecutor),
		"求和":  NewZnNativeFunction("X+Y", addValueExecutor),
		"X-Y": NewZnNativeFunction("X-Y", subValueExecutor),
		"求差":  NewZnNativeFunction("X-Y", subValueExecutor),
	}
}
