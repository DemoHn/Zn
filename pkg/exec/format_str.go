package exec

import (
	"fmt"
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

// fmtStack.fmtType
const (
	fmtTypeLiteral   int = 1
	fmtTypeFormatter int = 2
)

func formatString(formatStr *value.String, params *value.Array) (*value.String, error) {
	const (
		sBegin   = 1
		sLiteral = 2
		sFormat  = 3
	)
	var formatStrRune = []rune(formatStr.GetValue())

	///// #1. generate format stack
	// repeat [fmtType1, startIdx1, endIdx1, fmtType2, startIdx2, endIdx2, ...]
	var fmtStack = []int{}
	var state = sBegin
	var formatterCount = 0

	for idx, ch := range formatStrRune {
		switch ch {
		case '{':
			switch state {
			case sBegin:
				state = sFormat
				fmtStack = append(fmtStack, []int{fmtTypeFormatter, idx + 1}...)
			case sLiteral:
				state = sFormat
				fmtStack = append(fmtStack, []int{idx, fmtTypeFormatter, idx + 1}...)
			default:
				return nil, zerr.InvalidFmtTemplate(formatStr.String())
			}
		case '}':
			switch state {
			case sFormat:
				formatterCount += 1
				state = sBegin
				fmtStack = append(fmtStack, idx)
			default:
				return nil, zerr.InvalidFmtTemplate(formatStr.String())
			}
		default:
			switch state {
			case sBegin:
				state = sLiteral
				fmtStack = append(fmtStack, []int{fmtTypeLiteral, idx}...)
			}
		}
	}

	// append last char as last endIdx
	if state == sLiteral {
		fmtStack = append(fmtStack, len(formatStrRune))
	}

	//// since fmtStack = [fmtType, startIdx, endIdx] * N
	if len(fmtStack)%3 != 0 {
		return nil, zerr.InvalidFmtTemplate(formatStr.String())
	}

	paramElemList := params.GetValue()
	if len(paramElemList) != formatterCount {
		return nil, zerr.UnmatchFmtParams(formatStr.String())
	}

	///// #2. fill string
	paramElemIdx := 0
	fmtRuneList := []string{}

	for i := 0; i < len(fmtStack); i = i + 3 {
		var fmtType = fmtStack[i]
		var startIdx = fmtStack[i+1]
		var endIdx = fmtStack[i+2]

		formatter := string(formatStrRune[startIdx:endIdx])
		if fmtType == fmtTypeLiteral {
			fmtRuneList = append(fmtRuneList, formatter)
		} else if fmtType == fmtTypeFormatter {
			str, err := elementToString(formatter, paramElemList[paramElemIdx])
			if err != nil {
				return nil, err
			}

			fmtRuneList = append(fmtRuneList, str)
			paramElemIdx += 1
		}
	}

	return value.NewString(strings.Join(fmtRuneList, "")), nil
}

/*
* Formatter Rules:

 1. start with '#' - only Numbers are allowed to format
    1a. '#' -> format numbers with 6 significant digits (SD). The precise rules are as follows (copied from Python's `%g` format):

    suppose that the result formatted with presentation type 'e' 6 SD would have exponent exp. Then, if -4 <= exp < 6, the number is formatted with presentation type 'f' and precision 5-exp. Otherwise, the number is formatted with presentation type 'e' and 6 SD.

    Example:
    1234 --> "1234"
    12.3456789 --> "12.3457"
    1b. '#.N' -> format numbers, where N is the number of digits after the decimal point.

    1c. '#+' -> format numbers, add a '+' sign for positive numbers and 0
*/
func elementToString(formatter string, elem r.Element) (string, error) {
	if formatter == "" {
		switch elem.(type) {
		case *value.String, *value.Number, *value.Bool, *value.Array, *value.HashMap, *value.Null:
			return elem.String(), nil
		default:
			return "", zerr.InvalidParamType("")
		}
	}

	// if formatter starts from #
	if strings.HasPrefix(formatter, "#") {
		num, ok := elem.(*value.Number)
		if !ok {
			return "", zerr.NewErrorSLOT("格式化字符串只能用于数字")
		}

		return parseNumberFormatter(formatter[1:], num)
	}

	return "", zerr.NewErrorSLOT("无效的格式化字符串")
}

func parseNumberFormatter(formatter string, value *value.Number) (string, error) {
	// formatter: [+][.precision][E|%]
	const (
		sBegin          = 1
		sPositiveSign   = 2
		sFixedSign      = 3
		sScientificSign = 4
		sPercentSign    = 5
	)
	var (
		numFixedPrecision = 0
		flagPositive      = false
		flagFixed         = false
		flagScientific    = false
		flagPercent       = false
	)

	var state = sBegin

	// 1. parse formatter
	for _, ch := range formatter {
		switch ch {
		case '+':
			switch state {
			case sBegin:
				state = sPositiveSign
				flagPositive = true
			default:
				return "", zerr.NewErrorSLOT("无效的格式化字符串")
			}
		case '.':
			switch state {
			case sBegin, sPositiveSign:
				state = sFixedSign
				flagFixed = true
			default:
				return "", zerr.NewErrorSLOT("无效的格式化字符串")
			}
		case 'E':
			switch state {
			case sBegin, sPositiveSign, sFixedSign:
				state = sScientificSign
				flagScientific = true
			default:
				return "", zerr.NewErrorSLOT("无效的格式化字符串")
			}
		case '%':
			switch state {
			case sBegin, sPositiveSign, sFixedSign:
				state = sPercentSign
				flagPercent = true
			default:
				return "", zerr.NewErrorSLOT("无效的格式化字符串")
			}
		default:
			if ch >= '0' && ch <= '9' {
				switch state {
				case sFixedSign:
					numFixedPrecision = numFixedPrecision*10 + int(ch-'0')
				default:
					return "", zerr.NewErrorSLOT("无效的格式化字符串")
				}
			} else {
				return "", zerr.NewErrorSLOT("无效的格式化字符串")
			}
		}
	}

	// 2. get format number string (e.g. "%+.1E" / "%+.1f" / "%.6g")
	fmtStr := "%"
	if flagPositive {
		fmtStr += "+"
	}
	if flagFixed {
		fmtStr += fmt.Sprintf(".%d", numFixedPrecision)
	}
	// add E / f / .6g
	if flagScientific {
		fmtStr += "E"
	} else if flagFixed {
		fmtStr += "f"
	} else { // default case: NO FIXED & NO SCIENTIFIC_FLAG
		fmtStr += ".6g"
	}

	// 3. stringify number
	if flagPercent { // multiply 100 for percentage, then add "%"
		return fmt.Sprintf(fmtStr, value.GetValue()*100) + "%", nil
	}
	return fmt.Sprintf(fmtStr, value.GetValue()), nil
}
