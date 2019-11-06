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
	Tokens      []*tokens.Token
	lineScanner *LineScanner
	current     int
	code        []rune // source code
}

// NewLexer - new lexer
func NewLexer(code []rune) *Lexer {
	return &Lexer{
		Tokens:      []*tokens.Token{},
		lineScanner: NewLineScanner(),
		current:     0,
		code:        append(code, EOF),
	}
}

// Next - return current rune, and move forward the cursor for 1 character.
func (l *Lexer) Next() rune {
	data := l.code[l.current]
	if data == EOF {
		return EOF
	}

	l.current++
	return data
}

// Peek - get the character of the cursor
func (l *Lexer) Peek() rune {
	data := l.code[l.current]
	if data == EOF {
		return EOF
	}

	return data
}

// AppendToken - append one token to tokens
func (l *Lexer) AppendToken(token *tokens.Token) {
	l.Tokens = append(l.Tokens, token)
}

// GetIndex - get cursor value of lexer
func (l *Lexer) GetIndex() int {
	return l.current
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
	ch := l.Next()
	for ch != EOF {
		// parse indents
		l.lineScanner.NewLine(l.GetIndex() - 1)
		switch ch {
		case tokens.SP, tokens.TAB:
			curr := ch
			count := 0
			tokenType := TAB
			for {
				ch = l.Next()
				if ch != curr {
					if curr == tokens.SP {
						tokenType = SPACE
					}
					l.lineScanner.PushIndent(uint8(count), tokenType)
					break
				}
				count++
			}
		case tokens.CR, tokens.LF:
			// for CRLF <windows type>
			if ch == tokens.CR && l.Peek() == tokens.LF {
				l.lineScanner.EndLine(l.GetIndex() - 1)
				l.Next()
			}
			// for LFCR <no such type currently>
			if ch == tokens.LF && l.Peek() == tokens.CR {
				l.lineScanner.EndLine(l.GetIndex() - 1)
				l.Next()
			}

			// for LF or CR only
			// LF: <linux>, CR:<old mac>
			l.lineScanner.EndLine(l.GetIndex() - 1)
		default:
			// no other action, just move the cursor
			l.Next()
		}
	}
	return nil
}
