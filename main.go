package main

import (
	"fmt"
	"os"

	"github.com/DemoHn/Zn/cmd/zinc"
)

func main() {
	if err := zinc.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
