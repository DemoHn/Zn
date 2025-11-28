package http

import (
	libURL "net/url"

	"github.com/DemoHn/Zn/pkg/common"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

var httpLIB *r.Library

// HTTPClient, HTTPRequest, HTTPResponse

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
	var reqBody ReqBody
	if body.IsA() { // A:*String, B:*HashMap
		reqBody.ContentType = "application/x-www-form-urlencoded"
		reqBody.Value = body.GetA().GetValue()
	} else {
		bodyStr, err := common.HashMapToJSONString(body.GetB())
		if err != nil {
			return nil, err
		}
		reqBody.ContentType = "application/json"
		reqBody.Value = bodyStr.GetValue()
	}

	allowRedirectValue := allowRedicrect.GetValue()
	timeoutValue := timeout.GetValue()

	// build request
	req, err := buildBaseHttpRequest(
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
	resp, data, err := sendHttpRequest(req, allowRedirectValue, int(timeoutValue))
	if err != nil {
		return nil, err
	}

	return buildOBJ_HttpResponse(resp, data), nil
}

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
	emptyBody := ReqBody{}

	req, err := buildBaseHttpRequest(
		url,
		method,
		emptyHeader,
		emptyQuery,
		emptyBody,
	)
	if err != nil {
		return nil, value.ThrowException(err.Error())
	}

	resp, data, err := sendHttpRequest(req, true, defaultTimeout)
	if err != nil {
		return nil, value.ThrowException(err.Error())
	}

	return buildOBJ_HttpResponse(resp, data), nil
}

func Export() *r.Library {
	return httpLIB
}

func init() {
	var STDLIB_HTTP_NAME = "@HTTP"
	httpLIB = r.NewLibrary(STDLIB_HTTP_NAME)

	httpLIB.RegisterClass("HTTP请求", common.CLASS_HttpRequest).
		RegisterClass("HTTP响应", common.CLASS_HttpResponse).
		RegisterFunction("发送HTTP请求", value.NewFunction(FN_sendHTTPRequest))
}
