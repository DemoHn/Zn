package error

import "fmt"

// InvalidIndentType -
func InvalidIndentType(expect uint8, got uint8) *Error {
	findName := func(idetType uint8) string {
		name := "「空格」"
		if idetType == uint8(9) { // TAB
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
