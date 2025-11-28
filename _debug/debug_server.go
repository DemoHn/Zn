package main

import (
	zinc "github.com/DemoHn/Zn"
)

const TARGET_FILE2 = "./draft/http.zn"

func main() {
	znt := zinc.NewInterpreter()

	znt.SetMainServer(zinc.NewThreadServer(), zinc.NewHttpHandler(znt, TARGET_FILE2))

	znt.Listen("tcp://127.0.0.1:3862")
}
