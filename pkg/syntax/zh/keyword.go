package zh

import (
	"github.com/DemoHn/Zn/pkg/syntax"
)

//// keyword character (ideographs) definition
// keywords are all ideographs that its length varies from its definitions.
// so here we define all possible chars that may be an element of one keyword.
const (
	// GlyphBU - 不 - 不等于，不小于，不大于
	GlyphBU rune = 0x4E0D
	// GlyphQIE - 且 - 且
	GlyphQIE rune = 0x4E14
	// GlyphWEI - 为 - 成为，恒为，何为，为
	GlyphWEI rune = 0x4E3A
	// GlyphYIy - 义 - 定义
	GlyphYIy rune = 0x4E49
	// GlyphZHI - 之 - 之
	GlyphZHI rune = 0x4E4B
	// GlyphYU - 于 - 等于，小于，大于，不等于，不小于，不大于
	GlyphYU rune = 0x4E8E
	// GlyphLING - 令 - 令
	GlyphLING rune = 0x4EE4
	// GlyphYIi - 以 - 以
	GlyphYIi rune = 0x4EE5
	// GlyphHE - 何 - 如何，何为
	GlyphHE rune = 0x4F55
	// GlyphRUy - 入 - 输入，导入
	GlyphRUy rune = 0x5165
	// GlyphQI - 其 - 其
	GlyphQI rune = 0x5176
	// GlyphZAI - 再 - 再如
	GlyphZAI rune = 0x518D
	// GlyphCHU - 出 - 输出，抛出
	GlyphCHU rune = 0x51FA
	// GlyphZE - 则 - 否则
	GlyphZE rune = 0x5219
	// GlyphDAOy - 到 - 得到
	GlyphDAOy rune = 0x5230
	// GlyphLI - 历 - 遍历
	GlyphLI rune = 0x5386
	// GlyphQU - 取 -
	GlyphQU rune = 0x53D6
	// GlyphFOU - 否 - 否则
	GlyphFOU rune = 0x5426
	// GlyphDA - 大 - 大于，不大于
	GlyphDA rune = 0x5927
	// GlyphRU - 如 - 如果，如何，再如
	GlyphRU rune = 0x5982
	// GlyphDING - 定 - 定义
	GlyphDING rune = 0x5B9A
	// GlyphDUI - 对 -
	GlyphDUI rune = 0x5BF9
	// GlyphDAO - 导 - 导入
	GlyphDAO rune = 0x5BFC
	// GlyphXIAO - 小 - 小于，不小于
	GlyphXIAO rune = 0x5C0F
	// GlyphYI - 已 - 已知
	GlyphYI rune = 0x5DF2
	// GlyphDANG - 当 - 每当
	GlyphDANG rune = 0x5F53
	// GlyphDEy - 得 - 得到
	GlyphDEy rune = 0x5F97
	// GlyphXUN - 循 - 继续循环，结束循环
	GlyphXUN rune = 0x5FAA
	// GlyphHENG - 恒 - 恒为
	GlyphHENG rune = 0x6052
	// GlyphCHENG - 成 - 成为
	GlyphCHENG rune = 0x6210
	// GlyphHUO - 或 - 或
	GlyphHUO rune = 0x6216
	// GlyphPAO - 抛 - 抛出
	GlyphPAO rune = 0x629B
	// GlyphSHI - 是 -
	GlyphSHI rune = 0x662F
	// GlyphSHUy - 束 - 结束循环
	GlyphSHUy rune = 0x675F
	// GlyphGUO - 果 - 如果
	GlyphGUO rune = 0x679C
	// GlyphMEI - 每 - 每当
	GlyphMEI rune = 0x6BCF
	// GlyphZHU - 注 -
	GlyphZHU rune = 0x6CE8
	// GlyphHUAN - 环 - 继续循环，结束循环
	GlyphHUAN rune = 0x73AF
	// GlyphDE - 的 - 的
	GlyphDE rune = 0x7684
	// GlyphZHIy - 知 - 已知
	GlyphZHIy rune = 0x77E5
	// GlyphDENG - 等 - 等于，不等于
	GlyphDENG rune = 0x7B49
	// GlyphJIE - 结 - 结束循环
	GlyphJIE rune = 0x7ED3
	// GlyphJI - 继 - 继续循环
	GlyphJI rune = 0x7EE7
	// GlyphXU - 续 - 继续循环
	GlyphXU rune = 0x7EED
	// GlyphSHU - 输 - 输出，输入
	GlyphSHU rune = 0x8F93
	// GlyphBIAN - 遍 - 遍历
	GlyphBIAN rune = 0x904D
)

// Keyword token types
const (
	TypeDeclareW     uint8 = 40 // 令
	TypeLogicYesW    uint8 = 41 // 为
	TypeAssignConstW uint8 = 42 // 恒为
	TypeCondOtherW   uint8 = 43 // 再如
	TypeCondW        uint8 = 44 // 如果
	TypeFuncW        uint8 = 45 // 如何
	TypeGetterW      uint8 = 46 // 何为
	TypeParamAssignW uint8 = 47 // 已知
	TypeReturnW      uint8 = 48 // 输出
	TypeLogicNotEqW  uint8 = 51 // 不等于
	TypeLogicLteW    uint8 = 52 // 不大于
	TypeLogicGteW    uint8 = 53 // 不小于
	TypeLogicLtW     uint8 = 54 // 小于
	TypeLogicGtW     uint8 = 55 // 大于
	TypeVarOneW      uint8 = 56 // 以
	TypeCondElseW    uint8 = 59 // 否则
	TypeWhileLoopW   uint8 = 60 // 每当
	TypeObjNewW      uint8 = 61 // 成为
	TypeObjDefineW   uint8 = 63 // 定义
	TypeObjThisW     uint8 = 65 // 其
	TypeLogicOrW     uint8 = 69 // 或
	TypeLogicAndW    uint8 = 70 // 且
	TypeObjDotW      uint8 = 71 // 之
	TypeObjDotIIW    uint8 = 72 // 的
	TypeLogicEqualW  uint8 = 74 // 等于
	TypeInputW       uint8 = 75 // 输入
	TypeIteratorW    uint8 = 76 // 遍历
	TypeImportW      uint8 = 77 // 导入
	TypeGetResultW   uint8 = 78 // 得到
	TypeThrowErrorW  uint8 = 79 // 抛出
	TypeContinueW    uint8 = 80 // 继续循环
	TypeBreakW       uint8 = 81 // 结束循环
)

