package exec

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

type funcExecutor = func(*r.Context, []r.Element) (r.Element, error)

// global consts
var (
	ZnConstBoolTrue         = value.NewBool(true)
	ZnConstBoolFalse        = value.NewBool(false)
	ZnConstNull             = value.NewNull()
	ZnConstExceptionClass   = newExceptionModel()
	ZnConstHTTPRequestClass = newHTTPRequestModel()
	ZnConstDisplayFunc      = newDisplayFunc()
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
		"真":  ZnConstBoolTrue,
		"假":  ZnConstBoolFalse,
		"空":  ZnConstNull,
		"异常": ZnConstExceptionClass,
		"显示": ZnConstDisplayFunc,
	}

	GlobalValues = globalValues
}

func newExceptionModel() *value.ClassModel {
	constructorFunc := value.NewFunction(nil, func(c *r.Context, values []r.Element) (r.Element, error) {
		if err := value.ValidateExactParams(values, "string"); err != nil {
			return nil, err
		}

		message := values[0].(*value.String)
		return value.NewException(message.String()), nil
	})

	return value.NewClassModel("异常", nil).
		SetConstructorFunc(constructorFunc).
		DefineProperty("内容", value.NewString(""))
}

func newHTTPRequestModel() *value.ClassModel {
	readRequestBodyFunc := func(c *r.Context, values []r.Element) (r.Element, error) {
		if thisValue, ok := c.GetThisValue().(*value.Object); ok {
			if goValue, err := thisValue.GetProperty(c, "-goHttpRequest-"); err == nil {
				if t, ok := goValue.(*value.GoValue); ok {
					if t.GetTag() == "*http.Request" {
						reqV := t.GetValue().(*http.Request)
						data, err := ioutil.ReadAll(reqV.Body)
						if err != nil {
							return nil, zerr.NewRuntimeException("读取请求内容出现异常！")
						}

						return value.NewString(string(data)), nil
					}
				}
			}
		}
		return nil, zerr.InvalidParamType("goType: *http.Request")
	}

	return value.NewClassModel("HTTP请求", nil).
		DefineProperty("路径", value.NewString("")).
		DefineProperty("方法", value.NewString("")).
		DefineProperty("头部", value.NewHashMap([]value.KVPair{})).
		DefineProperty("查询参数", value.NewHashMap([]value.KVPair{})).
		DefineProperty("-goHttpRequest-", value.NewGoValue("*http.Request", nil)).
		DefineMethod("读取内容", value.NewFunction(nil, readRequestBodyFunc))
}

func newDisplayFunc() *value.Function {
	displayExecutor := func(c *r.Context, params []r.Element) (r.Element, error) {
		// display format string
		var items = []string{}
		for _, param := range params {
			// if param is a string, display its value (without 「 」 quotes) directly
			if str, ok := param.(*value.String); ok {
				items = append(items, str.String())
			} else {
				items = append(items, value.StringifyValue(param))
			}
		}

		os.Stdout.Write([]byte(fmt.Sprintf("%s\n", strings.Join(items, " "))))
		return value.NewNull(), nil
	}

	return value.NewFunction(nil, displayExecutor)
}
