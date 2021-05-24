package val

import (
	"regexp"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
)

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
			break // "break" HERE is to break the outer for-loop!
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
	if valid == false {
		return error.InvalidParamType(typeStr)
	}
	return nil
}
