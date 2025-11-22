package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/util"
	"github.com/DemoHn/Zn/pkg/value"
	"github.com/DemoHn/Zn/pkg/value/ext"
)

type ZnHttpHandler struct {
	interpreter *exec.Interpreter
	entryFile   string
}

func (h *ZnHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqObj, err := ConstructHTTPRequestObject(r)
	if err != nil {
		SendHTTPResponse(nil, err, w)
		return
	}

	varInput := runtime.ElementMap{
		"当前请求": reqObj,
	}
	// execute code
	rtnValue, err := h.interpreter.LoadFile(h.entryFile).Execute(varInput)
	SendHTTPResponse(rtnValue, err, w)

}

/** construct 当前请求 object for HTTP handler input
NOTE: Definition of class `HTTP请求`:

定义HTTP请求：
	其路径 = “/create”
	其方法 = “POST”
	其头部 =【
		“Content-Length” = “20”
		“Content-Type” = “application/json”
	】
	其查询参数 = 【
		“A” = “20”
		“B” =【“30”、“40”】
	】

	// 根据 content-type 智能解析HTTP请求内容，得到对应的对象
	如何解析内容，得到内容对象

NOTE2: Definition of class `HTTP响应`:
定义HTTP响应：
	其状态码 = 200
	其头部 = 【
		“Content-Type” =「text/plain; charset="utf-8"」
	】
	其内容 = “Hello World”
*/

func ConstructHTTPRequestObject(r *http.Request) (runtime.Element, error) {
	return ext.NewHTTPRequest(r), nil
}

func SendHTTPResponse(result runtime.Element, err error, w http.ResponseWriter) {
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
			respondOK(w, value.StringifyValue(v))
		}
	}
}
