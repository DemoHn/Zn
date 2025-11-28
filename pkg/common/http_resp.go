package common

import (
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

var CLASS_HttpResponse = value.NewClassModel("HTTP响应").
	DefineProperty("状态码", value.NewNumber(200)).
	DefineProperty("头部", value.NewEmptyHashMap()).
	DefineProperty("内容", value.NewString("")).
	SetConstructor(httpResponseConstructor)

func httpResponseConstructor(self r.Element, values []r.Element) (r.Element, error) {
	if err := value.ValidateLeastParams(values, "number", "any", "hashmap?"); err != nil {
		return nil, err
	}

	self.SetProperty("状态码", values[0])

	// set content & default header (content-type)
	switch v := values[1].(type) {
	case *value.String:
		self.SetProperty("头部", value.NewHashMap([]value.KVPair{
			{Key: "Content-Type", Value: value.NewString("text/plain")},
		}))
		self.SetProperty("内容", values[1])
	case *value.HashMap, *value.Array:
		bodyStr, err := ElementToJSONString(v)
		if err != nil {
			return nil, err
		}
		self.SetProperty("头部", value.NewHashMap([]value.KVPair{
			{Key: "Content-Type", Value: value.NewString("application/json")},
		}))
		self.SetProperty("内容", bodyStr)
	default:
		bodyStr, err := ElementToJSONString(v)
		if err != nil {
			return nil, err
		}
		self.SetProperty("头部", value.NewHashMap([]value.KVPair{
			{Key: "Content-Type", Value: value.NewString("text/plain")},
		}))
		self.SetProperty("内容", bodyStr)
	}

	// override headers
	if len(values) == 3 {
		self.SetProperty("头部", values[2])
	}
	return self, nil
}
