package val

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

// LoopCtl - a sgValue type used to control loop flow (aka. continue/break)
// and record values of current iteration.
type LoopCtl struct {
	currentIndex ctx.Value
	currentValue ctx.Value
}

// NewLoopCtl -
func NewLoopCtl() *LoopCtl {
	return &LoopCtl{
		currentIndex: NewNull(),
		currentValue: NewNull(),
	}
}

// GetProperty -
func (lc *LoopCtl) GetProperty(c *ctx.Context, name string) (ctx.Value, *error.Error) {
	switch name {
	case "值":
		return lc.currentValue, nil
	case "索引":
		return lc.currentIndex, nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (lc *LoopCtl) SetProperty(c *ctx.Context, name string, value ctx.Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (lc *LoopCtl) ExecMethod(c *ctx.Context, name string, values []ctx.Value) (ctx.Value, *error.Error) {
	switch name {
	case "结束":
		return NewNull(), error.BreakBreakError()
	case "继续":
		return NewNull(), error.ContinueBreakError()
	}
	return nil, error.MethodNotFound(name)
}

// SetCurrentKeyValue - internal usage to set current value
func (lc *LoopCtl) SetCurrentKeyValue(index ctx.Value, value ctx.Value) {
	lc.currentIndex = index
	lc.currentValue = value
}
