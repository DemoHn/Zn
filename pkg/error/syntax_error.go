package error

import "fmt"

type SyntaxError struct {
	Code    int
	Message string
	Cursor  int
}

func (e *SyntaxError) Error() string {
	return e.Message
}

const (
	// 20-25
	ErrInvalidSyntax           = 20
	ErrUnexpectedIndent        = 21
	ErrMustTypeID              = 22
	ErrInvalidIndent           = 23
	ErrInvalidIndentSpaceCount = 24
	ErrInvalidChar             = 25
	ErrEscapeStringFailed      = 26
)

// InvalidSyntax -
func InvalidSyntax(startIdx int) *SyntaxError {
	return &SyntaxError{
		Code:    ErrInvalidSyntax,
		Message: "当前语法不符合要求",
		Cursor:  startIdx,
	}
}

func UnexpectedIndent(startIdx int) *SyntaxError {
	return &SyntaxError{
		Code:    ErrUnexpectedIndent,
		Message: "出现不符合预期的缩进",
		Cursor:  startIdx,
	}
}

// ExprMustTypeID -
func ExprMustTypeID(startIdx int) *SyntaxError {
	return &SyntaxError{
		Code:    ErrMustTypeID,
		Message: "表达式须为「标识符」",
		Cursor:  startIdx,
	}
}

// InvalidIndentType -
func InvalidIndentType(expect uint8, got uint8, startIdx int) *SyntaxError {
	findName := func(item uint8) string {
		name := "「空格」"
		if item == uint8(9) { // TAB
			name = "「TAB」"
		}
		return name
	}

	return &SyntaxError{
		Code:    ErrInvalidIndent,
		Message: fmt.Sprintf("本行的缩进类型为%s，与此前缩进类型%s不符", findName(got), findName(expect)),
		Cursor:  startIdx,
	}
}

// InvalidIndentSpaceCount -
func InvalidIndentSpaceCount(count int, startIdx int) *SyntaxError {
	return &SyntaxError{
		Code:    ErrInvalidIndentSpaceCount,
		Message: fmt.Sprintf("当缩进类型为「空格」，其所列字符数应为4之倍数：当前空格字符数为%d", count),
		Cursor:  startIdx,
	}
}

// InvalidChar -
func InvalidChar(ch rune, startIdx int) *SyntaxError {
	return &SyntaxError{
		Code:    ErrInvalidChar,
		Message: fmt.Sprintf("未能识别字符「%c」", ch),
		Cursor:  startIdx,
	}
}

func EscapeStringFailed(startIdx int) *SyntaxError {
	return &SyntaxError{
		Code:    ErrEscapeStringFailed,
		Message: "解析文本中的特殊字符时出现异常",
		Cursor:  startIdx,
	}
}
