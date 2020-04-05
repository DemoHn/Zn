package lex

import (
	"github.com/DemoHn/Zn/error"
)

// LineStack - store line source and its indent info
type LineStack struct {
	IndentType
	CurrentLine int
	lines       []LineInfo
	scanCursor
	lineBuffer []rune
}

// LineInfo - stores the (absolute) start & end index of this line
// this should be added to the scanner after all parsing is done
type LineInfo struct {
	// the indent number (at the beginning) of this line.
	// all lines should have indents to differentiate scopes.
	indents int
	// startIdx - start index of lineBuffer
	startIdx int
	// endIdx - end index of lineBuffer
	endIdx int
}

type scanCursor struct {
	startIdx int
	indents  int
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
		IndentType:  IdetUnknown,
		lines:       []LineInfo{},
		CurrentLine: 1,
		scanCursor:  scanCursor{0, 0, scanIndent},
		lineBuffer:  []rune{},
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
	switch t {
	case IdetUnknown:
		if count > 0 && ls.IndentType != t {
			return error.InvalidIndentType(ls.IndentType, t)
		}
	case IdetSpace, IdetTab:
		// init ls.IndentType
		if ls.IndentType == IdetUnknown {
			ls.IndentType = t
		}
		// when t = space, the character count must be 4 * N
		if t == IdetSpace && count%4 != 0 {
			return error.InvalidIndentSpaceCount(count)
		}
		// when t does not match indentType, throw error
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
		indents:  idets,
		startIdx: ls.scanCursor.startIdx + count,
		endIdx:   lastIndex,
	}

	ls.lines = append(ls.lines, line)
	ls.scanCursor.scanState = scanEnd
}

// NewLine - reset scanCurosr
// change scanState from 2 -> 0
func (ls *LineStack) NewLine(index int) {
	// reset start index
	ls.scanCursor = scanCursor{index, 0, scanInit}

	// add CurrentLine
	ls.CurrentLine++
}

// onIndentStage - if the incoming SPACE/TAB should be regarded as indents
// or normal spaces.
func (ls *LineStack) onIndentStage() bool {
	return ls.scanState == scanInit
}

// AppendLineBuffer - push data to lineBuffer
func (ls *LineStack) AppendLineBuffer(data []rune) {
	ls.lineBuffer = append(ls.lineBuffer, data...)
}

// GetLineBufferSize -
func (ls *LineStack) GetLineBufferSize() int {
	return ls.getLineBufferSize()
}

// GetLineBuffer -
func (ls *LineStack) GetLineBuffer() []rune {
	return ls.lineBuffer
}

// GetLineIndent - get current lineNum indent
// NOTE: lineNum starts from 1
// NOTE2: if lineNum not found, return -1
func (ls *LineStack) GetLineIndent(lineNum int) int {
	// lineNum exceeds current range
	if lineNum > ls.CurrentLine {
		return -1
	}
	if lineNum == ls.CurrentLine {
		return ls.indents
	}

	if lineNum > 0 {
		lineInfo := ls.lines[lineNum-1]
		return lineInfo.indents
	}

	return -1
}

// GetParsedLineText - get line content of parsed lines (ie. not including current line)
// NOTICE: when lineNum >= $currentLine, this will return an empty string!
//
func (ls *LineStack) GetParsedLineText(lineNum int) []rune {
	if lineNum >= ls.CurrentLine {
		return []rune{}
	}

	if lineNum > 0 {
		lineInfo := ls.lines[lineNum-1]
		sIdx := lineInfo.startIdx
		eIdx := lineInfo.endIdx + 1
		return ls.lineBuffer[sIdx : eIdx+1]
	}

	return []rune{}
}

//// private helpers
func (ls *LineStack) getLineBufferSize() int {
	return len(ls.lineBuffer)
}

// getChar - get value from lineBuffer
func (ls *LineStack) getChar(idx int) rune {
	if idx >= len(ls.lineBuffer) {
		return EOF
	}
	return ls.lineBuffer[idx]
}
