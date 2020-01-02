package main

import (
	"flag"
	"os"

	"github.com/DemoHn/Zn/cmd"
)

// version
var (
	versionFlag = false
)

func main() {
	flag.Parse()

	// args := flag.Args()
	// show flags
	if versionFlag {
		cmd.ShowVersion()
		os.Exit(0)
	}

	cmd.EnterREPL()
	os.Exit(0)
}

func init() {
	flag.BoolVar(&versionFlag, "v", false, "显示Zn语言当前版本")
}
