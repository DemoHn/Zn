package val

import (
	"fmt"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

type ivTypeE uint8

// declare IVTypes
const (
	IVTypeArray   ivTypeE = 1
	IVTypeHashMap ivTypeE = 2
	IVTypeMember  ivTypeE = 3
)

// IV stands for Intermediate ctx.Value, which is a compound state of an expression.
// There're two elements inside an IV:
//   1. Root - root is a ctx.Value Object from which can set/get members.
//   2. Member - a string / int that represents for the "key" of root object
//
// IV could be reduced, either yields a result or the member value be chaned.
// When IV on Left-Hand Side (LHS), e.g. A之B 为 1, it could be reduced by ReduceLHS(), where memeber is updated;
// WHen IV on Right-Hand Side (RHS), e.g. 令C 为 A之B, it could be reduced by ReduceRHS(), and transfer result value to left-hand side.
type IV struct {
	// reduceType - value type
	reduceType ivTypeE
	// root Object
	root ctx.Value
	// member value
	member string
	// index is used ONLY when ivType = IVTypeArray,
	// it's a "Member", but in integer type.
	index int
}

// NewMemberIV -
func NewMemberIV(root ctx.Value, member string) *IV {
	return &IV{
		reduceType: IVTypeMember,
		root:       root,
		member:     member,
	}
}

// NewArrayIV -
func NewArrayIV(root ctx.Value, index int) *IV {
	return &IV{
		reduceType: IVTypeArray,
		root:       root,
		index:      index,
	}
}

// NewHashMapIV -
func NewHashMapIV(root ctx.Value, member string) *IV {
	return &IV{
		reduceType: IVTypeHashMap,
		root:       root,
		member:     member,
	}
}

// ReduceLHS - Reduce IV to value when IV on left-hand side
// usually for setters
func (iv *IV) ReduceLHS(c *ctx.Context, input ctx.Value) *error.Error {
	switch iv.reduceType {
	case IVTypeArray:
		arr, ok := iv.root.(*Array)
		if !ok {
			return error.InvalidExprType("array")
		}
		if iv.index < 0 || iv.index >= len(arr.value) {
			return error.IndexOutOfRange()
		}
		// set array value
		arr.value[iv.index] = input
		return nil
	case IVTypeHashMap:
		hm, ok := iv.root.(*HashMap)
		if !ok {
			return error.InvalidExprType("hashmap")
		}
		hm.value[iv.member] = input
	case IVTypeMember:
		return iv.root.SetProperty(c, iv.member, input)
	}
	return error.UnExpectedCase("IVReduceType", fmt.Sprintf("%d", iv.reduceType))
}

// ReduceRHS -
func (iv *IV) ReduceRHS(c *ctx.Context) (ctx.Value, *error.Error) {
	switch iv.reduceType {
	case IVTypeArray:
		arr, ok := iv.root.(*Array)
		if !ok {
			return nil, error.InvalidExprType("array")
		}
		if iv.index < 0 || iv.index >= len(arr.value) {
			return nil, error.IndexOutOfRange()
		}
		// set array value
		return arr.value[iv.index], nil
	case IVTypeHashMap:
		hm, ok := iv.root.(*HashMap)
		if !ok {
			return nil, error.InvalidExprType("hashmap")
		}
		result, ok := hm.value[iv.member]
		if !ok {
			return nil, error.IndexKeyNotFound(iv.member)
		}
		return result, nil
	case IVTypeMember:
		return iv.root.GetProperty(c, iv.member)
	}
	return nil, error.UnExpectedCase("IVReduceType", fmt.Sprintf("%d", iv.reduceType))
}
