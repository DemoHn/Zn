package server

import (
	"fmt"
	"net/http"
	"net/http/fcgi"
	"path/filepath"
	"strings"

	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

const EnvScriptFilename = "SCRIPT_FILENAME"

// if exec code ok - return 200 with result (ignore outputs from「显示」)
func RespondAsExecutor(w http.ResponseWriter, r *http.Request) {
	// exec program
	c := runtime.NewContext(exec.GlobalValues)
	// get filename
	cgiEnvs := fcgi.ProcessEnv(r)
	if scriptFile, ok := cgiEnvs[EnvScriptFilename]; ok {
		rootDir, fileName := filepath.Split(scriptFile)
		rootModule := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		c.SetRootDir(rootDir)
		rtnValue, err := exec.ExecuteModule(c, rootModule)
		if err != nil {
			writeErrorData(w, err)
			return
		}
		writeSuccessData(w, rtnValue)
		return
	}
	// else: no SCRIPT_FILENAME found
	writeErrorData(w, fmt.Errorf("待执行的文件参数未指定"))
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
