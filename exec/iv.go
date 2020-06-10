package exec

import "github.com/DemoHn/Zn/error"

// ZnIV - Zn Intermediate Value
type ZnIV interface {
	// Reduce - reduce an IV to a real ZnValue
	// NOTICE: results may differ from whether it's on LHS or RHS
	Reduce(input ZnValue, lhs bool) (ZnValue, *error.Error)
}

// ZnArrayIV - a structure for intermediate-expression of an array or a hashmap
//
// 'IV' just stands for 'intermediate value'
// For an IV, its value can be both retrieved or set, the only difference
// is whether on the left side or right side.
//
// For example:
// 令 B 为 【10，20，30】#0  => when IV is on RHS, it will assign the value (10) to variable B;
// 【10，20，30】#0 为 75   => when IV is on LHS, set the 0-th slot of array to 75
//
type ZnArrayIV struct {
	List  *ZnArray
	Index *ZnDecimal
}

// ZnHashMapIV - similar to ZnArrayIV, see above for details
type ZnHashMapIV struct {
	List  *ZnHashMap
	Index *ZnString
}

// Reduce -
func (iv *ZnArrayIV) Reduce(input ZnValue, lhs bool) (ZnValue, *error.Error) {
	// check data
	idx, err := iv.Index.asInteger()
	if err != nil {
		return nil, error.InvalidExprType("integer")
	}
	if idx < 0 || idx >= len(iv.List.Value) {
		return nil, error.IndexOutOfRange()
	}

	// iv is on LHS, that means its index will be assigned from a new value
	if lhs == true {
		iv.List.Value[idx] = input
		return input, nil
	}
	return iv.List.Value[idx], nil
}

// Reduce -
func (iv *ZnHashMapIV) Reduce(input ZnValue, lhs bool) (ZnValue, *error.Error) {
	// check data
	key := iv.Index.Value
	vr, ok := iv.List.Value[key]
	if !ok {
		return nil, error.IndexKeyNotFound(key)
	}

	if lhs == true {
		iv.List.Value[key] = input
		return input, nil
	}
	return vr, nil
}
