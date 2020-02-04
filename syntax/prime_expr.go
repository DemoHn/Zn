package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// PrimeExpr - primitive expression
type PrimeExpr struct {
	literal string
}

// IsPrimitive - a primeExpr must be primitive, that is, no longer additional
// calculation required.
func (pe *PrimeExpr) IsPrimitive() bool {
	return true
}

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
func ParsePrimeExpr(p *Parser) (Expression, *error.Error) {
	l := string(p.current().Literal)

	switch p.current().Type {
	case lex.TypeIdentifier, lex.TypeVarQuote:
		tk := new(ID)
		tk.SetLiteral(l)
		return tk, nil
	case lex.TypeNumber:
		tk := new(Number)
		tk.SetLiteral(l)
		return tk, nil
	case lex.TypeString:
		tk := new(String)
		tk.SetLiteral(l)
		return tk, nil
	default:
		return nil, error.NewErrorSLOT("no such prime expr!")
	}
}
