package tokens

// TokenType - general token type
type TokenType int

// Position - locate the position of a token
type Position struct {
	Line   int
	Column int // column index (deduct indent)
	Index  int // absolute index
}

// Token - general token type
type Token struct {
	Type    TokenType
	Literal []rune
	Start   Position
	End     Position
}

func (tk *Token) String(detailed bool) string {
	return ""
}
