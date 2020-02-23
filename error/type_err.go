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
		text: fmt.Sprintf("表达式不符合期望之%s类型", strings.Join(labels, "，")),
	})
}
