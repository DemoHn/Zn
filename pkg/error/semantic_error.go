package error

import "fmt"

type SemanticError struct {
	Code    int
	Message string
}

func (e *SemanticError) Error() string {
	return e.Message
}

const (
	// 30-39
	ErrInvalidIDFormat    = 30
	ErrIDNumberONLY       = 31
	ErrIDNameONLY         = 32
	ErrInvalidFmtTemplate = 33
	ErrUnmatchFmtParams   = 34
)

func InvalidIDFormat(idStr string) *SemanticError {
	return &SemanticError{
		Code:    ErrInvalidIDFormat,
		Message: fmt.Sprintf("标识符「%s」格式不符合要求", idStr),
	}
}

func IDNumberONLY(idStr string) *SemanticError {
	return &SemanticError{
		Code:    ErrIDNumberONLY,
		Message: "标识符只允许「数值」格式",
	}
}

func IDNameONLY(idStr string) *SemanticError {
	return &SemanticError{
		Code:    ErrIDNameONLY,
		Message: "标识符只允许「名称」格式",
	}
}

func InvalidFmtTemplate(template string) *SemanticError {
	return &SemanticError{
		Code:    ErrInvalidFmtTemplate,
		Message: fmt.Sprintf("解析文本拼接模板「%s」出现异常", template),
	}
}

func UnmatchFmtParams(template string) *SemanticError {
	return &SemanticError{
		Code:    ErrUnmatchFmtParams,
		Message: fmt.Sprintf("文本拼接模板「%s」所需数量与参数不匹配", template),
	}
}
