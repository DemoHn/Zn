package exec

import (
	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/syntax/zh"
)

// ExecuteModule - execute program from input Zn code (whether input source is a file or REPL)
func ExecuteModule(c *r.Context, name string) (r.Value, error) {
	// #1. find filepath of current module
	path, err := c.GetModulePath(name)
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
	p := syntax.NewParser(lexer, zh.NewParserZH())

	// #3. parse program
	program, err := p.Parse()
	if err != nil {
		return nil, WrapSyntaxError(lexer, name, err)
	}

	// #4. create module
	module := r.NewModule(name)
	// build module cache to dep tree
	c.BuildModuleCache(module)
	// create new scope (and pop the scope after execution)
	c.PushScope(module, lexer)
	defer c.PopScope()

	if err := evalProgram(c, program); err != nil {
		return nil, err
	}

	returnValue := c.GetCurrentScope().GetReturnValue()
	return returnValue, nil
}

// ExecuteREPLCode - execute program from input Zn code (whether input source is a file or REPL)
func ExecuteREPLCode(c *r.Context, lexer *syntax.Lexer) (r.Value, error) {
	// #1. construct parser
	p := syntax.NewParser(lexer, zh.NewParserZH())

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