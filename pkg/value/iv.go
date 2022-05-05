package value

import (
	"fmt"
	r "github.com/DemoHn/Zn/pkg/runtime"
)


// declare IVTypes
const (
	IVTypeArray   uint8 = 1
	IVTypeHashMap uint8 = 2
	IVTypeMember  uint8 = 3
)

// IV stands for Intermediate r.Value, which is a compound state of an expression.
// There're two elements inside an IV:
//   1. Root - root is a r.Value Object from which can set/get members.
//   2. Member - a string / int that represents for the "key" of root object
//
// IV could be reduced, either yields a result or the member value be chaned.
// When IV on Left-Hand Side (LHS), e.g. A之B 为 1, it could be reduced by ReduceLHS(), where memeber is updated;
// WHen IV on Right-Hand Side (RHS), e.g. 令C 为 A之B, it could be reduced by ReduceRHS(), and transfer result value to left-hand side.
type IV struct {
	// reduceType - value type
	reduceType uint8
	// root Object
	root r.Value
	// member value
	member string
	// index is used ONLY when ivType = IVTypeArray,
	// it's a "Member", but in integer type.
	index int
}

// NewMemberIV -
func NewMemberIV(root r.Value, member string) *IV {
	return &IV{
		reduceType: IVTypeMember,
		root:       root,
		member:     member,
	}
}

// NewArrayIV -
func NewArrayIV(root r.Value, index int) *IV {
	return &IV{
		reduceType: IVTypeArray,
		root:       root,
		index:      index,
	}
}

// NewHashMapIV -
func NewHashMapIV(root r.Value, member string) *IV {
	return &IV{
		reduceType: IVTypeHashMap,
		root:       root,
		member:     member,
	}
}

// ReduceLHS - Reduce IV to value when IV on left-hand side
// usually for setters
func (iv *IV) ReduceLHS(c *r.Context, input r.Value) error {
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
		return nil
	case IVTypeMember:
		return iv.root.SetProperty(c, iv.member, input)
	}
	return error.UnExpectedCase("IVReduceType", fmt.Sprintf("%d", iv.reduceType))
}

// ReduceRHS -
func (iv *IV) ReduceRHS(c *r.Context) (r.Value, error) {
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
