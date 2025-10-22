package main

import (
	"fmt"
	eio "io"
	"os"

	zinc "github.com/DemoHn/Zn"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"

	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/peterh/liner"
)

// EnterREPL - enter REPL to handle data
func EnterREPL() {
	linerR := liner.NewLiner()
	linerR.SetCtrlCAborts(true)

	interpreter := zinc.NewInterpreter()
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

		varInput := map[string]r.Element{}
		result, err := interpreter.LoadScript([]rune(text)).Execute(varInput)
		if err != nil {
			prettyPrintError(os.Stdout, err)
		} else if result != nil {
			prettyDisplayValue(result, os.Stdout)
		}
	}
}

// ExecProgram - exec program from file directly
func ExecProgram(file string) {
	znInterpreter := zinc.NewInterpreter()
	rtnValue, err := znInterpreter.LoadFile(file).Execute(map[string]r.Element{})

	if err != nil {
		prettyPrintError(os.Stdout, err)
		return
	}

	// print return value
	switch rtnValue.(type) {
	case *value.Null:
		return
	default:
		os.Stdout.Write([]byte(value.StringifyValue(rtnValue)))
		os.Stdout.Write([]byte{'\n'})
	}
}

// ShowVersion - show version

func ShowVersion() {
	znInterpreter := zinc.NewInterpreter()
	fmt.Printf("Zn语言版本：%s\n", znInterpreter.GetVersion())
}

// // display helpers
func prettyDisplayValue(v r.Element, w eio.Writer) {
	var displayData = ""
	var valStr = value.StringifyValue(v)
	switch v.(type) {
	case *value.Number:
		// FG color: Cyan (lightblue)
		displayData = fmt.Sprintf("\x1b[38;5;147m%s\x1b[0m\n", valStr)
	case *value.String:
		// FG color: Green
		// NOTE: string value quoted with 「 & 」
		displayData = fmt.Sprintf("\x1b[38;5;184m「%s」\x1b[0m\n", valStr)
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
