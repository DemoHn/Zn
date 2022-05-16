package exec

import (
	"fmt"
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"path/filepath"
	"strings"
)

// ExecuteModule - execute program from input Zn code (whether input source is a file or REPL)
func ExecuteModule(c *r.Context, in *io.FileStream, rootPath string) (r.Value, error) {
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

	moduleName, err := findModuleName(in.GetPath(), rootPath)
	if err != nil {
		return nil, err
	}
	// #3. create module
	module := r.NewModule(moduleName, lexer)
	c.PushScope(module)

	if err := evalProgram(c, program); err != nil {
		return nil, err
	}

	return c.GetCurrentScope().GetReturnValue(), nil
}

// ExecuteREPLCode - execute program from input Zn code (whether input source is a file or REPL)
func ExecuteREPLCode(c *r.Context, in io.InputStream) (r.Value, error) {
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

	if err := evalProgram(c, program); err != nil {
		return nil, err
	}

	return c.GetCurrentScope().GetReturnValue(), nil
}

// findModuleName - currentPath MUST BE in the directory and its child directories
// of rootPath
// Since we only support absolute path right now
func findModuleName(currentPath string, rootPath string) (string, error) {
	relFile, err := filepath.Rel(filepath.Dir(rootPath), currentPath)
	if err != nil {
		return "", zerr.NewErrorSLOT("定位模块位置出错")
	}

	relDir, relFileName := filepath.Split(relFile)
	// remove .zn   eg. A.zn -> A
	relFileName = strings.TrimSuffix(relFileName, filepath.Ext(relFileName))
	if relDir == "" {
		return relFileName, nil
	}

	relDir = strings.ReplaceAll(relDir, "/", "-")
	relDir = strings.ReplaceAll(relDir, "\\", "-")

	return fmt.Sprintf("%s-%s", relDir, relFileName), nil
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
