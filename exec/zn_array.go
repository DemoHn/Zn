package exec

import (
	"fmt"
	"strings"

	"github.com/DemoHn/Zn/error"
)

var defaultArrayClassRef *ClassRef

// ZnArray - Zn array type 「元组」型
type ZnArray struct {
	*ZnObject
	Value []ZnValue
}

// NewZnArray -
func NewZnArray(values []ZnValue) *ZnArray {
	return &ZnArray{
		Value:    values,
		ZnObject: NewZnObject(defaultArrayClassRef),
	}
}

func (za *ZnArray) String() string {
	strs := []string{}
	for _, item := range za.Value {
		strs = append(strs, item.String())
	}

	return fmt.Sprintf("【%s】", strings.Join(strs, "，"))
}

func init() {
	var arraySumGetter = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
		this, ok := scope.GetTargetThis().(*ZnArray)
		if !ok {
			return nil, error.InvalidParamType("array")
		}
		return addValueExecutor(ctx, scope, this.Value)
	}

	var arraySubGetter = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
		this, ok := scope.GetTargetThis().(*ZnArray)
		if !ok {
			return nil, error.InvalidParamType("array")
		}
		return subValueExecutor(ctx, scope, this.Value)
	}

	var arrayMulGetter = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
		this, ok := scope.GetTargetThis().(*ZnArray)
		if !ok {
			return nil, error.InvalidParamType("array")
		}
		return mulValueExecutor(ctx, scope, this.Value)
	}

	var arrayDivGetter = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
		this, ok := scope.GetTargetThis().(*ZnArray)
		if !ok {
			return nil, error.InvalidParamType("array")
		}
		return divValueExecutor(ctx, scope, this.Value)
	}

	var arrayFirstGetter = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
		this, ok := scope.GetTargetThis().(*ZnArray)
		if !ok {
			return nil, error.InvalidParamType("array")
		}
		if len(this.Value) == 0 {
			return NewZnNull(), nil
		}
		return this.Value[0], nil
	}

	var arrayLastGetter = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
		this, ok := scope.GetTargetThis().(*ZnArray)
		if !ok {
			return nil, error.InvalidParamType("array")
		}
		if len(this.Value) == 0 {
			return NewZnNull(), nil
		}
		return this.Value[len(this.Value)-1], nil
	}

	var arrayCountGetter = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
		this, ok := scope.GetTargetThis().(*ZnArray)
		if !ok {
			return nil, error.InvalidParamType("array")
		}
		return NewZnDecimalFromInt(len(this.Value), 0), nil
	}

	// array's getter list
	var getterMap = map[string]funcExecutor{
		"和":  arraySumGetter,
		"差":  arraySubGetter,
		"积":  arrayMulGetter,
		"商":  arrayDivGetter,
		"首":  arrayFirstGetter,
		"尾":  arrayLastGetter,
		"数目": arrayCountGetter,
		"长度": arrayCountGetter,
	}

	defaultArrayClassRef = NewClassRef("数组")
	bindClassGetters(defaultArrayClassRef, getterMap)
}
