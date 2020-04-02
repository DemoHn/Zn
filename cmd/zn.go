package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/DemoHn/Zn/exec"
	"github.com/DemoHn/Zn/lex"
	"github.com/peterh/liner"
)

const version = "rv2"

// EnterREPL - enter REPL to handle data
func EnterREPL() {
	linerR := liner.NewLiner()
	linerR.SetCtrlCAborts(true)
	ctx := exec.NewContext()

	// REPL loop
	for {
		ctx.ResetLastValue()
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

		// execute program
		in := lex.NewTextStream(text)
		env := exec.NewProgramEnv()
		result := ctx.ExecuteCode(in, env)
		if !result.HasError {
			if result.Value != nil {
				prettyDisplayValue(result.Value, os.Stdout)
			}
		} else {
			fmt.Println(result.Error.Display())
		}
	}
}

// ExecProgram - exec program from file directly
func ExecProgram(file string) {
	ctx := exec.NewContext()
	env := exec.NewProgramEnv()
	in, errF := lex.NewFileStream(file)
	if errF != nil {
		fmt.Println(errF.Display())
		return
	}

	result := ctx.ExecuteCode(in, env)
	// when exec program, unlike REPL, it's not necessary to print last executed value
	if result.HasError {
		fmt.Println(result.Error.Display())
	}
}

// ShowVersion - show version
func ShowVersion() {
	fmt.Printf("Zn语言版本：%s\n", version)
}

//// display helpers
func prettyDisplayValue(val exec.ZnValue, w io.Writer) {
	var displayData = ""

	switch v := val.(type) {
	case *exec.ZnDecimal:
		// FG color: Cyan (lightblue)
		displayData = fmt.Sprintf("\x1b[38;5;147m%s\x1b[0m\n", v.String())
	case *exec.ZnString:
		// FG color: Green
		displayData = fmt.Sprintf("\x1b[38;5;184m%s\x1b[0m\n", v.String())
	case *exec.ZnBool:
		// FG color: White
		displayData = fmt.Sprintf("\x1b[38;5;231m%s\x1b[0m\n", v.String())
	case *exec.ZnNull, *exec.ZnFunction:
		displayData = fmt.Sprintf("‹\x1b[38;5;80m%s\x1b[0m›\n", v.String())
	default:
		displayData = fmt.Sprintf("%s\n", v.String())
	}

	w.Write([]byte(displayData))
}
