package cmd

import (
	"fmt"
	"os"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec"
	"github.com/DemoHn/Zn/lex"
	"github.com/DemoHn/Zn/syntax"
	"github.com/peterh/liner"
)

const version = "rv1"

// ExecuteProgram - read file and execute
func execProgram(stream *lex.InputStream, inpt *exec.Interpreter) (string, *error.Error) {
	var nInpt *exec.Interpreter = inpt
	if inpt == nil {
		nInpt = exec.NewInterpreter()
	}

	p := syntax.NewParser(lex.NewLexer(stream))
	programNode, err := p.Parse()
	if err != nil {
		return "", err
	}
	fmt.Printf("node = \x1b[33m%s\x1b[0m\n", syntax.StringifyAST(programNode))

	// return with green color
	return fmt.Sprintf("\x1b[32m%s\x1b[0m\n", nInpt.Execute(programNode)), nil
}

// EnterREPL - enter REPL to handle data
func EnterREPL() {
	linerR := liner.NewLiner()
	linerR.SetCtrlCAborts(true)

	inpt := exec.NewInterpreter()

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

		rtn, errE := execProgram(lex.NewTextStream(text), inpt)
		if errE != nil {
			fmt.Printf("%s\n", errE.Display())
			continue
		}

		fmt.Println(rtn)
	}
}

// ExecProgram -
func ExecProgram(file string) {
	s := lex.Source{}
	in, errF := lex.NewFileStream(file)
	if errF != nil {
		fmt.Println(errF.Display())
		return
	}
	s.AddStream(in)

	rtn, errE := execProgram(s.Streams[0], nil)
	if errE != nil {
		fmt.Println(errE.Display())
		return
	}

	fmt.Println(rtn)
}

// ShowVersion - show version
func ShowVersion() {
	fmt.Printf("Zn语言版本：%s\n", version)
}
