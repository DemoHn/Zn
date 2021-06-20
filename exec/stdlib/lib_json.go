package stdlib

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/exec/val"
)

var libjsonValueMap = map[string]ctx.Value{
	"解析JSON": val.NewFunction("解析JSON", parseJsonFunc),
	"生成JSON": val.NewFunction("生成JSON", generateJsonFunc),
}

// parseJsonFunc - 解析JSON
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

	return buildHashMapItem(result), nil
}

// generateJsonFunc - 生成JSON
func generateJsonFunc(c *ctx.Context, values []ctx.Value) (ctx.Value, *error.Error) {
	if err := val.ValidateExactParams(values, "hashmap"); err != nil {
		return nil, err
	}
	data, err := json.Marshal(buildPlainStrItem(values[0]))
	if err != nil {
		return nil, error.NewErrorSLOT("生成JSON失败：" + err.Error())
	}
	return val.NewString(string(data)), nil
}

// buildHashMapItem -
func buildHashMapItem(item interface{}) ctx.Value {
	if item == nil { // nil for json value "null"
		return val.NewNull()
	}
	switch vv := item.(type) {
	case float64:
		return val.NewDecimalFromFloat64(vv)
	case string:
		return val.NewString(vv)
	case bool:
		return val.NewBool(vv)
	case map[string]interface{}:
		target := val.NewHashMap([]val.KVPair{})
		for k, v := range vv {
			finalValue := buildHashMapItem(v)
			target.AppendKVPair(val.KVPair{
				Key:   k,
				Value: finalValue,
			})
		}
		return target
	case []interface{}:
		varr := val.NewArray([]ctx.Value{})
		for _, vitem := range vv {
			varr.AppendValue(buildHashMapItem(vitem))
		}
		return varr
	default:
		return val.NewString(fmt.Sprintf("%v", vv))
	}
}

// buildPlainStrItem - from ctx.Value -> plain interface{} value
func buildPlainStrItem(item ctx.Value) interface{} {
	switch vv := item.(type) {
	case *val.Null:
		return nil
	case *val.String:
		return vv.String()
	case *val.Bool:
		return vv.GetValue()
	case *val.Decimal:
		valStr := vv.String()
		valStr = strings.Replace(valStr, "*10^", "e", 1)
		// replace *10^ -> e
		result, err := strconv.ParseFloat(valStr, 64)
		// Sometimes parseFloat may fail due to overflow, underflow etc.
		// For those invalid numbers, return NaN instead.
		if err != nil {
			return math.NaN()
		}
		return result
	case *val.Array:
		resultList := []interface{}{}
		for _, vi := range vv.GetValue() {
			resultList = append(resultList, buildPlainStrItem(vi))
		}
		return resultList
	case *val.HashMap:
		resultMap := map[string]interface{}{}
		for k, vi := range vv.GetValue() {
			resultMap[k] = buildPlainStrItem(vi)
		}
		return resultMap
	}
	return nil
}
