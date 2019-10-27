package error

// DecodeUTF8Fail -
// TODO: add params
func DecodeUTF8Fail() *Error {
	err := inputError.NewError(0x01, Error{
		text:   "解析UTF-8输入失败，请确认输入的编码是UTF-8！",
		cursor: nil,
		info:   nil,
	})

	return &err
}
