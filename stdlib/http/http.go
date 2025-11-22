package http

import (
	libURL "net/url"

	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/util"
	"github.com/DemoHn/Zn/pkg/value"
)

var httpLIB *r.Library

/*
*
定义HTTP请求：

	其路径 = “”
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
	DefineProperty("路径", value.NewString("")).
	DefineProperty("方法", value.NewString("GET")).
	DefineProperty("头部", value.NewEmptyHashMap()).
	DefineProperty("查询参数", value.NewEmptyHashMap()).
	DefineProperty("内容", value.NewString("")).
	DefineProperty("允许重定向", value.NewBool(true)).
	DefineProperty("请求时限", value.NewNumber(30)).
	DefineMethod("发送请求", value.NewFunction(methodSendRequest))

func methodSendRequest(receiver r.Element, values []r.Element) (r.Element, error) {
	receiver.GetProperty("路径")
	path, err1 := value.AssertPropertyElement[*value.String](receiver, "路径")
	method, err2 := value.AssertPropertyElement[*value.String](receiver, "方法")
	headers, err3 := value.AssertPropertyElement[*value.HashMap](receiver, "头部")
	// assert String or HashMap
	queryParam, err4 := value.BuildEitherPropertyElement[*value.String, *value.HashMap](receiver, "查询参数")
	// assert String or HashMap
	body, err5 := value.BuildEitherPropertyElement[*value.String, *value.HashMap](receiver, "内容")
	allowRedicrect, err6 := value.AssertPropertyElement[*value.Bool](receiver, "允许重定向")
	timeout, err7 := value.AssertPropertyElement[*value.Number](receiver, "请求时限")

	// assert errors
	for _, e := range []error{err1, err2, err3, err4, err5, err6, err7} {
		if e != nil {
			return nil, e
		}
	}

	// get params
	pathValue := path.GetValue()
	methodValue := method.GetValue()

	// get headers
	headersValue := make([][2]string, 0)
	for k, v := range headers.GetValue() {
		// only string Value is allowed to append into header set
		// others will be filted without any warning wor~
		if str, ok := v.(*value.String); ok {
			headersValue = append(headersValue, [2]string{k, str.GetValue()})
		}
	}

	// get query param
	queryParamValue := make([][2]string, 0)
	if queryParam.IsA() { // A:*String, B:*HashMap
		qs := queryParam.GetA().GetValue()
		qsValue, err := libURL.ParseQuery(qs)
		if err != nil {
			return nil, err
		}
		for k, v := range qsValue {
			for _, item := range v {
				queryParamValue = append(queryParamValue, [2]string{k, item})
			}
		}
	} else {
		// from HashMap -> map[string][]string
		for k, v := range queryParam.GetB().GetValue() {
			switch vv := v.(type) {
			case *value.String:
				queryParamValue = append(queryParamValue, [2]string{k, vv.GetValue()})
			// handle Array of String
			case *value.Array:
				for _, inItem := range vv.GetValue() {
					if inItemStr, ok := inItem.(*value.String); ok {
						queryParamValue = append(queryParamValue, [2]string{k, inItemStr.GetValue()})
					}
				}
			}
		}
	}

	// convert reqBody
	var reqBody util.ReqBody
	if body.IsA() { // A:*String, B:*HashMap
		reqBody.ContentType = "application/x-www-form-urlencoded"
		reqBody.Value = body.GetA().GetValue()
	} else {
		bodyStr, err := util.HashMapToJSONString(body.GetB())
		if err != nil {
			return nil, err
		}
		reqBody.ContentType = "application/json"
		reqBody.Value = bodyStr.GetValue()
	}

	allowRedirectValue := allowRedicrect.GetValue()
	timeoutValue := timeout.GetValue()

	// build request
	req, err := util.BuildBaseHttpRequest(
		pathValue,
		methodValue,
		headersValue,
		queryParamValue,
		reqBody,
	)
	if err != nil {
		return nil, err
	}

	// sendRequest
	resp, data, err := util.SendHttpRequest(req, allowRedirectValue, int(timeoutValue))
	if err != nil {
		return nil, err
	}

	// build HTTP响应 object
	initProps := map[string]r.Element{
		"状态码": value.NewNumber(float64(resp.StatusCode)),
		"内容":  value.NewString(string(data)),
	}
	return value.NewObject(CLASS_HttpResposne, initProps), nil
}

// HTTP响应类型
/**
定义HTTP响应:
	其状态码 = 200
	其头部 = [
		“Content-Type” = “application/json”
	]
	其内容 = “”
*/
var CLASS_HttpResposne = value.NewClassModel("HTTP响应").
	DefineProperty("状态码", value.NewNumber(200)).
	DefineProperty("头部", value.NewEmptyHashMap()).
	DefineProperty("内容", value.NewString(""))

// 发送HTTP请求方法
func FN_sendHTTPRequest(receiver r.Element, values []r.Element) (r.Element, error) {
	if err := value.ValidateExactParams(values, "string", "string"); err != nil {
		return nil, err
	}

	// #1. get exact type of params
	const defaultTimeout = 30

	method := values[0].(*value.String).GetValue()
	url := values[1].(*value.String).GetValue()
	emptyHeader := [][2]string{}
	emptyQuery := [][2]string{}
	emptyBody := util.ReqBody{}

	req, err := util.BuildBaseHttpRequest(
		url,
		method,
		emptyHeader,
		emptyQuery,
		emptyBody,
	)
	if err != nil {
		return nil, value.ThrowException(err.Error())
	}

	resp, body, err := util.SendHttpRequest(req, true, defaultTimeout)
	if err != nil {
		return nil, value.ThrowException(err.Error())
	}

	// 构造 HTTP响应 对象
	initProps := map[string]r.Element{
		"状态码": value.NewNumber(float64(resp.StatusCode)),
		"内容":  value.NewString(string(body)),
	}
	return value.NewObject(CLASS_HttpResposne, initProps), nil
}

func Export() *r.Library {
	return httpLIB
}

func init() {
	var STDLIB_HTTP_NAME = "@HTTP"
	httpLIB = r.NewLibrary(STDLIB_HTTP_NAME)

	httpLIB.RegisterClass("HTTP请求", CLASS_HttpRequest).
		RegisterClass("HTTP响应", CLASS_HttpResposne).
		RegisterFunction("发送HTTP请求", value.NewFunction(FN_sendHTTPRequest))
}
