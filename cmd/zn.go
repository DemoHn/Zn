package cmd

import (
	"bufio"
	"fmt"
	"os"
)

const version = "rv1"

// ExecuteProgram - read file and execute
func ExecuteProgram() {
	// TODO
}

// EnterREPL - enter REPL to handle data
func EnterREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Zn> ")
	for scanner.Scan() {
		fmt.Println(scanner.Text())
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
