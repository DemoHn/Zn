package syntax

// Parser - parse source file into syntax tree for further execution.
type Parser struct {
	*Lexer
	ASTBuilder
}

// ASTBuilder - build AST from tokens. Its logic varies from different languages.
// Currently, only Chinese ASTBuilder is supported
type ASTBuilder interface {
	ParseAST(lexer *Lexer) (*BlockStmt, error)
}

// NewParser - create a new parser from source
func NewParser(lexer *Lexer, astBuilder ASTBuilder) *Parser {
	return &Parser{
		Lexer:        lexer,
		ASTBuilder:   astBuilder,
	}
}

// Parser - parse all tokens into syntax tree
// TODO: in the future we'll parse it into bytecodes directly, instead.
func (p *Parser) Parse() (ast *BlockStmt, err error) {
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

	ast, err = p.ParseAST(p.Lexer)
	return
}