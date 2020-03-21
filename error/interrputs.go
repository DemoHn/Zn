package error

// ReturnValueInterrupt - triggers when a value is returned
func ReturnValueInterrupt() *Error {
	return interrupts.NewError(0x01, Error{})
}
