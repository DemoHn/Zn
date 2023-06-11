package zinc

import (
	"fmt"
	eio "io"
	"os"
	"path/filepath"
	"strings"

	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/value"

	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/peterh/liner"
)

const version = "rev06"

// EnterREPL - enter REPL to handle data
func EnterREPL() {
	linerR := liner.NewLiner()
	linerR.SetCtrlCAborts(true)
	c := r.NewContext("", exec.GlobalValues)

	// init global module and scope
	module := r.NewModule("REPL", nil)
	c.EnterModule(module)

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

		// #1. read source code
		source, err := in.ReadAll()
		if err != nil {
			prettyPrintError(c, os.Stdout, err)
			return
		}

		lexer := syntax.NewLexer(source)
		c.GetCurrentModule().SetLexer(lexer)

		// execute code
		result, err2 := exec.ExecuteREPLCode(c, lexer)
		if err2 != nil {
			prettyPrintError(c, os.Stdout, err2)
			return
		}

		if result != nil {
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

	c := r.NewContext(rootDir, exec.GlobalValues)

	// when exec program, unlike REPL, it's not necessary to print last executed value
	rtnValue, err := exec.ExecuteModule(c, rootModule)
	if err != nil {
		prettyPrintError(c, os.Stdout, err)
		return
	}

	// print return value
	switch rtnValue.(type) {
	case *value.Null:
		return
	default:
		if c.GetHasPrinted() {
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

func prettyPrintError(c *r.Context, w eio.Writer, err error) {
	_, _ = w.Write([]byte(exec.DisplayError(err)))
}
