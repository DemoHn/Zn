package http

import (
	libHTTP "net/http"
	libURL "net/url"

	"github.com/DemoHn/Zn/pkg/common"
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

var httpLIB *r.Library

// build *http.Request from HTTP请求 Element
func digestBaseRequest(reqElement r.Element) (*libHTTP.Request, error) {
	reqObj, err0 := value.AssertElement[*value.Object](reqElement)
	if err0 != nil {
		return nil, err0
	}
	if !reqObj.IsInstanceOf(common.CLASS_HttpRequest) {
		return nil, value.ThrowException("不匹配的参数类型：HTTP请求")
	}

	path, err1 := value.AssertPropertyElement[*value.String](reqObj, "路径")
	method, err2 := value.AssertPropertyElement[*value.String](reqObj, "方法")
	headers, err3 := value.AssertPropertyElement[*value.HashMap](reqObj, "头部")
	// assert String or HashMap
	queryParam, err4 := value.BuildEitherPropertyElement[*value.String, *value.HashMap](reqObj, "查询参数")
	// assert String or HashMap
	body, err5 := value.BuildEitherPropertyElement[*value.String, *value.HashMap](reqObj, "内容")

	// assert errors
	for _, e := range []error{err1, err2, err3, err4, err5} {
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

	// build request
	return buildBaseHttpRequest(
		pathValue,
		methodValue,
		headersValue,
		queryParamValue,
		reqBody,
	)
}

// (发送HTTP请求：请求@HTTP请求)
func FN_sendHTTPRequest(receiver r.Element, values []r.Element) (r.Element, error) {
	const TOTAL_PARAMS = 1
	if len(values) != TOTAL_PARAMS {
		return nil, zerr.ExactParamsError(TOTAL_PARAMS)
	}
	req, err1 := digestBaseRequest(values[0])
	if err1 != nil {
		return nil, err1
	}

	const defaultTimeout = 30
	// #1. go! sendRequest
	resp, data, err := sendHttpRequest(req, true, defaultTimeout)
	if err != nil {
		return nil, value.ThrowException(err.Error())
	}

	return buildOBJ_HttpResponse(resp, data), nil
}

// (发送GET请求：URL)
func FN_sendHTTPRequest_GET(receiver r.Element, values []r.Element) (r.Element, error) {
	if err := value.ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}

	// #1. get exact type of params
	const defaultTimeout = 30
	const method = "GET"

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

// (发送POST请求：URL@文本、内容@字典/文本)
func FN_sendHTTPRequest_POST(receiver r.Element, values []r.Element) (r.Element, error) {
	const TOTAL_PARAMS = 2
	if len(values) != TOTAL_PARAMS {
		return nil, zerr.ExactParamsError(2)
	}

	url, err1 := value.AssertElement[*value.String](values[0])
	body := value.BuildEitherElement[*value.String, *value.HashMap](values[1])

	if err1 != nil {
		return nil, err1
	}
	// #1. get exact type of params
	const defaultTimeout = 30
	const method = "GET"

	emptyHeader := [][2]string{}
	emptyQuery := [][2]string{}

	// #2. handle body
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

	req, err := buildBaseHttpRequest(
		url.String(),
		method,
		emptyHeader,
		emptyQuery,
		reqBody,
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
		RegisterFunction("发送HTTP请求", value.NewFunction(FN_sendHTTPRequest)).
		RegisterFunction("发送GET请求", value.NewFunction(FN_sendHTTPRequest_GET)).
		RegisterFunction("发送POST请求", value.NewFunction(FN_sendHTTPRequest_POST))
}
