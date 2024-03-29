package main

import (
	"fmt"
	"os"

	"github.com/DemoHn/Zn/cmd/zinc-devtool/cmds"
	"github.com/spf13/cobra"
)

// ToolCommand - tool command
var rootCommand = &cobra.Command{
	Use:   "znt",
	Short: "znt - Zn Tools 辅助开发工具",
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCommand.AddCommand(cmds.GenCodeImageCmd)
	rootCommand.AddCommand(cmds.GenKeywordCmd)
}
