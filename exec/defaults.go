package exec

import (
	"fmt"
	"math/big"
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

	sum, _ := NewZnDecimal("0")
	for _, decimal := range decimals {
		r1, r2 := rescalePair(sum, decimal)
		newco := new(big.Int).Add(r1.co, r2.co)

		sum.co = newco
		sum.exp = r1.exp
	}

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

	sum, _ := NewZnDecimal("0")
	for _, decimal := range decimals {
		r1, r2 := rescalePair(sum, decimal)
		negco := new(big.Int).Neg(r2.co)
		newco := new(big.Int).Add(r1.co, negco)

		sum.co = newco
		sum.exp = r1.exp
	}

	return sum, nil
}

// init function
func init() {
	//// predefined values - those variables (symbols) are defined before
	//// any execution procedure.
	//// NOTICE: those variables are all constants!
	predefinedValues = map[string]ZnValue{
		"真":  NewZnBool(true),
		"假":  NewZnBool(false),
		"空":  NewZnNull(),
		"显示": NewZnNativeFunction("显示", displayExecutor),
		"递增": NewZnNativeFunction("递增", addValueExecutor),
		"递加": NewZnNativeFunction("递增", addValueExecutor),
		"递减": NewZnNativeFunction("递减", subValueExecutor),
	}
}
