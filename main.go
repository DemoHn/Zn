package main

import (
	"fmt"
	"os"

	"github.com/DemoHn/Zn/cmd"
	"github.com/spf13/cobra"
)

var (
	versionFlag bool
	rootCmd     = &cobra.Command{
		Use:   "Zn",
		Short: "Zn语言解释器",
		Long:  "Zn语言解释器",
		Run: func(c *cobra.Command, args []string) {
			// -v, --version
			if versionFlag {
				cmd.ShowVersion()
				return
			}
			// if len(args) > 0, execute file
			if len(args) > 0 {
				filename := args[0]
				cmd.ExecProgram(filename)
				return
			}
			// by default, enter REPL
			cmd.EnterREPL()
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "显示Zn语言版本")
}
