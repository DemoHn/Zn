package tokens

import (
	"fmt"
)

// declare marks
const (
	Comma          rune = 0xFF0C //，
	Colon          rune = 0xFF1A //：
	Semicolon      rune = 0xFF1B //；
	QuestionMark   rune = 0xFF1F //？
	RefMark        rune = 0x0026 // &
	BangMark       rune = 0xFF01 // ！
	AnnotationMark rune = 0x0040 // @
	HashMark       rune = 0x0023 // #
	EmMark         rune = 0x2014 // —
	EllipsisMark   rune = 0x2026 // …
)

// MarkerTokenType -
type MarkerTokenType int

// MarkerToken should implement all Token's functions
type MarkerToken struct {
	Type  MarkerTokenType
	Start int
	End   int
}

func (mk MarkerToken) String(detailed bool) string {
	aliasNameMap := map[MarkerTokenType]string{
		CommaType:      "Comma<,>",
		ColonType:      "Colon<:>",
		SemiColonType:  "SemiColon<;>",
		QuestionType:   "Question<?>",
		RefType:        "Ref<&>",
		BangType:       "Bang<!>",
		AnnotationType: "Annotation<@>",
		HashType:       "Hash<#>",
		EmType:         "Em<—>",
		EllipsisType:   "Ellipsis<…>",
		EmptyType:      "NIL",
	}

	raw := aliasNameMap[mk.Type]
	if detailed {
		return fmt.Sprintf("%s[%d,%d]", raw, mk.Start, mk.End)
	}
	return raw
}

// Position -
func (mk MarkerToken) Position() (int, int) {
	return mk.Start, mk.End
}

// declare markerToken types
const (
	EmptyType MarkerTokenType = iota
	CommaType
	ColonType
	SemiColonType
	QuestionType
	RefType
	BangType
	AnnotationType
	HashType
	EmType
	EllipsisType
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

// ConstructMarkerToken - consturct token from current character
// ch: current char
// idx: the index of current char
/**
func ConstructMarkerToken(l *lex.Lexer, ch rune, idx int) (lex.Token, *error.Error) {
	switch ch {
	case Comma:
		return MarkerToken{CommaType, idx, idx}, nil
	case Colon:
		return MarkerToken{ColonType, idx, idx}, nil
	case Semicolon:
		return MarkerToken{SemiColonType, idx, idx}, nil
	case QuestionMark:
		return MarkerToken{QuestionType, idx, idx}, nil
	case RefMark:
		return MarkerToken{RefType, idx, idx}, nil
	case BangMark:
		return MarkerToken{BangType, idx, idx}, nil
	case AnnotationMark:
		return MarkerToken{AnnotationType, idx, idx}, nil
	case HashMark:
		return MarkerToken{HashType, idx, idx}, nil
	case EmMark:
		// notice - only two consecutive em marks are valid!
		if l.Next() == EmMark {
			return MarkerToken{EmType, idx, idx + 1}, nil
		}
		return MarkerToken{EmptyType, 0, 0}, error.InvalidSingleEm(idx)
	case EllipsisMark:
		if l.Next() == EllipsisMark {
			return MarkerToken{EllipsisType, idx, idx + 1}, nil
		}
		return MarkerToken{EmptyType, 0, 0}, error.InvalidSingleEllipsis(idx)
	}
	// others ignore them
	return MarkerToken{EmptyType, 0, 0}, nil
}
*/
