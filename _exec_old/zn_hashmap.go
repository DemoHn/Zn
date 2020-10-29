package exec

import (
	"fmt"
	"strings"

	"github.com/DemoHn/Zn/error"
)

var defaultHashMapClassRef *ClassRef

// ZnHashMap -
type ZnHashMap struct {
	*ZnObject
	// now only support string as key
	Value map[string]ZnValue
	// The order of key is (delibrately) random when iterating a hashmap.
	// Thus, we preserve the (insertion) order of key using "KeyOrder" field.
	KeyOrder []string
}

// KVPair - key-value pair, used for ZnHashMap
type KVPair struct {
	Key   string
	Value ZnValue
}

func (zh *ZnHashMap) String() string {
	strs := []string{}
	for _, key := range zh.KeyOrder {
		value := zh.Value[key]
		strs = append(strs, fmt.Sprintf("%s == %s", key, value.String()))
	}
	return fmt.Sprintf("【%s】", strings.Join(strs, "，"))
}

// NewZnHashMap -
func NewZnHashMap(kvPairs []KVPair) *ZnHashMap {
	hm := &ZnHashMap{
		Value:    map[string]ZnValue{},
		KeyOrder: []string{},
		ZnObject: NewZnObject(defaultHashMapClassRef),
	}

	for _, kvPair := range kvPairs {
		hm.Value[kvPair.Key] = kvPair.Value
		hm.KeyOrder = append(hm.KeyOrder, kvPair.Key)
	}

	return hm
}

func init() {
	var mapCountGetter = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
		this, ok := scope.GetTargetThis().(*ZnHashMap)
		if !ok {
			return nil, error.InvalidParamType("hashmap")
		}
		return NewZnDecimalFromInt(len(this.Value), 0), nil
	}

	var getterMap = map[string]funcExecutor{
		"长度": mapCountGetter,
	}

	defaultHashMapClassRef = NewClassRef("列表")
	bindClassGetters(defaultHashMapClassRef, getterMap)
}
