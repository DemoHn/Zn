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

		if fmtType == fmtTypeLiteral {
			fmtRuneList = append(fmtRuneList, string(formatStrRune[startIdx:endIdx]))
		} else if fmtType == fmtTypeFormatter {
			fmtRuneList = append(fmtRuneList, elementToString(
				formatStrRune[startIdx:endIdx],
				paramElemList[paramElemIdx],
			))

			paramElemIdx += 1
		}
	}

	return value.NewString(strings.Join(fmtRuneList, "")), nil
}

func elementToString(formatter []rune, elem r.Element) string {
	switch v := elem.(type) {
	case *value.String:
		return v.GetValue()
	case *value.Number:
		return fmt.Sprintf("%.2f", v.GetValue())
	default:
		return ""
	}
}
