package main

import (
	"fmt"

	zinc "github.com/DemoHn/Zn"
)

const TARGET_FILE = "./draft/test.zn"

func main() {
	znt := zinc.NewInterpreter()
	res, err := znt.LoadFile(TARGET_FILE).Execute(map[string]zinc.Element{})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(res)
	}
}
