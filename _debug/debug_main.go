package main

import (
	"fmt"

	zinc "github.com/DemoHn/Zn"
)

const TARGET_FILE = "./doc/zh-cn/snippets/example/冒泡排序.zn"
const VAR_INPUT = `数组文本=“123,456，-2,3，DD”`

func main() {
	znt := zinc.NewInterpreter()

	varInput, err := znt.ExecuteVarInputText(VAR_INPUT)
	if err != nil {
		panic(err)
	}

	res, err := znt.LoadFile(TARGET_FILE).Execute(varInput)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(res)
	}
}
