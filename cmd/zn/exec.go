package zn

import (
	"fmt"
	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/value"
	eio "io"
	"os"
	"path/filepath"
	"strings"

	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/peterh/liner"
)

const version = "rev06"

// EnterREPL - enter REPL to handle data
func EnterREPL() {
	linerR := liner.NewLiner()
	linerR.SetCtrlCAborts(true)
	c := r.NewContext(exec.GlobalValues)

	// REPL loop
	for {
		text, err := linerR.Prompt("Zn> ")
		if err != nil {
			if err == liner.ErrPromptAborted {
				os.Exit(0)
			} else if err.Error() == "EOF" {
				os.Exit(0)
			} else {
				fmt.Printf("未知错误：%s\n", err.Error())
				os.Exit(0)
			}
		}
		// append history
		linerR.AppendHistory(text)
		// add special command
		if text == ".exit" {
			break
		}

		// execute program
		in := io.NewByteStream([]byte(text))
		module := r.NewModule("REPL")

		// #1. read source code
		source, err := in.ReadAll()
		if err != nil {
			prettyPrintError(c, err)
			return
		}

		lexer := syntax.NewLexer(source)
		// global execution scope
		// TODO: fix global scope
		c.PushScope(module, lexer)

		// execute code
		result, err2 := exec.ExecuteREPLCode(c, lexer)
		if err2 != nil {
			prettyPrintError(c, err2)
			return
		}

		if result != nil {
			prettyDisplayValue(result, os.Stdout)
		}
	}
}

// ExecProgram - exec program from file directly
func ExecProgram(file string) {
	c := r.NewContext(exec.GlobalValues)

	rootDir := filepath.Dir(file)
	// get module name
	_, fileName := filepath.Split(file)
	rootModule := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	c.SetRootDir(rootDir)
	// when exec program, unlike REPL, it's not necessary to print last executed value
	if _, err := exec.ExecuteModule(c, rootModule); err != nil {
		prettyPrintError(c, err)
	}
}

// ShowVersion - show version
func ShowVersion() {
	fmt.Printf("Zn语言版本：%s\n", version)
}

//// display helpers
func prettyDisplayValue(v r.Value, w eio.Writer) {
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

func prettyPrintError(c *r.Context, err error) {
	_, _ = os.Stdout.Write([]byte(exec.DisplayError(err)))
}
