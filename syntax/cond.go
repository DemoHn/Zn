package syntax

import "github.com/DemoHn/Zn/error"

// ConditionStmt - conditional (if-else) statement
type ConditionStmt struct {
	// if
	IfTrueExpr  Expression
	IfTrueBlock BlockStmt
	// else
	IfFalseExpr  Expression
	IfFalseBlock BlockStmt
	// else if
	OtherExprs  []Expression
	OtherBlocks []BlockStmt
	// else-branch exists or not
	hasIfFalse bool
}

// BlockStmt - block stmts consists of statements that share
// same indent number
type BlockStmt struct {
	Content []Statement
}

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
func ParseCondStmt(p *Parser) (*ConditionStmt, *error.Error) {

}
