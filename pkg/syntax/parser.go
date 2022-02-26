package syntax

// Parser - parse source file into syntax tree for further execution.
type Parser struct {
	Source []rune
	TokenBuilder
	// TokenP1: Peek1 token
	TokenP1 Token
	// TokenP2: Peek2 token
	TokenP2 Token
}

// TokenBuilder - build tokens varies from different supporting languages.
// Currently, only Chinese TokenBuilder is supported
type TokenBuilder interface {
	NextToken() (Token, error)
}

// NewParser - create a new parser from source
func NewParser(source []rune, builder TokenBuilder) *Parser {
	return &Parser{
		Source: source,
		TokenBuilder: builder,
	}
}

// Parser - parse all tokens into syntax tree
// TODO: in the future we'll parse it into bytecodes directly, instead.
func (p *Parser) Parse() (*AST, error) {
	err := p.ParseBeginLex()
	if err != nil {
		return nil, err
	}

	return nil, nil
}