package error

// ArithDivZeroError - for A/B, when B = 0
func ArithDivZeroError() *Error {
	return arithError.NewError(0x01, Error{
		text: "被除数不得为0",
	})
}
