package exec

import (
	"github.com/DemoHn/Zn/error"
)

var defaultBoolClassRef *ClassRef

// ZnBool - (bool) 「二象」型
type ZnBool struct {
	*ZnObject
	Value bool
}

// String - show displayed value
func (zb *ZnBool) String() string {
	data := "真"
	if zb.Value == false {
		data = "假"
	}
	return data
}

// NewZnBool -
func NewZnBool(value bool) *ZnBool {
	return &ZnBool{
		Value:    value,
		ZnObject: NewZnObject(defaultBoolClassRef),
	}
}

func init() {
	var toStringGetter = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
		this, ok := scope.GetTargetThis().(*ZnBool)
		if !ok {
			return nil, error.InvalidParamType("bool")
		}
		return NewZnString(this.String()), nil
	}

	var getterMap = map[string]funcExecutor{
		"文本*": toStringGetter,
	}

	defaultBoolClassRef = NewClassRef("二象")
	bindClassGetters(defaultBoolClassRef, getterMap)
}
