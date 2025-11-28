package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/DemoHn/Zn/pkg/common"
	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/DemoHn/Zn/pkg/runtime"
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
		return common.JSONStringToElement(value.NewString(string(body)))
	}
	return value.NewString(string(body)), nil
}

func buildIncomingRequest(r *http.Request) (runtime.Element, error) {
	headerDict := value.NewEmptyHashMap()
	for k, v := range r.Header {
		if len(v) > 0 {
			headerDict.AppendKVPair(value.KVPair{
				Key:   k,
				Value: value.NewString(v[0]),
			})
		}
	}

	qsDict := value.NewEmptyHashMap()
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			qsDict.AppendKVPair(value.KVPair{
				Key:   k,
				Value: value.NewString(v[0]),
			})
		}
	}

	body, err := buildIncomingRequestBody(r)
	if err != nil {
		return nil, err
	}

	reqObj := value.NewObject(common.CLASS_HttpRequest, runtime.ElementMap{
		"URL":  value.NewString(r.URL.String()),
		"路径":   value.NewString(r.URL.Path),
		"方法":   value.NewString(r.Method),
		"头部":   headerDict,
		"查询参数": qsDict,
		"内容":   body,
	})

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
			jsonStr, err := common.ElementToJSONString(v)
			// write resp body
			if err != nil {
				respondError(w, err)
			} else {
				respondJSON(w, []byte(jsonStr.GetValue()))
			}
		case *value.Object:
			if v.IsInstanceOf(common.CLASS_HttpResponse) {
				content, _ := v.GetProperty("内容")
				respHeader, _ := v.GetProperty("头部")
				statusCode, _ := v.GetProperty("状态码")

				// write resp body directly
				var contentStr string
				if v, ok := content.(*value.String); ok {
					contentStr = v.GetValue()
				} else {
					jsonStr, err := common.ElementToJSONString(content)
					if err != nil {
						respondError(w, err)
						return
					}
					contentStr = jsonStr.String()
				}
				// write to response directly
				for k, v := range respHeader.(*value.HashMap).GetValue() {
					w.Header().Add(k, v.String())
				}
				w.WriteHeader(int(statusCode.(*value.Number).GetValue()))
				w.Write([]byte(contentStr))
			}
		default:
			respondOK(w, v.String())
		}
	}
}
