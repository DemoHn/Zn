package tokens

import "fmt"

// SingleOp - operators that has only one unicode character
type SingleOp = rune

// DualOp - operator that has two unicode characters
type DualOp = [2]rune

const (
	// LeftQuoteI - left quote Type I (highest precedence)
	LeftQuoteI SingleOp = 0x300A //《

	// RightQuoteI - paired with LeftQuoteI
	RightQuoteI SingleOp = 0x300B //》

	// LeftQuoteII - left quote Type II
	LeftQuoteII SingleOp = 0x300C // 「

	// RightQuoteII - paired with LeftQuoteII
	RightQuoteII SingleOp = 0x300D // 」

	// LeftQuoteIII - left quote Type III
	LeftQuoteIII SingleOp = 0x3000E // 『

	// RightQuoteIII - paired with RightQuoteIII
	RightQuoteIII SingleOp = 0x3000F // 』

	// LeftQuoteIV - left quote Type IV
	LeftQuoteIV SingleOp = 0x201C // “

	// RightQuoteIV -
	RightQuoteIV SingleOp = 0x201D // ”

	// LeftQuoteV - left quote Type V (should be paired with right quote Type V)
	LeftQuoteV SingleOp = 0x2018 // ‘

	// RightQuoteV -
	RightQuoteV SingleOp = 0x2019 // ’

	// LeftBracket -
	LeftBracket SingleOp = 0x3010 // 【

	// RightBracket -
	RightBracket SingleOp = 0x3011 // 】

	// LeftParen -
	LeftParen SingleOp = 0xFF08 // （

	// RightParen -
	RightParen SingleOp = 0xFF09 // ）

	// TypeMark - declare variable (or type) is non-nullable
	TypeMark SingleOp = 0xFF01 // ！

	// VarRemark - the middle dot
	VarRemark SingleOp = 0x00B7 //
	// more marks...
)

// QuoterToken (引号) - assign all quote operators for quoting strings. This is a Token Type.
type QuoterToken struct {
	// Currently there are 5 levels. Higher level quoter could quote lower level quotes.
	// e.g. 《 “” 》 will regard the inner “”　as raw strings,
	// while “《》”　won't (and may throw syntax error)
	// All levels:
	//
	// Level 1 - 《》
	// Level 2 - 「」
	// Level 3 - 『』
	// Level 4 - “”
	// Level 5 - ‘’
	Level int
	// if this mark is a left quote or a right quote
	Left bool
	// original literal
	Literal []rune
	// start position
	Start int
	// end position
	End int
}

func (qt QuoterToken) String(detailed bool) string {
	var raw string
	if qt.Left {
		raw = fmt.Sprintf("LQuoter<%d>", qt.Level)
	}
	raw = fmt.Sprintf("RQuoter<%d>", qt.Level)

	if detailed {
		return fmt.Sprintf("%s:(%d,%d)", raw, qt.Start, qt.End)
	}
	return raw
}

// VarRemarkerToken (间隔号) - madatory assign the runes betweeen first and second as an identifier (variable)
// regardless whether there're keywords inside.
// It's a middle dot (·)
type VarRemarkerToken struct {
	start int
	end   int
}

func (vr VarRemarkerToken) String(detailed bool) string {
	raw := "VarRemarker"
	if detailed {
		return fmt.Sprintf("%s:(%d,%d)", raw, vr.start, vr.end)
	}
	return raw
}

// ArrayMapBracketToken (方括号) - use block bracket to quote a sequence as an array
// or a map.
type ArrayMapBracketToken struct {
	Left    bool
	MapSign bool
	start   int
	end     int
}

func (am ArrayMapBracketToken) String(detailed bool) string {
	b := "Array"
	if am.MapSign {
		b = "Map"
	}
	l := "L"
	if !am.Left {
		l = "R"
	}
	raw := fmt.Sprintf("%s%sSign", l, b)
	if detailed {
		return fmt.Sprintf("%s:(%d,%d)", raw, am.start, am.end)
	}
	return raw
}

// MapKeyQuoterToken (花括号) is a quick-sign for quoting the key of a map.
type MapKeyQuoterToken struct {
	Left  bool
	start int
	end   int
}
