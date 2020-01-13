package error

import (
	"fmt"
)

// InvalidSingleEllipsis -
// TODO: add cursor
func InvalidSingleEllipsis(idx int) *Error {
	return lexError.NewError(0x02, Error{
		text: "未能识别单个「…」字符，请注意需要有连续两个方可有效！",
		info: idx,
	})
}

// DecodeUTF8Fail - decode error
func DecodeUTF8Fail(ch byte) *Error {
	return lexError.NewError(0x20, Error{
		text: fmt.Sprintf("前方有无法解析成UTF-8编码之异常字符'\\x%x'，请确认文件编码之正确性及完整性！", ch),
		info: ch,
	})
}

// InvalidIndentType -
func InvalidIndentType(expect uint8, got uint8) *Error {
	findName := func(idetType uint8) string {
		name := "「空格」"
		if idetType == uint8(9) { // TAB
			name = "「TAB」"
		}
		return name
	}
	return lexError.NewError(0x21, Error{
		text: fmt.Sprintf("此行现行缩进类型为%s，与前设缩进类型%s不符！", findName(got), findName(expect)),
		info: struct {
			expect uint8
			got    uint8
		}{
			expect, got,
		},
	})
}

// InvalidIndentSpaceCount -
func InvalidIndentSpaceCount(count int) *Error {
	return lexError.NewError(0x22, Error{
		text: fmt.Sprintf("当缩进类型为「空格」，其所列字符数应为4之倍数：当前空格字符数为%d", count),
		info: count,
	})
}
