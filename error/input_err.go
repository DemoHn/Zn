package error

import "fmt"

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

// FileNotFound -
func FileNotFound(filePath string) *Error {
	err := inputError.NewError(0x02, Error{
		text:   fmt.Sprintf("未找到文件：%s", filePath),
		cursor: nil,
		info:   filePath,
	})

	return &err
}

// FileOpenError -
func FileOpenError(filePath string, oriError error) *Error {
	err := inputError.NewError(0x03, Error{
		text:   fmt.Sprintf("无法打开文件：%s！", filePath),
		cursor: nil,
		info: struct {
			path string
			err  error
		}{filePath, oriError},
	})

	return &err
}
