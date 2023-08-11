package zinc

import (
	"fmt"
	eio "io"
	"os"
	"path/filepath"
	"strings"

	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"

	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/peterh/liner"
)

const version = "rev07"

// EnterREPL - enter REPL to handle data
func EnterREPL() {
	linerR := liner.NewLiner()
	linerR.SetCtrlCAborts(true)

	replExecutor := exec.NewREPLExecutor()
	// REPL loop
	for {
		text, err := linerR.Prompt("Zn> ")
		if err != nil {
			if err == liner.ErrPromptAborted {
				break
			} else if err.Error() == "EOF" {
				break
			} else {
				fmt.Printf("未知错误：%s\n", err.Error())
				break
			}
		}
		// append history
		linerR.AppendHistory(text)
		// add special command
		if text == "退出" {
			break
		}

		result, err := replExecutor.RunCode(text)
		if err != nil {
			prettyPrintError(os.Stdout, err)
		} else if result != nil {
			prettyDisplayValue(result, os.Stdout)
		}
	}
}

// ExecProgram - exec program from file directly
func ExecProgram(file string) {
	rootDir := filepath.Dir(file)
	// get module name
	_, fileName := filepath.Split(file)
	rootModule := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	fileExecutor := exec.NewFileExecutor(rootDir, rootModule)
	// when exec program, unlike REPL, it's not necessary to print last executed value
	rtnValue, err := fileExecutor.RunModule()
	if err != nil {
		prettyPrintError(os.Stdout, err)
		return
	}

	// print return value
	switch rtnValue.(type) {
	case *value.Null:
		return
	default:
		if fileExecutor.HasPrinted() {
			os.Stdout.Write([]byte{'\n'})
		}
		os.Stdout.Write([]byte(value.StringifyValue(rtnValue)))
		os.Stdout.Write([]byte{'\n'})
	}
}

// ShowVersion - show version
func ShowVersion() {
	fmt.Printf("Zn语言版本：%s\n", version)
}

//// display helpers
func prettyDisplayValue(v r.Element, w eio.Writer) {
	var displayData = ""
	var valStr = value.StringifyValue(v)
	switch v.(type) {
	case *value.Number:
		// FG color: Cyan (lightblue)
		displayData = fmt.Sprintf("\x1b[38;5;147m%s\x1b[0m\n", valStr)
	case *value.String:
		// FG color: Green
		displayData = fmt.Sprintf("\x1b[38;5;184m%s\x1b[0m\n", valStr)
	case *value.Bool:
		// FG color: White
		displayData = fmt.Sprintf("\x1b[38;5;231m%s\x1b[0m\n", valStr)
	case *value.Null, *value.Function:
		displayData = fmt.Sprintf("‹\x1b[38;5;80m%s\x1b[0m›\n", valStr)
	default:
		displayData = fmt.Sprintf("%s\n", valStr)
	}

	_, _ = w.Write([]byte(displayData))
}

func prettyPrintError(w eio.Writer, err error) {
	_, _ = w.Write([]byte(exec.DisplayError(err)))
}
