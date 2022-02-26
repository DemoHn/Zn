package zh

import "github.com/DemoHn/Zn/pkg/syntax"

// TokenBuilderZH
type TokenBuilderZH struct {
	*syntax.Lexer
}

// NextToken -
func (tb *TokenBuilderZH) NextToken() (syntax.Token, error) {
	return syntax.Token{}, nil
}
