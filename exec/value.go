package exec

import (
	"fmt"
	"strings"
)

// ZnValue - general value interface
type ZnValue interface {
	String() string
}

//// Primitive Types Definition

// ZnString - string
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

// ZnBool - (bool)
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

// ZnArray - Zn array type
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
