package lex

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex/tokens"
)

// EOF - mark as end of file, should only exists at the end of sequence
const EOF rune = 0

// Lexer is a structure that pe provides a set of tools to help tokenizing the code.
type Lexer struct {
	lines      *LineScanner
	currentPos int
	readPos    int
	code       []rune // source code
}

// NewLexer - new lexer
func NewLexer(code []rune) *Lexer {
	return &Lexer{
		lines:      NewLineScanner(),
		currentPos: 0,
		readPos:    0,
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

// isWhiteSpace - if a character belongs to white space (including tabs, full-width spaces, etc.)
// see 'the draft' for details
func (l *Lexer) isWhiteSpace(ch rune) bool {
	spaceList := []rune{
		0x0009, 0x000B, 0x000C, 0x0020, 0x00A0,
		0x2000, 0x2001, 0x2002, 0x2003, 0x2004,
		0x2005, 0x2006, 0x2007, 0x2008, 0x2009,
		0x200A, 0x200B, 0x202F, 0x205F, 0x3000,
	}

	for _, s := range spaceList {
		if ch == s {
			return true
		}
	}

	return false
}

// NextToken - parse and generate the next token (including comments)
func (l *Lexer) NextToken() (*tokens.Token, *error.Error) {
	ch := l.next()

	// for EOF, done quickly
	if ch == EOF {
		return &tokens.Token{
			Type:    tokens.EOF,
			Literal: []rune{},
		}, nil
	}

	switch ch {
	case SP, TAB:
		if err := l.parseIndents(ch); err != nil {
			return nil, err
		}
	case CR, LF:
		l.parseCRLF(ch)
	default:
		l.parseContent(ch)
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
	case SPACE:
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

func (l *Lexer) parseContent(ch rune) {
	for l.peek() != CR && l.peek() != LF && l.peek() != EOF {
		l.next()
	}
}
