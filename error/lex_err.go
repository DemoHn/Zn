package error

// InvalidSingleEm - only single em exists is meaningless!
// TODO: add cursor
func InvalidSingleEm(idx int) *Error {
	err := lexError.NewError(0x01, Error{
		text:   "未能识别单个「—」字符，请注意需要有连续两个方可有效！",
		cursor: nil,
		info:   idx,
	})

	return &err
}

// InvalidSingleEllipsis -
// TODO: add cursor
func InvalidSingleEllipsis(idx int) *Error {
	err := lexError.NewError(0x02, Error{
		text:   "未能识别单个「…」字符，请注意需要有连续两个方可有效！",
		cursor: nil,
		info:   idx,
	})

	return &err
}

// InvalidIndent -
func InvalidIndent(idx int) *Error {
	err := lexError.NewError(0x03, Error{
		text:   "Indent Error",
		cursor: nil,
		info:   idx,
	})

	return &err
}
