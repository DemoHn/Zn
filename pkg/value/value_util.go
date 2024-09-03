package value

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

// Define compareVerbs, for details of each verb, check the following comments
// on compareValues() function.
const (
	CmpEq uint8 = 1
	CmpLt uint8 = 2
	CmpGt uint8 = 3
)

// compareValues - some ZnValues are comparable from specific types of right value
// otherwise it will throw zerr.
//
// There are three types of compare verbs (actions): Eq, Lt and Gt.
//
// Eq - compare if two values are "equal". Usually there are two rules:
// 1. types of left and right value are same. A number MUST BE equals to a number, that means
// (string) “2” won't be equals to (number) 2;
// 2. each item SHOULD BE identical, even for composited types (i.e. array, hashmap)
//
// Lt - for two decimals ONLY. If leftValue < rightValue.
//
// Gt - for two decimals ONLY. If leftValue > rightValue.
//
func CompareValues(left r.Element, right r.Element, verb uint8) (bool, error) {
	switch vl := left.(type) {
	case *Null:
		if _, ok := right.(*Null); ok {
			return true, nil
		}
		return false, nil
	case *Number:
		// compare right value - decimal only
		if vr, ok := right.(*Number); ok {
			cmpResult := false
			switch verb {
			case CmpEq:
				cmpResult = vl.value == vr.value
			case CmpLt:
				cmpResult = vl.value < vr.value
			case CmpGt:
				cmpResult = vl.value > vr.value
			default:
				return false, zerr.UnexpectedCase("比较类型", strconv.Itoa(int(verb)))
			}
			return cmpResult, nil
		}
		// if vert == CmbEq and rightValue is not decimal type
		// then return `false` directly
		if verb == CmpEq {
			return false, nil
		}
		return false, zerr.InvalidCompareRType("number")
	case *String:
		// Only CmpEq is valid for comparison
		if verb != CmpEq {
			return false, zerr.InvalidCompareLType("number", "string", "bool", "array", "hashmap")
		}
		// compare right value - string only
		if vr, ok := right.(*String); ok {
			cmpResult := strings.Compare(vl.value, vr.value) == 0
			return cmpResult, nil
		}
		return false, nil
	case *Bool:
		if verb != CmpEq {
			return false, zerr.InvalidCompareLType("number", "string", "bool", "array", "hashmap")
		}
		// compare right value - bool only
		if vr, ok := right.(*Bool); ok {
			cmpResult := vl.value == vr.value
			return cmpResult, nil
		}
		return false, nil
	case *Array:
		if verb != CmpEq {
			return false, zerr.InvalidCompareLType("number", "string", "bool", "array", "hashmap")
		}

		if vr, ok := right.(*Array); ok {
			if len(vl.value) != len(vr.value) {
				return false, nil
			}
			// cmp each item
			for idx := range vl.value {
				cmpVal, err := CompareValues(vl.value[idx], vr.value[idx], CmpEq)
				if err != nil {
					return false, err
				}
				// break the loop only when cmpVal = false
				if !cmpVal {
					return false, nil
				}
			}
			return true, nil
		}
		return false, nil
	case *HashMap:
		if verb != CmpEq {
			return false, zerr.InvalidCompareLType("number", "string", "bool", "array", "hashmap")
		}

		if vr, ok := right.(*HashMap); ok {
			if len(vl.value) != len(vr.value) {
				return false, nil
			}
			// cmp each item
			for idx := range vl.value {
				// ensure the key exists on vr
				vrr, ok := vr.value[idx]
				if !ok {
					return false, nil
				}
				cmpVal, err := CompareValues(vl.value[idx], vrr, CmpEq)
				if err != nil {
					return false, err
				}
				return cmpVal, nil
			}
			return true, nil
		}
		return false, nil
	}
	return false, zerr.InvalidCompareLType("number", "string", "bool", "array", "hashmap")
}

// DuplicateValue - deepcopy values' structure, including bool, string, decimal, array, hashmap
// for function or object or null, pass the original reference instead.
// This is due to the 'copycat by default' policy
func DuplicateValue(in r.Element) r.Element {
	switch v := in.(type) {
	case *Bool:
		return NewBool(v.value)
	case *String:
		return NewString(v.value)
	case *Number:
		return NewNumber(v.value)
	case *Null:
		return in // no need to copy since all "NULL" values are same
	case *Array:
		newArr := []r.Element{}
		for _, val := range v.value {
			newArr = append(newArr, DuplicateValue(val))
		}
		return NewArray(newArr)
	case *HashMap:
		kvPairs := []KVPair{}
		for _, key := range v.keyOrder {
			dupVal := DuplicateValue(v.value[key])
			kvPairs = append(kvPairs, KVPair{key, dupVal})
		}
		return NewHashMap(kvPairs)
	case *Function: // function itself is immutable, so return directly
		return in
	case *Object: // we don't copy object value at all
		return in
	}
	return in
}

