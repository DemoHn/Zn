package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// ParseStatement - a program consists of statements -
func (p *Parser) ParseStatement() (Statement, *error.Error) {
	switch p.current().Type {
	case lex.TypeDeclareW:
		return p.ParseVarDeclare()
	}

	return nil, nil
}
