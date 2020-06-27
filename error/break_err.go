package error

// ReturnBreakError - breaks when return statement is executed
func ReturnBreakError(extra interface{}) *Error {
	return breakError.NewError(0x01, Error{
		text:  "「返回」中断",
		extra: extra,
	})
}
