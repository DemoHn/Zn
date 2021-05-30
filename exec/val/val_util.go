package val

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

type compareVerb uint8

// Define compareVerbs, for details of each verb, check the following comments
// on compareValues() function.
const (
	CmpEq compareVerb = 1
	CmpLt compareVerb = 2
	CmpGt compareVerb = 3
)

// compareValues - some ZnValues are comparable from specific types of right value
// otherwise it will throw error.
//
// There are three types of compare verbs (actions): Eq, Lt and Gt.
//
// Eq - compare if two values are "equal". Usually there are two rules:
// 1. types of left and right value are same. A number MUST BE equals to a number, that means
// (string) “2” won't be equals to (number) 2;
// 2. each items SHOULD BE identical, even for composited types (i.e. array, hashmap)
//
// Lt - for two decimals ONLY. If leftValue < rightValue.
//
// Gt - for two decimals ONLY. If leftValue > rightValue.
//
func CompareValues(left ctx.Value, right ctx.Value, verb compareVerb) (bool, *error.Error) {
	switch vl := left.(type) {
	case *Null:
		if _, ok := right.(*Null); ok {
			return true, nil
		}
		return false, nil
	case *Decimal:
		// compare right value - decimal only
		if vr, ok := right.(*Decimal); ok {
			r1, r2 := rescalePair(vl, vr)
			cmpResult := false
			switch verb {
			case CmpEq:
				cmpResult = (r1.co.Cmp(r2.co) == 0)
			case CmpLt:
				cmpResult = (r1.co.Cmp(r2.co) < 0)
			case CmpGt:
				cmpResult = (r1.co.Cmp(r2.co) > 0)
			default:
				return false, error.UnExpectedCase("比较原语", strconv.Itoa(int(verb)))
			}
			return cmpResult, nil
		}
		// if vert == CmbEq and rightValue is not decimal type
		// then return `false` directly
		if verb == CmpEq {
			return false, nil
		}
		return false, error.InvalidCompareRType("decimal")
	case *String:
		// Only CmpEq is valid for comparison
		if verb != CmpEq {
			return false, error.InvalidCompareLType("decimal", "string", "bool", "array", "hashmap")
		}
		// compare right value - string only
		if vr, ok := right.(*String); ok {
			cmpResult := (strings.Compare(vl.value, vr.value) == 0)
			return cmpResult, nil
		}
		return false, nil
	case *Bool:
		if verb != CmpEq {
			return false, error.InvalidCompareLType("decimal", "string", "bool", "array", "hashmap")
		}
		// compare right value - bool only
		if vr, ok := right.(*Bool); ok {
			cmpResult := vl.value == vr.value
			return cmpResult, nil
		}
		return false, nil
	case *Array:
		if verb != CmpEq {
			return false, error.InvalidCompareLType("decimal", "string", "bool", "array", "hashmap")
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
			return false, error.InvalidCompareLType("decimal", "string", "bool", "array", "hashmap")
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
	return false, error.InvalidCompareLType("decimal", "string", "bool", "array", "hashmap")
}

// DuplicateValue - deepcopy values' structure, including bool, string, decimal, array, hashmap
// for function or object or null, pass the original reference instead.
// This is due to the 'copycat by default' policy
func DuplicateValue(in ctx.Value) ctx.Value {
	switch v := in.(type) {
	case *Bool:
		return NewBool(v.value)
	case *String:
		return NewString(v.value)
	case *Decimal:
		x := new(big.Int)
		return &Decimal{
			co:  x.Set(v.co),
			exp: v.exp,
		}
	case *Null:
		return in // no need to copy since all "NULL" values are same
	case *Array:
		newArr := []ctx.Value{}
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

// StringifyValue - yield a string from ctx.Value
func StringifyValue(value ctx.Value) string {
	switch v := value.(type) {
	case *String:
		return fmt.Sprintf("「%s」", v.value)
	case *Decimal:
		return v.String()
	case *Array:
		strs := []string{}
		for _, item := range v.value {
			strs = append(strs, StringifyValue(item))
		}

		return fmt.Sprintf("【%s】", strings.Join(strs, "，"))
	case *Bool:
		data := "真"
		if !v.value {
			data = "假"
		}
		return data
	case *Function:
		return "[方法]"
	case *Null:
		return "空"
	case *Object:
		return "[对象]"
	case *HashMap:
		strs := []string{}
		for _, key := range v.keyOrder {
			value := v.value[key]
			strs = append(strs, fmt.Sprintf("%s == %s", key, StringifyValue(value)))
		}
		return fmt.Sprintf("【%s】", strings.Join(strs, "，"))
	}
	return ""
}

//// param validators

// validateExactParams is a helper function that asserts input params type where each parameter
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
func validateExactParams(values []ctx.Value, typeStr ...string) *error.Error {
	if len(values) != len(typeStr) {
		return error.ExactParamsError(len(typeStr))
	}
	for idx, v := range values {
		if err := validateOneParam(v, typeStr[idx]); err != nil {
			return err
		}
	}
	return nil
}

// validateLeastParams is a helper function similar to validateExactParams(),
// however, it handles the situation that target params are variadic, like
// [int, int, string...], the length of target params varies from 0~N.
//
// To validate variadic params, we introduce wildcard (* or +) at the end of typeStr to
// mark variadic part, like "string+", "bool*"
//
// e.g.:
// ["decimal", "string+"] means the FIRST param is a decimal, and the FOLLOWING params
// are all strings (must have ONE string param)
//
// ["decimal", "bool", "string*"] means the FIRST param is a decimal, the SECOND param is a bool, and the FOLLOWING params
// are all strings (allow 0 string params)
func validateLeastParams(values []ctx.Value, typeStr ...string) *error.Error {
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
				return error.NewErrorSLOT("通配符需要至少一个参数")
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

// validateAllParams doesn't limit the length of input values; instead, it requires all the parameters
// to have same value type denoted by `typeStr`
// e.g. validateAllParams([]Value{“1”, “2”, “3”}, "string")
func validateAllParams(values []ctx.Value, typeStr string) *error.Error {
	for _, v := range values {
		if err := validateOneParam(v, typeStr); err != nil {
			return err
		}
	}
	return nil
}

func validateOneParam(v ctx.Value, typeStr string) *error.Error {
	valid := true
	switch typeStr {
	case "decimal":
		if _, ok := v.(*Decimal); !ok {
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
		return error.InvalidParamType(typeStr)
	}
	return nil
}

//// Program Executors

// （显示） 方法的执行逻辑
var DisplayExecutor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	// display format string
	var items = []string{}

	for _, param := range params {
		// if param is a string, display its value (without 「 」 quotes) directly
		if str, ok := param.(*String); ok {
			items = append(items, str.String())
		} else {
			items = append(items, StringifyValue(param))
		}
	}
	fmt.Printf("%s\n", strings.Join(items, " "))
	return NewNull(), nil
}

// （递增）方法的执行逻辑
var AddValueExecutor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	var decimals = []*Decimal{}

	if err := validateLeastParams(params, "decimal+"); err != nil {
		return nil, err
	}
	// validate types
	for _, param := range params {
		vparam := param.(*Decimal)
		decimals = append(decimals, vparam)
	}

	if len(decimals) == 1 {
		return decimals[0], nil
	}

	sum := decimals[0].Add(decimals[1:]...)
	return sum, nil
}

// （递减）方法的执行逻辑
var SubValueExecutor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	var decimals = []*Decimal{}

	if err := validateLeastParams(params, "decimal+"); err != nil {
		return nil, err
	}
	// validate types
	for _, param := range params {
		vparam := param.(*Decimal)
		decimals = append(decimals, vparam)
	}

	if len(decimals) == 1 {
		return decimals[0], nil
	}

	sum := decimals[0].Sub(decimals[1:]...)
	return sum, nil
}

var MulValueExecutor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	var decimals = []*Decimal{}

	if err := validateLeastParams(params, "decimal+"); err != nil {
		return nil, err
	}
	// validate types
	for _, param := range params {
		vparam := param.(*Decimal)
		decimals = append(decimals, vparam)
	}

	if len(decimals) == 1 {
		return decimals[0], nil
	}

	sum := decimals[0].Mul(decimals[1:]...)
	return sum, nil
}

var DivValueExecutor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	var decimals = []*Decimal{}

	if err := validateLeastParams(params, "decimal+"); err != nil {
		return nil, err
	}
	// validate types
	for _, param := range params {
		vparam := param.(*Decimal)
		decimals = append(decimals, vparam)
	}
	if len(decimals) == 1 {
		return decimals[0], nil
	}

	res, err := decimals[0].Div(decimals[1:]...)
	return res, err
}

var ProbeExecutor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	if len(params) != 2 {
		return nil, error.ExactParamsError(2)
	}

	vtag, ok := params[0].(*String)
	if !ok {
		return nil, error.InvalidParamType("string")
	}
	// add probe data to log
	c.GetProbe().AddLog(vtag.value, params[1])
	return params[1], nil
}
