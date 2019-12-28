package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// Node - general node type
type Node interface {
	getType() nodeType
}

// Statement - a program consists of statements
type Statement interface {
	Node
	statementNode()
}

// Expression - a special type of statement
type Expression interface {
	Node
	expressionNode()
}

type nodeType int

// ProgramNode - the syntax tree of a program
type ProgramNode struct {
	Children []Statement
}

func (p *ProgramNode) getType() nodeType {
	return TypeProgram
}

// declare node types
const (
	TypeProgram    nodeType = 0
	TypeVarAssign  nodeType = 1 // 令...为...
	TypeArrayExpr  nodeType = 3 // 【1，2，3，4】
	TypeIdentifier nodeType = 5
	TypeNumber     nodeType = 6
	TypeString     nodeType = 7
)

// Parser - parse all nodes
type Parser struct {
	lexer        *lex.Lexer
	currentToken *lex.Token
	peekToken    *lex.Token
}

// NewParser -
func NewParser(l *lex.Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}
	// read current and peek token
	p.next()
	p.next()
	return p
}

func (p *Parser) next() *error.Error {
	tk, err := p.lexer.NextToken()
	if err != nil {
		return err
	}

	p.currentToken = p.peekToken
	p.peekToken = tk
	return nil
}

func (p *Parser) current() *lex.Token {
	return p.currentToken
}

func (p *Parser) peek() *lex.Token {
	return p.peekToken
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
