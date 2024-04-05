package zinc

import (
	"github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
)

//// This file provides ALL types for customizing the Zinc Language and services.
type Element = runtime.Element

/// collect some interfaces we are currently using to help users extend more implementations!
type ZnExecutor interface {
	RunCode(varInputs map[string]Element)
}

type ZnServer interface {
	Listen(connUrl string)
}

type ZnASTBuilder interface {
	ParseAST(lexer *syntax.Lexer) (*syntax.Program, error)
}

type ZnElement interface {
	GetProperty(*runtime.Context, string) (Element, error)
	SetProperty(*runtime.Context, string, Element) error
	ExecMethod(*runtime.Context, string, []Element) (Element, error)
}

type ZnError interface {
	Error() string
}
