package http

import (
	"github.com/DemoHn/Zn/pkg/common"
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

var httpLIB *r.Library

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
