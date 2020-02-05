package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

/**
// Node - general node type
type Node interface{}

// ProgramNode - the syntax tree of a program
type ProgramNode struct {
	Children []Statement
}
*/

// Parser - parse all nodes
type Parser struct {
	*lex.Lexer
	tokens     [3]*lex.Token
	mockTokens mockTokens
	lineMask   uint16
}

/**
// Expression - a special type of statement
type Expression interface {
	Node
	IsPrimitive() bool
}

// Statement - a program consists of statements
type Statement interface {
	Node
	statementNode()
}
*/
type mockTokens struct {
	tokens       []lex.Token
	cursor       int
	useMockToken bool
}

const (
	modeInline uint16 = 0x01
	modeBlock  uint16 = 0x02
)

// NewParser -
func NewParser(l *lex.Lexer) *Parser {
	return &Parser{
		Lexer:      l,
		mockTokens: mockTokens{}, // mockTokens data, for unit testing
	}
}

// Parse - parse all tokens into an AST (stored as ProgramNode)
func (p *Parser) Parse() (pg *Program, err *error.Error) {
	defer func() {
		// recover error from next()
		if r := recover(); r != nil {
			err, _ = r.(*error.Error)
		}
		if err != nil {
			// if subcode >= 0x50, that means this error is generated from
			// parser, i.e. we have to set cursor manually by retrieving the start cursor
			// of current() token
			if (err.GetCode() & 0xff) >= uint16(0x50) {
				if tk := p.current(); tk != nil {
					line := tk.Range.StartLine
					col := tk.Range.StartCol

					p.moveAndSetCursor(line, col, err)
				}
			}
		}
	}()

	// advance tokens TWICE
	p.next()
	p.next()

	pg = new(Program)
	var block *BlockStmt

	peekIndent := p.getPeekIndent()
	block, err = ParseBlockStmt(p, peekIndent)
	if err == nil {
		pg.Content = block
	}
	return
}

// InitMockToken - init mockToken for parser
// after that, tokens will be retrieved directly from provided token list instead of lexer.
// this API is used actively for unit testing.
func (p *Parser) InitMockToken(tokens []lex.Token) {
	p.mockTokens = mockTokens{
		tokens:       tokens,
		cursor:       0,
		useMockToken: true,
	}
}

func (p *Parser) next() *lex.Token {
	var tk *lex.Token
	var err *error.Error

	// use pre-load token list (for mock tests)
	if p.mockTokens.useMockToken {
		if p.mockTokens.cursor >= len(p.mockTokens.tokens) {
			tk = lex.NewTokenEOF(0, 0)
		} else {
			tk = &(p.mockTokens.tokens[p.mockTokens.cursor])
			p.mockTokens.cursor = p.mockTokens.cursor + 1
		}
	} else { // normal way
		tk, err = p.NextToken()
		if err != nil {
			panic(err)
		}
	}

	// after retrieving next token successfully, check if current token has
	// violate lineMasks
	// check the comment of validateLineMask() for details
	if p.tokens[0] != nil && p.tokens[1] != nil {
		if err = p.validateLineMask(p.tokens[0], p.tokens[1]); err != nil {
			panic(err)
		}
	}

	// move advanced token buffer
	p.tokens[0] = p.tokens[1]
	p.tokens[1] = p.tokens[2]
	p.tokens[2] = tk

	return p.tokens[0]
}

func (p *Parser) current() *lex.Token {
	return p.tokens[0]
}

func (p *Parser) peek() *lex.Token {
	return p.tokens[1]
}

func (p *Parser) peek2() *lex.Token {
	return p.tokens[2]
}

func (p *Parser) setLineMask(loc uint16) {
	p.lineMask = p.lineMask | loc
}

func (p *Parser) unsetLineMask(loc uint16) {
	p.lineMask = p.lineMask & (^loc)
}

func (p *Parser) validateLineMask(lastToken *lex.Token, newToken *lex.Token) *error.Error {

	var line1 = lastToken.Range.EndLine
	var line2 = newToken.Range.StartLine
	// if new token has entered a new line
	if line2 > line1 {
		// for modeInline, all tokens should have no explicit newline (CRLF)
		if (p.lineMask & modeInline) > 0 {
			return error.NewErrorSLOT("prohibited newline! mostly because the statement doesn't finish yet!")
		}
	}
	return nil
}

