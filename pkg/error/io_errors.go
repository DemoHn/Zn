package error

import (
	"fmt"
	"io"
)

func FileNotFound (path string) *Error {
	return &Error{
		Code: 0x2010,
		Message: fmt.Sprintf("未能找到文件 %s，请检查它是否存在！", path),
		Extra: extraMap{
			"path": path,
		},
	}
}

// ReadFileError - Read I/O Stream failed
func ReadFileError(err error) *Error {
	errTextMap := map[error]string{
		io.ErrShortBuffer:   "需要更大的缓冲区",
		io.ErrUnexpectedEOF: "未知文件结束符",
		io.ErrNoProgress:    "多次尝试读取，皆无数据或返回错误",
		io.ErrShortWrite:    "操作写入的数据比提供的少",
	}

	errText := err.Error()
	if v, ok := errTextMap[err]; ok {
		errText = fmt.Sprintf("%s (%s)", v, err.Error())
	}

	return &Error{
		Code: 0x2012,
		Message: fmt.Sprintf("读取I/O流失败：%s", errText),
		Extra: extraMap{
			"error": errText,
		},
	}
}