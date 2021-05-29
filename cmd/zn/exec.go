package zn

import (
	"fmt"
	"io"
	"os"

	"github.com/DemoHn/Zn/exec"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/exec/val"
	"github.com/DemoHn/Zn/lex"
	"github.com/peterh/liner"
)

const version = "rev04"

// EnterREPL - enter REPL to handle data
func EnterREPL() {
	linerR := liner.NewLiner()
	linerR.SetCtrlCAborts(true)
	c := ctx.NewContext(nil) // TODO:

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
		if text == ".print" {
			printSymbols(c)
			continue
		} else if text == ".exit" {
			break
		}

		// execute program
		in := lex.NewTextStream(text)
		result, err2 := exec.ExecuteCode(c, in)
		if err2 != nil {
			fmt.Println(err2.Display())
		} else {
			if result != nil {
				prettyDisplayValue(result, os.Stdout)
			}
		}
	}
}

// ExecProgram - exec program from file directly
func ExecProgram(file string) {
	c := ctx.NewContext(nil)
	in, errF := lex.NewFileStream(file)
	if errF != nil {
		fmt.Println(errF.Display())
		return
	}

	_, err := exec.ExecuteCode(c, in)
	// when exec program, unlike REPL, it's not necessary to print last executed value
	if err != nil {
		fmt.Println(err.Display())
	}
}

// ShowVersion - show version
func ShowVersion() {
	fmt.Printf("Zn语言版本：%s\n", version)
}

//// display helpers
func prettyDisplayValue(v ctx.Value, w io.Writer) {
	var displayData = ""
	var valStr = val.StringifyValue(v)
	switch v.(type) {
	case *val.Decimal:
		// FG color: Cyan (lightblue)
		displayData = fmt.Sprintf("\x1b[38;5;147m%s\x1b[0m\n", valStr)
	case *val.String:
		// FG color: Green
		displayData = fmt.Sprintf("\x1b[38;5;184m%s\x1b[0m\n", valStr)
	case *val.Bool:
		// FG color: White
		displayData = fmt.Sprintf("\x1b[38;5;231m%s\x1b[0m\n", valStr)
	case *val.Null, *val.Function:
		displayData = fmt.Sprintf("‹\x1b[38;5;80m%s\x1b[0m›\n", valStr)
	default:
		displayData = fmt.Sprintf("%s\n", valStr)
	}

	w.Write([]byte(displayData))
}

// printSymbols -
func printSymbols(c *ctx.Context) {
	/** TODO
	strs := []string{}
	for k, symArr := range ctx.GetSymbols() {
		if symArr != nil {
			for _, symItem := range symArr {
				symStr := "ε"
				if symItem.Value != nil {
					symStr = symItem.Value.String()
				}
				strs = append(strs, fmt.Sprintf("‹%s, %d› => %s", k, symItem.NestLevel, symStr))
			}
		}
	}

	data := strings.Join(strs, "\n")
	fmt.Println(data)
	*/
}
