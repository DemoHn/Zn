package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// PrimeExpr - primitive expression
type PrimeExpr struct {
	literal string
}

func (pe *PrimeExpr) expressionNode() {}

// SetLiteral - set literal for primeExpr
func (pe *PrimeExpr) SetLiteral(literal string) {
	pe.literal = literal
}

// GetLiteral -
func (pe *PrimeExpr) GetLiteral() string {
	return pe.literal
}

//// following some concrete PrimeExpr types
// including String, Number, ID

// ID - Identifier type
type ID struct {
	PrimeExpr
}

// Number -
type Number struct {
	PrimeExpr
}

// String -
type String struct {
	PrimeExpr
}

// ParsePrimeExpr -
func (p *Parser) ParsePrimeExpr() (Expression, *error.Error) {
	l := string(p.current().Literal)

	switch p.current().Type {
	case lex.TypeIdentifier, lex.TypeVarQuote:
		tk := new(ID)
		tk.SetLiteral(l)
		p.next()
		return tk, nil
	case lex.TypeNumber:
		tk := new(Number)
		tk.SetLiteral(l)
		p.next()
		return tk, nil
	case lex.TypeString:
		tk := new(String)
		tk.SetLiteral(l)
		p.next()
		return tk, nil
	default:
		return nil, error.NewErrorSLOT("no such prime expr!")
	}
}