// StringifyValue - yield a string from r.Value
func StringifyValue(value r.Element) string {
	switch v := value.(type) {
	case *String:
		return fmt.Sprintf("「%s」", v.value)
	case *Number:
		return v.String()
	case *Array:
		strs := []string{}
		for _, item := range v.value {
			strs = append(strs, StringifyValue(item))
		}
		return fmt.Sprintf("【%s】", strings.Join(strs, "、"))
	case *Bool:
		data := "真"
		if !v.value {
			data = "假"
		}
		return data
	case *Function:
		return fmt.Sprintf("[某方法]")
	case *Null:
		return "空"
	case *Object:
		// show object as "[对象: <name>]"
		return fmt.Sprintf("[对象: %s]", v.GetObjectName())
	case *HashMap:
		strs := []string{}
		for _, key := range v.keyOrder {
			value := v.value[key]
			strs = append(strs, fmt.Sprintf("%s = %s", key, StringifyValue(value)))
		}
		return fmt.Sprintf("【%s】", strings.Join(strs, "，"))
	case *GoValue:
		// show go value as string like "[GoValue type:<tag>]"
		return fmt.Sprintf("[GoValue type:%s]", v.GetTag())
	}
	return ""
}

//// param validators

// ValidateExactParams is a helper function that asserts input params type where each parameter
// should exactly match the list of typeStr.
// valid typeStr are one of the followings:
//
// 1. decimal  --> *Decimal
// 2. string   --> *String
// 3. array    --> *Array
// 4. hashmap  --> *HashMap
// 5. bool     --> *Bool
// 6. object   --> *Object
// 7. function --> *Function
// 8. any      --> any value type
func ValidateExactParams(values []r.Element, typeStr ...string) error {
	if len(values) != len(typeStr) {
		return zerr.ExactParamsError(len(typeStr))
	}
	for idx, v := range values {
		if err := validateOneParam(v, typeStr[idx]); err != nil {
			return err
		}
	}
	return nil
}

// ValidateLeastParams is a helper function similar to ValidateExactParams(),
// however, it handles the situation that target params are variadic, like
// [int, int, string...], the length of target params varies from 0~N.
//
// To validate variadic params, we introduce wildcard (* or +) at the end of typeStr to
// mark variadic part, like "string+", "bool*"
//
// e.g.:
// ["number", "string+"] means the FIRST param is a decimal, and the FOLLOWING params
// are all strings (must have ONE more string param)
//
// ["number", "bool", "string*"] means the FIRST param is a decimal, the SECOND param is a bool, and the FOLLOWING params
// are all strings (allow 0 string params)
func ValidateLeastParams(values []r.Element, typeStr ...string) error {
Loop:
	for idx, t := range typeStr {
		// find if there's wildcard
		re := regexp.MustCompile(`(\w+)(\*|\+)?`)
		matches := re.FindStringSubmatch(t)
		// match: [_, name, wildcard]
		switch matches[2] {
		// matches 0 or more params
		case "*", "+":
			if matches[2] == "+" && idx > len(values) {
				return zerr.UnexpectedParamWildcard()
			}
			for i := idx; i < len(values); i++ {
				if err := validateOneParam(values[i], matches[1]); err != nil {
					return err
				}
			}
			break Loop
		default:
			if err := validateOneParam(values[idx], t); err != nil {
				return err
			}
		}
	}
	return nil
}

// ValidateAllParams doesn't limit the length of input values; instead, it requires all the parameters
// to have same value type denoted by `typeStr`
// e.g. ValidateAllParams([]Value{“1”, “2”, “3”}, "string")
func ValidateAllParams(values []r.Element, typeStr string) error {
	for _, v := range values {
		if err := validateOneParam(v, typeStr); err != nil {
			return err
		}
	}
	return nil
}

func validateOneParam(v r.Element, typeStr string) error {
	valid := true
	switch typeStr {
	case "number":
		if _, ok := v.(*Number); !ok {
			valid = false
		}
	case "string":
		if _, ok := v.(*String); !ok {
			valid = false
		}
	case "array":
		if _, ok := v.(*Array); !ok {
			valid = false
		}
	case "hashmap":
		if _, ok := v.(*HashMap); !ok {
			valid = false
		}
	case "bool":
		if _, ok := v.(*Bool); !ok {
			valid = false
		}
	case "object":
		if _, ok := v.(*Object); !ok {
			valid = false
		}
	case "function":
		if _, ok := v.(*Function); !ok {
			valid = false
		}
	case "any":
		valid = true
	}
	if !valid {
		return zerr.InvalidParamType(typeStr)
	}
	return nil
}
