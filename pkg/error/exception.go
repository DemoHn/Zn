package error

type Exception struct {
	Name string
	Message string
}

func (e *Exception) Error() string {
	return e.Message
}

func NewRuntimeException(message string) *Exception {
	return &Exception{
		Name:    "运行异常",
		Message: message,
	}
}

