package exec

import (
	"fmt"
	"strings"

	"github.com/DemoHn/Zn/error"
)

// ZnValue - general value interface
type ZnValue interface {
	String() string
}

// ZnComparable - to make the value comparable
type ZnComparable interface {
	Equals(val ZnComparable) (*ZnBool, *error.Error)      // A 等于 B
	Is(val ZnComparable) (*ZnBool, *error.Error)          // 此 A 为 B
	LessThan(val ZnComparable) (*ZnBool, *error.Error)    // A 小于 B
	GreaterThan(val ZnComparable) (*ZnBool, *error.Error) // A 大于 B
}

var predefinedValues map[string]ZnValue

//// Primitive Types Definition

// ZnString - string 「文本」型
type ZnString struct {
	Value string
}

func (zs *ZnString) String() string {
	return fmt.Sprintf("「%s」", zs.Value)
}

// NewZnString -
func NewZnString(value string) *ZnString {
	return &ZnString{
		Value: value,
	}
}

// Equals -
func (zs *ZnString) Equals(val ZnComparable) (*ZnBool, *error.Error) {
	v, ok := val.(*ZnString)
	if !ok {
		return nil, error.NewErrorSLOT("Right value must be ZnString")
	}
	if res := strings.Compare(zs.Value, v.Value); res == 0 {
		return NewZnBool(true), nil
	}
	return NewZnBool(false), nil
}

// Is -
func (zs *ZnString) Is(val ZnComparable) (*ZnBool, *error.Error) {
	v, ok := val.(*ZnString)
	if !ok {
		return nil, error.NewErrorSLOT("Right value must be ZnString")
	}
	if res := strings.Compare(zs.Value, v.Value); res == 0 {
		return NewZnBool(true), nil
	}
	return NewZnBool(false), nil
}

// LessThan -
func (zs *ZnString) LessThan(val ZnComparable) (*ZnBool, *error.Error) {
	v, ok := val.(*ZnString)
	if !ok {
		return nil, error.NewErrorSLOT("Right value must be ZnString")
	}
	if res := strings.Compare(zs.Value, v.Value); res == -1 {
		return NewZnBool(true), nil
	}
	return NewZnBool(false), nil
}

// GreaterThan -
func (zs *ZnString) GreaterThan(val ZnComparable) (*ZnBool, *error.Error) {
	v, ok := val.(*ZnString)
	if !ok {
		return nil, error.NewErrorSLOT("Right value must be ZnString")
	}
	if res := strings.Compare(zs.Value, v.Value); res == 1 {
		return NewZnBool(true), nil
	}
	return NewZnBool(false), nil
}

// ZnBool - (bool) 「二象」型
type ZnBool struct {
	Value bool
}

func (zb *ZnBool) String() string {
	data := "真"
	if zb.Value == false {
		data = "假"
	}
	return data
}

// Rev - reverse value; i.e. from TRUE -> FALSE or from FALSE -> TRUE
func (zb *ZnBool) Rev() *ZnBool {
	zb.Value = !zb.Value
	return zb
}

// Equals -
func (zb *ZnBool) Equals(val ZnComparable) (*ZnBool, *error.Error) {
	v, ok := val.(*ZnBool)
	if !ok {
		return nil, error.NewErrorSLOT("Right value must be ZnBool")
	}
	if zb.Value == v.Value {
		return NewZnBool(true), nil
	}
	return NewZnBool(false), nil
}

// Is -
func (zb *ZnBool) Is(val ZnComparable) (*ZnBool, *error.Error) {
	v, ok := val.(*ZnBool)
	if !ok {
		return nil, error.NewErrorSLOT("Right value must be ZnBool")
	}
	if zb.Value == v.Value {
		return NewZnBool(true), nil
	}
	return NewZnBool(false), nil
}

// LessThan -
func (zb *ZnBool) LessThan(val ZnComparable) (*ZnBool, *error.Error) {
	return nil, error.NewErrorSLOT("not supported for ZnBool")
}

// GreaterThan -
func (zb *ZnBool) GreaterThan(val ZnComparable) (*ZnBool, *error.Error) {
	return nil, error.NewErrorSLOT("not supported for ZnBool")
}

// NewZnBool -
func NewZnBool(value bool) *ZnBool {
	return &ZnBool{
		Value: value,
	}
}

