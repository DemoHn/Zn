package exec

import (
	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
)

// ExecuteCode - execute program from input Zn code (whether input source is a file or REPL)
func ExecuteCode(c *r.Context, in io.InputStream) (r.Value, error) {
	// #1. read source code
	source, err := in.ReadAll()
	if err != nil {
		return nil, err
	}

	lexer := syntax.NewLexer(source)
	p := syntax.NewParser(lexer, c.GetBuilder())

	// #2. parse program
	program, err := p.Parse()
	if err != nil {
		return nil, err
	}

	// #3. create module
	module := r.NewModule("BOOM", lexer)
	c.PushScope(module)

	if err := evalProgram(c, program); err != nil {
		return nil, err
	}

	return c.GetCurrentScope().GetReturnValue(), nil
}


func execProgram(c *r.Context, program *syntax.Program) (r.Value, error) {
	m := r.NewModule("BOOM", nil)
	c.PushScope(m)

	err := evalProgram(c, program)
	if err != nil {
		return nil, err
	}
	/** TODO: complete error display system

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
	 */
	return c.GetCurrentScope().GetReturnValue(), nil
}
