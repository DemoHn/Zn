package lex

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/util"
)

// Lexer is a structure that pe provides a set of tools to help tokenizing the code.
type Lexer struct {
	lines      *LineScanner
	currentPos int
	peekPos    int
	lexScope   lexScope
	quoteStack *util.RuneStack
	lexError   *error.Error // if a lexError occurs, mark it
	chBuffer   []rune
	code       []rune // source code
}

// lexScope - lexer parsing scope
// lexScope indicates the next step for lexer to recognize.
// For example, if shiftLevel = LvQuoteVAR, all strings (including numbers & keywords) should be
// recognized as normal variable name
type lexScope uint8

// define lexScopes
const (
	LvNormalVAR         lexScope = 40
	LvKeyword           lexScope = 38
	LvQuoteVAR          lexScope = 36
	LvQuoteSTRING       lexScope = 34
	LvSingleLineComment lexScope = 32
	LvMultiLineComment  lexScope = 31
)

// NewLexer - new lexer
func NewLexer(code []rune) *Lexer {
	return &Lexer{
		lines:      NewLineScanner(),
		currentPos: 0,
		peekPos:    0,
		lexScope:   LvNormalVAR,
		quoteStack: util.NewRuneStack(32),
		lexError:   nil,
		chBuffer:   []rune{},
		code:       append(code, EOF),
	}
}

// next - return current rune, and move forward the cursor for 1 character.
func (l *Lexer) next() rune {
	if l.peekPos >= len(l.code) {
		return EOF
	}

	data := l.code[l.peekPos]

	l.currentPos = l.peekPos
	l.peekPos++
	return data
}

// peek - get the character of the cursor
func (l *Lexer) peek() rune {
	if l.peekPos >= len(l.code) {
		return EOF
	}
	data := l.code[l.peekPos]

	return data
}

// peek2 - get the next next character without moving the cursor
func (l *Lexer) peek2() rune {
	if l.peekPos >= len(l.code)-1 {
		return EOF
	}
	data := l.code[l.peekPos+1]

	return data
}

// peek3 - get the next next next character without moving the cursor
func (l *Lexer) peek3() rune {
	if l.peekPos >= len(l.code)-2 {
		return EOF
	}
	data := l.code[l.peekPos+2]

	return data
}

// current - get cursor value of lexer
func (l *Lexer) current() int {
	return l.currentPos
}

// rebase - rebase cursor
func (l *Lexer) rebase(cursor int) {
	l.currentPos = cursor
	l.peekPos = cursor + 1
}

func (l *Lexer) clearBuffer() {
	l.chBuffer = []rune{}
}

func (l *Lexer) pushBuffer(ch rune) {
	l.chBuffer = append(l.chBuffer, ch)
}

// pushBufferRange - push buffer for a range
// @returns bool - if push success
func (l *Lexer) pushBufferRange(start int, end int) bool {
	if end >= len(l.code) || start > end {
		return false
	}

	l.chBuffer = append(l.chBuffer, l.code[start:end+1]...)
	return true
}

// NextToken - parse and generate the next token (including comments)
func (l *Lexer) NextToken() (*Token, *error.Error) {
	var ch = l.next()
	var token *Token

	switch ch {
	case EOF:
		l.lines.PushLine(l.current())
		return TokenEOF(), nil
	case SP, TAB:
		// if indent has been scanned, it should be regarded as whitespaces
		// (it's totally ignored)
		if !l.lines.HasScanIndent() {
			l.parseIndents(ch)
		}
	case CR, LF:
		l.parseCRLF(ch)
		l.parseIndents(l.peek())
	// meet with 注, it may be possibly a lead character of a comment block
	// notice: it would also be a normal identifer (if 注[number]：) does not satisfy.
	case GlyphZHU:
		if l.lexScope == LvNormalVAR {
			cursor := l.current()
			isComment, isMultiLine := l.validateComment(ch)
			if isComment {
				token = l.parseComment(l.code[l.currentPos], isMultiLine)
			} else {
				l.rebase(cursor)
				// handle the char to next to prevent dead lock
				l.pushBuffer(ch)
				l.next()
			}
		}
	// left quotes
	case LeftQuoteI, LeftQuoteII, LeftQuoteIII, LeftQuoteIV, LeftQuoteV:
		token = l.parseString(ch)
	case MiddleDot:

	// right quotes
	case RightQuoteI, RightQuoteII, RightQuoteIII, RightQuoteIV, RightQuoteV:

	}

	if l.lexError != nil {
		return nil, l.lexError
	}
	return token, nil
}

//// parsing logics
func (l *Lexer) parseIndents(ch rune) *error.Error {
	count := 0
	for l.peek() == ch {
		count++
		l.next()
	}
	// determine indentType
	indentType := IdetUnknown
	switch ch {
	case TAB:
		indentType = IdetTab
	case SP:
		indentType = IdetSpace
	}

	return l.lines.SetIndent(count, indentType, l.current()+1)
}

