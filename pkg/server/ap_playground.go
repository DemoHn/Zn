package server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/DemoHn/Zn/pkg/value"
)

// if exec code ok - return 200 with result (ignore outputs from「显示」)
func RespondAsPlayground(r *http.Request, w *RespWriter) error {
	// exec program
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return w.WriteErrorData(fmt.Errorf("读取程序内容出现异常：%s", err.Error()))
	}

	pExecutor := exec.NewPlaygroundExecutor(body)
	rtnValue, err := pExecutor.RunCode()
	if err != nil {
		return w.WriteErrorData(err)
	}

	// write return value as resp body
	switch rtnValue.(type) {
	case *value.Null:
		return w.WriteSuccessData("")
	default:
		return w.WriteSuccessData(value.StringifyValue(rtnValue))
	}
}
