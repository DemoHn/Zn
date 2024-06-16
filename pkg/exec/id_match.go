package exec

import (
	"strconv"
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
)

/**
Unlike most languages, the token _Identifier_ (type = 5) has a large variety of
representation types, including variable name, number, datetime, unit value, even
currency value!

This module aims to breakdown the identifier token to match more specific idTypes for
further code execution usages: For example, if an identifier is matched as a "number", then
it's impossible to be assigned to another value!

So far, we support only two idTypes, which corresponds to its specific format; For other formats of this identifier token, we throw an error as "SemanticError"!

1. idNumber

ID Format: all char shall be a number ('0-9') or decimal point (.) or exponential indicator (E / *10^) or number flag (+/-); and all char follows the Number Construction Regexp to make it represent ONE number!

Example:
```
+12.5
114514
0.12
-12.8*10^15
-6.4
````

INVALID Example: (will throw SemanticError)
```
-1/2
2.3.5
++18
```

2. idName

ID Format: the leading chars SHOULD NOT be numbers chars (any valid chars but not 0-9)

Example:
```
公交车站-数目
自税/包税计费重
IR80
数值+
非常7+1*
-某不公开属性-
_window_handler2
```
*/

func MatchIDType(id *syntax.ID) (r.IDType, error) {
	idStr := id.GetLiteral()
	// #1. match number
	isNumber, err := tryParseNumber(id)
	if err != nil {
		return nil, err
	}
	if isNumber {
		return &r.IDNumber{
			Literal:  idStr,
			NumValue: parseIDNumberToFloat64(idStr),
		}, nil
	}

	// #2. match name
	return &r.IDName{
		Literal: idStr,
	}, nil
}

func MatchIDName(id *syntax.ID) (*r.IDName, error) {
	v, err := MatchIDType(id)
	if err != nil {
		return nil, err
	}

	if idName, ok := v.(*r.IDName); ok {
		return idName, nil
	} else {
		return nil, zerr.IDNameONLY(id.GetLiteral())
	}
}

func MatchIDNumber(id *syntax.ID) (*r.IDNumber, error) {
	v, err := MatchIDType(id)
	if err != nil {
		return nil, err
	}

	if idNumber, ok := v.(*r.IDNumber); ok {
		return idNumber, nil
	} else {
		return nil, zerr.IDNumberONLY(id.GetLiteral())
	}
}

// regex: ^[-+]?[0-9]+\.?[0-9]+((([eE][-+])|(\*(10)?\^[-+]?))[0-9]+)?$
// ref: https://github.com/DemoHn/Zn/issues/4
func tryParseNumber(id *syntax.ID) (bool, error) {
	charArr := []rune(id.GetLiteral())

	// hand-written regex parser
	// ref: https://cyberzhg.github.io/toolbox/min_dfa?regex=Rj9QP0QqLj9EKygoKEVQKXwocygxMCk/dVA/KSlEKyk/
	// hand-drawn min-DFA:
	// https://github.com/DemoHn/Zn/issues/6
	const (
		sBegin      = 1
		sDot        = 2
		sIntEnd     = 3
		sIntPMFlag  = 5
		sDotDecEnd  = 6
		sEFlag      = 7
		sSFlag      = 8
		sExpPMFlag  = 9
		sSciI       = 10
		sSciEndFlag = 11
		sExpEnd     = 12
		sSciII      = 13
	)
	var state = sBegin
	var endStates = []int{sIntEnd, sDotDecEnd, sExpEnd}
	var parsedChars = 0
	for _, ch := range charArr {
		switch ch {
		case 'e', 'E':
			switch state {
			case sDotDecEnd, sIntEnd:
				state = sEFlag
			default:
				goto end
			}
		case '.':
			switch state {
			case sIntEnd:
				state = sDot
			default:
				goto end
			}
		case '-', '+':
			switch state {
			case sBegin:
				state = sIntPMFlag
			case sEFlag, sSciEndFlag:
				state = sExpPMFlag
			default:
				goto end
			}
		case '*':
			switch state {
			case sDotDecEnd, sIntEnd:
				state = sSFlag
			default:
				goto end
			}
		case '1':
			switch state {
			case sSFlag:
				state = sSciI
				// same with other numbers
			case sBegin, sIntEnd, sIntPMFlag:
				state = sIntEnd
			case sDot, sDotDecEnd:
				state = sDotDecEnd
			case sExpPMFlag, sSciEndFlag, sExpEnd:
				state = sExpEnd
			default:
				goto end
			}
		case '0':
			switch state {
			case sSciI:
				state = sSciII
			case sBegin, sIntEnd, sIntPMFlag:
				state = sIntEnd
			case sDot, sDotDecEnd:
				state = sDotDecEnd
			case sExpPMFlag, sSciEndFlag, sExpEnd:
				state = sExpEnd
			default:
				goto end
			}
		case '2', '3', '4', '5', '6', '7', '8', '9':
			switch state {
			case sBegin, sIntEnd, sIntPMFlag:
				state = sIntEnd
			case sDot, sDotDecEnd:
				state = sDotDecEnd
			case sExpPMFlag, sSciEndFlag, sExpEnd:
				state = sExpEnd
			default:
				goto end
			}
		case '^':
			switch state {
			case sSFlag, sSciII:
				state = sSciEndFlag
			default:
				goto end
			}
		default:
			goto end
		}
		parsedChars += 1
	}

end:
	// NO chars parsed, maybe the token is idName (instead of idNumber)
	if parsedChars == 0 {
		return false, nil
	}
	// only + or - is parsed
	if state == sIntPMFlag {
		return false, nil
	}
	// Parsing flow NOT FINISH: e.g. `15.` got error where no number after decimal point
	if !syntax.ContainsInt(state, endStates) {
		return false, zerr.InvalidIDFormat(id.GetLiteral())
	}
	// There're still characters after number: e.g. `128kg`
	if parsedChars < len(charArr) {
		return false, zerr.InvalidIDFormat(id.GetLiteral())
	}

	return true, nil
}

func parseIDNumberToFloat64(idStr string) float64 {
	v := strings.Replace(idStr, "*^", "e", 1)
	v = strings.Replace(v, "*10^", "e", 1)

	f, _ := strconv.ParseFloat(v, 64)
	return f
}
