package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/util"
	"github.com/DemoHn/Zn/pkg/value"
)

type ZnHttpHandler struct {
	interpreter *exec.Interpreter
	entryFile   string
}

func NewZnHttpHandler(interpreter *exec.Interpreter, entryFile string) *ZnHttpHandler {
	return &ZnHttpHandler{interpreter: interpreter, entryFile: entryFile}
}

func (h *ZnHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqObj, err := buildIncomingRequest(r)
	if err != nil {
		sendHTTPResponse(nil, err, w)
		return
	}

	varInput := runtime.ElementMap{
		"当前请求": reqObj,
	}
	// execute code
	rtnValue, err := h.interpreter.LoadFile(h.entryFile).Execute(varInput)
	sendHTTPResponse(rtnValue, err, w)

}

/** construct 传入请求 (incoming request) object for HTTP handler input
NOTE: the 传入请求 class is different from @HTTP-HTTP请求 class,
  1. the properties of 传入请求 class is READONLY
  2. no other method (like 发送请求) for the class

Actually, the 传入请求 class is just a dataclass to map incoming *http.Request data!

定义传入请求：
	其URL = “http://127.0.0.1:3862/create”
	其路径 = “/create”
	其方法 = “POST”
	其头部 = [
		“Content-Length” = “20”
		“Content-Type” = “application/json”
	]
	其查询参数 = [
		“A” = “20”
		“B” = [“30”、“40”]
	]
*/

type IncomingRequest struct {
	Request *http.Request
	URL     *value.String
	Path    *value.String
	Method  *value.String
	Headers *value.HashMap
	Query   *value.HashMap
}

func (iq *IncomingRequest) String() string {
	return fmt.Sprintf("‹对象·传入请求 (URL=%s)›", iq.URL.String())
}

func (iq *IncomingRequest) GetProperty(name string) (runtime.Element, error) {
	switch name {
	case "URL":
		return iq.URL, nil
	case "路径":
		return iq.Path, nil
	case "方法":
		return iq.Method, nil
	case "头部":
		return iq.Headers, nil
	case "查询参数":
		return iq.Query, nil
	default:
		return nil, zerr.PropertyNotFound(name)
	}
}

func (iq *IncomingRequest) SetProperty(name string, value runtime.Element) error {
	return zerr.PropertyNotFound(name)
}

func (iq *IncomingRequest) ExecMethod(name string, values []runtime.Element) (runtime.Element, error) {
	if name == "读取内容" {
		return buildIncomingRequestBody(iq.Request)
	}

	return nil, zerr.MethodNotFound(name)
}

func buildIncomingRequestBody(req *http.Request) (runtime.Element, error) {
	if req.Body == nil {
		return value.NewString(""), nil
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, value.ThrowException("读取请求内容出现异常")
	}

	contentType := req.Header.Get("Content-Type")
	if contentType == "application/json" {
		return util.JSONStringToElement(value.NewString(string(body)))
	}
	return value.NewString(string(body)), nil
}

func buildIncomingRequest(r *http.Request) (runtime.Element, error) {
	reqObj := &IncomingRequest{
		Request: r,
		URL:     value.NewString(r.URL.String()),
		Path:    value.NewString(r.URL.Path),
		Method:  value.NewString(r.Method),
		Headers: value.NewEmptyHashMap(),
		Query:   value.NewEmptyHashMap(),
	}
	for k, v := range r.Header {
		reqObj.Headers.AppendKVPair(value.KVPair{
			Key:   k,
			Value: value.NewString(v[0]),
		})
	}
	for k, v := range r.URL.Query() {
		var pairValue runtime.Element
		if len(v) == 0 {
			continue
		} else if len(v) == 1 {
			pairValue = value.NewString(v[0])
		} else {
			pairValue = value.NewEmptyArray()
			for _, vItem := range v {
				pairValue.(*value.Array).AppendValue(value.NewString(vItem))
			}
		}
		reqObj.Query.AppendKVPair(value.KVPair{
			Key:   k,
			Value: pairValue,
		})
	}
	return reqObj, nil
}

func sendHTTPResponse(result runtime.Element, err error, w http.ResponseWriter) {
	if err != nil {
		respondError(w, err)
	} else {
		switch v := result.(type) {
		case *value.String:
			respondOK(w, v.GetValue())
		case *value.Number:
			respondOK(w, fmt.Sprintf("%v", v.GetValue()))
		case *value.HashMap, *value.Array:
			plainValue := util.ElementToPlainValue(v)
			jsonBytes, _ := json.Marshal(plainValue)
			// write resp body
			respondJSON(w, jsonBytes)
		default:
			respondOK(w, v.String())
		}
	}
}
