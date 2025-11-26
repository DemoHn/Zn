package common

import (
	"fmt"
	"net/http"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

/*
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

type HttpResponse struct {
	StatusCode *value.Number
	Header     *value.HashMap
	Body       *value.String
}

//// implement Element's interface

func (r *HttpResponse) String() string {
	return fmt.Sprintf("‹对象·HTTP响应 (%d，“%s”)›", int(r.StatusCode.GetValue()), r.Body.GetValue())
}

// (新建HTTP响应：200、“response text”) or
// (新建HTTP响应：200、[1,2,3]) or
// (新建HTTP响应: 404、[A=1，B=2])
func (r *HttpResponse) Construct(params []r.Element) (r.Element, error) {
	// validate params first (at least 2 params, the first is status code, the second is header map, the third is body string)
	// 3rd param is optional, represents for additional resp headers
	err := value.ValidateLeastParams(params, "number", "any", "hashmap?")
	if err != nil {
		return nil, err
	}
	p0 := params[0].(*value.Number)
	p1, err := ElementToJSONString(params[1])
	if err != nil {
		return nil, err
	}

	resp := &HttpResponse{
		StatusCode: p0,
		Header:     value.NewEmptyHashMap(),
		Body:       p1,
	}
	if len(params) > 2 {
		p2 := params[2].(*value.HashMap)
		resp.Header = p2
	}
	return resp, nil
}

func (r *HttpResponse) GetProperty(name string) (r.Element, error) {
	switch name {
	case "状态码":
		return r.StatusCode, nil
	case "头部":
		return r.Header, nil
	case "内容":
		return r.Body, nil
	}
	return nil, zerr.PropertyNotFound(name)
}

// ALL
func (r *HttpResponse) SetProperty(name string, v r.Element) error {
	switch name {
	case "状态码":
		r.StatusCode = v.(*value.Number)
	case "头部":
		r.Header = v.(*value.HashMap)
	case "内容":
		r.Body = v.(*value.String)
	}
	return zerr.PropertyNotFound(name)
}

func (r *HttpResponse) ExecMethod(name string, params []r.Element) (r.Element, error) {
	return nil, zerr.MethodNotFound(name)
}

func BuildHttpResponse(resp *http.Response, body []byte) *value.Object {
	initProps := map[string]r.Element{
		"状态码": value.NewNumber(float64(resp.StatusCode)),
		"头部":  value.NewEmptyHashMap(),
		"内容":  value.NewString(string(body)),
	}

	headerValue := value.NewEmptyHashMap()
	for k, v := range resp.Header {
		if len(v) > 0 {
			headerValue.AppendKVPair(value.KVPair{
				Key:   k,
				Value: value.NewString(v[0]),
			})
		}

	}
	initProps["头部"] = headerValue

	return value.NewObject(CLASS_HttpResposne, initProps)
}
