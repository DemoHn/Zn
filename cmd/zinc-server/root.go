package main

import (
	"fmt"

	zinc "github.com/DemoHn/Zn"
	"github.com/spf13/cobra"
)

const defaultConnUrl = "tcp://127.0.0.1:3862"

var (
	connUrl   string
	entryFile string

	rootCmd = &cobra.Command{
		Use:   "zinc-server",
		Short: "Zn HTTP服务器",
		Long:  "Zn HTTP服务器 - 处理上游传过来的HTTP请求，并返回相应的结果 - 请注意和 zinc-playground 不同，这里每接受一次请求时都会创建一个新的 goroutine，所以请小心别把服务器搞挂了！",
		Run: func(c *cobra.Command, args []string) {
			znServer := zinc.NewServer()

			// listen and handle
			if err := znServer.SetHTTPHandler(entryFile).Launch(connUrl); err != nil {
				fmt.Printf("启动服务器时发生错误：%v\n", err)
				return
			}
		},
	}
)

func main() {
	rootCmd.Flags().StringVarP(&connUrl, "listen", "l", defaultConnUrl, "设置服务器监听的URL 如 tcp://127.0.0.1:3862 或 unix:///tmp/zinc.sock")
	rootCmd.Flags().StringVarP(&entryFile, "file", "f", "", "执行入口文件")
	rootCmd.Execute()
}
