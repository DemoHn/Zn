package error

import (
	"fmt"
	"strings"
)

type RuntimeError struct {
	Code    int
	Message string
	Extra   interface{}
}

func (r *RuntimeError) Error() string {
	return r.Message
}

const (
	ErrIndexOutOfRange          = 30
	ErrIndexKeyNotFound         = 31
	ErrNameNotDefined           = 32
	ErrNameRedeclared           = 33
	ErrAssignToConstant         = 34
	ErrPropertyNotFound         = 35
	ErrMethodNotFound           = 36
	ErrClassNotOnRoot           = 37
	ErrThisValueNotFound        = 38
	ErrLeastParamsError         = 40
	ErrMismatchParamLengthError = 41
	ErrMostParamsError          = 42
	ErrExactParamsError         = 43
	// module error
	ErrModuleNotFound           = 50
	ErrImportSameModule         = 51
	ErrDuplicateModule          = 52
	ErrModuleCircularDependency = 53
	// internal error
	ErrUnexpectedCase           = 60
	ErrUnexpectedEmptyExecLogic = 61
	ErrUnexpectedAssign         = 62
	ErrUnexpectedParamWildcard  = 63
	// type error
	ErrInvalidExprType            = 70
	ErrInvalidFuncVariable        = 71
	ErrInvalidParamType           = 72
	ErrInvalidCompareLType        = 73
	ErrInvalidCompareRType        = 74
	ErrInvalidExceptionType       = 75
	ErrInvalidExceptionObjectType = 76
	// arith error
	ErrArithDivZero          = 80
	ErrArithRootLessThanZero = 81
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

// IndexOutOfRange -
func IndexOutOfRange() *RuntimeError {
	return &RuntimeError{
		Code:    ErrIndexOutOfRange,
		Message: "索引超出此对象可用范围",
		Extra:   nil,
	}
}

// IndexKeyNotFound - used in hashmap
func IndexKeyNotFound(key string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrIndexKeyNotFound,
		Message: fmt.Sprintf("索引「%s」并不存在于此对象中", key),
		Extra:   key,
	}
}

// NameNotDefined -
func NameNotDefined(name string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrNameNotDefined,
		Message: fmt.Sprintf("标识「%s」未有定义", name),
		Extra:   name,
	}
}

// NameRedeclared -
func NameRedeclared(name string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrNameRedeclared,
		Message: fmt.Sprintf("标识「%s」被重复定义", name),
		Extra:   name,
	}
}

// AssignToConstant -
func AssignToConstant() *RuntimeError {
	return &RuntimeError{
		Code:    ErrAssignToConstant,
		Message: "不允许赋值给常量",
		Extra:   nil,
	}
}

// PropertyNotFound -
func PropertyNotFound(name string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrPropertyNotFound,
		Message: fmt.Sprintf("属性「%s」不存在", name),
		Extra:   name,
	}
}

// MethodNotFound -
func MethodNotFound(name string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrMethodNotFound,
		Message: fmt.Sprintf("方法「%s」不存在", name),
		Extra:   name,
	}
}

// ClassNotOnRoot -
func ClassNotOnRoot(name string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrClassNotOnRoot,
		Message: fmt.Sprintf("只能在模块主层级定义「%s」类", name),
		Extra:   name,
	}
}

// ThisValueNotFound -
func ThisValueNotFound() *RuntimeError {
	return &RuntimeError{
		Code:    ErrThisValueNotFound,
		Message: "未找到此方法/属性对应的主对象 (thisValue)",
		Extra:   nil,
	}
}

// LeastParamsError -
func LeastParamsError(minParams int) *RuntimeError {
	return &RuntimeError{
		Code:    ErrLeastParamsError,
		Message: fmt.Sprintf("需要输入至少%d个参数", minParams),
		Extra:   minParams,
	}
}

// MismatchParamLengthError -
func MismatchParamLengthError(expect int, got int) *RuntimeError {
	return &RuntimeError{
		Code:    ErrMismatchParamLengthError,
		Message: fmt.Sprintf("此方法定义了%d个参数，而实际输入%d个参数", expect, got),
		Extra:   []int{expect, got},
	}
}

// MostParamsError -
func MostParamsError(maxParams int) *RuntimeError {
	return &RuntimeError{
		Code:    ErrMostParamsError,
		Message: fmt.Sprintf("至多需要%d个参数", maxParams),
		Extra:   maxParams,
	}
}

// ExactParamsError -
func ExactParamsError(exactParams int) *RuntimeError {
	return &RuntimeError{
		Code:    ErrExactParamsError,
		Message: fmt.Sprintf("需要正好%d个参数", exactParams),
		Extra:   exactParams,
	}
}

// ModuleNotFound -
func ModuleNotFound(name string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrModuleNotFound,
		Message: fmt.Sprintf("未找到「%s」模块", name),
		Extra:   name,
	}
}

