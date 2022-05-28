package stdlib

import (
	"encoding/json"
	"fmt"
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
	"math"
	"strconv"
	"strings"
)

var moduleName = "JSON"
var jsonModule = r.NewModule(moduleName)


// parseJsonFunc - 解析JSON
func parseJsonFunc(c *r.Context, values []r.Value) (r.Value, error) {
	// validate string ONLY
	if err := value.ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	v := values[0].(*value.String)
	vd := []byte(v.String())

	result := map[string]interface{}{}
	// core logic (via Go's stdlib)
	if err := json.Unmarshal(vd, &result); err != nil {
		return nil, zerr.NewErrorSLOT("解析JSON失败：" + err.Error())
	}

	return buildHashMapItem(result), nil
}

// generateJsonFunc - 生成JSON
func generateJsonFunc(c *r.Context, values []r.Value) (r.Value, error) {
	if err := value.ValidateExactParams(values, "hashmap"); err != nil {
		return nil, err
	}
	data, err := json.Marshal(buildPlainStrItem(values[0]))
	if err != nil {
		return nil, zerr.NewErrorSLOT("生成JSON失败：" + err.Error())
	}
	return value.NewString(string(data)), nil
}

// buildHashMapItem -
func buildHashMapItem(item interface{}) r.Value {
	if item == nil { // nil for json value "null"
		return value.NewNull()
	}
	switch vv := item.(type) {
	case float64:
		return value.NewNumber(vv)
	case string:
		return value.NewString(vv)
	case bool:
		return value.NewBool(vv)
	case map[string]interface{}:
		target := value.NewHashMap([]value.KVPair{})
		for k, v := range vv {
			finalValue := buildHashMapItem(v)
			target.AppendKVPair(value.KVPair{
				Key:   k,
				Value: finalValue,
			})
		}
		return target
	case []interface{}:
		varr := value.NewArray([]r.Value{})
		for _, vitem := range vv {
			varr.AppendValue(buildHashMapItem(vitem))
		}
		return varr
	default:
		return value.NewString(fmt.Sprintf("%v", vv))
	}
}

// buildPlainStrItem - from r.Value -> plain interface{} value
func buildPlainStrItem(item r.Value) interface{} {
	switch vv := item.(type) {
	case *value.Null:
		return nil
	case *value.String:
		return vv.String()
	case *value.Bool:
		return vv.GetValue()
	case *value.Number:
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
	case *value.Array:
		var resultList []interface{}
		for _, vi := range vv.GetValue() {
			resultList = append(resultList, buildPlainStrItem(vi))
		}
		return resultList
	case *value.HashMap:
		resultMap := map[string]interface{}{}
		for k, vi := range vv.GetValue() {
			resultMap[k] = buildPlainStrItem(vi)
		}
		return resultMap
	}
	return nil
}

func init() {
	jsonModule.RegisterFunction("解析JSON", parseJsonFunc)
	jsonModule.RegisterFunction("生成JSON", generateJsonFunc)

	RegisterModule(moduleName, jsonModule)
}