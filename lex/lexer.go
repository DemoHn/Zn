package lex

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/util"
)

// Lexer is a structure that pe provides a set of tools to help tokenizing the code.
type Lexer struct {
	lines      *LineScanner
	currentPos int
	readPos    int
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
		readPos:    0,
		lexScope:   LvNormalVAR,
		quoteStack: util.NewRuneStack(32),
		lexError:   nil,
		chBuffer:   []rune{},
		code:       append(code, EOF),
	}
}

// next - return current rune, and move forward the cursor for 1 character.
func (l *Lexer) next() rune {
	if l.readPos >= len(l.code) {
		return EOF
	}

	data := l.code[l.readPos]

	l.currentPos = l.readPos
	l.readPos++
	return data
}

// peek - get the character of the cursor
func (l *Lexer) peek() rune {
	if l.readPos >= len(l.code) {
		return EOF
	}
	data := l.code[l.readPos]

	return data
}

// peek2 - get the next next character without moving the cursor
func (l *Lexer) peek2() rune {
	if l.readPos >= len(l.code)-1 {
		return EOF
	}
	data := l.code[l.readPos+1]

	return data
}

// peek3 - get the next next next character without moving the cursor
func (l *Lexer) peek3() rune {
	if l.readPos >= len(l.code)-2 {
		return EOF
	}
	data := l.code[l.readPos+2]

	return data
}

// current - get cursor value of lexer
func (l *Lexer) current() int {
	return l.currentPos
}

// rebase - rebase cursor
func (l *Lexer) rebase(cursor int) {
	l.currentPos = cursor
	l.readPos = cursor + 1
}

func (l *Lexer) clearBuffer() {
	l.chBuffer = []rune{}
}

func (l *Lexer) pushBuffer(ch rune) {
	l.chBuffer = append(l.chBuffer, ch)
}

// NextToken - parse and generate the next token (including comments)
func (l *Lexer) NextToken() (*Token, *error.Error) {
	var ch = l.peek()
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
	// meet with 注, it may be possibly a lead character of a comment block
	// notice: it would also be a normal identifer (if 注[number]：) does not satisfy.
	case GlyphZHU:
		if l.lexScope == LvNormalVAR {
			cursor := l.current()
			isComment, isMultiLine := l.validateComment(ch)
			if isComment {
				token = l.parseComment(l.peek(), isMultiLine)
			} else {
				l.rebase(cursor)
				// handle the char to next to prevent dead lock
				l.pushBuffer(ch)
				l.next()
			}
		}
	// left quotes
	case LeftQuoteI, LeftQuoteII, LeftQuoteIII, LeftQuoteIV, LeftQuoteV:
		l.lexScope = LvQuoteSTRING
		l.parseString(ch)
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
	// It's clear that the line has been ended, whether it's CR or LF
	l.next()

	// for CRLF <windows type> or LFCR
	if (ch == CR && l.peek() == LF) ||
		(ch == LF && l.peek() == CR) {

		// skip one char since we have judge two chars
		l.next()
		l.lines.PushLine(l.current() + 1)
		return
	}
	// for LF or CR only
	// LF: <linux>, CR:<old mac>
	l.lines.PushLine(l.current())
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
	// read next char
	l.next()
	// if next char is a number or whitespace, skip it
	for isNumber(l.peek()) || isWhiteSpace(l.peek()) {
		l.next()
	}
	// match pattern 1, 2
	if l.peek() == Colon {
		l.next()
		// “ or 「
		lquotes := []rune{LeftQuoteV, LeftQuoteII}
		// match pattern 3, 4
		if util.Contains(l.peek(), lquotes) {
			return true, true
		}

		return true, false
	}
	return false, false
}

// parseComment until its end
func (l *Lexer) parseComment(ch rune, isMultiLine bool) *Token {
	// tear-down operations before return
	defer func() {
		l.clearBuffer()
	}()

	for ch != EOF {
		// CR, LF
		if util.Contains(ch, []rune{CR, LF}) {
			// parse CR,LF first
			l.parseCRLF(ch)
			if !isMultiLine {
				return NewCommentToken(l.chBuffer, isMultiLine)
			}
		} else {
			// normal string
			l.pushBuffer(ch)
			// for mutli-line comment, calculate quotes is necessary.
			if isMultiLine {
				// push left quotes
				if util.Contains(ch, LeftQuotes) {
					if !l.quoteStack.Push(ch) {
						l.lexError = error.NewErrorSLOT("quote stack if full")
						return nil
					}
					// pop right quotes if possible
				} else if util.Contains(ch, RightQuotes) {
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
		}

		l.next()
		ch = l.peek()
	}

	// meeting with EOF
	// whenever it's a multi-line or single-line comment block,
	// whenever the multi-line comment block (if so) has been enclosed or not,
	// return the token directly.
	l.lines.PushLine(l.current())
	return NewCommentToken(l.chBuffer, isMultiLine)
}

func (l *Lexer) parseString(ch rune) {

}
