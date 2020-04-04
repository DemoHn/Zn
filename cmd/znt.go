package cmd

import (
	"fmt"
	"os"
)

// znt - Zn tools

func main() {
	args := os.Args

	if len(args) > 0 {
		directive := args[0]
		switch directive {
		case "md:pretty":
			rtn := mdPretty(args[1])
			os.Exit(rtn)
		}
	} else {
		fmt.Println("需要输入合适的指令。")
		fmt.Println("可能合适的指令有：")
		fmt.Println("    md:pretty - 语法高亮 markdown 文件中含有 Zn 语言的代码片段，使之更易读")

		os.Exit(1)
	}
}

func mdPretty(file string) int {
	return 0
}
