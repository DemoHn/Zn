package error

import (
	"fmt"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

const (
	SigTypeReturn uint8 = 1
	SigTypeContinue uint8 = 2
	SigTypeBreak uint8 = 3
)

// signal name
var sigTypeMap = map[uint8]string{
	SigTypeReturn: "返回",
	SigTypeContinue: "继续",
	SigTypeBreak: "结束",
}
type Signal struct {
	SigType uint8
	Extra interface{}
}

func (s *Signal) Error() string {
	return fmt.Sprintf("发现「%s」中断", sigTypeMap[s.SigType])
}

// return XXX
func NewReturnSignal(val r.Value) *Signal {
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