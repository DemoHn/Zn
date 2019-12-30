package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// some concrete prime expression types

// ID - Identifier type
type ID struct {
	literal string
}

func (x *ID) getType() nodeType {
	return TypeIdentifier
}

func (x *ID) expressionNode() {}

// Number -
type Number struct {
	literal string
}

func (n *Number) getType() nodeType {
	return TypeNumber
}

func (n *Number) expressionNode() {}

// String -
type String struct {
	literal string
}

func (s *String) getType() nodeType {
	return TypeString
}

func (s *String) expressionNode() {}

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
	case lex.TypeIdentifier, lex.TypeVarQuote:
		tk = &ID{
			literal: string(p.current().Literal),
		}
		p.next()
	case lex.TypeNumber:
		tk = &Number{
			literal: string(p.current().Literal),
		}
		p.next()
	case lex.TypeString:
		tk = &String{
			literal: string(p.current().Literal),
		}
		p.next()
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
