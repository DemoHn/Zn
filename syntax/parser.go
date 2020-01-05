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
	lexer        *lex.Lexer
	currentToken *lex.Token
	peekToken    *lex.Token
	peek2Token   *lex.Token
}

// NewParser -
func NewParser(l *lex.Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}
	// read current and peek token
	p.next()
	p.next()
	p.next()
	return p
}

// Parse - parse all tokens into an AST (stored as ProgramNode)
func (p *Parser) Parse() (*ProgramNode, *error.Error) {
	pg := &ProgramNode{
		Children: []Statement{},
	}
	for p.current().Type != lex.TypeEOF {
		err := p.ParseStatement(pg)
		if err != nil {
			return nil, err
		}
	}
	return pg, nil
}

func (p *Parser) next() *error.Error {
	tk, err := p.lexer.NextToken()
	if err != nil {
		return err
	}

	p.currentToken = p.peekToken
	p.peekToken = p.peek2Token
	p.peek2Token = tk
	return nil
}

func (p *Parser) current() *lex.Token {
	return p.currentToken
}

func (p *Parser) peek() *lex.Token {
	return p.peekToken
}

func (p *Parser) peek2() *lex.Token {
	return p.peek2Token
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
	return error.NewErrorSLOT("syntax error")
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
	return error.NewErrorSLOT("syntax error")
}
