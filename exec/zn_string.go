package exec

import (
	"fmt"
	"unicode/utf8"

	"github.com/DemoHn/Zn/error"
)

var defaultStringClassRef *ClassRef

// ZnString - string 「文本」型
type ZnString struct {
	*ZnObject
	Value string
}

// String() - display those types
func (zs *ZnString) String() string {
	return fmt.Sprintf("「%s」", zs.Value)
}

// NewZnString -
func NewZnString(value string) *ZnString {
	return &ZnString{
		ZnObject: NewZnObject(defaultStringClassRef),
		Value:    value,
	}
}

func init() {
	var stringCountGetter = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
		this, ok := scope.GetTargetThis().(*ZnString)
		if !ok {
			return nil, error.InvalidParamType("string")
		}
		return NewZnDecimalFromInt(utf8.RuneCountInString(this.Value), 0), nil
	}

	var getterMap = map[string]funcExecutor{
		"长度": stringCountGetter,
	}

	defaultStringClassRef = NewClassRef("文本")
	bindClassGetters(defaultStringClassRef, getterMap)
}
