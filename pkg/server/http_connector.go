package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

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
	headerMap := map[string]string{}
	for k, v := range r.Header {
		if len(v) > 0 {
			headerMap[k] = v[0]
		}
	}

	// build query params
	queryPair := []value.KVPair{}
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			if len(v) > 1 {
				qsArr := []runtime.Element{}
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

	initialProps := map[string]runtime.Element{
		"方法":   value.NewString(r.Method),
		"路径":   value.NewString(r.URL.Path),
		"头部":   buildHashMapItem(headerMap),
		"查询参数": value.NewHashMap(queryPair),
		"-goHttpRequest-": value.NewGoValue("*http.Request", r),
	}

	return value.NewObject(exec.ZnConstHTTPRequestClass, initialProps), nil
}

func ConstructHTTPResponseObject() (runtime.Element, error) {
	return nil, nil
}

// TODO - use ConstructHTTPResponseObject() to send response
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
			jsonBytes, _ := json.Marshal(buildPlainStrItem(v))
			// write resp body
			respondJSON(w, jsonBytes)
		default:
			respondOK(w, value.StringifyValue(v))
		}
	}
}
