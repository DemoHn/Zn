package error

import "fmt"

// InvalidIndentType -
func InvalidIndentType(expect uint8, got uint8) *Error {
	findName := func(indentType uint8) string {
		name := "「空格」"
		if expect == uint8(9) { // TAB
			name = "「TAB」"
		}
		return name
	}

	return &Error{
		Code:    0x2021,
		Message: fmt.Sprintf("本行的缩进类型为%s，与此前缩进类型%s不符", findName(got), findName(expect)),
		Extra:   []uint8{expect, got},
	}
}

// InvalidIndentSpaceCount -
func InvalidIndentSpaceCount(count int) *Error {
	return &Error{
		Code:    0x2022,
		Message: fmt.Sprintf("当缩进类型为「空格」，其所列字符数应为4之倍数：当前空格字符数为%d", count),
		Extra:   count,
	}
}

// InvalidChar -
func InvalidChar(ch rune) *Error {
	return &Error{
		Code:    0x2026,
		Message: fmt.Sprintf("未能识别字符「%c」", ch),
		Extra:   ch,
	}
}

//// syntax error
// InvalidSyntax -
func InvalidSyntax(startIdx int) *Error {
	return &Error{
		Code:    0x2250,
		Message: "当前语法不符合规范",
		Extra:   startIdx,
	}
}

func UnexpectedIndent(startIdx int) *Error {
	return &Error{
		Code:    0x2251,
		Message: "出现不符合预期的缩进",
		Extra:   startIdx,
	}
}

// ExprMustTypeID -
func ExprMustTypeID(startIdx int) *Error {
	return &Error{
		Code: 0x2253,
		Message: "表达式须为「标识符」",
		Extra: startIdx,
	}
}

