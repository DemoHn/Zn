package cmd

import (
	"fmt"
	"os"

	"github.com/DemoHn/Zn/exec"
	"github.com/DemoHn/Zn/lex"
	"github.com/DemoHn/Zn/syntax"
	"github.com/peterh/liner"
)

const version = "rv1"

// ExecuteProgram - read file and execute
func execProgram(text string, inpt *exec.Interpreter) (string, error) {
	var nInpt *exec.Interpreter = inpt
	if inpt == nil {
		nInpt = exec.NewInterpreter()
	}

	data := []rune(text)
	p := syntax.NewParser(lex.NewLexer(data))
	programNode, err := p.Parse()
	if err != nil {
		return "", err
	}

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

		rtn, errE := execProgram(text, inpt)
		if errE != nil {
			fmt.Printf("[语法错误] %s\n", errE.Error())
			continue
		}

		fmt.Println(rtn)
	}
}

// ShowVersion - show version
func ShowVersion() {
	fmt.Printf("Zn语言版本：%s\n", version)
}
