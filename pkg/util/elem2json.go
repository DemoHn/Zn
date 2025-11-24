package util

import (
	"encoding/json"
	"fmt"

	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

func HashMapToJSONString(hm *value.HashMap) (*value.String, error) {
	data, err := json.Marshal(buildPlainValueFromElement(hm))
	if err != nil {
		return nil, value.ThrowException("生成JSON失败 - " + err.Error())
	}
	return value.NewString(string(data)), nil
}

func JSONStringToElement(jsonStr *value.String) (r.Element, error) {
	plainMap := map[string]any{}
	vdata := []byte(jsonStr.GetValue())
	if err := json.Unmarshal(vdata, &plainMap); err != nil {
		return nil, value.ThrowException("解析JSON失败 - " + err.Error())
	}

	return buildElementFromPlainValue(plainMap), nil
}

func ElementToJSONString(elem r.Element) (*value.String, error) {
	plainValue := buildPlainValueFromElement(elem)
	jsonStr, err := json.Marshal(plainValue)
	if err != nil {
		return nil, value.ThrowException("生成JSON失败 - " + err.Error())
	}
	return value.NewString(string(jsonStr)), nil
}

func buildPlainValueFromElement(elem r.Element) any {
	switch vv := elem.(type) {
	case *value.Null:
		return nil
	case *value.String:
		return vv.String()
	case *value.Bool:
		return vv.GetValue()
	case *value.Number:
		return vv.GetValue()
	case *value.Array:
		var resultList []interface{}
		for _, vi := range vv.GetValue() {
			resultList = append(resultList, buildPlainValueFromElement(vi))
		}
		return resultList
	case *value.HashMap:
		resultMap := map[string]any{}
		for k, vi := range vv.GetValue() {
			resultMap[k] = buildPlainValueFromElement(vi)
		}
		return resultMap
	}
	return nil
}

func buildElementFromPlainValue(item any) r.Element {
	if item == nil {
		return value.NewNull()
	}
	switch vv := item.(type) {
	//// case#1: numbers
	// may cause precision lose !!
	case uint:
	case uint8:
	case uint16:
	case uint32:
	case uint64:
	case int:
	case int8:
	case int16:
	case int32:
	case int64:
		return value.NewNumber(float64(vv))
	case float64:
		return value.NewNumber(vv)
	//// case#2: strings
	case []rune:
		return value.NewString(string(vv))
	case string:
		return value.NewString(vv)
	//// case#3: booleans
	case bool:
		return value.NewBool(vv)
	case map[string]any:
		target := value.NewEmptyHashMap()
		for k, v := range vv {
			finalValue := buildElementFromPlainValue(v)
			target.AppendKVPair(value.KVPair{
				Key:   k,
				Value: finalValue,
			})
		}
		return target
	case []any:
		varr := value.NewEmptyArray()
		for _, vitem := range vv {
			varr.AppendValue(buildElementFromPlainValue(vitem))
		}
		return varr
	}
	// default fallback logic
	return value.NewString(fmt.Sprintf("%v", item))
}
