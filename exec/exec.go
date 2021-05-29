package exec

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/lex"
	"github.com/DemoHn/Zn/syntax"
)

// ExecuteCode - execute program from input Zn code (whether from file or REPL)
func ExecuteCode(c *ctx.Context, in *lex.InputStream) (ctx.Value, *error.Error) {
	program, err := parseCode(in)
	if err != nil {
		return nil, err
	}
	// init fileInfo
	fileInfo := ctx.InitFileInfo(program.Lexer)
	c.SetFileInfo(fileInfo)

	// eval program
	return execProgram(c, program)
}

// parseCode - lex & parse code text
func parseCode(in *lex.InputStream) (*syntax.Program, *error.Error) {
	l := lex.NewLexer(in)
	p := syntax.NewParser(l)
	// start
	block, err := p.Parse()
	if err != nil {
		return nil, err
	}

	return syntax.NewProgramNode(block, l), nil
}

func execProgram(c *ctx.Context, program *syntax.Program) (ctx.Value, *error.Error) {
	err := evalProgram(c, program)
	if err != nil {
		cursor := err.GetCursor()

		// wrapError if lineInfo is missing (mostly for non-syntax errors)
		// If lineInfo missing, then we will add current execution line and hide some part to
		// display errors properly.
		if cursor.LineNum == 0 {
			fileInfo := c.GetFileInfo()
			newCursor := error.Cursor{
				File:    fileInfo.File,
				LineNum: fileInfo.CurrentLine,
				Text:    fileInfo.LineStack.GetLineText(fileInfo.CurrentLine, false),
			}
			err.SetCursor(newCursor)
		}
		return nil, err
	}
	return c.GetScope().GetReturnValue(), nil
}
