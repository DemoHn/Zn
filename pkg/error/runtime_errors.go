package error

import "fmt"

// IndexOutOfRange -
func IndexOutOfRange() *Error {
	return &Error{
		Code:    0x2401,
		Message: "索引超出此对象可用范围",
		Extra:   nil,
	}
}

// IndexKeyNotFound - used in hashmap
func IndexKeyNotFound(key string) *Error {
	return &Error{
		Code:    0x2402,
		Message: fmt.Sprintf("索引「%s」并不存在于此对象中", key),
		Extra:   key,
	}
}

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
		Code:    0x2505,
		Message: fmt.Sprintf("方法「%s」不存在", name),
		Extra:   name,
	}
}

// LeastParamsError -
func LeastParamsError(minParams int) *Error {
	return &Error{
		Code:    0x2701,
		Message: fmt.Sprintf("需要输入至少%d个参数", minParams),
		Extra:   minParams,
	}
}

// MismatchParamLengthError -
func MismatchParamLengthError(expect int, got int) *Error {
	return &Error{
		Code:    0x2702,
		Message: fmt.Sprintf("此方法定义了%d个参数，而实际输入%d个参数", expect, got),
		Extra:   []int{expect, got},
	}
}

// MostParamsError -
func MostParamsError(maxParams int) *Error {
	return &Error{
		Code:    0x2703,
		Message: fmt.Sprintf("至多需要%d个参数", maxParams),
		Extra:   maxParams,
	}
}

// ExactParamsError -
func ExactParamsError(exactParams int) *Error {
	return &Error{
		Code:    0x2704,
		Message: fmt.Sprintf("需要正好%d个参数", exactParams),
		Extra:   exactParams,
	}
}

// ModuleNotFound -
func ModuleNotFound(name string) *Error {
	return &Error{
		Code: 0x2801,
		Message: fmt.Sprintf("未找到「%s」模块", name),
		Extra: name,
	}
}

// Internal Error Class, for Zn Internal exception (rare to happen)
// e.g. Unexpected switch-case

// UnExpectedCase -
func UnExpectedCase(tag string, value string) *Error {
	return &Error{
		Code:    0x6001,
		Message: fmt.Sprintf("未定义的条件项：「%s」的值为「%s」", tag, value),
		Extra:   nil,
	}
}

//// SLOT
func NewErrorSLOT(info string) error {
	return fmt.Errorf(info)
}
