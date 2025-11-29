package http

import (
	"fmt"
	"io"
	libHTTP "net/http"
	libURL "net/url"
	"strings"
	"time"

	"github.com/DemoHn/Zn/pkg/common"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

type ReqBody struct {
	// ContentType - describe the content-type of
	ContentType string
	Value       string
}

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

// accept HashMap or String
func digestReqBody(reqBodyElement r.Element) (ReqBody, error) {

}

func buildBaseHttpRequest(
	url string,
	method string,
	headers [][2]string, // [][key=xxx, value=xxx]
	queryParams [][2]string,
	body ReqBody,
) (*libHTTP.Request, error) {
	// #1. parse url
	urlObj, err := libURL.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("解析URL时出现异常 - %s", err.Error())
	}
	// #2. normalize method
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

	// #3. if queryParams content has params, then override original urlObj's rawQuery
	if len(queryParams) > 0 {
		qValues := libURL.Values{}
		for _, qpair := range queryParams {
			k := qpair[0]
			v := qpair[1]

			if qValues.Get(k) == "" {
				qValues.Set(k, v)
			} else {
				qValues.Add(k, v)
			}
		}
		urlObj.RawQuery = qValues.Encode()
	}

	finalHeaders := libHTTP.Header{}
	// #4. set body
	if body.ContentType == "" {
		// set default content-type header
		finalHeaders.Set("Content-Type", "text/plain")
	} else {
		finalHeaders.Set("Content-Type", body.ContentType)
	}
	bodyReadCloser := io.NopCloser(strings.NewReader(body.Value))

	// #5. set headers (may override content-type)
	for _, headerPair := range headers {
		k := headerPair[0]
		v := headerPair[1]
		if finalHeaders.Get(k) == "" {
			finalHeaders.Set(k, v)
		} else {
			finalHeaders.Add(k, v)
		}
	}

	return &libHTTP.Request{
		Method:     method,
		URL:        urlObj,
		Header:     finalHeaders,
		Body:       bodyReadCloser,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}, nil
}

func sendHttpRequest(req *libHTTP.Request, allowRedicrect bool, timeout int) (*libHTTP.Response, []byte, error) {
	client := &libHTTP.Client{
		CheckRedirect: func(req *libHTTP.Request, via []*libHTTP.Request) error {
			if !allowRedicrect {
				return fmt.Errorf("此请求不允许自动重定向")
			}
			return nil
		},
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, []byte{}, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, []byte{}, err
	}

	return resp, content, nil
}

func buildOBJ_HttpResponse(resp *libHTTP.Response, body []byte) *value.Object {
	headerHashMap := value.NewEmptyHashMap()
	for k, v := range resp.Header {
		if len(v) > 0 {
			headerHashMap.AppendKVPair(value.KVPair{
				Key:   k,
				Value: value.NewString(v[0]),
			})
		}
	}

	return value.NewObject(common.CLASS_HttpResponse, map[string]r.Element{
		"状态码": value.NewNumber(float64(resp.StatusCode)),
		"内容":  value.NewString(string(body)),
		"头部":  headerHashMap,
	})
}
