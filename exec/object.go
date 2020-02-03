package exec

import (
	"fmt"
	"strings"
)

// ZnObject - the global Zn Object interface
type ZnObject interface {
	IsNull() bool
	// String - print string
	String() string
}

// ZnNullable - if the value of object is null
type ZnNullable struct {
	nullFlag bool
}

// IsNull -
func (zl *ZnNullable) IsNull() bool {
	return zl.nullFlag
}

// SetNull -
func (zl *ZnNullable) SetNull() {
	zl.nullFlag = true
}

// UnsetNull -
func (zl *ZnNullable) UnsetNull() {
	zl.nullFlag = false
}

//// primitive types

// ZnString - Zn string concrete type
type ZnString struct {
	ZnNullable
	Value string
}

func (zs *ZnString) String() string {
	return fmt.Sprintf("「%s」", zs.Value)
}

// SetValue -
func (zs *ZnString) SetValue(v string) bool {
	zs.Value = v
	return true
}

// ZnArray - Zn array type
type ZnArray struct {
	ZnNullable
	Items []ZnObject
}

// Init - initialize an array with predefined seqeuencial data of ZnObject
func (zs *ZnArray) Init(objs []ZnObject) {
	zs.Items = objs
}

func (zs *ZnArray) String() string {
	strs := []string{}
	for _, item := range zs.Items {
		strs = append(strs, item.String())
	}

	return fmt.Sprintf("【%s】", strings.Join(strs, "，"))
}