func (l *Lexer) parseCRLF(ch rune) {
	cursor := l.current()
	// for CRLF <windows type> or LFCR
	if (ch == CR && l.peek() == LF) ||
		(ch == LF && l.peek() == CR) {

		// skip one char since we have judge two chars
		l.next()
		l.lines.PushLine(cursor - 1)
		return
	}
	// for LF or CR only
	// LF: <linux>, CR:<old mac>
	l.lines.PushLine(cursor - 1)
}

// validate if the coming block is a comment block
// valid comment block are listed below:
// (single-line)
// 1. 注：
// 2. 注123456：
//
// (multi-line)
//
// 3. 注：“
// 4. 注123456：“
//
// @returns (isValid, isMultiLine)
func (l *Lexer) validateComment(ch rune) (bool, bool) {
	// “ or 「
	lquotes := []rune{LeftQuoteIV, LeftQuoteII}
	for {
		ch = l.next()
		// match pattern 1, 2
		if ch == Colon {
			// match pattern 3, 4
			if util.Contains(l.peek(), lquotes) {
				l.next()
				return true, true
			}
			return true, false
		}
		if isNumber(ch) || isWhiteSpace(ch) {
			continue
		}
		return false, false
	}
}

// parseComment until its end
func (l *Lexer) parseComment(ch rune, isMultiLine bool) *Token {
	// setup
	l.clearBuffer()
	if isMultiLine {
		if !l.quoteStack.Push(ch) {
			l.lexError = error.NewErrorSLOT("push stack is full")
			return nil
		}
	}
	// iterate
	for {
		ch = l.next()
		switch ch {
		case EOF:
			l.lines.PushLine(l.current() - 1)
			return NewCommentToken(l.chBuffer, isMultiLine)
		case CR, LF:
			c1 := l.current()
			// parse CR,LF first
			l.parseCRLF(ch)
			if !isMultiLine {
				return NewCommentToken(l.chBuffer, isMultiLine)
			}
			// for multi-line comment blocks, CRLF is also included
			l.pushBufferRange(c1, l.current())
			// manually set no indents
			l.lines.SetIndent(0, IdetUnknown, l.current()+1)
		default:
			// for mutli-line comment, calculate quotes is necessary.
			if isMultiLine {
				// push left quotes
				if util.Contains(ch, LeftQuotes) {
					if !l.quoteStack.Push(ch) {
						l.lexError = error.NewErrorSLOT("quote stack if full")
						return nil
					}
				}
				// pop right quotes if possible
				if util.Contains(ch, RightQuotes) {
					currentL, hasValue := l.quoteStack.Current()
					if hasValue {
						if QuoteMatchMap[currentL] == ch {
							l.quoteStack.Pop()
						}
						// stop quoting
						if l.quoteStack.IsEmpty() {
							l.next()
							return NewCommentToken(l.chBuffer, isMultiLine)
						}
					}
				}
			}
			l.pushBuffer(ch)
		}
	}
}

// parseString -
func (l *Lexer) parseString(ch rune) *Token {
	// start up
	l.lexScope = LvQuoteSTRING
	l.clearBuffer()
	l.quoteStack.Push(ch)
	l.next()

	// parse string
	for {
		pch := l.peek()
		switch pch {
		// push quotes
		case LeftQuoteI, LeftQuoteII, LeftQuoteIII, LeftQuoteIV, LeftQuoteV:
			l.pushBuffer(pch)
			if !l.quoteStack.Push(pch) {
				l.lexError = error.NewErrorSLOT("quote stack is full")
				return nil
			}
		// pop quotes if match
		case RightQuoteI, RightQuoteII, RightQuoteIII, RightQuoteIV, RightQuoteV:
			currentL, hasValue := l.quoteStack.Current()
			if hasValue {
				if QuoteMatchMap[currentL] == pch {
					l.quoteStack.Pop()
				}
				// stop quoting
				if l.quoteStack.IsEmpty() {
					l.next()
					return NewStringToken(l.chBuffer, ch)
				}
			}
			l.pushBuffer(pch)
		case CR, LF:
			c1 := l.current()
			l.parseCRLF(pch)
			// push buffer & mark new line
			l.pushBufferRange(c1+1, l.current())
			l.lines.SetIndent(0, IdetUnknown, l.current()+1)
		case EOF:
			// after meeting with EOF
			l.lines.PushLine(l.current())
			return NewStringToken(l.chBuffer, ch)
		default:
			l.pushBuffer(pch)
		}
		l.next()
	}
}

// parseVarRemark -
func (l *Lexer) parseVarRemark(ch rune) *Token {
	// set up
	l.lexScope = LvQuoteVAR
	l.clearBuffer()
	l.next()
	// iterate
	for {
		pch := l.peek()
		switch pch {
		case EOF:
		case MiddleDot:
		}
		l.next()
	}
}
