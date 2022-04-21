package syntax

// Parser - parse source file into syntax tree for further execution.
type Parser struct {
	*Lexer
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

// ASTBuilder - build AST from tokens. Its logic varies from different languages.
// Currently, only Chinese ASTBuilder is supported
type ASTBuilder interface {
	ParseAST(parser *Parser) (*AST, error)
}

// NewParser - create a new parser from source
func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		Lexer:        lexer,
	}
}

// Parser - parse all tokens into syntax tree
// TODO: in the future we'll parse it into bytecodes directly, instead.
func (p *Parser) Parse(tkBuilder TokenBuilder, astBuilder ASTBuilder) (ast *AST, err error) {
	// handle panics
	defer func() {
		var ok bool
		if r := recover(); r != nil {
			err, ok = r.(error)
			if !ok {
				panic(r)
			}
			// handleDeferError(p, err)
		}
	}()

	// advance tokens TWICE
	p.Next(tkBuilder)
	p.Next(tkBuilder)

	ast, err = astBuilder.ParseAST(p)
	return
}

func (p *Parser) Next(tkBuilder TokenBuilder) Token {
	tk, err := tkBuilder.NextToken(p.Lexer)
	if err != nil {
		panic(err)
	}

	p.TokenP1 = p.TokenP2
	p.TokenP2 = tk

	return p.TokenP1
}

func (p *Parser) Current() Token {
	return p.TokenP1
}

func (p *Parser) Peek() Token {
	return p.TokenP2
}