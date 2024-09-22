package syntax

// Parser - parse source file into syntax tree for further execution.
type Parser struct {
	*Lexer
	ASTBuilder
}

// ASTBuilder - build AST from tokens. Its logic varies from different languages.
// Currently, only Chinese ASTBuilder is supported
type ASTBuilder interface {
	ParseAST(lexer *Lexer) (*Program, error)
	ParseVarInputs(lexer *Lexer) (*VarDeclareStmt, error)
}

func NewParser(source []rune, astBuilder ASTBuilder) *Parser {
	lexer := NewLexer(source)
	return &Parser{
		Lexer:      lexer,
		ASTBuilder: astBuilder,
	}
}

func (p *Parser) GetLexer() *Lexer {
	return p.Lexer
}

// Parser - parse all tokens into syntax tree
func (p *Parser) Parse() (ast *Program, err error) {
	// handle panics
	defer func() {
		var ok bool
		if r := recover(); r != nil {
			err, ok = r.(error)
			if !ok {
				panic(r)
			}
		}
	}()

	ast, err = p.ParseAST(p.Lexer)
	return
}

func (p *Parser) ParseVarInputs() (vdStmt *VarDeclareStmt, err error) {
	// handle panics
	defer func() {
		var ok bool
		if r := recover(); r != nil {
			err, ok = r.(error)
			if !ok {
				panic(r)
			}
		}
	}()

	vdStmt, err = p.ASTBuilder.ParseVarInputs(p.Lexer)
	return
}
