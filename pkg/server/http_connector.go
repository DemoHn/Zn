package server

import (
	"net/http"

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
	ref := value.NewClassModel("HTTP请求", nil)

	return nil, nil
}

func ConstructHTTPResponseObject() (runtime.Element, error) {
	return nil, nil
}
