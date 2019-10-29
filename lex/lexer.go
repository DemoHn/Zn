package lex

import "unicode/utf8"

// Lexer is a structure that provides a set of tools to help tokenizing the code.
type Lexer struct {
	Tokens  []Token
	current int
	code    []rune // source code
}

// Token - an universal and abstract type as the smallest unit of code syntax
type Token interface {
	String(detailed bool) string
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
	return l.end()
}

// Current - get current rune string
func (l *Lexer) Current() rune {
	return l.code[l.current]
}

// Next - return current rune, and move forward the cursor for 1 character.
func (l *Lexer) Next() rune {
	data := l.code[l.current]
	l.current++
	return data
}

// Peek - get the character of the next cursor, while the cursor doesn't move.
func (l *Lexer) Peek() rune {
	if l.current+1 < len(l.code) {
		return l.code[l.current+1]
	}
	return utf8.RuneError
}

// AppendToken - append one token to tokens
func (l *Lexer) AppendToken(token Token) {
	l.Tokens = append(l.Tokens, token)
}

// GetIndex - get cursor value of lexer
func (l *Lexer) GetIndex() int {
	return l.current
}

// private functions
func (l *Lexer) end() bool {
	return (l.current >= len(l.code))
}
