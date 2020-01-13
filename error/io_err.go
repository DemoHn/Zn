package error

import "fmt"

// FileNotFound - file not found
func FileNotFound(path string) *Error {
	return lexError.NewError(0x10, Error{
		text: fmt.Sprintf("未能找到文件：%s，请检查它是否存在！", path),
		info: path,
	})
}

// FileOpenError -
func FileOpenError(filePath string, oriError error) *Error {
	return lexError.NewError(0x11, Error{
		text: fmt.Sprintf("未能读取文件：%s，请检查其是否存在及有无读取权限！", filePath),
		info: struct {
			path string
			err  error
		}{filePath, oriError},
	})
}
