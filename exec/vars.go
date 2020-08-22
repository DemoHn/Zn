package exec

import (
	"fmt"
	"strings"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/syntax"
)

//// General Value types

// ZnValue - general value interface
type ZnValue interface {
	String() string
}

type funcExecutor func(ctx *Context, scope Scope, params []ZnValue) (ZnValue, *error.Error)

//////// ZnObject Definition

// ZnObject -
type ZnObject struct {
	// defines all properties (readable and writable)
	PropList map[string]ZnValue
	// defines all methods: A 之 （B：XX）
	MethodList map[string]*ZnFunction
}

//////// Primitive Types Definition

// ZnString - string 「文本」型
type ZnString struct {
	ZnObject
	Value string
}

// ZnBool - (bool) 「二象」型
type ZnBool struct {
	ZnObject
	Value bool
}

// ZnArray - Zn array type 「元组」型
type ZnArray struct {
	ZnObject
	Value []ZnValue
}

// ZnNull - Zn null type - a special marker indicates that
// this value has neither type nor value
type ZnNull struct {
	ZnObject
}

// ZnFunction -
type ZnFunction struct {
	ZnObject
	Node     *syntax.FunctionDeclareStmt
	Executor funcExecutor
}

// ZnHashMap -
type ZnHashMap struct {
	ZnObject
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

//////// Variable Type Implementation

// String() - display those types
func (zs *ZnString) String() string {
	return fmt.Sprintf("「%s」", zs.Value)
}

func (zb *ZnBool) String() string {
	data := "真"
	if zb.Value == false {
		data = "假"
	}
	return data
}

func (za *ZnArray) String() string {
	strs := []string{}
	for _, item := range za.Value {
		strs = append(strs, item.String())
	}

	return fmt.Sprintf("【%s】", strings.Join(strs, "，"))
}

func (zn *ZnNull) String() string {
	return "空"
}

func (zf *ZnFunction) String() string {
	return fmt.Sprintf("方法： %s", zf.Node.FuncName.GetLiteral())
}

func (zh *ZnHashMap) String() string {
	strs := []string{}
	for _, key := range zh.KeyOrder {
		value := zh.Value[key]
		strs = append(strs, fmt.Sprintf("%s == %s", key, value.String()))
	}
	return fmt.Sprintf("【%s】", strings.Join(strs, "，"))
}

// Rev - ZnBool
func (zb *ZnBool) Rev() *ZnBool {
	zb.Value = !zb.Value
	return zb
}

//////// New[Type] Constructors

// NewZnString -
func NewZnString(value string) *ZnString {
	return &ZnString{
		Value: value,
	}
}

// NewZnBool -
func NewZnBool(value bool) *ZnBool {
	return &ZnBool{
		Value: value,
	}
}

// NewZnArray -
func NewZnArray(values []ZnValue) *ZnArray {
	return &ZnArray{
		Value: values,
	}
}

// NewZnNull - null value
func NewZnNull() *ZnNull {
	t := &ZnNull{}
	return t
}

// NewZnFunction -
func NewZnFunction(node *syntax.FunctionDeclareStmt) *ZnFunction {
	return &ZnFunction{
		Node:     node,
		Executor: nil,
	}
}

// NewZnNativeFunction - new Zn native function
func NewZnNativeFunction(name string, executor funcExecutor) *ZnFunction {
	id := new(syntax.ID)
	id.SetLiteral([]rune(name))

	return &ZnFunction{
		Node: &syntax.FunctionDeclareStmt{
			FuncName: id,
		},
		Executor: executor,
	}
}

// NewZnHashMap -
func NewZnHashMap(kvPairs []KVPair) *ZnHashMap {
	hm := &ZnHashMap{
		Value:    map[string]ZnValue{},
		KeyOrder: []string{},
	}

	for _, kvPair := range kvPairs {
		hm.Value[kvPair.Key] = kvPair.Value
		hm.KeyOrder = append(hm.KeyOrder, kvPair.Key)
	}

	return hm
}

// NewZnObject -
func NewZnObject() *ZnObject {
	obj := &ZnObject{
		PropList:   map[string]ZnValue{},
		MethodList: map[string]*ZnFunction{},
	}

	return obj
}
