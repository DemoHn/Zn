package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/DemoHn/Zn/pkg/exec"
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

func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: impl LOGIC HERE!
	w.WriteHeader(200)
	io.WriteString(w, "Hello World")
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

func respondError(w http.ResponseWriter, reason error) {
	// declare content-type = text/plain (not json/file/...)
	w.Header().Add("Content-Type", "text/plain; charset=\"utf-8\"")
	// http status: 500 Internal Error
	w.WriteHeader(http.StatusInternalServerError)
	// write resp body
	io.WriteString(w, reason.Error())
}
