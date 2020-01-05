package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// ArrayExpr - array expression
type ArrayExpr struct {
	PrimeExpr
	Items []Expression
}

// ParseArrayExpr - yield ArrayExpr node
// CFG:
// ArrayExpr -> 【 ItemList 】
// ItemList  -> Expr ExprTail
//           ->
// ExprTail  -> ， Expr ExprTail
//           ->
//
// Expr      -> PrimaryExpr
//
// PrimaryExpr -> Number
//             -> String
//             -> ID
//             -> ArrayExpr
//             -> （ Expr ）
func (p *Parser) ParseArrayExpr() (*ArrayExpr, *error.Error) {
	// #0. consume left brancket
	if err := p.consume(lex.TypeArrayQuoteL); err != nil {
		return nil, err
	}

	ar := &ArrayExpr{
		Items: make([]Expression, 0),
	}
	// #1. consume item list
	if err := parseItemList(p, ar); err != nil {
		return nil, err
	}
	// #2. consume right brancket
	if err := p.consume(lex.TypeArrayQuoteR); err != nil {
		return nil, err
	}
	return ar, nil
}

func parseItemList(p *Parser, ar *ArrayExpr) *error.Error {
	// #0. parse expression
	expr, err := p.ParseExpression()
	if err != nil {
		return err
	}
	ar.Items = append(ar.Items, expr)

	// #1. parse list tail
	return parseItemListTail(p, ar)
}

func parseItemListTail(p *Parser, ar *ArrayExpr) *error.Error {
	// skip parsing
	if p.current().Type != lex.TypeCommaSep {
		return nil
	}
	// #0. consume comma
	if err := p.consume(lex.TypeCommaSep); err != nil {
		return err
	}
	// #1. parse expression
	expr, err := p.ParseExpression()
	if err != nil {
		return err
	}
	ar.Items = append(ar.Items, expr)

	// #2. parse tail nested again
	return parseItemListTail(p, ar)
}
