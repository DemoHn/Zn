package runtime

import "github.com/DemoHn/Zn/pkg/syntax"

type ExecCursor struct {
	*syntax.Lexer
	currentLine int
}