// parseKeyword -
// @return bool matchKeyword
// @return *Token token
//
// when matchKeyword = true, a keyword token will be generated
// matchKeyword = false, regard it as normal identifier
// and return directly.
func parseKeyword(l *syntax.Lexer, moveForward bool) (bool, syntax.Token, error) {
	var tk syntax.Token
	var wordLen = 1

	startIdx := l.GetCursor()
	ch := l.GetCurrentChar()

	tk.StartIdx = startIdx
	// manual matching one or consecutive keywords
	switch ch {
	case GlyphBU:
		if l.Peek() == GlyphDA && l.Peek2() == GlyphYU {
			wordLen = 3
			tk.Type = TypeLogicLteW
		} else if l.Peek() == GlyphDENG && l.Peek2() == GlyphYU {
			wordLen = 3
			tk.Type = TypeLogicNotEqW
		} else if l.Peek() == GlyphXIAO && l.Peek2() == GlyphYU {
			wordLen = 3
			tk.Type = TypeLogicGteW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphQIE:
		tk.Type = TypeLogicAndW
	case GlyphWEI:
		tk.Type = TypeLogicYesW
	case GlyphZHI:
		tk.Type = TypeObjDotW
	case GlyphLING:
		tk.Type = TypeDeclareW
	case GlyphYIi:
		tk.Type = TypeVarOneW
	case GlyphHE:
		if l.Peek() == GlyphWEI {
			wordLen = 2
			tk.Type = TypeGetterW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphQI:
		tk.Type = TypeObjThisW
	case GlyphZAI:
		if l.Peek() == GlyphRU {
			wordLen = 2
			tk.Type = TypeCondOtherW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphFOU:
		if l.Peek() == GlyphZE {
			wordLen = 2
			tk.Type = TypeCondElseW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphDA:
		if l.Peek() == GlyphYU {
			wordLen = 2
			tk.Type = TypeLogicGtW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphRU:
		if l.Peek() == GlyphGUO {
			wordLen = 2
			tk.Type = TypeCondW
		} else if l.Peek() == GlyphHE {
			wordLen = 2
			tk.Type = TypeFuncW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphDING:
		if l.Peek() == GlyphYIy {
			wordLen = 2
			tk.Type = TypeObjDefineW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphDAO:
		if l.Peek() == GlyphRUy {
			wordLen = 2
			tk.Type = TypeImportW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphXIAO:
		if l.Peek() == GlyphYU {
			wordLen = 2
			tk.Type = TypeLogicLtW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphYI:
		if l.Peek() == GlyphZHIy {
			wordLen = 2
			tk.Type = TypeParamAssignW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphDEy:
		if l.Peek() == GlyphDAOy {
			wordLen = 2
			tk.Type = TypeGetResultW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphHENG:
		if l.Peek() == GlyphWEI {
			wordLen = 2
			tk.Type = TypeAssignConstW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphCHENG:
		if l.Peek() == GlyphWEI {
			wordLen = 2
			tk.Type = TypeObjNewW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphHUO:
		tk.Type = TypeLogicOrW
	case GlyphPAO:
		if l.Peek() == GlyphCHU {
			wordLen = 2
			tk.Type = TypeThrowErrorW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphMEI:
		if l.Peek() == GlyphDANG {
			wordLen = 2
			tk.Type = TypeWhileLoopW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphDE:
		tk.Type = TypeObjDotIIW
	case GlyphDENG:
		if l.Peek() == GlyphYU {
			wordLen = 2
			tk.Type = TypeLogicEqualW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphJIE:
		if l.Peek() == GlyphSHUy && l.Peek2() == GlyphXUN && l.Peek3() == GlyphHUAN {
			wordLen = 4
			tk.Type = TypeBreakW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphJI:
		if l.Peek() == GlyphXU && l.Peek2() == GlyphXUN && l.Peek3() == GlyphHUAN {
			wordLen = 4
			tk.Type = TypeContinueW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphSHU:
		if l.Peek() == GlyphCHU {
			wordLen = 2
			tk.Type = TypeReturnW
		} else if l.Peek() == GlyphRUy {
			wordLen = 2
			tk.Type = TypeInputW
		} else {
			return false, syntax.Token{}, nil
		}
	case GlyphBIAN:
		if l.Peek() == GlyphLI {
			wordLen = 2
			tk.Type = TypeIteratorW
		} else {
			return false, syntax.Token{}, nil
		}
	}

	// tk not empty
	if tk.Type != 0 {
		if moveForward {
			for i := 1; i <= wordLen; i++ {
				l.Next()
			}
		}
		tk.EndIdx = l.GetCursor()
		return true, tk, nil
	}
	return false, syntax.Token{}, nil
}
