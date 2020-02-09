package exec

import (
	"fmt"
	"strings"
)

// ZnValue - general value interface
type ZnValue interface {
	String() string
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

func (zs *ZnArray) String() string {
	strs := []string{}
	for _, item := range zs.Value {
		strs = append(strs, item.String())
	}

	return fmt.Sprintf("【%s】", strings.Join(strs, "，"))
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
