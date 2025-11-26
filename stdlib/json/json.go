package json

import (
	"github.com/DemoHn/Zn/pkg/common"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

const STDLIB_JSON_NAME = "@JSON"

var jsonLIB *r.Library

/*

testing...
如何解析JSON？
	输入文本

	输出[[ parseJson $1 ]]

	拦截异常：
	    xxxx
*/
// parseJsonFunc - 解析JSON
func FN_parseJson(receiver r.Element, values []r.Element) (r.Element, error) {
	// validate string ONLY
	if err := value.ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	// get exact type of params
	p1 := values[0].(*value.String)

	return common.JSONStringToElement(p1)
}

// generateJsonFunc - 生成JSON
func FN_generateJson(receiver r.Element, values []r.Element) (r.Element, error) {
	if err := value.ValidateExactParams(values, "hashmap"); err != nil {
		return nil, err
	}
	p1 := values[0].(*value.HashMap)
	return common.HashMapToJSONString(p1)
}

func Export() *r.Library {
	return jsonLIB
}

func init() {
	jsonLIB = r.NewLibrary(STDLIB_JSON_NAME)
	jsonLIB.RegisterFunction("解析JSON", value.NewFunction(FN_parseJson)).
		RegisterFunction("生成JSON", value.NewFunction(FN_generateJson))
}
