package lex

import "github.com/DemoHn/Zn/error"

// LineScanner provides a structure to store all indents info and its start cursor
type LineScanner struct {
	indentType IndentType
	lines      []LineInfo
	lineCache  *LineInfo
}

// LineInfo - stores the (absolute) start & end index of this line
// it should be added to the scanner after all parsing is done
type LineInfo struct {
	// the indent (at the beginning) of this line
	// all lines should have indents to differentiate scopes.
	Indents uint8
	// the first character (exclude indents) of this line
	// notice when emptyLine = true, the value will be -1
	StartIndex int
	// the last character (exclude CRLF or LF) of this line
	// notice when emptyLine = true, the value will be -1
	EndIndex int
	// if the line contains no effective characters
	EmptyLine bool
}

// IndentType - only TAB or SPACE (U+0020) are supported
type IndentType uint8

// define IndentTypes
const (
	UNKNOWN IndentType = 0
	TAB     IndentType = 9
	SPACE   IndentType = 32
)

// NewLineScanner -
func NewLineScanner() *LineScanner {
	return &LineScanner{
		indentType: UNKNOWN,
		lines:      []LineInfo{},
		lineCache:  nil,
	}
}

// NewLine - init and stash new LineInfo to lineCache
func (ls *LineScanner) NewLine(idx int) {
	ls.lineCache = &LineInfo{
		Indents:    0,
		StartIndex: idx,
		EndIndex:   idx,
		EmptyLine:  true,
	}
}

// PushIndent - push indent (for counting the consecutive intent chars, it's the task of lexer)
// notice for IndentType = SPACE, only 4 * N chars as indent is valid!
//
// possible errors:
// 1. inconsist indentType
// 2. when IndentType = SPACE, the count is not 4 * N chars
func (ls *LineScanner) PushIndent(count uint8, t IndentType) *error.Error {
	if t == SPACE && count%4 != 0 {
		// TODO: add error
		return error.InvalidIndent(0)
	}

	if ls.indentType == UNKNOWN {
		ls.indentType = t
	}
	if ls.indentType != t {
		// TODO: add error
		return error.InvalidIndent(0)
	}

	indents := count
	if ls.indentType == SPACE {
		indents = count / 4
	}

	ls.lineCache.Indents = indents
	ls.lineCache.StartIndex = ls.lineCache.StartIndex + int(count)
	return nil
}

// EndLine - when meet CRLF, that means the line is going to end
// It's time to push cache into line info
func (ls *LineScanner) EndLine() {

}
