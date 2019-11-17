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
}

func (tk *Token) String(detailed bool) string {
	return ""
}

// token types
const (
	None TokenType = 1
	EOF  TokenType = 0
)
