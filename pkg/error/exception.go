package error

import "fmt"

type Exception struct {
	Name string
	Message string
}

func (e *Exception) Error() string {
	return fmt.Sprintf("%s：%s", e.Name, e.Message)
}

func NewRuntimeException(message string) *Exception {
	return &Exception{
		Name:    "运行异常",
		Message: message,
	}
}

