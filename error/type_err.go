package error

import "fmt"

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
func InvalidExprType(assertType string) *Error {
	label := assertType
	if v, ok := typeNameMap[assertType]; ok {
		label = v
	}
	return typeError.NewError(0x01, Error{
		text: fmt.Sprintf("表达式不符合期望之「%s」类型", label),
	})
}
