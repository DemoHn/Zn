package lex

import (
	"fmt"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex/tokens"
)

// EOF - mark as end of file, should only exists at the end of sequence
const EOF rune = 0

// Lexer is a structure that pe provides a set of tools to help tokenizing the code.
type Lexer struct {
	Tokens     []*tokens.Token
	lineScan   *LineScanner
	currentPos int
	readPos    int
	code       []rune // source code
}

// NewLexer - new lexer
func NewLexer(code []rune) *Lexer {
	return &Lexer{
		Tokens:     []*tokens.Token{},
		lineScan:   NewLineScanner(),
		currentPos: 0,
		readPos:    0,
		code:       append(code, EOF),
	}
}

// Next - return current rune, and move forward the cursor for 1 character.
func (l *Lexer) Next() rune {
	if l.readPos >= len(l.code) {
		return EOF
	}

	data := l.code[l.readPos]

	l.currentPos = l.readPos
	l.readPos++
	return data
}

// Peek - get the character of the cursor
func (l *Lexer) Peek() rune {
	if l.readPos >= len(l.code) {
		return EOF
	}
	data := l.code[l.readPos]

	return data
}

// AppendToken - append one token to tokens
func (l *Lexer) AppendToken(token *tokens.Token) {
	l.Tokens = append(l.Tokens, token)
}

// CurrentPos - get cursor value of lexer
func (l *Lexer) CurrentPos() int {
	return l.currentPos
}

// DisplayTokens - display tokens, usually used for debugging
func (l *Lexer) DisplayTokens() string {
	result := ""
	for idx, tk := range l.Tokens {
		if idx == 0 {
			result = tk.String(true)
		} else {
			result = fmt.Sprintf("%s %s", result, tk.String(true))
		}
	}

	return result
}

// IsWhiteSpace - if a character belongs to white space (including tabs, full-width spaces, etc.)
// see 'the draft' for details
func (l *Lexer) IsWhiteSpace(ch rune) bool {
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

// Tokenize - the main logic that transforms codes into tokens
func (l *Lexer) Tokenize() *error.Error {
	var ch rune
	// read first char
	ch = l.Next()

	l.lineScan.NewLine(l.CurrentPos())
	for ch != EOF {
		switch ch {
		case tokens.SPACE, tokens.TAB:
			if err := l.parseIndents(ch); err != nil {
				return err
			}
		case tokens.CR, tokens.LF:
			l.parseCRLF(ch)
		default:
			l.parseContent(ch)
		}
		ch = l.Peek()
	}

	// submit end line
	l.lineScan.EndLine(l.CurrentPos())
	return nil
}

func (l *Lexer) parseIndents(ch rune) *error.Error {
	count := 0
	for l.Peek() == ch {
		count++
		l.Next()
	}

	// determine indentType
	indentType := IdetUnknown
	switch ch {
	case tokens.TAB:
		indentType = IdetTab
	case tokens.SPACE:
		indentType = IdetSpace
	}

	return l.lineScan.SetIndent(count, indentType)
}

func (l *Lexer) parseCRLF(ch rune) {
	// It's clear that the line has been ended, whether it's CR or LF
	l.lineScan.EndLine(l.CurrentPos())
	l.Next()

	// for CRLF <windows type> or LFCR
	if (ch == tokens.CR && l.Peek() == tokens.LF) ||
		(ch == tokens.LF && l.Peek() == tokens.CR) {

		// skip one char since we have judge two chars
		l.Next()
		l.lineScan.NewLine(l.CurrentPos() + 1)
		return
	}
	// for LF or CR only
	// LF: <linux>, CR:<old mac>
	l.lineScan.NewLine(l.CurrentPos() + 1)
}

func (l *Lexer) parseContent(ch rune) {
	for l.Peek() != tokens.CR && l.Peek() != tokens.LF && l.Peek() != EOF {
		l.Next()
	}
}
