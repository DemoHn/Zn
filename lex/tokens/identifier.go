package tokens

import (
	"fmt"
)

// declare identifier ranges
const (
	// minimum length of one identifier
	MinLength = 1
	// maximum length of one identifier
	MaxLength = 32
)

// IdentifierToken -
type IdentifierToken struct {
	Literal []rune
	// actual literal length - notice since there may have spaces inside
	// the original string, LiteralLen may not be (End - Start + 1)!
	LiteralLen int
	Start      int
	End        int
}

func (idf IdentifierToken) String(detailed bool) string {
	raw := fmt.Sprintf("IDF{%d}<%s>", idf.LiteralLen, string(idf.Literal))

	if detailed {
		return fmt.Sprintf("%s[%d,%d]", raw, idf.Start, idf.End)
	}
	return raw
}

// Position -
func (idf IdentifierToken) Position() (int, int) {
	return idf.Start, idf.End
}

// ConstructIdentifierToken - make an identifier
/**
func ConstructIdentifierToken(l *lex.Lexer) (lex.Token, *error.Error) {
	literal := []rune{}
	gotFirst := false
	for !l.End() {
		ch := l.Next()
		// skip white spaces
		if l.IsWhiteSpace(ch) {
			continue
		}

		if !gotFirst && !inFirstChar(ch) {
			// throw error
		}

	}
}
*/
