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
		context: r.NewContext(globalValues, r.NewMainModule(nil)),
	}
}

// RunCode - ONE-TIME DEAL! run code text from request's input body
func (pl *PlaygroundExecutor) RunCode(varInputs map[string]r.Element) (r.Element, error) {
	in := io.NewByteStream(pl.source)

	// #1. read source code
	source, err := in.ReadAll()
	if err != nil {
		return nil, err
	}

	// #2. get lexer
	module := pl.context.GetCurrentModule()

	// #3. parse program
	parser := syntax.NewParserFromSource(source, zh.NewParserZH())
	program, err := parser.Parse()
	if err != nil {
		return nil, WrapSyntaxError(parser, module.GetName(), err)
	}
	// set source lines
	module.SetSourceLines(program.Lines)

	// #4. eval code
	// #4.1 first set init input value
	pl.context.SetVarInputs(varInputs)
	// #4.2 then eval the program
	if err := evalProgram(pl.context, program); err != nil {
		return nil, WrapRuntimeError(pl.context, err)
	}

	return pl.context.GetCurrentScope().GetReturnValue(), nil
}
