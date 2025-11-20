package exec

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

// global consts
var (
	ZnConstBoolTrue       = value.NewBool(true)
	ZnConstBoolFalse      = value.NewBool(false)
	ZnConstNull           = value.NewNull()
	ZnConstExceptionClass = newExceptionModel()
	ZnConstDisplayFunc    = newDisplayFunc()
	ZnConstGetRandomFloat = newGetRandomFloatFunc()
)

// globalValues -
var globalValues map[string]r.Element

var GlobalValues map[string]r.Element

// init function
func init() {

	//// predefined values - those variables (symbols) are defined before
	//// any execution procedure.
	//// NOTICE: those variables are all constants!
	globalValues = map[string]r.Element{
		"真":    ZnConstBoolTrue,
		"假":    ZnConstBoolFalse,
		"空":    ZnConstNull,
		"异常":   ZnConstExceptionClass,
		"显示":   ZnConstDisplayFunc,
		"取随机数": ZnConstGetRandomFloat,
	}

	GlobalValues = globalValues
}

func newExceptionModel() *value.ClassModel {
	constructorFunc := func(receiver r.Element, values []r.Element) (r.Element, error) {
		if err := value.ValidateExactParams(values, "string"); err != nil {
			return nil, err
		}

		message := values[0].(*value.String)
		return value.NewException(message.String()), nil
	}

	return value.NewClassModel("异常").SetConstructor(constructorFunc)
}

func newDisplayFunc() *value.Function {
	displayExecutor := func(receiver r.Element, params []r.Element) (r.Element, error) {
		// display format string
		var items = []string{}
		for _, param := range params {
			items = append(items, value.StringifyValue(param))
		}

		os.Stdout.Write([]byte(fmt.Sprintf("%s\n", strings.Join(items, " "))))
		return value.NewNull(), nil
	}

	return value.NewFunction(displayExecutor)
}

func newGetRandomFloatFunc() *value.Function {
	getRandomFloatExecutor := func(receiver r.Element, params []r.Element) (r.Element, error) {
		return value.NewNumber(rand.Float64()), nil
	}

	return value.NewFunction(getRandomFloatExecutor)
}