// ZnArray - Zn array type 「元组」型
type ZnArray struct {
	Value []ZnValue
}

func (za *ZnArray) String() string {
	strs := []string{}
	for _, item := range za.Value {
		strs = append(strs, item.String())
	}

	return fmt.Sprintf("【%s】", strings.Join(strs, "，"))
}

// Equals -
func (za *ZnArray) Equals(val ZnComparable) (*ZnBool, *error.Error) {
	v, ok := val.(*ZnArray)
	if !ok {
		return nil, error.NewErrorSLOT("right value must be ZnArray")
	}
	if len(za.Value) != len(v.Value) {
		return NewZnBool(false), nil
	}
	// cmp each item
	for idx, item := range za.Value {
		vitemL, okL := item.(ZnComparable)
		if !okL {
			return nil, error.NewErrorSLOT("item must be comparable")
		}
		vitemR, okR := v.Value[idx].(ZnComparable)
		if !okR {
			return nil, error.NewErrorSLOT("item must be comparable")
		}
		cmpVal, err := vitemL.Equals(vitemR)
		if err != nil {
			return nil, err
		}
		if cmpVal.Value == false {
			return NewZnBool(false), nil
		}
	}
	return NewZnBool(true), nil
}

// Is -
func (za *ZnArray) Is(val ZnComparable) (*ZnBool, *error.Error) {
	v, ok := val.(*ZnArray)
	if !ok {
		return nil, error.NewErrorSLOT("right value must be ZnArray")
	}
	if len(za.Value) != len(v.Value) {
		return NewZnBool(false), nil
	}
	// cmp each item
	for idx, item := range za.Value {
		vitemL, okL := item.(ZnComparable)
		if !okL {
			return nil, error.NewErrorSLOT("item must be comparable")
		}
		vitemR, okR := v.Value[idx].(ZnComparable)
		if !okR {
			return nil, error.NewErrorSLOT("item must be comparable")
		}
		cmpVal, err := vitemL.Is(vitemR)
		if err != nil {
			return nil, err
		}
		if cmpVal.Value == false {
			return NewZnBool(false), nil
		}
	}
	return NewZnBool(true), nil
}

// LessThan -
func (za *ZnArray) LessThan(val ZnComparable) (*ZnBool, *error.Error) {
	return nil, error.NewErrorSLOT("not supported for ZnArray")
}

// GreaterThan -
func (za *ZnArray) GreaterThan(val ZnComparable) (*ZnBool, *error.Error) {
	return nil, error.NewErrorSLOT("not supported for ZnArray")
}

// NewZnArray -
func NewZnArray(values []ZnValue) *ZnArray {
	return &ZnArray{
		Value: values,
	}
}

// ZnNull - Zn null type - a special marker indicates that
// this value has neither type nor value
type ZnNull struct{}

func (zn *ZnNull) String() string {
	return "‹空›"
}

// Equals -
func (zn *ZnNull) Equals(val ZnComparable) (*ZnBool, *error.Error) {
	_, ok := val.(*ZnNull)
	if !ok {
		return NewZnBool(true), nil
	}
	return NewZnBool(false), nil
}

// Is -
func (zn *ZnNull) Is(val ZnComparable) (*ZnBool, *error.Error) {
	_, ok := val.(*ZnNull)
	if !ok {
		return NewZnBool(true), nil
	}
	return NewZnBool(false), nil
}

// LessThan -
func (zn *ZnNull) LessThan(val ZnComparable) (*ZnBool, *error.Error) {
	return nil, error.NewErrorSLOT("not supported for ZnNull")
}

// GreaterThan -
func (zn *ZnNull) GreaterThan(val ZnComparable) (*ZnBool, *error.Error) {
	return nil, error.NewErrorSLOT("not supported for ZnNull")
}

// NewZnNull - null value
func NewZnNull() *ZnNull {
	return &ZnNull{}
}

// init function
func init() {
	//// predefined values - those variables (symbols) are defined before
	//// any execution procedure.
	//// NOTICE: those variables are all constants!
	predefinedValues = map[string]ZnValue{
		"真": NewZnBool(true),
		"假": NewZnBool(false),
		"空": NewZnNull(),
	}
}
