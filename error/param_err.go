package error

import "fmt"

// LeastParamsError -
func LeastParamsError(minParams int) *Error {
	return paramError.NewError(0x01, Error{
		text: fmt.Sprintf("需要输入至少%d个参数", minParams),
		info: fmt.Sprintf("minParams=(%d)", minParams),
	})
}

// MismatchParamLengthError -
func MismatchParamLengthError(expect int, got int) *Error {
	return paramError.NewError(0x02, Error{
		text: fmt.Sprintf("此方法定义了%d个参数，而实际输入%d个参数", expect, got),
		info: fmt.Sprintf("expect=(%d) got=(%d)", expect, got),
	})
}
