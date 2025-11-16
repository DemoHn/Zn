package stdlib

import (
	"fmt"
	"io"
	"net/http"

	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

/*
*
定义HTTP请求：

	其URL = “”
	其方法 = “”

	注：头部 (headers) 是请求头部的信息；在进行 HTTP 请求时，会自动添加一些默认的头部信息，如
		如 Content-Type、Host、Content-Length等；HTTP请求#头部即是在此基础上添加额外的头部信息.
	其头部 = [=]

	注：查询参数 (queryString) 是URL中的 ? 后面的部分；
		如 URL: https://abc.com?a=1&b=2&c=3，对应的查询参数是 ‘a=1&b=2&c=3’

		在定义HTTP请求的‘查询参数’属性时，可以是文本 (String)，也可以是字典 (HashMap) 类型，
		- 如果是文本类型，那就是原始的查询参数字符串
		- 如果是字典类型，那就是将字典转换成的查询参数字符串
	其查询参数 = [=]

	注：内容 (body) 是请求体的内容，主要用在 POST 等需要传输大量数据的地方

		内容可以是文本 (String)，也可以是字典 (HashMap)
		- 若是文本，本次请求的 Content-Type 为 application/x-www-form-urlencoded，其内容即是诸如 a=1&b=2&c=3 这样的文本
		- 若是字典，本次请求的 Content-Type 为 application/json，实际传输内容会转换成JSON字符串，诸如 {"a": 1, "b": 2, "c": 3}
	其内容 = “”

	其允许重定向 = 是

	注：请求时限 (timeout) 指的是在处理一次请求时，所能够等待的最长时间 - 如果超过这个时间还没有反应，那么此请求将会自动中断并抛出异常
		- 默认是 30 秒
	其请求时限 = 30
*/
var CLASS_HttpRequest = value.NewClassModel("HTTP请求").
	DefineProperty("URL", value.NewString("")).
	DefineProperty("方法", value.NewString("")).
	DefineProperty("头部", value.NewEmptyHashMap()).
	DefineProperty("查询参数", value.NewEmptyHashMap()).
	DefineProperty("内容", value.NewString("")).
	DefineProperty("允许重定向", value.NewBool(true)).
	DefineProperty("请求时限", value.NewNumber(30)).
	DefineMethod("TESTING", value.NewFunction(func(receiver r.Element, values []r.Element) (r.Element, error) {
		fmt.Println("======")
		fmt.Println(receiver)
		fmt.Println(values)

		return nil, nil
	}))

// HTTP响应类型
var CLASS_HttpResposne = value.NewClassModel("HTTP响应").
	DefineProperty("代码", value.NewNumber(200)).
	DefineProperty("内容", value.NewString(""))

// 发送HTTP请求方法
func sendHTTPRequestFunc(receiver r.Element, values []r.Element) (r.Element, error) {
	if err := value.ValidateExactParams(values, "string", "string"); err != nil {
		return nil, err
	}

	// #1. get exact type of params
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, value.ThrowException("读取HTTP响应内容失败：" + err.Error())
	}

	// 构造 HTTP响应 对象
	initProps := map[string]r.Element{
		"代码": value.NewNumber(float64(resp.StatusCode)),
		"内容": value.NewString(string(body)),
	}
	return value.NewObject(CLASS_HttpResposne, initProps), nil
}

func init() {
	var STDLIB_HTTP_NAME = "@HTTP"
	var httpLIB = NewLibrary(STDLIB_HTTP_NAME)

	RegisterClassForLibrary(httpLIB, "HTTP请求", CLASS_HttpRequest)
	// 注册 HTTP响应 类型
	RegisterClassForLibrary(httpLIB, "HTTP响应", CLASS_HttpResposne)
	// 注册 发送HTTP请求 方法
	RegisterFunctionForLibrary(httpLIB, "发送HTTP请求", sendHTTPRequestFunc)
	// 注册模块
	RegisterLibrary(httpLIB)
}
