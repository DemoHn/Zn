package error

import (
	"fmt"
	"io"
)

type IOError struct {
	Code    int
	Message string
	Path    string
}

// errCode: 1X
const (
	ErrFileNotFound = 10
	ErrReadFile     = 11
	ErrReadVarInput = 12
)

func (e *IOError) Error() string {
	return fmt.Sprintf("读取%s失败，%s", e.Path, e.Message)
}

func FileNotFound(path string) *IOError {
	return &IOError{
		Code:    ErrFileNotFound,
		Message: "未能找到文件，请检查它是否存在！",
		Path:    fmt.Sprintf("文件「%s」", path),
	}
}

// ReadFileError - Read I/O Stream failed
func ReadFileError(err error, path string) *IOError {
	errTextMap := map[error]string{
		io.ErrShortBuffer:   "需要更大的缓冲区",
		io.ErrUnexpectedEOF: "未知文件结束符",
		io.ErrNoProgress:    "多次尝试读取，皆没有获取到数据",
		io.ErrShortWrite:    "数据没有完全写入",
	}

	errText := err.Error()
	if v, ok := errTextMap[err]; ok {
		errText = fmt.Sprintf("%s (%s)", v, err.Error())
	}

	return &IOError{
		Code:    ErrReadFile,
		Message: fmt.Sprintf("读取I/O流失败：%s", errText),
		Path:    path,
	}
}

func ReadVarInputError(err error) *IOError {
	return &IOError{
		Code:    ErrReadVarInput,
		Message: "无法读取内容",
		Path:    "预定义变量",
	}
}
