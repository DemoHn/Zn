package main

import "fmt"

type Number struct {
	value float64
}

type String struct {
	value string
}

func getTypeString[T any]() string {
	var t T
	switch any(t).(type) {
	case *Number:
		return "number"
	case *String:
		return "string"
	}
	return "unknown"
}

func main() {
	fmt.Println(getTypeString[*String]())
}
