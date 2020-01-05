package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DemoHn/Zn/exec"
	"github.com/DemoHn/Zn/lex"
	"github.com/DemoHn/Zn/syntax"
)

const version = "rv1"

// ExecuteProgram - read file and execute
func ExecuteProgram() {
	// TODO
}

// EnterREPL - enter REPL to handle data
func EnterREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	inpt := exec.NewInterpreter()

	fmt.Printf("Zn> ")
	for scanner.Scan() {
		// exect
		data := []rune(scanner.Text())
		p := syntax.NewParser(lex.NewLexer(data))
		programNode, err := p.Parse()
		if err != nil {
			fmt.Printf("[SyntaxError] %s\n", err.Error())

			fmt.Printf("Zn> ")
			continue
		}

		fmt.Printf("%s\n", inpt.Execute(programNode))
		fmt.Printf("Zn> ")
	}

	if scanner.Err() != nil {
		// handle error.
	}
}

// ShowVersion - show version
func ShowVersion() {
	fmt.Printf("Zn语言版本：%s\n", version)
}
