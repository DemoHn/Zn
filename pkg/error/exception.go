package error

type Exception struct {
	name    string
	message string
}

func (e *Exception) Error() string {
	return e.message
}

func (e *Exception) GetName() string {
	return e.name
}

func NewRuntimeException(message string) *Exception {
	return &Exception{
		name:    "运行异常",
		message: message,
	}
}
