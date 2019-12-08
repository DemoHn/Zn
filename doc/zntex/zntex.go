package zntex

import (
	"fmt"
	"unicode/utf8"
)

// ZnTex - main type of ZnTex
type ZnTex struct {
	input  []rune
	tokens []Token
}

// New ZnTex instance
func New() *ZnTex {
	return &ZnTex{
		input: []rune{},
	}
}

// ReadInput -
func (zt *ZnTex) ReadInput(rawData []byte) error {
	b := rawData
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if r == utf8.RuneError {
			return fmt.Errorf("invalid utf-8 byte")
		}
		zt.input = append(zt.input, r)
		b = b[size:]
	}

	// append EOF
	zt.input = append(zt.input, EOF)
	return nil
}

// Parse - get all tokens
func (zt *ZnTex) Parse() error {
	var parser = &Parser{
		currentPos: 0,
		peekPos:    0,
		quoteStack: NewRuneStack(64),
		chBuffer:   []rune{},
		input:      []rune{},
	}

	for {
		tk := parser.NextToken()
		if tk == nil {
			break
		}

		zt.tokens = append(zt.tokens, tk)
	}
	return nil
}
