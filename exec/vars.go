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
	Compare(val ZnValue, cmpType znCompareType) (*ZnBool, *error.Error)
}

type znCompareType uint8

const (
	compareTypeEq = 1
	compareTypeIs = 2
	compareTypeLt = 3
	compareTypeGt = 4
)

type funcExecutor func(params []ZnValue, template *syntax.FunctionDeclareStmt, ctx *Context) (ZnValue, *error.Error)

//////// Primitive Types Definition

// ZnString - string 「文本」型
type ZnString struct {
	Value string
}

// ZnBool - (bool) 「二象」型
type ZnBool struct {
	Value bool
}

// ZnArray - Zn array type 「元组」型
type ZnArray struct {
	Value []ZnValue
}

// ZnNull - Zn null type - a special marker indicates that
// this value has neither type nor value
type ZnNull struct{}

// ZnFunction -
type ZnFunction struct {
	Node     *syntax.FunctionDeclareStmt
	Executor funcExecutor
}

// ZnHashMap -
type ZnHashMap struct {
	// now only support string as key
	Value map[string]ZnValue
}

// KVPair - key-value pair, used for ZnHashMap
type KVPair struct {
	Key   string
	Value ZnValue
}

//////// Variable Type Implementation

// String() - display those types
func (zs *ZnString) String() string {
	return fmt.Sprintf("%s", zs.Value)
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
	for key, value := range zh.Value {
		strs = append(strs, fmt.Sprintf("%s == %s", key, value.String()))
	}
	return fmt.Sprintf("【%s】", strings.Join(strs, "，"))
}

// Compare - compare data
func (zs *ZnString) Compare(val ZnValue, cmpType znCompareType) (*ZnBool, *error.Error) {
	var valR *ZnString
	var targetResult int

	switch v := val.(type) {
	case *ZnString:
		valR = v
	case *ZnNull:
		return NewZnBool(false), nil
	default:
		if cmpType == compareTypeEq || cmpType == compareTypeIs {
			return NewZnBool(false), nil
		}
		return nil, error.InvalidExprType("string")
	}

	switch cmpType {
	case compareTypeEq, compareTypeIs:
		targetResult = 0
	case compareTypeGt:
		targetResult = 1
	case compareTypeLt:
		targetResult = -1
	}

	if res := strings.Compare(zs.Value, valR.Value); res == targetResult {
		return NewZnBool(true), nil
	}
	return NewZnBool(false), nil
}

// Compare - ZnBool
func (zb *ZnBool) Compare(val ZnValue, cmpType znCompareType) (*ZnBool, *error.Error) {
	var valR *ZnBool
	switch v := val.(type) {
	case *ZnBool:
		valR = v
	case *ZnNull:
		return NewZnBool(false), nil
	default:
		if cmpType == compareTypeEq || cmpType == compareTypeIs {
			return NewZnBool(false), nil
		}
		return nil, error.InvalidExprType("bool")
	}

	switch cmpType {
	case compareTypeEq, compareTypeIs:
		if zb.Value == valR.Value {
			return NewZnBool(true), nil
		}
		return NewZnBool(false), nil
	default:
		return nil, error.NewErrorSLOT("not supported for ZnBool")
	}
}

// Compare - ZnArray
func (za *ZnArray) Compare(val ZnValue, cmpType znCompareType) (*ZnBool, *error.Error) {
	var valR *ZnArray
	switch v := val.(type) {
	case *ZnArray:
		valR = v
	case *ZnNull:
		return NewZnBool(false), nil
	default:
		if cmpType == compareTypeEq || cmpType == compareTypeIs {
			return NewZnBool(false), nil
		}
		return nil, error.InvalidExprType("array")
	}

	switch cmpType {
	case compareTypeEq, compareTypeIs:
		if len(za.Value) != len(valR.Value) {
			return NewZnBool(false), nil
		}
		// cmp each item
		for idx, item := range za.Value {
			cmpVal, err := item.Compare(valR.Value[idx], cmpType)
			if err != nil {
				return nil, err
			}
			if cmpVal.Value == false {
				return NewZnBool(false), nil
			}
		}
		return NewZnBool(true), nil
	default:
		return nil, error.NewErrorSLOT("not supported for ZnArray")
	}
}

// Compare - ZnNull
func (zn *ZnNull) Compare(val ZnValue, cmpType znCompareType) (*ZnBool, *error.Error) {
	switch val.(type) {
	case *ZnNull:
		return NewZnBool(true), nil
	default:
		return NewZnBool(false), nil
	}
}

// Compare - ZnFunction
func (zf *ZnFunction) Compare(val ZnValue, cmpType znCompareType) (*ZnBool, *error.Error) {
	return nil, error.NewErrorSLOT("function is incomparable!")
}

// Compare - ZnHashMap
func (zh *ZnHashMap) Compare(val ZnValue, cmpType znCompareType) (*ZnBool, *error.Error) {
	// TODO -
	return nil, nil
}

// Rev - ZnBool
func (zb *ZnBool) Rev() *ZnBool {
	zb.Value = !zb.Value
	return zb
}

// Exec - ZnFunction exec function
func (zf *ZnFunction) Exec(params []ZnValue, ctx *Context, env Env) (ZnValue, *error.Error) {
	// TODO1: add new env
	// st -> global symbol table
	// if executor = nil, then use default function executor
	if zf.Executor == nil {
		// enter scope to add values
		ctx.EnterScope()
		defer ctx.ExitScope()
		// check param length
		if len(params) != len(zf.Node.ParamList) {
			return nil, error.NewErrorSLOT("param list length mismatch!")
		}

		// set id
		for idx, param := range params {
			paramID := zf.Node.ParamList[idx]
			err := ctx.Bind(paramID.GetLiteral(), param, false)
			if err != nil {
				return nil, err
			}
		}

		result := ctx.ExecuteBlockAST(zf.Node.ExecBlock, env)
		if result.HasError {
			return nil, result.Error
		}
		return result.Value, nil
	}

	return zf.Executor(params, zf.Node, ctx)
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
	return &ZnNull{}
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
		Value: map[string]ZnValue{},
	}

	for _, kvPair := range kvPairs {
		hm.Value[kvPair.Key] = kvPair.Value
	}

	return hm
}
