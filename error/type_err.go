package error

import (
	"fmt"
	"strings"
)

var typeNameMap = map[string]string{
	"string":   "文本",
	"decimal":  "数值",
	"integer":  "整数",
	"function": "方法",
	"bool":     "二象",
	"null":     "空",
	"array":    "元组",
	"hashmap":  "列表",
	"id":       "标识",
}

// InvalidExprType -
func InvalidExprType(assertType ...string) *Error {
	labels := []string{}
	for _, at := range assertType {
		label := at
		if v, ok := typeNameMap[at]; ok {
			label = v
		}
		labels = append(labels, fmt.Sprintf("「%s」", label))
	}
	return typeError.NewError(0x01, Error{
		text: fmt.Sprintf("表达式不符合期望之%s类型", strings.Join(labels, "、")),
	})
}

// InvalidFuncVariable -
func InvalidFuncVariable(tag string) *Error {
	return typeError.NewError(0x02, Error{
		text: fmt.Sprintf("「%s」须为一个方法", tag),
		info: fmt.Sprintf("tag=(%s)", tag),
	})
}

// InvalidCaseType - general error for default branch of switch
// Theoratically, it should NOT be triggered at all!
func InvalidCaseType() *Error {
	return typeError.NewError(0x03, Error{
		text: "不符合期望之类型",
	})
}
