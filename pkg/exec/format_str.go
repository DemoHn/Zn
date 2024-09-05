package exec

import (
	"fmt"
	"strconv"
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
				return nil, zerr.NewErrorSLOT("解析模板异常")
			}
		case '}':
			switch state {
			case sFormat:
				formatterCount += 1
				state = sBegin
				fmtStack = append(fmtStack, idx)
			default:
				return nil, zerr.NewErrorSLOT("解析模板异常")
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
		return nil, zerr.NewErrorSLOT("解析模板异常")
	}

	paramElemList := params.GetValue()
	if len(paramElemList) != formatterCount {
		return nil, zerr.NewErrorSLOT("参数与模板数量不匹配")
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

/** Formatter Rules:

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
	if len(formatter) == 0 {
		switch elem.(type) {
		case *value.String, *value.Number, *value.Bool, *value.Array, *value.HashMap, *value.Null:
			return value.StringifyValue(elem), nil
		default:
			return "", zerr.NewErrorSLOT("无效的元素类型")
		}
	}

	// if formatter starts from #
	if strings.HasPrefix(formatter, "#") {
		num, ok := elem.(*value.Number)
		if !ok {
			return "", zerr.NewErrorSLOT("格式化字符串只能用于数字")
		}

		numValue := num.GetValue()
		switch {
		case formatter == "#":
			return fmt.Sprintf("%.6g", numValue), nil
		case formatter == "#+":
			return fmt.Sprintf("%+.6g", numValue), nil
		case strings.HasPrefix(formatter, "#."):
			precision, err := strconv.Atoi(formatter[2:])
			if err != nil {
				return "", zerr.NewErrorSLOT("无效的格式化字符串")
			}
			formatStr := fmt.Sprintf("%%.%df", precision)
			return fmt.Sprintf(formatStr, numValue), nil
		case strings.HasPrefix(formatter, "#+."):
			precision, err := strconv.Atoi(formatter[3:])
			if err != nil {
				return "", zerr.NewErrorSLOT("无效的格式化字符串")
			}
			formatStr := fmt.Sprintf("%%+.%df", precision)
			return fmt.Sprintf(formatStr, numValue), nil
		default:
			return "", zerr.NewErrorSLOT("无效的格式化字符串")
		}
	}

	return "", zerr.NewErrorSLOT("无效的格式化字符串")
}
