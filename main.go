package main

import (
	"fmt"
	"os"

	"github.com/DemoHn/Zn/cmd/zn"
)

func main() {
	if err := zn.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
