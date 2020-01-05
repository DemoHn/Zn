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

// Expression - a special type of statement
type Expression interface {
	Node
	expressionNode()
}

// Statement - a program consists of statements
type Statement interface {
	Node
	statementNode()
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
		return nil, error.NewErrorSLOT("no match expression")
	}
	return tk, nil
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
