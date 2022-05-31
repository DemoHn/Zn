package error

import (
	"fmt"
	"strings"
)

var typeNameMap = map[string]string{
	"string":   "文本",
	"number":   "数值",
	"integer":  "整数",
	"function": "方法",
	"bool":     "逻辑",
	"null":     "空",
	"array":    "元组",
	"hashmap":  "列表",
	"id":       "标识",
}

// InvalidExprType -
func InvalidExprType(assertType ...string) *Error {
	var labels []string
	for _, at := range assertType {
		label := at
		if v, ok := typeNameMap[at]; ok {
			label = v
		}
		labels = append(labels, fmt.Sprintf("「%s」", label))
	}

	return &Error{
		Code:    0x2301,
		Message: fmt.Sprintf("表达式不符合期望的%s类型", strings.Join(labels, "、")),
		Extra:   labels,
	}
}

// InvalidFuncVariable -
func InvalidFuncVariable(tag string) *Error {
	return &Error{
		Code:    0x2302,
		Message: fmt.Sprintf("「%s」须为一个方法", tag),
		Extra:   tag,
	}
}

// InvalidParamType -
func InvalidParamType(assertType ...string) *Error {
	labels := []string{}
	for _, at := range assertType {
		label := at
		if v, ok := typeNameMap[at]; ok {
			label = v
		}
		labels = append(labels, fmt.Sprintf("「%s」", label))
	}
	return &Error{
		Code:    0x2303,
		Message: fmt.Sprintf("输入参数不符合期望之%s类型", strings.Join(labels, "、")),
		Extra:   labels,
	}
}

// InvalidCompareLType - 比较的值的类型
func InvalidCompareLType(assertType ...string) *Error {
	labels := []string{}
	for _, at := range assertType {
		label := at
		if v, ok := typeNameMap[at]; ok {
			label = v
		}
		labels = append(labels, fmt.Sprintf("「%s」", label))
	}
	return &Error{
		Code:    0x2304,
		Message: fmt.Sprintf("比较值的类型应为%s", strings.Join(labels, "、")),
		Extra:   labels,
	}
}

// InvalidCompareRType - 被比较的值的类型
func InvalidCompareRType(assertType ...string) *Error {
	labels := []string{}
	for _, at := range assertType {
		label := at
		if v, ok := typeNameMap[at]; ok {
			label = v
		}
		labels = append(labels, fmt.Sprintf("「%s」", label))
	}
	return &Error{
		Code:    0x2305,
		Message: fmt.Sprintf("被比较值的类型应为%s", strings.Join(labels, "、")),
		Extra:   labels,
	}
}

