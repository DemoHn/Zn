package server

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

type RespWriter struct {
	writer *bufio.Writer
	header http.Header
}

func NewRespWriter(conn net.Conn) *RespWriter {
	return &RespWriter{
		writer: bufio.NewWriter(conn),
		header: http.Header{},
	}
}

//// fit http.ResponseWriter interface
func (w *RespWriter) Header() http.Header {
	return w.header
}

func (w *RespWriter) WriteHeader(statusCode int) {
	headString := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))

	w.writer.WriteString(headString)
	// write header (e.g. Content-Type: text/plain; charset="utf-8")
	w.header.Write(w.writer)
	// write \r\n
	w.writer.WriteString("\r\n")
}

func (w *RespWriter) Write(data []byte) (int, error) {
	nn, err := w.writer.Write(data)
	if err != nil {
		return nn, err
	}

	if err := w.writer.Flush(); err != nil {
		return nn, err
	}

	return nn, nil
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

//// parseConnUrl -> network + address
// currently support: tcp://, unix://
func parseConnUrl(connUrl string) (string, string, error) {
	u, err := url.Parse(connUrl)
	if err != nil {
		return "", "", err
	}
	switch u.Scheme {
	case TypeUnixSocket:
		return TypeUnixSocket, u.Path, nil
	case TypeTcp:
		return TypeTcp, u.Host, nil
	default:
		return "", "", fmt.Errorf("不支持的协议：%s", u.Scheme)
	}
}

//// object <-> value.Hashmap

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
