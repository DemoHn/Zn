package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

type playgroundReq struct {
	VarInput   string
	SourceCode string
}

func ReadRequestForPlayground(r *http.Request) ([]rune, map[string]runtime.Element, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("读取请求内容出现异常：%s", err.Error())
	}

	var reqInfo playgroundReq
	if err := json.Unmarshal(body, &reqInfo); err != nil {
		return nil, nil, fmt.Errorf("解析请求格式不符合要求！")
	}

	if reqInfo.VarInput != "" {
		varInputs, err := exec.ExecVarInputText(reqInfo.VarInput)
		if err != nil {
			return nil, nil, err
		}
		return []rune(reqInfo.SourceCode), varInputs, nil
	}
	return []rune(reqInfo.SourceCode), map[string]runtime.Element{}, nil
}

func WriteResponseForPlayground(w http.ResponseWriter, rtnValue runtime.Element, err error) {
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
