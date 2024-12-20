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
	ErrIndexOutOfRange          = 40
	ErrIndexKeyNotFound         = 41
	ErrNameNotDefined           = 42
	ErrNameRedeclared           = 43
	ErrAssignToConstant         = 44
	ErrPropertyNotFound         = 45
	ErrMethodNotFound           = 46
	ErrClassNotOnRoot           = 47
	ErrThisValueNotFound        = 48
	ErrInvalidExceptionClass    = 49
	ErrLeastParamsError         = 50
	ErrMismatchParamLengthError = 51
	ErrMostParamsError          = 52
	ErrExactParamsError         = 53
	// module error
	ErrModuleNotFound           = 60
	ErrImportSameModule         = 61
	ErrDuplicateModule          = 62
	ErrModuleCircularDependency = 63
	// internal error
	ErrUnexpectedCase           = 70
	ErrUnexpectedEmptyExecLogic = 71
	ErrUnexpectedAssign         = 72
	ErrUnexpectedParamWildcard  = 73
	// type error
	ErrInvalidExprType            = 80
	ErrInvalidFuncVariable        = 81
	ErrInvalidParamType           = 82
	ErrInvalidCompareLType        = 83
	ErrInvalidCompareRType        = 84
	ErrInvalidExceptionType       = 85
	ErrInvalidExceptionObjectType = 86
	ErrInvalidClassType           = 87
	// arith error
	ErrArithDivZero          = 90
	ErrArithRootLessThanZero = 91
	// input error
	ErrInputValueNotFound = 95
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
		Message: "此变量的值不允许更改",
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

func InvalidExceptionClass(name string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrInvalidExceptionClass,
		Message: fmt.Sprintf("当前接到的异常类型并非「%s」", name),
		Extra:   name,
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
		Message: fmt.Sprintf("未定义的分支逻辑：「%s」的值为「%s」", tag, value),
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
		if label != "" {
			labels = append(labels, fmt.Sprintf("「%s」", label))
		}
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

// InvalidClassType -
func InvalidClassType(tag string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrInvalidClassType,
		Message: fmt.Sprintf("「%s」须为一个定义类型", tag),
		Extra:   tag,
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

func InputValueNotFound(tag string) *RuntimeError {
	return &RuntimeError{
		Code:    ErrInputValueNotFound,
		Message: fmt.Sprintf("没有设置输入变量「%s」的值", tag),
		Extra:   tag,
	}
}

//// SLOT
func NewErrorSLOT(info string) error {
	return fmt.Errorf(info)
}
