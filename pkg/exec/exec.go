package exec

import (
	"errors"
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"os"
	"path/filepath"
	"strings"
)

// ExecuteModule - execute program from input Zn code (whether input source is a file or REPL)
func ExecuteModule(c *r.Context, name string, rootDir string) (r.Value, error) {
	// #1. find filepath of current module
	path, err := getModulePath(name, rootDir)
	if err != nil {
		return nil, err
	}

	// #2. read source code
	in, err := io.NewFileStream(path)
	if err != nil {
		return nil, err
	}
	source, err := in.ReadAll()
	if err != nil {
		return nil, err
	}

	lexer := syntax.NewLexer(source)
	p := syntax.NewParser(lexer, c.GetBuilder())

	// #3. parse program
	program, err := p.Parse()
	if err != nil {
		return nil, err
	}

	// #4. create module
	module := r.NewModule(name)
	c.PushScope(module)
	defer c.PopScope()

	if err := evalProgram(c, program); err != nil {
		return nil, err
	}

	returnValue := c.GetCurrentScope().GetReturnValue()
	return returnValue, nil
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

// getModulePath - get filepath of current module relative to rootPath
func getModulePath(name string, rootDir string) (string, error) {
	dirs := strings.Split(name, "-")
	// add .zn for last item
	dirs[len(dirs)-1] = dirs[len(dirs)-1] + ".zn"

	path := filepath.Join(rootDir, filepath.Join(dirs...))
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return "", zerr.ModuleNotFound(name)
	}

	return path, nil
}