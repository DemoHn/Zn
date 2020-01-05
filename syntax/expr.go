package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

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
