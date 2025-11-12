package main

import (
	"strings"

	"github.com/spf13/cobra"
)

var (
	versionFlag  bool
	varInputFlag []string
	rootCmd      = &cobra.Command{
		Use:   "Zn",
		Short: "Zn语言解释器",
		Long:  "Zn语言解释器",
		Run: func(c *cobra.Command, args []string) {
			// -v, --version
			if versionFlag {
				ShowVersion()
				return
			}
			// if len(args) > 0, execute file
			if len(args) > 0 {
				filename := args[0]
				varInputBlock := strings.Join(varInputFlag, "\n")
				ExecProgram(filename, varInputBlock)
				return
			}
			// by default, enter REPL
			EnterREPL()
		},
	}
)

func main() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "显示Zn语言版本")
	rootCmd.Flags().StringArrayVarP(&varInputFlag, "input", "i", []string{}, "定义输入变量(支持多个变量)，格式为 <变量名>=<表达式>，如：‘./zinc xx.zn -i 客单价=28.25 -i 销量=300’")
	rootCmd.Execute()
}
