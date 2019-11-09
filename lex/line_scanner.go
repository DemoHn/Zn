package lex

import (
	"fmt"
	"strings"

	"github.com/DemoHn/Zn/error"
)

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
	Indents int
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
	IdetUnknown IndentType = 0
	IdetTab     IndentType = 9
	IdetSpace   IndentType = 32
)

// NewLineScanner -
func NewLineScanner() *LineScanner {
	return &LineScanner{
		indentType: IdetUnknown,
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

// SetIndent - set current line's indent (for counting the consecutive intent chars, it's the task of lexer)
// notice for IndentType = SPACE, only 4 * N chars as indent is valid!
//
// possible errors:
// 1. inconsist indentType
// 2. when IndentType = SPACE, the count is not 4 * N chars
func (ls *LineScanner) SetIndent(count int, t IndentType) *error.Error {
	if t == IdetSpace && count%4 != 0 {
		return error.NewErrorSLOT("SPACE count should be 4 times!")
	}

	if ls.indentType == IdetUnknown {
		ls.indentType = t
	} else {
		if ls.indentType != t {
			return error.NewErrorSLOT("前后indent char 不同！")
		}
	}

	// when indentType = TAB, count = indents
	indents := count
	if ls.indentType == IdetSpace {
		indents = count / 4
	}

	ls.lineCache.Indents = indents
	ls.lineCache.StartIndex = ls.lineCache.StartIndex + int(count)
	return nil
}

// EndLine - when meet CRLF, that means the line is going to end
// It's time to push cache into line info
func (ls *LineScanner) EndLine(endIndex int) {
	// for some reason (e.g. empty code), ls.lineCache maybe nil
	// therefore we won't submit data to lineInfo
	if ls.lineCache == nil {
		return
	}

	ls.lineCache.EndIndex = endIndex
	// handle non-empty line
	if endIndex > ls.lineCache.StartIndex {
		ls.lineCache.EmptyLine = false
	}

	ls.lines = append(ls.lines, *(ls.lineCache))
	ls.clearCache()
}

// String shows all lines info for testing.
// format::=
//   {lineInfo1} {lineInfo2} {lineInfo3} ...
//
// lineInfo ::=
//   Space<2>[23,45] or
//   Tab<4>[0,1] or
//   Empty<0>
func (ls *LineScanner) String() string {
	ss := []string{}
	var indentChar string
	// get indent type
	switch ls.indentType {
	case IdetUnknown:
		indentChar = "Unknown"
	case IdetSpace:
		indentChar = "Space"
	case IdetTab:
		indentChar = "Tab"
	}

	for _, line := range ls.lines {
		if line.EmptyLine {
			ss = append(ss, "Empty<0>")
		} else {
			ss = append(ss, fmt.Sprintf(
				"%s<%d>[%d,%d]",
				indentChar, line.Indents,
				line.StartIndex, line.EndIndex,
			))
		}
	}

	return strings.Join(ss, " ")
}

//// private function
func (ls *LineScanner) clearCache() {
	ls.lineCache = nil
}
