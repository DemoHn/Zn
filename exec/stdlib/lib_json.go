package stdlib

import (
	"encoding/json"
	"fmt"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/exec/val"
)

var libjsonMethodMap = map[string]ctx.Value{
	"解析JSON": val.NewFunction("解析JSON", parseJsonFunc),
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

	return buildHashMap(result), nil
}

// buildHashMap - build hash map
func buildHashMap(source map[string]interface{}) *val.HashMap {
	target := val.NewHashMap([]val.KVPair{})
	for k, v := range source {
		finalValue := buildHashMapItem(v)
		target.AppendKVPair(val.KVPair{
			Key:   k,
			Value: finalValue,
		})
	}
	return target
}

// buildHashMapItem -
func buildHashMapItem(item interface{}) ctx.Value {
	var finalValue ctx.Value
	switch vv := item.(type) {
	case float64:
		finalValue = val.NewDecimalFromInt(int(vv), 0)
	case string:
		finalValue = val.NewString(vv)
	case bool:
		finalValue = val.NewBool(vv)
	case map[string]interface{}:
		finalValue = buildHashMap(vv)
	case []interface{}:
		varr := val.NewArray([]ctx.Value{})
		for _, vitem := range vv {
			varr.AppendValue(buildHashMapItem(vitem))
		}
		finalValue = varr
	default:
		finalValue = val.NewString(fmt.Sprintf("%v", vv))
	}

	return finalValue
}
