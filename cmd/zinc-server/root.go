package main

import (
	"fmt"
	"os"

	"github.com/DemoHn/Zn/pkg/server"
	"github.com/spf13/cobra"
)

const defaultConnUrl = "tcp://127.0.0.1:3862"

var (
	connUrl         string
	childWorkerFlag bool
	maxProcs        int
	initProcs       int

	rootCmd = &cobra.Command{
		Use:   "zinc-server",
		Short: "Zn FastCGI 服务器",
		Long:  "Zn FastCGI 服务器 - 监听上游传过来的请求，执行指定代码并输出结果；同时也可以当作 HTTP 处理器搭配 nginx 等使用，如 PHP-FPM 一般",
		Run: func(c *cobra.Command, args []string) {
			///// run child worker if  --child-worker = true & preForkChild env is "OK"
			if childWorkerFlag && os.Getenv(server.EnvPreforkChildKey) == server.EnvPreforkChildVal {
				// start child worker to handle requests
				if err := server.StartWorker(); err != nil {
					fmt.Printf("启动子进程时发生错误：%v\n", err)
					return
				}
			} else {
				//// otherwise, just listen to the server
				zns, err := server.NewFromURL(connUrl)
				if err != nil {
					fmt.Printf("启动服务器时发生错误：%v\n", err)
					return
				}

				cfg := server.ZnServerConfig{
					InitProcs: initProcs,
					MaxProcs:  maxProcs,
				}
				if err := zns.StartMaster(cfg); err != nil {
					fmt.Printf("启动主进程时发生错误：%v\n", err)
					return
				}
			}
		},
	}
)

func main() {
	rootCmd.Flags().StringVarP(&connUrl, "listen", "l", defaultConnUrl, "设置服务器监听的URL 如 tcp://127.0.0.1:3862 或 unix:///tmp/zinc.sock")
	rootCmd.Flags().BoolVar(&childWorkerFlag, "child-worker", false, "[仅限内部使用] 启动用于处理子请求的Worker进程")

	rootCmd.Flags().IntVar(&maxProcs, "max-procs", 100, "限制最大可创建进程数量")
	rootCmd.Flags().IntVar(&initProcs, "init-procs", 20, "初始创建子进程数量")
	rootCmd.Execute()
}
