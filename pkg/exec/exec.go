package exec

import (
	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/syntax/zh"
)

// ExecuteCode - execute program from input Zn code (whether input source is a file or REPL)
func ExecuteCode(c *r.Context, in io.InputStream) (r.Value, error) {
	program, err := parseCode(in)
	if err != nil {
		return nil, err
	}
	// eval program
	return execProgram(c, program)
}


// parseCode - lex & parse code text
func parseCode(in io.InputStream) (*syntax.Program, error) {
	// TODO: import more builders in more languages
	source, err := in.ReadAll()
	if err != nil {
		return nil, err
	}
	zhBuilder := zh.NewParserZH()
	l := syntax.NewLexer(source)
	p := syntax.NewParser(l, zhBuilder)

	return p.Parse()
}

func execProgram(c *r.Context, program *syntax.Program) (r.Value, error) {
	/** TODO: complete error display system
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
	 */
	return c.GetScope().GetReturnValue(), nil
}
