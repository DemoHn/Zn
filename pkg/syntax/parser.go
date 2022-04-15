package syntax

// Parser - parse source file into syntax tree for further execution.
type Parser struct {
	*Lexer
	TokenBuilder
	// TokenP1: Peek1 token
	TokenP1 Token
	// TokenP2: Peek2 token
	TokenP2 Token
}

// TokenBuilder - build tokens varies from different supporting languages.
// Currently, only Chinese TokenBuilder is supported
type TokenBuilder interface {
	NextToken(lexer *Lexer) (Token, error)
}

// NewParser - create a new parser from source
func NewParser(lexer *Lexer, builder TokenBuilder) *Parser {
	return &Parser{
		Lexer:        lexer,
		TokenBuilder: builder,
	}
}

// Parser - parse all tokens into syntax tree
// TODO: in the future we'll parse it into bytecodes directly, instead.
func (p *Parser) Parse() (*AST, error) {
	return nil, nil
}

func (p *Parser) next() Token {
	tk, err := p.NextToken(p.Lexer)
	if err != nil {
		panic(err)
	}

	p.TokenP1 = p.TokenP2
	p.TokenP2 = tk

	return p.TokenP1
}

func (p *Parser) current() Token {
	return p.TokenP1
}

func (p *Parser) peek() Token {
	return p.TokenP2
}