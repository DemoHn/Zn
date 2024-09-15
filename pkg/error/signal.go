package error

import (
	"fmt"
)

const (
	SigTypeReturn    uint8 = 1
	SigTypeContinue  uint8 = 2
	SigTypeBreak     uint8 = 3
	SigTypeException uint8 = 4
)

// signal name
var sigTypeMap = map[uint8]string{
	SigTypeReturn:    "返回",
	SigTypeContinue:  "继续",
	SigTypeBreak:     "结束",
	SigTypeException: "异常",
}

type Signal struct {
	SigType uint8
	Extra   interface{}
}

func (s *Signal) Error() string {
	return fmt.Sprintf("收到「%s」中断信号", sigTypeMap[s.SigType])
}

// return XXX
func NewReturnSignal(val interface{}) *Signal {
	return &Signal{SigTypeReturn, val}
}

// continue in for loop
func NewContinueSignal() *Signal {
	return &Signal{SigTypeContinue, nil}
}

// break in for loop
func NewBreakSignal() *Signal {
	return &Signal{SigTypeBreak, nil}
}

// throw exception
func NewExceptionSignal(val interface{}) *Signal {
	return &Signal{SigTypeException, val}
}
