package stdlib

import (
	"io/ioutil"
	"net/http"

	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

var httpModuleName = "HTTP"
var httpModule = r.NewInternalModule(httpModuleName)

// HTTP响应类型
var httpResponseClass = value.NewClassModel("HTTP响应", httpModule).
	DefineProperty("代码", value.NewNumber(200)).
	DefineProperty("内容", value.NewString(""))

// 发送HTTP请求方法
func sendHTTPRequestFunc(values []r.Element) (r.Element, error) {
	if err := value.ValidateExactParams(values, "string", "string"); err != nil {
		return nil, err
	}
	method := values[0].(*value.String).GetValue()
	url := values[1].(*value.String).GetValue()

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, value.ThrowException("创建HTTP请求失败：" + err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, value.ThrowException("发送HTTP请求失败：" + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, value.ThrowException("读取HTTP响应内容失败：" + err.Error())
	}

	// 构造 HTTP响应 对象
	initProps := map[string]r.Element{
		"代码": value.NewNumber(float64(resp.StatusCode)),
		"内容": value.NewString(string(body)),
	}
	return value.NewObject(httpResponseClass, initProps), nil
}

func init() {
	// 注册 HTTP响应 类型
	RegisterClassForModule(httpModule, "HTTP响应", httpResponseClass)
	// 注册 发送HTTP请求 方法
	RegisterFunctionForModule(httpModule, "发送HTTP请求", sendHTTPRequestFunc)
	// 注册模块
	RegisterModule(httpModuleName, httpModule)
}
