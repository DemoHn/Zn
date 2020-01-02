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

func (bs *BlockStatement) getType() nodeType {
	return TypeBlockStmt
}

func (bs *BlockStatement) statementNode() {}

// ParseStatement - a program consists of statements
//
// CFG:
// Statement -> StmtItem StmtTail
// StmtTail  -> ；StmtItem StmtTail
//           ->
// StmtItem  -> VarDeclareStmt
//           -> VarAssignStmt
func (p *Parser) ParseStatement(pg *ProgramNode) *error.Error {
	// #0. parse statement item
	stmt, err := parseStmtItem(p)
	if err != nil {
		return err
	}
	pg.Children = append(pg.Children, stmt)
	// #1. parse statement tail
	return parseStmtTail(p, pg)
}

func parseStmtItem(p *Parser) (Statement, *error.Error) {
	switch p.current().Type {
	case lex.TypeDeclareW:
		return p.ParseVarDeclare()
	default:
		return p.ParseVarAssignStmt()
	}
}

func parseStmtTail(p *Parser, pg *ProgramNode) *error.Error {
	t := p.current().Type
	if t != lex.TypeStmtSep {
		return nil
	}
	p.consume(lex.TypeStmtSep)

	// #1. parse stmt item
	stmt, err := parseStmtItem(p)
	if err != nil {
		return err
	}
	pg.Children = append(pg.Children, stmt)
	// #2. parse stmt tail
	return parseStmtTail(p, pg)
}
