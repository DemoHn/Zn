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
		context: r.NewContext(globalValues, r.NewMainModule(nil)),
	}
}

// RunCode - run code as one input line
func (pl *REPLExecutor) RunCode(text string) (r.Element, error) {
	in := io.NewByteStream([]byte(text))
	c := pl.context

	// #1. read source code
	source, err := in.ReadAll()
	if err != nil {
		return nil, err
	}
	// #2. parse code
	parser := syntax.NewParserFromSource(source, zh.NewParserZH())
	program, err := parser.Parse()
	if err != nil {
		return nil, WrapSyntaxError(parser, "", err)
	}

	// #3. bind new source text lines
	c.GetCurrentModule().SetSourceLines(program.Lines)

	// #4. execute code
	if err := evalProgram(c, program); err != nil {
		return nil, WrapRuntimeError(c, err)
	}

	return c.GetCurrentScope().GetReturnValue(), nil
}
