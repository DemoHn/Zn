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
	rowPos     int
}

// LineInfo - stores the (absolute) start & end index of this line
// it should be added to the scanner after all parsing is done
type LineInfo struct {
	// the indent (at the beginning) of this line
	// all lines should have indents to differentiate scopes.
	IndentNum int
	// the first character (exclude indents) of this line
	// notice when emptyLine = true, the value will be -1
	Start int
	// the last character (exclude CRLF or LF) of this line
	// notice when emptyLine = true, the value will be -1
	End int
	// if the line contains no effective characters
	EmptyLine bool
	scanState
}

// IndentType - only TAB or SPACE (U+0020) are supported
type IndentType uint8

// scanState - internal line scan state
type scanState uint8

// define IndentTypes
const (
	IdetUnknown IndentType = 0
	IdetTab     IndentType = 9
	IdetSpace   IndentType = 32
)

// define scanStates
const (
	scanInit   scanState = 0
	scanIndent scanState = 1
	scanEnd    scanState = 2
)

// NewLineScanner -
func NewLineScanner() *LineScanner {
	return &LineScanner{
		indentType: IdetUnknown,
		rowPos:     0,
		lines: []LineInfo{
			{
				IndentNum: 0,
				Start:     0,
				End:       0,
				EmptyLine: true,
				scanState: scanInit,
			},
		},
	}
}

// SetIndent - set current line's indent (for counting the consecutive intent chars, it's the task of lexer)
// notice for IndentType = SPACE, only 4 * N chars as indent is valid!
//
// possible errors:
// 1. inconsist indentType
// 2. when IndentType = SPACE, the count is not 4 * N chars
func (ls *LineScanner) SetIndent(count int, t IndentType, start int) *error.Error {
	if t == IdetSpace && count%4 != 0 {
		return error.NewErrorSLOT("SPACE count should be 4 times!")
	}

	if ls.indentType == IdetUnknown && t != IdetUnknown {
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

	// set current info
	info := ls.lines[ls.rowPos]
	info.Start = start
	info.IndentNum = indents
	info.scanState = scanIndent
	// write back data
	ls.lines[ls.rowPos] = info
	return nil
}

// PushLine - when meet CRLF/LF/CR/LFCR, that means the line is going to end
// so update the current one and create new one
func (ls *LineScanner) PushLine(endIndex int) {
	info := ls.lines[ls.rowPos]

	info.End = endIndex
	info.scanState = scanEnd
	// handle non-empty line
	if endIndex > info.Start {
		info.EmptyLine = false
	}
	// write back to data
	ls.lines[ls.rowPos] = info

	// add new template
	ls.lines = append(ls.lines, LineInfo{
		IndentNum: 0,
		Start:     0,
		End:       0,
		EmptyLine: true,
		scanState: scanInit,
	})
	ls.rowPos++
}

// HasScanIndent - determines if indents has been scanned properly
// If so, further SPs and TABs would be regarded as normal whitespaces and will
// be neglacted as usual.
func (ls *LineScanner) HasScanIndent() bool {
	state := ls.lines[ls.rowPos].scanState
	return state == scanIndent
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
		if line.scanState == scanEnd {
			if line.EmptyLine {
				ss = append(ss, "Empty<0>")
			} else {
				ss = append(ss, fmt.Sprintf(
					"%s<%d>[%d,%d]",
					indentChar, line.IndentNum,
					line.Start, line.End,
				))
			}
		}
	}
	return strings.Join(ss, " ")
}
