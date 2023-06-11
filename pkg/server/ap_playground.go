package server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

// if exec code ok - return 200 with result (ignore outputs from「显示」)
func RespondAsPlayground(w http.ResponseWriter, r *http.Request) {
	// exec program
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorData(w, fmt.Errorf("读取程序内容出现异常：%s", err.Error()))
		return
	}

	pExecutor := exec.NewPlaygroundExecutor(body)
	rtnValue, err := pExecutor.RunCode()
	if err != nil {
		writeErrorData(w, err)
		return
	}
	writeSuccessData(w, rtnValue)
}

func writeErrorData(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Header().Add("Content-Type", "text/plain; charset=\"utf-8\"")
	w.Write([]byte(err.Error()))
}

func writeSuccessData(w http.ResponseWriter, rtnValue runtime.Value) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain; charset=\"utf-8\"")

	// print return value
	switch rtnValue.(type) {
	case *value.Null:
		return
	default:
		w.Write([]byte(value.StringifyValue(rtnValue)))
	}
}
