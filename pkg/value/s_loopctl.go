package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

// LoopCtl - a loopctl type used to control loop flow (aka. continue/break)
// and record values of current iteration.
type LoopCtl struct {
	currentIndex r.Value
	currentValue r.Value
}

// NewLoopCtl -
func NewLoopCtl() *LoopCtl {
	return &LoopCtl{
		currentIndex: NewNull(),
		currentValue: NewNull(),
	}
}

// GetProperty -
func (lc *LoopCtl) GetProperty(c *r.Context, name string) (r.Value, error) {
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (lc *LoopCtl) SetProperty(c *r.Context, name string, value r.Value) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (lc *LoopCtl) ExecMethod(c *r.Context, name string, values []r.Value) (r.Value, error) {
	switch name {
	case "结束循环":
		return NewNull(), zerr.NewBreakSignal()
	case "继续循环":
		return NewNull(), zerr.NewContinueSignal()
	}
	return nil, zerr.MethodNotFound(name)
}

// SetCurrentKeyValue - internal usage to set current value
func (lc *LoopCtl) SetCurrentKeyValue(index r.Value, value r.Value) {
	lc.currentIndex = index
	lc.currentValue = value
}

