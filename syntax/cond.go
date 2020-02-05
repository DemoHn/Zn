package syntax

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
)

// ConditionStmt - conditional (if-else) statement
type ConditionStmt struct {
	// if
	IfTrueExpr  Expression
	IfTrueBlock *BlockStmt
	// else
	IfFalseBlock *BlockStmt
	// else if
	OtherExprs  []Expression
	OtherBlocks []*BlockStmt
	// else-branch exists or not
	HasIfFalse bool
}

// BlockStmt - block stmts consists of statements that share
// same indent number
type BlockStmt struct {
	Content []Statement
}

func (cs *ConditionStmt) statementNode() {}

// ParseCondStmt - yield ConditionStmt node
// CFG:
// CondStmt -> 如果 IfTrueExpr ：
//         ...     IfTrueBlock
//
//          -> 如果 IfTrueExpr ：
//         ...     IfTrueBlock
//         ... 否则 ：
//         ...     IfFalseBlock
//
//          -> 如果 IfTrueExpr ：
//         ...     IfTrueBlock
//         ... 再如 OtherExpr1 ：
//         ...     OtherBlock1
//         ... 再如 OtherExpr2 ：
//         ...     OtherBlock2
//         ... ....
//             否则 ：
//         ...     IfFalseBlock
func ParseCondStmt(p *Parser, mainIndent int) (*ConditionStmt, *error.Error) {
	var condExpr Expression
	var err *error.Error

	var stmt = new(ConditionStmt)

	var condKeywords = []lex.TokenType{
		lex.TypeCondElseW,
		lex.TypeCondOtherW,
	}
	// by definition, the first Branch (if-branch) is required,
	// and the 如果 (if) keyword has been consumed before this function call.
	//
	// thus for other branches (like else-branch and elseif-branch),
	// we should consume the corresponding keyword explicitly. (否则，再如)
	var firstBranch = true
	var branchTokenType = lex.TypeCondW

	for p.peek().Type != lex.TypeEOF {
		// when firstBranch = true, we are parsing (if-branch) by default
		if !firstBranch {
			// check indent
			if p.getPeekIndent() != mainIndent {
				return stmt, nil
			}
			// consume other condKeywords (否则，再如)
			if match, tk := p.tryConsume(condKeywords); match {
				branchTokenType = tk.Type
			} else {
				// for other keywords, this parsing task has
				return stmt, nil
			}
		}
		p.setLineMask(modeInline)
		// #1 parse expression
		if branchTokenType != lex.TypeCondElseW {
			condExpr, err = ParseExpression(p)
			if err != nil {
				return nil, err
			}
		}
		// #2. parse colon
		if err := p.consume(lex.TypeFuncCall); err != nil {
			return nil, err
		}
		// #3 parse block stmts
		p.unsetLineMask(modeInline)
		ok, blockIndent := p.expectBlockIndent()
		if !ok {
			return nil, error.NewErrorSLOT("unexpected indent")
		}
		blockStmt, err := ParseBlockStmt(p, blockIndent)
		if err != nil {
			return nil, err
		}

		// #4. merge data
		switch branchTokenType {
		case lex.TypeCondW:
			stmt.IfTrueExpr = condExpr
			stmt.IfTrueBlock = blockStmt
			firstBranch = false
		case lex.TypeCondElseW:
			// set false-branch flag
			stmt.HasIfFalse = true
			stmt.IfFalseBlock = blockStmt
		case lex.TypeCondOtherW:
			stmt.OtherExprs = append(stmt.OtherExprs, condExpr)
			stmt.OtherBlocks = append(stmt.OtherBlocks, blockStmt)
		}
	}
	return stmt, nil
}

// ParseBlockStmt -
func ParseBlockStmt(p *Parser, blockIndent int) (*BlockStmt, *error.Error) {
	// TODO: in the future, we will merge ProgramNode & BlockStmt into one type
	// THIS IS A TMP SOLUTION!!
	bStmt := new(BlockStmt)
	pg := new(ProgramNode)

	for (p.peek().Type != lex.TypeEOF) && p.getPeekIndent() == blockIndent {
		if err := ParseStatement(p, pg); err != nil {
			return nil, err
		}
	}

	// copy data from pg to bSttmt
	for _, pcs := range pg.Children {
		bStmt.Content = append(bStmt.Content, pcs)
	}

	return bStmt, nil
}
