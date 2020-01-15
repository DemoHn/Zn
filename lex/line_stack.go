package lex

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/util"
)

// LineStack - store line source and its indent info
type LineStack struct {
	IndentType
	lines []LineInfo
	scanCursor
	lineBuffer []rune
}

// LineInfo - stores the (absolute) start & end index of this line
// this should be added to the scanner after all parsing is done
type LineInfo struct {
	// the indent number (at the beginning) of this line.
	// all lines should have indents to differentiate scopes.
	Indents int
	// source data of the line (without indentation chars)
	Source []rune
}

type scanCursor struct {
	indents int
	scanState
}

// IndentType - only TAB or SPACE (U+0020) are supported
type IndentType = uint8

// scanState - internal line scan state
type scanState uint8

// define IndentTypes
const (
	IdetUnknown IndentType = 0
	IdetTab     IndentType = 9
	IdetSpace   IndentType = 32
)

// define scanStates
// example:
//
// [ IDENTS ] [ TEXT TEXT TEXT ] [ CR LF ]
// ^         ^                  ^
// 0         1                  2
//
// 0: scanInit
// 1: scanIndent
// 2: scanEnd
const (
	scanInit   scanState = 0
	scanIndent scanState = 1
	scanEnd    scanState = 2
)

// NewLineStack - new line stack
func NewLineStack() *LineStack {
	return &LineStack{
		IndentType: IdetUnknown,
		lines:      []LineInfo{},
		scanCursor: scanCursor{
			indents:   0,
			scanState: scanInit,
		},
		lineBuffer: []rune{},
	}
}

// SetIndent - set current line's indent (for counting the consecutive intent chars, it's the task of lexer)
// notice: for IndentType = SPACE, only 4 * N chars as indent is valid!
// and change scanState from 0 -> 1
//
// possible errors:
// 1. inconsist indentType
// 2. when IndentType = SPACE, the count is not 4 * N chars
func (ls *LineStack) SetIndent(count int, t IndentType) *error.Error {
	if t == IdetSpace && count%4 != 0 {
		return error.InvalidIndentSpaceCount(count)
	}

	if ls.IndentType == IdetUnknown && t != IdetUnknown {
		ls.IndentType = t
	} else {
		if ls.IndentType != t {
			return error.InvalidIndentType(ls.IndentType, t)
		}
	}

	// when indentType = TAB, count = indents
	// otherwise, count = indents * 4
	indentNum := count
	if ls.IndentType == IdetSpace {
		indentNum = count / 4
	}

	// set scanCursor
	ls.scanCursor.indents = indentNum
	ls.scanCursor.scanState = scanIndent

	return nil
}

// PushLine - push effective line text into line info
// change scanState from 1 -> 2
func (ls *LineStack) PushLine(lastIndex int) {
	idets := ls.scanCursor.indents
	count := idets
	if ls.IndentType == IdetSpace {
		count = idets * 4
	}

	// push index
	line := LineInfo{
		Indents: idets,
		Source:  util.Copy(ls.lineBuffer[count : lastIndex+1]),
	}
	ls.lines = append(ls.lines, line)
	ls.scanCursor.scanState = scanEnd
}

// NewLine - reset scanCurosr
// change scanState from 2 -> 0
func (ls *LineStack) NewLine(index int) {
	// reset start index
	ls.lineBuffer = ls.lineBuffer[index:]
	ls.scanCursor = scanCursor{
		indents:   0,
		scanState: scanInit,
	}
}

// HasScanIndent - determines if indents has been scanned properly
// If so, further SPs and TABs would be regarded as normal whitespaces and will
// be neglacted as usual.
func (ls *LineStack) HasScanIndent() bool {
	state := ls.scanCursor.scanState
	return state == scanIndent
}

// AppendLineBuffer - push data to lineBuffer
func (ls *LineStack) AppendLineBuffer(data []rune) {
	ls.lineBuffer = append(ls.lineBuffer, data...)
}

// GetColIndex - get value from lineBuffer
func (ls *LineStack) GetColIndex(idx int) rune {
	if idx >= len(ls.lineBuffer) {
		return EOF
	}
	return ls.lineBuffer[idx]
}

// GetLineBufferSize -
func (ls *LineStack) GetLineBufferSize() int {
	return len(ls.lineBuffer)
}

// GetLineBuffer -
func (ls *LineStack) GetLineBuffer() []rune {
	return ls.lineBuffer
}

// CurrentLineNum -
func (ls *LineStack) CurrentLineNum() int {
	if ls.scanCursor.scanState == scanEnd {
		return len(ls.lines)
	}
	return len(ls.lines) + 1
}
