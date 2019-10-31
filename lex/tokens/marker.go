package lex

import "github.com/DemoHn/Zn/error"

const (
	// Comma -
	Comma rune = 0xFF0C //，

	// Colon -
	Colon rune = 0xFF1A //：

	// Semicolon -
	Semicolon rune = 0xFF1B //；

	// QuestionMark -
	QuestionMark rune = 0xFF1F //？

	// RefMark -
	RefMark rune = 0x0026 // &

	// BangMark - declare variable (or type) is non-nullable
	BangMark rune = 0xFF01 // ！

	// AnnotationMark -
	AnnotationMark rune = 0x0040 // @

	// HashMark -
	HashMark rune = 0x0023 // #

	// EmMark - notice only two consecutive markers are valid
	EmMark rune = 0x2014 // —

	// EllipsisMark - notice only two consecutive markers are valid
	EllipsisMark rune = 0x2026 // …
)

// Markers - get marker character list
var Markers = []rune{
	Comma,
	Colon,
	Semicolon,
	QuestionMark,
	RefMark,
	BangMark,
	AnnotationMark,
	HashMark,
	EmMark,
	EllipsisMark,
}

// MarkerToken - a subtype of token
type MarkerToken interface {
	String(detailed bool) string
}

// FilterMarker - for an input character, this functions detects whether
// it belongs to a marker or not
func FilterMarker(ch rune) bool {
	for _, marker := range Markers {
		if marker == ch {
			return true
		}
	}
	return false
}

func ConstructMarkerToken(ch rune, idx int) (MarkerToken, *error.Error) {

	return emptyMarker{}, nil
}

//// helpers
type emptyMarker struct{}

func (tk emptyMarker) String(detailed bool) string {
	return "emptyMarker"
}