// ImportSameModule -
func ImportSameModule(name string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrImportSameModule,
		Message: "导入模块与当前模块相同",
		Extra:   name,
	}
}

// DuplicateModule -
func DuplicateModule(name string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrDuplicateModule,
		Message: fmt.Sprintf("重复导入「%s」模块", name),
		Extra:   name,
	}
}

// DuplicateModule -
func ModuleCircularDependency() *RuntimeError {
	return &RuntimeError{
		Code:    ErrModuleCircularDependency,
		Message: "导入模块时出现循环依赖，无法进行下一步操作",
		Extra:   nil,
	}
}

// Internal Error Class, for Zn Internal exception (rare to happen)
// e.g. Unexpected switch-case

// UnexpectedCase -
func UnexpectedCase(tag string, value string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrUnexpectedCase,
		Message: fmt.Sprintf("未定义的条件分支：「%s」的值为「%s」", tag, value),
		Extra:   nil,
	}
}

func UnexpectedEmptyExecLogic() *RuntimeError {
	return &RuntimeError{
		Code:    ErrUnexpectedEmptyExecLogic,
		Message: "执行逻辑不能为空",
		Extra:   nil,
	}
}

func UnexpectedAssign() *RuntimeError {
	return &RuntimeError{
		Code:    ErrUnexpectedAssign,
		Message: "方法不能被赋值",
		Extra:   nil,
	}
}

func UnexpectedParamWildcard() *RuntimeError {
	return &RuntimeError{
		Code:    ErrUnexpectedParamWildcard,
		Message: "无效的参数通配符",
		Extra:   nil,
	}
}

//// type errors

// InvalidExprType -
func InvalidExprType(assertType ...string) *RuntimeError {
	var labels []string
	for _, at := range assertType {
		label := at
		if v, ok := typeNameMap[at]; ok {
			label = v
		}
		labels = append(labels, fmt.Sprintf("「%s」", label))
	}

	return &RuntimeError{
		Code:    ErrInvalidExprType,
		Message: fmt.Sprintf("表达式不符合期望的%s类型", strings.Join(labels, "、")),
		Extra:   labels,
	}
}

// InvalidFuncVariable -
func InvalidFuncVariable(tag string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrInvalidFuncVariable,
		Message: fmt.Sprintf("「%s」须为一个方法", tag),
		Extra:   tag,
	}
}

// InvalidParamType -
func InvalidParamType(assertType ...string) *RuntimeError {
	var labels []string
	for _, at := range assertType {
		label := at
		if v, ok := typeNameMap[at]; ok {
			label = v
		}
		labels = append(labels, fmt.Sprintf("「%s」", label))
	}
	return &RuntimeError{
		Code:    ErrInvalidParamType,
		Message: fmt.Sprintf("输入参数不符合期望之%s类型", strings.Join(labels, "、")),
		Extra:   labels,
	}
}

// InvalidCompareLType - 比较的值的类型
func InvalidCompareLType(assertType ...string) *RuntimeError {
	var labels []string
	for _, at := range assertType {
		label := at
		if v, ok := typeNameMap[at]; ok {
			label = v
		}
		labels = append(labels, fmt.Sprintf("「%s」", label))
	}
	return &RuntimeError{
		Code:    ErrInvalidCompareLType,
		Message: fmt.Sprintf("比较值的类型应为%s", strings.Join(labels, "、")),
		Extra:   labels,
	}
}

// InvalidCompareRType - 被比较的值的类型
func InvalidCompareRType(assertType ...string) *RuntimeError {
	var labels []string
	for _, at := range assertType {
		label := at
		if v, ok := typeNameMap[at]; ok {
			label = v
		}
		labels = append(labels, fmt.Sprintf("「%s」", label))
	}
	return &RuntimeError{
		Code:    ErrInvalidCompareRType,
		Message: fmt.Sprintf("被比较值的类型应为%s", strings.Join(labels, "、")),
		Extra:   labels,
	}
}

func InvalidExceptionType(name string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrInvalidExceptionType,
		Message: fmt.Sprintf("「%s」必须是一个类型！", name),
		Extra:   nil,
	}
}

func InvalidExceptionObjectType(name string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrInvalidExceptionObjectType,
		Message: fmt.Sprintf("「%s」构造出来的对象须是一个异常类型的值！", name),
		Extra:   nil,
	}
}

func ArithDivZero() *RuntimeError {
	return &RuntimeError{
		Code:    ErrArithDivZero,
		Message: "被除数不得为0",
		Extra:   nil,
	}
}

func ArithRootLessThanZero() *RuntimeError {
	return &RuntimeError{
		Code:    ErrArithRootLessThanZero,
		Message: "计算平方根时，底数须大于0",
		Extra:   nil,
	}
}

//// SLOT
func NewErrorSLOT(info string) error {
	return fmt.Errorf(info)
}
