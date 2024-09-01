package main

import (
	"fmt"

	zinc "github.com/DemoHn/Zn"
	"github.com/DemoHn/Zn/pkg/server"
	"github.com/spf13/cobra"
)

const defaultConnUrl = "tcp://127.0.0.1:3862"

var (
	connUrl         string
	childWorkerFlag bool
	maxProcs        int
	initProcs       int
	timeout         int

	rootCmd = &cobra.Command{
		Use:   "zinc-playground",
		Short: "zinc playground",
		Long:  "zinc playground - 在启动服务器之后，用户发送HTTP请求并提交代码后即可执行，并返回对应的结果；这样用户可以线上编写并运行代码",
		Run: func(c *cobra.Command, args []string) {
			err := zinc.NewServer().
				SetPlaygroundHandler().
				SetPMServerConfig(server.ZnPMServerConfig{
					InitProcs: initProcs,
					MaxProcs:  maxProcs,
					Timeout:   timeout,
				}).
				Launch(connUrl)

			if err != nil {
				fmt.Println("启动服务器时发生异常：%s\n", err)
				return
			}
		},
	}
)

func main() {
	rootCmd.Flags().StringVarP(&connUrl, "listen", "l", defaultConnUrl, "设置服务器监听的URL 如 tcp://127.0.0.1:3862 或 unix:///tmp/zinc.sock")
	rootCmd.Flags().BoolVar(&childWorkerFlag, "child-worker", false, "[仅限内部使用] 启动用于处理子请求的Worker进程")

	rootCmd.Flags().IntVar(&maxProcs, "max-procs", 100, "限制最大可创建进程数量")
	rootCmd.Flags().IntVar(&initProcs, "init-procs", 20, "初始创建子进程数量")

	rootCmd.Flags().IntVar(&timeout, "timeout", 60, "执行超时时间，单位为秒")
	rootCmd.Execute()
}