// consume one token with denoted validTypes
// if not, return syntaxError
func (p *Parser) consume(validTypes ...lex.TokenType) *error.Error {
	tkType := p.peek().Type
	for _, item := range validTypes {
		if item == tkType {
			p.next()
			return nil
		}
	}
	return error.InvalidSyntax()
}

// trying to consume one token. if the token is valid in the given range of tokenTypes,
// will return its tokenType; if not, then nothing will happen.
//
// returns (matched, tokenType)
func (p *Parser) tryConsume(validTypes []lex.TokenType) (bool, *lex.Token) {
	tk := p.peek()
	for _, vt := range validTypes {
		if vt == tk.Type {
			p.next()
			return true, tk
		}
	}

	return false, nil
}

// expectBlockIndent - detect if the Indent(peek) == Indent(current) + 1
// returns (validBlockIndent, newIndent)
func (p *Parser) expectBlockIndent() (bool, int) {
	var peekLine = p.peek().Range.StartLine
	var currLine = p.current().Range.StartLine

	var peekIndent = p.GetLineIndent(peekLine)
	var currIndent = p.GetLineIndent(currLine)

	if peekIndent == currIndent+1 {
		return true, peekIndent
	}
	return false, 0
}

// getPeekIndent -
func (p *Parser) getPeekIndent() int {
	var peekLine = p.peek().Range.StartLine

	return p.GetLineIndent(peekLine)
}

/**
//// parse element functions

// ParseStatement - a program consists of statements
//
// CFG:
// Statement -> VarDeclareStmt
//           -> VarAssignStmt
//           -> ；
func ParseStatement(p *Parser, pg *ProgramNode) *error.Error {
	validTypes := []lex.TokenType{
		lex.TypeStmtSep,
		lex.TypeDeclareW,
		lex.TypeCondW,
	}
	match, tk := p.tryConsume(validTypes)
	if match {
		switch tk.Type {
		case lex.TypeStmtSep:
			// skip
			return nil
		case lex.TypeDeclareW:
			stmt, err := ParseVarDeclare(p)
			if err != nil {
				return err
			}
			pg.Children = append(pg.Children, stmt)
			return nil
		case lex.TypeCondW:
			mainIndent := p.getPeekIndent()

			stmt, err := ParseCondStmt(p, mainIndent)
			if err != nil {
				return err
			}
			pg.Children = append(pg.Children, stmt)
			return nil
		}
	}

	stmt, err := p.ParseVarAssignStmt()
	if err != nil {
		return err
	}
	pg.Children = append(pg.Children, stmt)
	return nil
}

// ParseExpression - parse general expression (abstract expression type)
//
// currently, expression only contains
// ID
// Number
// String
// ArrayExpr
func ParseExpression(p *Parser) (Expression, *error.Error) {
	var validTypes = []lex.TokenType{
		lex.TypeIdentifier, lex.TypeVarQuote, lex.TypeNumber, lex.TypeString,
		lex.TypeArrayQuoteL,
	}

	match, tk := p.tryConsume(validTypes)
	if match {
		switch tk.Type {
		case lex.TypeIdentifier, lex.TypeVarQuote, lex.TypeNumber, lex.TypeString:
			return ParsePrimeExpr(p)
		case lex.TypeArrayQuoteL:
			token, err := ParseArrayExpr(p)
			if err != nil {
				return nil, err
			}
			return token, nil
		}
	}
	return nil, error.InvalidSyntax()
}
**/

// similar to lexer's version, but with given line & col
func (p *Parser) moveAndSetCursor(line int, col int, err *error.Error) {
	buf := p.GetLineBuffer()
	cursor := error.Cursor{
		File:    p.Lexer.InputStream.Scope,
		ColNum:  col,
		LineNum: line,
		Text:    string(buf),
	}

	defer func() {
		// recover but not handle it
		recover()
		err.SetCursor(cursor)
	}()

	endCursor := p.SlideToLineEnd()
	cursor.Text = string(buf[:endCursor+1])
	err.SetCursor(cursor)
}
