package exec

import (
	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/syntax/zh"
)

type PlaygroundExecutor struct {
	source  []byte
	context *r.Context
}

func NewPlaygroundExecutor(code []byte) *PlaygroundExecutor {
	return &PlaygroundExecutor{
		source:  code,
		context: r.NewContext(globalValues, r.NewAnonymousModule(nil)),
	}
}

// RunCode - ONE-TIME DEAL! run code text from fastCGI request's input body
func (pl *PlaygroundExecutor) RunCode() (r.Element, error) {
	in := io.NewByteStream(pl.source)

	// #1. read source code
	source, err := in.ReadAll()
	if err != nil {
		return nil, err
	}

	// #2. get lexer
	lexer := syntax.NewLexer(source)
	module := pl.context.GetCurrentModule()
	module.SetLexer(lexer)

	// #3. parse program
	parser := syntax.NewParser(lexer, zh.NewParserZH())
	program, err := parser.Parse()
	if err != nil {
		return nil, WrapSyntaxError(lexer, module, err)
	}

	// #4. eval code
	if err := evalProgram(pl.context, program); err != nil {
		return nil, WrapRuntimeError(pl.context, err)
	}

	return pl.context.GetCurrentScope().GetReturnValue(), nil
}
