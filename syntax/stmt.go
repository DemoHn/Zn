package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// Statement - a program consists of statements
type Statement interface {
	Node
	statementNode()
}

// BlockStatement - a block of statements
type BlockStatement struct {
	Children []Statement
}

func (bs *BlockStatement) statementNode() {}

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
