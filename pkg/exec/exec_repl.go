package exec

import (
	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/syntax/zh"
)

type REPLExecutor struct {
	context *r.Context
}

func NewREPLExecutor() *REPLExecutor {
	return &REPLExecutor{
		context: r.NewContext(GlobalValues, r.NewAnonymousModule(nil)),
	}
}

// RunCode - run code as one input line
func (pl *REPLExecutor) RunCode(text string) (r.Value, error) {
	in := io.NewByteStream([]byte(text))

	// #1. read source code
	source, err := in.ReadAll()
	if err != nil {
		return nil, err
	}

	// #2. get lexer
	lexer := syntax.NewLexer(source)
	module := pl.context.GetCurrentModule()
	module.SetLexer(lexer)

	// #3. execute code
	return pl.execREPLCode(lexer)
}

func (pl *REPLExecutor) execREPLCode(lexer *syntax.Lexer) (r.Value, error) {
	c := pl.context
	// #1. construct parser
	p := syntax.NewParser(lexer, zh.NewParserZH())

	// #2. parse program
	program, err := p.Parse()
	if err != nil {
		return nil, WrapSyntaxError(lexer, c.GetCurrentModule(), err)
	}

	if err := evalProgram(c, program); err != nil {
		return nil, WrapRuntimeError(c, err)
	}

	return c.GetCurrentScope().GetReturnValue(), nil
}
