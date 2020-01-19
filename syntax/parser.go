package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// Node - general node type
type Node interface{}

// ProgramNode - the syntax tree of a program
type ProgramNode struct {
	Children []Statement
}

// Parser - parse all nodes
type Parser struct {
	*lex.Lexer
	tokens     [3]*lex.Token
	mockTokens mockTokens
}

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

type mockTokens struct {
	tokens       []lex.Token
	cursor       int
	useMockToken bool
}

// NewParser -
func NewParser(l *lex.Lexer) *Parser {
	p := &Parser{
		Lexer:      l,
		mockTokens: mockTokens{}, // mockTokens data, for unit testing
	}
	// read current and peek token
	return p
}

// Parse - parse all tokens into an AST (stored as ProgramNode)
func (p *Parser) Parse() (pg *ProgramNode, err *error.Error) {
	defer func() {
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

	// pre-read tokens
	for i := 0; i < 3; i++ {
		err = p.next()
		if err != nil {
			return
		}
	}
	pg = &ProgramNode{
		Children: []Statement{},
	}
	for p.current().Type != lex.TypeEOF {
		err = p.ParseStatement(pg)
		if err != nil {
			return
		}
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

func (p *Parser) next() *error.Error {
	var tk *lex.Token
	var err *error.Error
	// use pre-load token list
	if p.mockTokens.useMockToken {
		if p.mockTokens.cursor >= len(p.mockTokens.tokens) {
			tk = lex.NewTokenEOF(0, 0)
		} else {
			tk = &(p.mockTokens.tokens[p.mockTokens.cursor])
			p.mockTokens.cursor = p.mockTokens.cursor + 1
		}
	} else {
		tk, err = p.NextToken()
		if err != nil {
			return err
		}
	}

	p.tokens[0] = p.tokens[1]
	p.tokens[1] = p.tokens[2]
	p.tokens[2] = tk
	return nil
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

// consume one token (without callback), will return error if the incoming token (p.currentToken)
// is not in validTypes
func (p *Parser) consume(validTypes ...lex.TokenType) *error.Error {
	tk := p.current()
	tkType := tk.Type
	for _, item := range validTypes {
		if item == tkType {
			return p.next()
		}
	}
	return error.InvalidSyntax()
}

// consume one token with error func
func (p *Parser) consumeFunc(callback func(*lex.Token), validTypes ...lex.TokenType) *error.Error {
	tk := p.current()
	tkType := tk.Type
	for _, item := range validTypes {
		if item == tkType {
			callback(tk)
			return p.next()
		}
	}
	return error.InvalidSyntax()
}

//// parse element functions

// ParseStatement - a program consists of statements
//
// CFG:
// Statement -> VarDeclareStmt
//           -> VarAssignStmt
//           -> ；
func (p *Parser) ParseStatement(pg *ProgramNode) *error.Error {
	switch p.current().Type {
	case lex.TypeStmtSep:
		p.consume(lex.TypeStmtSep)
		// skip
		return nil
	case lex.TypeDeclareW:
		stmt, err := p.ParseVarDeclare()
		if err != nil {
			return err
		}
		pg.Children = append(pg.Children, stmt)
		return nil
	default:
		stmt, err := p.ParseVarAssignStmt()
		if err != nil {
			return err
		}
		pg.Children = append(pg.Children, stmt)
		return nil
	}
}

// ParseExpression - parse general expression (abstract expression type)
//
// currently, expression only contains
// ID
// Number
// String
// ArrayExpr
// （ Expr ）
func (p *Parser) ParseExpression() (Expression, *error.Error) {
	var tk Expression
	switch p.current().Type {
	case lex.TypeIdentifier, lex.TypeVarQuote, lex.TypeNumber, lex.TypeString:
		return p.ParsePrimeExpr()
	case lex.TypeArrayQuoteL:
		token, err := p.ParseArrayExpr()
		if err != nil {
			return nil, err
		}
		tk = token
	case lex.TypeStmtQuoteL:
		token, err := parseParenExpr(p)
		if err != nil {
			return nil, err
		}
		tk = token
	default:
		return nil, error.InvalidSyntax()
	}
	return tk, nil
}

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

func parseParenExpr(p *Parser) (Expression, *error.Error) {
	// #0. left paren
	if err := p.consume(lex.TypeStmtQuoteL); err != nil {
		return nil, err
	}
	// #1. parse expr
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}

	// #2. right paren
	if err := p.consume(lex.TypeStmtQuoteR); err != nil {
		return nil, err
	}
	return expr, nil
}
