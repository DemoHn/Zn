package http

import (
	"fmt"
	"io"
	nHTTP "net/http"
	nURL "net/url"
	"strings"
	"time"

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
	DefineProperty("路径", value.NewString("")).
	DefineProperty("方法", value.NewString("GET")).
	DefineProperty("头部", value.NewEmptyHashMap()).
	DefineProperty("查询参数", value.NewEmptyHashMap()).
	DefineProperty("内容", value.NewString("")).
	DefineProperty("允许重定向", value.NewBool(true)).
	DefineProperty("请求时限", value.NewNumber(30)).
	DefineMethod("发送请求", value.NewFunction(func(receiver r.Element, values []r.Element) (r.Element, error) {
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
		headersValue := make(map[string]string)
		// filter element to string
		for k, v := range headers.GetValue() {
			if str, ok := v.(*value.String); ok {
				headersValue[k] = str.GetValue()
			}
		}

		// get query param
		queryParamU := UQueryParams{}
		if queryParam.IsA() { // A:*String, B:*HashMap
			queryParamU.isJSON = false
			queryParamU.queryString = queryParam.GetA().GetValue()
		} else {
			queryParamU.isJSON = true

			// from HashMap -> map[string][]string
			queryParamU.queryDict = make(map[string][]string)
			for k, v := range queryParam.GetB().GetValue() {
				if vstr, ok := v.(*value.String); ok {
					if _, ok := queryParamU.queryDict[k]; !ok {
						queryParamU.queryDict[k] = []string{}
					}
					queryParamU.queryDict[k] = append(queryParamU.queryDict[k], vstr.GetValue())
				}
			}
		}

		var bodyValue string
		if body.IsA() { // A:*String, B:*HashMap
			bodyValue = body.GetA().GetValue()
		} else {
			bodyStr, err := util.HashMapToJSONString(body.GetB())
			if err != nil {
				return nil, err
			}
			bodyValue = bodyStr
		}

		allowRedirectValue := allowRedicrect.GetValue()
		timeoutValue := timeout.GetValue()

		// build request
		req, err := buildHttpRequest(
			pathValue,
			methodValue,
			headersValue,
			queryParamU,
			bodyValue,
		)
		if err != nil {
			return nil, err
		}

		// sendRequest
		resp, data, err := sendHttpRequest(req, allowRedirectValue, int(timeoutValue))
		if err != nil {
			return nil, err
		}

		// build HTTP响应 object
		initProps := map[string]r.Element{
			"状态码": value.NewNumber(float64(resp.StatusCode)),
			"内容":  value.NewString(string(data)),
		}
		return value.NewObject(CLASS_HttpResposne, initProps), nil
	}))

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
	method := values[0].(*value.String).GetValue()
	url := values[1].(*value.String).GetValue()

	req, err := buildHttpRequest(
		url,
		method,
		make(map[string]string),
		UQueryParams{isJSON: false, queryString: ""},
		"",
	)
	if err != nil {
		return nil, value.ThrowException(err.Error())
	}

	resp, body, err := sendHttpRequest(req, true, 30)
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

// /////// http internal logic
// UQueryParams - union type for query params (if isJSON = true, then use queryDict; otherwise use queryString)
type UQueryParams struct {
	isJSON      bool
	queryString string
	queryDict   map[string][]string
}

func buildHttpRequest(
	url string,
	method string,
	headers map[string]string,
	queryParams UQueryParams,
	body string,
) (*nHTTP.Request, error) {
	// #1. 解析url
	urlObj, err := nURL.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("解析URL时出现异常 - %s", err.Error())
	}
	// #2. 匹配 method
	var methodMatch bool = false
	if method == "" {
		method = "GET"
	}

	for _, method := range []string{"get", "post", "put", "delete", "head", "options"} {
		if method == strings.ToLower(method) {
			methodMatch = true
			break
		}
	}
	if !methodMatch {
		return nil, fmt.Errorf("不支持的HTTP请求方法 - %s", method)
	}

	// #3. 匹配 queryParams
	if queryParams.isJSON {
		// for non-empty query params ONLY
		if len(queryParams.queryDict) > 0 {
			// construct query params value set
			qValue := make(nURL.Values)
			// flush old query params
			for k, v := range queryParams.queryDict {
				if len(v) > 0 {
					qValue.Set(k, v[0])
					// append reset value
					for _, vi := range v[1:] {
						qValue.Add(k, vi)
					}
				}
			}
			urlObj.RawQuery = qValue.Encode()
		}
	} else {
		// queryParam is string
		if queryParams.queryString != "" {
			urlObj.RawQuery = nURL.QueryEscape(queryParams.queryString)
		}

	}

	// #4. 设置body
	finalHeaders := nHTTP.Header{}

	// #5. 匹配 headers (may override content-type)
	for k, v := range headers {
		finalHeaders.Set(k, v)
	}

	return &nHTTP.Request{
		Method:     method,
		URL:        urlObj,
		Header:     finalHeaders,
		Body:       io.NopCloser(strings.NewReader(body)),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}, nil
}

// sendHttpRequest - suppose timeout
func sendHttpRequest(req *nHTTP.Request, allowRedicrect bool, timeout int) (*nHTTP.Response, []byte, error) {
	client := &nHTTP.Client{
		CheckRedirect: func(req *nHTTP.Request, via []*nHTTP.Request) error {
			if !allowRedicrect {
				return fmt.Errorf("此请求不允许自动重定向")
			}
			return nil
		},
		Timeout: time.Duration(timeout) * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, []byte{}, value.ThrowException("发送HTTP请求时出现异常 - " + err.Error())
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, []byte{}, value.ThrowException("读取HTTP响应内容时出现异常 - " + err.Error())
	}

	return resp, content, nil
}

func init() {
	var STDLIB_HTTP_NAME = "@HTTP"
	httpLIB = r.NewLibrary(STDLIB_HTTP_NAME)

	httpLIB.RegisterClass("HTTP请求", CLASS_HttpRequest).
		RegisterClass("HTTP响应", CLASS_HttpResposne).
		RegisterFunction("发送HTTP请求", value.NewFunction(FN_sendHTTPRequest))
}
