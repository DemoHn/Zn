package main

import (
	"fmt"

	"github.com/DemoHn/Zn/pkg/server"
	"github.com/spf13/cobra"
)

const defaultConnUrl = "tcp://127.0.0.1:3862"

var (
	connUrl string
	rootCmd = &cobra.Command{
		Use:   "zinc-server",
		Short: "Zn FastCGI 服务器",
		Long:  "Zn FastCGI 服务器 - 监听上游传过来的请求，执行指定代码并输出结果；同时也可以当作 HTTP 处理器搭配 nginx 等使用，如 PHP-FPM 一般",
		Run: func(c *cobra.Command, args []string) {
			zns, err := server.NewFromURL(connUrl)
			if err != nil {
				fmt.Printf("启动服务器时发生错误：%v\n", err)
				return
			}

			zns.Listen()
		},
	}
)

func main() {
	rootCmd.Flags().StringVarP(&connUrl, "listen", "l", defaultConnUrl, "设置服务器监听的URL 如 tcp://127.0.0.1:3862 或 unix:///tmp/zinc.sock")
	rootCmd.Execute()
}
