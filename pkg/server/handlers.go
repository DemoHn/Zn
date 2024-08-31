package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

type playgroundReq struct {
	VarInput   string
	SourceCode string
}

// if exec code ok - return 200 with result (ignore outputs from「显示」)
func PlaygroundHandler(w http.ResponseWriter, r *http.Request) {
	// exec program
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondError(w, fmt.Errorf("读取请求内容出现异常：%s", err.Error()))
		return
	}

	var reqInfo playgroundReq
	if err := json.Unmarshal(body, &reqInfo); err != nil {
		respondError(w, fmt.Errorf("解析请求格式不符合要求！"))
		return
	}

	pExecutor := exec.NewPlaygroundExecutor([]byte(reqInfo.SourceCode))
	// exec varInputs
	varInputs, err := exec.ExecVarInputs(reqInfo.VarInput)
	if err != nil {
		respondError(w, err)
		return
	}

	rtnValue, err := pExecutor.RunCode(varInputs)
	if err != nil {
		respondError(w, err)
		return
	}

	// write return value as resp body
	switch rtnValue.(type) {
	case *value.Null:
		respondOK(w, "")
	default:
		respondOK(w, value.StringifyValue(rtnValue))
	}
}

func HTTPHandlerWithEntry(entryFile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerObject := []value.KVPair{}
		for k, v := range r.Header {
			if len(v) > 0 {
				headerObject = append(headerObject, value.KVPair{
					Key: k,
					// only read *first* occured value of the header name
					Value: value.NewString(v[0]),
				})
			}
		}
		// TODO: handle binary data of request body (e.g. file)
		// the format of final request body varies from different Content-Type: by default, it's a plain text, but when Content-Type = application/json, it's a json hashmap
		var requestBody runtime.Element

		bodyData, _ := ioutil.ReadAll(r.Body)
		switch r.Header.Get("Content-Type") {
		case "application/json":
			var unmarshalMap map[string]interface{}
			// try unmarshal json - if failed then return string directly
			err := json.Unmarshal(bodyData, &unmarshalMap)
			if err != nil {
				requestBody = value.NewString(string(bodyData))
			} else {
				requestBody = buildHashMapItem(unmarshalMap)
			}
		default:
			requestBody = value.NewString(string(bodyData))
		}

		// construct http request: zinc hashMap
		httpReqObject := value.NewHashMap([]value.KVPair{
			{
				Key:   "请求方法",
				Value: value.NewString(r.Method),
			},
			{
				Key:   "路径",
				Value: value.NewString(r.URL.Path),
			},
			{
				Key:   "请求头",
				Value: value.NewHashMap(headerObject),
			},
			{
				Key:   "请求参数",
				Value: requestBody,
			},
		})

		executor := exec.NewFileExecutor(entryFile)
		result, err := executor.RunCode(map[string]runtime.Element{
			"当前请求": httpReqObject,
		})

		if err != nil {
			respondError(w, err)
		} else {
			switch v := result.(type) {
			case *value.String:
				respondOK(w, v.GetValue())
			case *value.Number:
				respondOK(w, fmt.Sprintf("%v", v.GetValue()))
			case *value.HashMap, *value.Array:
				jsonBytes, _ := json.Marshal(buildPlainStrItem(v))
				// write resp body
				respondJSON(w, jsonBytes)
			default:
				respondOK(w, value.StringifyValue(v))
			}
		}
	}
}

// buildHashMapItem - from plain object to HashMap Element
func buildHashMapItem(item interface{}) runtime.Element {
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
		varr := value.NewArray([]runtime.Element{})
		for _, vitem := range vv {
			varr.AppendValue(buildHashMapItem(vitem))
		}
		return varr
	default:
		return value.NewString(fmt.Sprintf("%v", vv))
	}
}

// buildPlainStrItem - from r.Value -> plain interface{} value
func buildPlainStrItem(item runtime.Element) interface{} {
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
	// TODO: stringify other objects
	return value.StringifyValue(item)
}

//// helpers

func respondOK(w http.ResponseWriter, body string) {
	// declare content-type = text/plain (not json/file/...)
	w.Header().Add("Content-Type", "text/plain; charset=\"utf-8\"")
	// http status: 200 OK
	w.WriteHeader(http.StatusOK)
	// write resp body
	io.WriteString(w, body)
}

func respondJSON(w http.ResponseWriter, body []byte) {
	// declare content-type = text/plain (not json/file/...)
	w.Header().Add("Content-Type", "application/json; charset=\"utf-8\"")
	// http status: 200 OK
	w.WriteHeader(http.StatusOK)
	// write resp body
	w.Write(body)
}

func respondError(w http.ResponseWriter, reason error) {
	// declare content-type = text/plain (not json/file/...)
	w.Header().Add("Content-Type", "text/plain; charset=\"utf-8\"")
	// http status: 500 Internal Error
	w.WriteHeader(http.StatusInternalServerError)
	// write resp body
	io.WriteString(w, reason.Error())
}
