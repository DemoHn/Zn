package lex

import (
	"fmt"
	"unicode/utf8"

	"github.com/DemoHn/Zn/error"
)

// TokenType - an alias of number to declare the type of tokens
type TokenType uint16

// Lexer is a structure that pe provides a set of tools to help tokenizing the code.
type Lexer struct {
	Tokens  []Token
	current int
	code    []rune // source code
}

// Token - general token model
type Token interface {
	String(detailed bool) string //
	Position() (int, int)
}

// NewLexer - new lexer
func NewLexer(code []rune) *Lexer {
	return &Lexer{
		Tokens:  []Token{},
		current: 0,
		code:    code,
	}
}

// End - if the lexer cursor has come to the end
func (l *Lexer) End() bool {
	return (l.current >= len(l.code))
}

// Next - return current rune, and move forward the cursor for 1 character.
func (l *Lexer) Next() rune {
	if l.End() {
		return utf8.RuneError
	}
	data := l.code[l.current]
	l.current++
	return data
}

// Peek - get the character of the cursor
func (l *Lexer) Peek() rune {
	if l.End() {
		return utf8.RuneError
	}
	return l.code[l.current]
}

// AppendToken - append one token to tokens
func (l *Lexer) AppendToken(token Token) {
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

// Tokenize - the main logic that transforms codes into tokens
func (l *Lexer) Tokenize() *error.Error {
	return nil
}
