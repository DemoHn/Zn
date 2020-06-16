package exec

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
	"github.com/DemoHn/Zn/syntax"
)

// Context - code lifecycle management
// TODO: this is a tmp solution. in the future, we will
// gradually obselete this tree-walk based interperter.
type Context struct {
	*SymbolTable
	*ArithInstance
	// lastValue is set during the execution, usually stands for 'the return value' of a function.
	lastValue ZnValue
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
	ctx := new(Context)
	ctx.SymbolTable = NewSymbolTable()
	ctx.ArithInstance = NewArithInstance(defaultPrecision)
	ctx.lastValue = NewZnNull()
	return ctx
}

// ExecuteCode - execute program from input Zn code (whether from file or REPL)
func (ctx *Context) ExecuteCode(in *lex.InputStream, scope Scope) Result {
	l := lex.NewLexer(in)
	p := syntax.NewParser(l)
	// start
	block, err := p.Parse()
	if err != nil {
		return Result{true, nil, err}
	}
	// After parsing, lines are split & cached completely.
	// It's time to initialize lexScope
	ctx.initLexScope(l)

	// construct root (program) node
	program := syntax.NewProgramNode(block)

	if err := EvalProgram(ctx, scope, program); err != nil {
		wrapError(ctx, err)
		return Result{true, nil, err}
	}
	return Result{false, ctx.lastValue, nil}
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

	lastValue := ctx.lastValue
	if ctx.lastValue == nil {
		lastValue = NewZnNull()
	}
	return Result{false, lastValue, nil}
}

// ResetLastValue - set ctx.lastValue -> nil
func (ctx *Context) ResetLastValue() {
	ctx.lastValue = nil
}

// lexScope helpers
func (ctx *Context) initLexScope(l *lex.Lexer) {
	ctx.lexScope = lexScope{
		file:        l.InputStream.Scope,
		currentLine: 0,
		lineStack:   l.LineStack,
	}
}

func (ctx *Context) setCurrentLine(line int) {
	ctx.lexScope.currentLine = line
}

func (ctx *Context) getFile() string {
	return ctx.lexScope.file
}

func (ctx *Context) getCurrentLine() int {
	return ctx.lexScope.currentLine
}

func (ctx *Context) getCurrentLineText() string {
	ls := ctx.lexScope
	txt := ls.lineStack.GetLineText(ctx.currentLine, false)
	return txt
}

//// Execute (Evaluate) statements

// EvalProgram - evaluate global program (root node)
func EvalProgram(ctx *Context, scope Scope, program *syntax.Program) *error.Error {
	return evalStmtBlock(ctx, scope, program.Content)
}

// wrapError if lineInfo is missing (mostly for non-syntax errors)
// If lineInfo missing, then we will add current execution line and hide some part to
// display errors properly.
func wrapError(ctx *Context, err *error.Error) {
	cursor := err.GetCursor()
	if cursor.LineNum == 0 {
		newCursor := error.Cursor{
			File:    ctx.getFile(),
			LineNum: ctx.getCurrentLine(),
			Text:    ctx.getCurrentLineText(),
		}
		err.SetCursor(newCursor)
	}
}
