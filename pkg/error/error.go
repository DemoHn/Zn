package error

// Error - Zn's general error model (with error stack)
type Error struct {
	Code int
	Message string
	Extra interface{}
}

func (e *Error) Error() string {
	return e.Message
}