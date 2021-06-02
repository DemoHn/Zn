package stdlib

import (
	"encoding/json"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/exec/val"
)

var libjsonMethodMap = map[string]funcExecutor{}

func init() {
	libjsonMethodMap = map[string]funcExecutor{
		"解析JSON": parseJsonFunc,
	}
}

func parseJsonFunc(c *ctx.Context, values []ctx.Value) (ctx.Value, *error.Error) {
	// validate string ONLY
	if err := val.ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	v := values[0].(*val.String)
	vd := []byte(v.String())

	result := map[string]interface{}{}
	// core logic (via Go's stdlib)
	if err := json.Unmarshal(vd, &result); err != nil {
		return nil, error.NewErrorSLOT("解析JSON失败：" + err.Error())
	}

	// TODO: add value type parser
}
