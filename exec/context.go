package exec

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
	"github.com/DemoHn/Zn/syntax"
)

// SymbolInfo - symbol info
type SymbolInfo struct {
	NestLevel  int
	Value      ZnValue
	IsConstant bool // if isConstant = true, the value of this symbol is prohibited from any modification.
}

// Context - GLOBAL execution context, usually create only once in one program.
type Context struct {
	symbols map[string][]SymbolInfo
	globals map[string]ZnValue
	arith   *Arith
	lexScope
}

// lexScope defines current lex info of context.
// NOTICE: a Context is eligible to call ExecuteCode() many times, which means lexScope may
// varies from different input stream even in same context!
type lexScope struct {
	file string
	// current execution line - it's continuously changing within the execution process
	currentLine int
	lineStack   *lex.LineStack
}

const defaultPrecision = 8

// Result - context execution result structure
// NOTICE: when HasError = true, Value = nil, while execution yields error
//         when HasError = false, Error = nil, Value = <result Value>
//
// Currently only one value is supported as return argument.
type Result struct {
	HasError bool
	Value    ZnValue
	Error    *error.Error
}

// NewContext - create new Zn Context for furthur execution
func NewContext() *Context {
	return &Context{
		symbols: map[string][]SymbolInfo{},
		globals: predefinedValues,
		arith:   NewArith(defaultPrecision),
	}
}

// GetSymbols -
func (ctx *Context) GetSymbols() map[string][]SymbolInfo {
	return ctx.symbols
}

// ExecuteCode - execute program from input Zn code (whether from file or REPL)
func (ctx *Context) ExecuteCode(in *lex.InputStream, scope *RootScope) Result {
	l := lex.NewLexer(in)
	p := syntax.NewParser(l)
	// start
	block, err := p.Parse()
	if err != nil {
		return Result{true, nil, err}
	}
	// init scope
	scope.Init(l)

	// construct root (program) node
	program := syntax.NewProgramNode(block)

	// eval program
	if err := evalProgram(ctx, scope, program); err != nil {
		wrapError(ctx, err)
		return Result{true, nil, err}
	}
	return Result{false, scope.GetLastValue(), nil}
}

// ExecuteBlockAST - execute blockStmt AST
// usually for executing function template
func (ctx *Context) ExecuteBlockAST(scope Scope, block *syntax.BlockStmt) Result {
	if err := evalStmtBlock(ctx, scope, block); err != nil {
		// handle returnValue Interrupts
		if err.GetErrorClass() != error.InterruptsClass {
			wrapError(ctx, err)
			return Result{true, nil, err}
		}
	}

	return Result{false, NewZnNull(), nil}
}

// wrapError if lineInfo is missing (mostly for non-syntax errors)
// If lineInfo missing, then we will add current execution line and hide some part to
// display errors properly.
func wrapError(ctx *Context, err *error.Error) {
	/**
	cursor := err.GetCursor()

	if cursor.LineNum == 0 {

		newCursor := error.Cursor{
			File:    ctx.getFile(),
			LineNum: ctx.getCurrentLine(),
			Text:    ctx.getCurrentLineText(),
		}
		err.SetCursor(newCursor)
	}
	*/
}
