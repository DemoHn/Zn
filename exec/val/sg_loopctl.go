package val

import "github.com/DemoHn/Zn/error"

// LoopCtl - a sgValue type used to control loop flow (aka. continue/break)
// and record values of current iteration.
type LoopCtl struct {
	currentIndex Value
	currentValue Value
}

// NewLoopCtl -
func NewLoopCtl() *LoopCtl {
	return &LoopCtl{
		currentIndex: NewNull(),
		currentValue: NewNull(),
	}
}

// GetProperty -
func (lc *LoopCtl) GetProperty(ctx *Context, name string) (Value, *error.Error) {
	switch name {
	case "值":
		return lc.currentValue, nil
	case "索引":
		return lc.currentIndex, nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (lc *LoopCtl) SetProperty(ctx *Context, name string, value Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (lc *LoopCtl) ExecMethod(ctx *Context, name string, values []Value) (Value, *error.Error) {
	switch name {
	case "结束":
		return NewNull(), error.BreakBreakError()
	case "继续":
		return NewNull(), error.ContinueBreakError()
	}
	return nil, error.MethodNotFound(name)
}

// SetCurrentKeyValue - internal usage to set current value
func (lc *LoopCtl) SetCurrentKeyValue(index Value, value Value) {
	lc.currentIndex = index
	lc.currentValue = value
}
