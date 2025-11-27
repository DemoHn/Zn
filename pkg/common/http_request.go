package common

import (
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

/*
定义HTTP请求：

	其路径 = “”
	其方法 = “”

	注：头部 (headers) 是请求头部的信息；在进行 HTTP 请求时，会自动添加一些默认的头部信息，如
		如 Content-Type、Host、Content-Length等；HTTP请求#头部即是在此基础上添加额外的头部信息.
	其头部 = [=]

	注：查询参数 (queryString) 是URL中的 ? 后面的部分；
		如 URL: https://abc.com?a=1&b=2&c=3，对应的查询参数是 ‘a=1&b=2&c=3’

		在定义HTTP请求的‘查询参数’属性时，可以是文本 (String)，也可以是字典 (HashMap) 类型，
		- 如果是文本类型，那就是原始的查询参数字符串
		- 如果是字典类型，那就是将字典转换成的查询参数字符串
	其查询参数 = [=]

	注：内容 (body) 是请求体的内容，主要用在 POST 等需要传输大量数据的地方

		内容可以是文本 (String)，也可以是字典 (HashMap)
		- 若是文本，本次请求的 Content-Type 为 application/x-www-form-urlencoded，其内容即是诸如 a=1&b=2&c=3 这样的文本
		- 若是字典，本次请求的 Content-Type 为 application/json，实际传输内容会转换成JSON字符串，诸如 {"a": 1, "b": 2, "c": 3}
	其内容 = “”
*/
var CLASS_HttpRequest = value.NewClassModel("HTTP请求").
	DefineProperty("URL", value.NewString("")).
	DefineProperty("路径", value.NewString("")).
	DefineProperty("方法", value.NewString("GET")).
	DefineProperty("头部", value.NewEmptyHashMap()).
	DefineProperty("查询参数", value.NewEmptyHashMap()).
	DefineProperty("内容", value.NewString("")).
	SetConstructor(httpRequestConsturctor)

func httpRequestConsturctor(self r.Element, values []r.Element) (r.Element, error) {
	// params: 方法、URL、内容 (文本 or 字典)
	if err := value.ValidateLeastParams(values, "string", "string", "any?"); err != nil {
		return nil, err
	}
	self.SetProperty("方法", values[0])
	self.SetProperty("URL", values[1])

	// IF there's 3rd param
	if len(values) == 3 {
		switch v := values[2].(type) {
		case *value.String:
			self.SetProperty("头部", value.NewHashMap([]value.KVPair{
				{Key: "Content-Type", Value: value.NewString("text/plain")},
			}))
			self.SetProperty("内容", values[2])
		case *value.HashMap:
			bodyStr, err := HashMapToJSONString(v)
			if err != nil {
				return nil, err
			}
			self.SetProperty("头部", value.NewHashMap([]value.KVPair{
				{Key: "Content-Type", Value: value.NewString("application/json")},
			}))
			self.SetProperty("内容", bodyStr)
		}
	}
	return self, nil
}
