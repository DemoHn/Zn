package value

import (
	"fmt"

	zerr "github.com/DemoHn/Zn/pkg/error"
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
//  1. Root - root is a r.Value Object from which can set/get members.
//  2. Member - a string / int that represents for the "key" of root object
//
// IV could be reduced, either yields a result or the member value be chaned.
// When IV on Left-Hand Side (LHS), e.g. A之B = 1, it could be reduced by ReduceLHS(), where memeber is updated;
// WHen IV on Right-Hand Side (RHS), e.g. 令C = A之B, it could be reduced by ReduceRHS(), and transfer result value to left-hand side.
type IV struct {
	// reduceType - value type
	reduceType uint8
	// root Object
	root r.Element
	// member value
	member string
	// index is used ONLY when ivType = IVTypeArray,
	// it's a "Member", but in integer type.
	index int
}

// NewMemberIV -
func NewMemberIV(root r.Element, member string) *IV {
	return &IV{
		reduceType: IVTypeMember,
		root:       root,
		member:     member,
	}
}

// NewArrayIV -
func NewArrayIV(root r.Element, index int) *IV {
	return &IV{
		reduceType: IVTypeArray,
		root:       root,
		index:      index,
	}
}

// NewHashMapIV -
func NewHashMapIV(root r.Element, member string) *IV {
	return &IV{
		reduceType: IVTypeHashMap,
		root:       root,
		member:     member,
	}
}

// ReduceLHS - Reduce IV to value when IV on left-hand side
// usually for setters
func (iv *IV) ReduceLHS(input r.Element) error {
	switch iv.reduceType {
	case IVTypeArray:
		arr, ok := iv.root.(*Array)
		if !ok {
			return zerr.InvalidExprType("array")
		}
		// in Zn, array index starts from 1, while in Golang, array index starts from 0,
		// thus, X#2 (Zn) <--> x[2-1] --> x[1] (Go)
		realIndex := iv.index - 1
		if realIndex < 0 || realIndex >= len(arr.value) {
			return zerr.IndexOutOfRange()
		}
		// set array value
		arr.value[realIndex] = input
		return nil
	case IVTypeHashMap:
		hm, ok := iv.root.(*HashMap)
		if !ok {
			return zerr.InvalidExprType("hashmap")
		}
		// use AppendKVPair instead of `hm.value[iv.member] = input` directly
		hm.AppendKVPair(KVPair{iv.member, input})
		return nil
	case IVTypeMember:
		return iv.root.SetProperty(c, iv.member, input)
	}
	return zerr.UnexpectedCase("IVReduceType", fmt.Sprintf("%d", iv.reduceType))
}

// ReduceRHS -
func (iv *IV) ReduceRHS(c *r.Context) (r.Element, error) {
	switch iv.reduceType {
	case IVTypeArray:
		arr, ok := iv.root.(*Array)
		if !ok {
			return nil, zerr.InvalidExprType("array")
		}
		// in Zn, array index starts from 1, while in Golang, array index starts from 0,
		// thus, X#2 (Zn) <--> x[2-1] --> x[1] (Go)
		realIndex := iv.index - 1
		if realIndex < 0 || realIndex >= len(arr.value) {
			return nil, zerr.IndexOutOfRange()
		}
		// set array value
		return arr.value[realIndex], nil
	case IVTypeHashMap:
		hm, ok := iv.root.(*HashMap)
		if !ok {
			return nil, zerr.InvalidExprType("hashmap")
		}
		result, ok := hm.value[iv.member]
		if !ok {
			return nil, zerr.IndexKeyNotFound(iv.member)
		}
		return result, nil
	case IVTypeMember:
		return iv.root.GetProperty(c, iv.member)
	}
	return nil, zerr.UnexpectedCase("IVReduceType", fmt.Sprintf("%d", iv.reduceType))
}
