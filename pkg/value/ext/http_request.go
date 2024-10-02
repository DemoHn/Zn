package ext

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

type httpRequestGetterFunc func(*HTTPRequest, *r.Context) (r.Element, error)
type httpRequestMethodFunc func(*HTTPRequest, *r.Context, []r.Element) (r.Element, error)

type HTTPRequest struct {
	request *http.Request
}

// NewHTTPRequest - new HTTPRequest from raw *http.Request
func NewHTTPRequest(req *http.Request) *HTTPRequest {
	return &HTTPRequest{request: req}
}

// GetProperty - get property of HTTPRequest
func (h *HTTPRequest) GetProperty(c *r.Context, name string) (r.Element, error) {
	httpRequestGetterMap := map[string]httpRequestGetterFunc{
		"路径":   httpRequestGetPath,
		"方法":   httpRequestGetMethod,
		"头部":   httpRequestGetHeaders,
		"查询参数": httpRequestGetQueryParams,
	}

	if fn, ok := httpRequestGetterMap[name]; ok {
		return fn(h, c)
	}
	return nil, zerr.PropertyNotFound(name)
}

func httpRequestGetPath(h *HTTPRequest, c *r.Context) (r.Element, error) {
	return value.NewString(h.request.URL.Path), nil
}

func httpRequestGetMethod(h *HTTPRequest, c *r.Context) (r.Element, error) {
	return value.NewString(h.request.Method), nil
}

func httpRequestGetHeaders(h *HTTPRequest, c *r.Context) (r.Element, error) {
	queryPair := []value.KVPair{}
	for k, v := range h.request.Header {
		if len(v) > 0 {
			queryPair = append(queryPair, value.KVPair{
				Key:   k,
				Value: value.NewString(v[0]),
			})
		}
	}
	return value.NewHashMap(queryPair), nil
}

func httpRequestGetQueryParams(h *HTTPRequest, c *r.Context) (r.Element, error) {
	queryPair := []value.KVPair{}
	for k, v := range h.request.URL.Query() {
		if len(v) > 0 {
			if len(v) > 1 {
				qsArr := []r.Element{}
				for _, qs := range v {
					qsArr = append(qsArr, value.NewString(qs))
				}
				queryPair = append(queryPair, value.KVPair{
					Key:   k,
					Value: value.NewArray(qsArr),
				})
			} else {
				queryPair = append(queryPair, value.KVPair{
					Key:   k,
					Value: value.NewString(v[0]),
				})
			}
		}
	}
	return value.NewHashMap(queryPair), nil
}

// SetProperty - set property of HTTPRequest
func (h *HTTPRequest) SetProperty(c *r.Context, name string, value r.Element) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod - execute method of HTTPRequest
func (h *HTTPRequest) ExecMethod(c *r.Context, name string, values []r.Element) (r.Element, error) {
	methodMap := map[string]httpRequestMethodFunc{
		"读取内容": httpRequestExecReadBody,
	}

	if method, ok := methodMap[name]; ok {
		return method(h, c, values)
	}
	return nil, zerr.MethodNotFound(name)
}

func httpRequestExecReadBody(h *HTTPRequest, c *r.Context, values []r.Element) (r.Element, error) {
	// impl GetBody function here
	body, err := ioutil.ReadAll(h.request.Body)
	if err != nil {
		return nil, value.ThrowException("读取请求内容出现异常")
	}

	contentType := h.request.Header.Get("Content-Type")
	if contentType == "application/json" {
		var jsonBody interface{}
		if err := json.Unmarshal(body, &jsonBody); err != nil {
			return nil, value.ThrowException("将请求内容解析成JSON格式时出现异常")
		}
		return buildHashMapItem(jsonBody), nil
	}
	return value.NewString(string(body)), nil
}

// buildHashMapItem - from plain object to HashMap Element
func buildHashMapItem(item interface{}) r.Element {
	if item == nil { // nil for json value "null"
		return value.NewNull()
	}
	switch vv := item.(type) {
	case float64:
		return value.NewNumber(vv)
	case string:
		return value.NewString(vv)
	case bool:
		return value.NewBool(vv)
	case map[string]interface{}:
		target := value.NewHashMap([]value.KVPair{})
		for k, v := range vv {
			finalValue := buildHashMapItem(v)
			target.AppendKVPair(value.KVPair{
				Key:   k,
				Value: finalValue,
			})
		}
		return target
	case []interface{}:
		varr := value.NewArray([]r.Element{})
		for _, vitem := range vv {
			varr.AppendValue(buildHashMapItem(vitem))
		}
		return varr
	default:
		return value.NewString(fmt.Sprintf("%v", vv))
	}
}
