package error

import "fmt"

// NameNotDefined -
func NameNotDefined(name string) *Error {
	return &Error{
		Code:    0x2501,
		Message: fmt.Sprintf("标识「%s」未有定义", name),
		Extra:   name,
	}
}

// NameRedeclared -
func NameRedeclared(name string) *Error {
	return &Error{
		Code:    0x2502,
		Message: fmt.Sprintf("标识「%s」被重复定义", name),
		Extra:   name,
	}
}

// AssignToConstant -
func AssignToConstant() *Error {
	return &Error{
		Code:    0x2503,
		Message: "不允许赋值给常量",
		Extra:   nil,
	}
}

// PropertyNotFound -
func PropertyNotFound(name string) *Error {
	return &Error{
		Code:    0x2504,
		Message: fmt.Sprintf("属性「%s」不存在", name),
		Extra:   name,
	}
}

// MethodNotFound -
func MethodNotFound(name string) *Error {
	return &Error{
		Code:    0x2504,
		Message: fmt.Sprintf("方法「%s」不存在", name),
		Extra:   name,
	}
}
