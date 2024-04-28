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
	ErrInvalidIDFormat = 30
	ErrIDNumberONLY    = 31
	ErrIDNameONLY      = 32
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
