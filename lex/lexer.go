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
	quoteStack *util.RuneStack
	lexError   *error.Error // if a lexError occurs, mark it
	chBuffer   []rune
	code       []rune // source code
}

// NewLexer - new lexer
func NewLexer(code []rune) *Lexer {
	return &Lexer{
		lines:      NewLineScanner(),
		currentPos: 0,
		peekPos:    0,
		quoteStack: util.NewRuneStack(32),
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

// current - get current char value
func (l *Lexer) current() rune {
	if l.peekPos >= len(l.code) {
		return EOF
	}

	return l.code[l.currentPos]
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

	switch ch {
	case EOF:
		l.lines.PushLine(l.currentPos)
		return NewTokenEOF(), nil
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
		cursor := l.currentPos
		isComment, isMultiLine := l.validateComment(ch)
		if isComment {
			return l.parseComment(l.current(), isMultiLine)
		}

		l.rebase(cursor)
		// handle the char to next to prevent dead lock
		l.pushBuffer(ch)
		l.next()
	// left quotes
	case LeftQuoteI, LeftQuoteII, LeftQuoteIII, LeftQuoteIV, LeftQuoteV:
		return l.parseString(ch)
	case MiddleDot:
		return l.parseVarQuote(ch)
	}
	return nil, nil
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

	return l.lines.SetIndent(count, indentType, l.currentPos+1)
}

func (l *Lexer) parseCRLF(ch rune) {
	cursor := l.currentPos
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
func (l *Lexer) parseComment(ch rune, isMultiLine bool) (*Token, *error.Error) {
	// setup
	l.clearBuffer()
	if isMultiLine {
		if !l.quoteStack.Push(ch) {
			return nil, error.NewErrorSLOT("push stack is full")
		}
	}
	// iterate
	for {
		ch = l.next()
		switch ch {
		case EOF:
			l.lines.PushLine(l.currentPos - 1)
			return NewCommentToken(l.chBuffer, isMultiLine), nil
		case CR, LF:
			c1 := l.currentPos
			// parse CR,LF first
			l.parseCRLF(ch)
			if !isMultiLine {
				return NewCommentToken(l.chBuffer, isMultiLine), nil
			}
			// for multi-line comment blocks, CRLF is also included
			l.pushBufferRange(c1, l.currentPos)
			// manually set no indents
			l.lines.SetIndent(0, IdetUnknown, l.currentPos+1)
		default:
			// for mutli-line comment, calculate quotes is necessary.
			if isMultiLine {
				// push left quotes
				if util.Contains(ch, LeftQuotes) {
					if !l.quoteStack.Push(ch) {
						return nil, error.NewErrorSLOT("quote stack if full")
					}
				}
				// pop right quotes if possible
				if util.Contains(ch, RightQuotes) {
					currentL, _ := l.quoteStack.Current()
					if QuoteMatchMap[currentL] == ch {
						l.quoteStack.Pop()
					}
					// stop quoting
					if l.quoteStack.IsEmpty() {
						l.next()
						return NewCommentToken(l.chBuffer, isMultiLine), nil
					}
				}
			}
			l.pushBuffer(ch)
		}
	}
}

// parseString -
func (l *Lexer) parseString(ch rune) (*Token, *error.Error) {
	// start up
	l.clearBuffer()
	l.quoteStack.Push(ch)
	firstChar := ch
	// iterate
	for {
		ch := l.next()
		switch ch {
		case EOF:
			// after meeting with EOF
			l.lines.PushLine(l.currentPos - 1)
			return NewStringToken(l.chBuffer, firstChar), nil
		// push quotes
		case LeftQuoteI, LeftQuoteII, LeftQuoteIII, LeftQuoteIV, LeftQuoteV:
			l.pushBuffer(ch)
			if !l.quoteStack.Push(ch) {
				return nil, error.NewErrorSLOT("quote stack is full")
			}
		// pop quotes if match
		case RightQuoteI, RightQuoteII, RightQuoteIII, RightQuoteIV, RightQuoteV:
			currentL, _ := l.quoteStack.Current()
			if QuoteMatchMap[currentL] == ch {
				l.quoteStack.Pop()
			}
			// stop quoting
			if l.quoteStack.IsEmpty() {
				l.next()
				return NewStringToken(l.chBuffer, firstChar), nil
			}
			l.pushBuffer(ch)
		case CR, LF:
			c1 := l.currentPos
			l.parseCRLF(ch)
			// push buffer & mark new line
			l.pushBufferRange(c1, l.currentPos)
			l.lines.SetIndent(0, IdetUnknown, l.currentPos+1)
		default:
			l.pushBuffer(ch)
		}
	}
}

func (l *Lexer) parseVarQuote(ch rune) (*Token, *error.Error) {
	// setup
	l.clearBuffer()
	isFirst := true
	// iterate
	count := 0
	for {
		ch = l.next()
		// we should ensure the following chars to satisfy the condition
		// of an identifier
		switch ch {
		case EOF:
			l.lines.PushLine(l.currentPos - 1)
			return NewVarQuoteToken(l.chBuffer), nil
		case MiddleDot:
			return NewVarQuoteToken(l.chBuffer), nil
		default:
			// ignore white-spaces!
			if isWhiteSpace(ch) {
				continue
			}
			if isIdentifierChar(ch, isFirst) {
				l.pushBuffer(ch)
				count++
				if count > maxIdentifierLength {
					return nil, error.NewErrorSLOT("invalid syntax: 变量名太长")
				}
				isFirst = false
			} else {
				return nil, error.NewErrorSLOT("invalid syntax: invalid var")
			}
		}
	}
}
