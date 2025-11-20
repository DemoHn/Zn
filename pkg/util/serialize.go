package util

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"

	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

func HashMapToJSONString(hm *value.HashMap) (string, error) {
	data, err := json.Marshal(buildPlainValueFromElement(hm))
	if err != nil {
		return "", value.ThrowException("生成JSON失败 - " + err.Error())
	}
	return string(data), nil
}

func ElementToPlainValue(elem r.Element) interface{} {
	return buildPlainValueFromElement(elem)
}

func buildPlainValueFromElement(elem r.Element) interface{} {
	switch vv := elem.(type) {
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
			resultList = append(resultList, buildPlainValueFromElement(vi))
		}
		return resultList
	case *value.HashMap:
		resultMap := map[string]interface{}{}
		for k, vi := range vv.GetValue() {
			resultMap[k] = buildPlainValueFromElement(vi)
		}
		return resultMap
	}
	return nil
}
