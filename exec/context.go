package exec

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/lex"
	"github.com/DemoHn/Zn/syntax"
)

// SymbolInfo - symbol info
type SymbolInfo struct {
	nestLevel  int
	value      ZnValue
	isConstant bool // if isConstant = true, the value of this symbol is prohibited from any modification.
}

// Context - GLOBAL execution context, usually create only once in one program.
type Context struct {
	symbols map[string][]SymbolInfo
	globals map[string]ZnValue
	arith *Arith
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
	return &Context{
		symbols: map[string][]SymbolInfo{},
		globals: predefinedValues,
		arith: NewArith(defaultPrecision),
	}
}


// BindSymbol - add value to symbol table
func (ctx *SymbolTable) BindSymbol(nestLevel int, id string, obj ZnValue, isConstant bool) *error.Error {
	newInfo := SymbolInfo{
		nestLevel:  nestLevel,
		value:      obj,
		isConstant: isConstant,
	}

	symArr, ok := ctx.symbols[id]
	if !ok {
		// init symbolInfo array
		ctx.symbols[id] = []SymbolInfo{newInfo}
		return nil
	}

	// check if there's variable re-declaration
	if len(symArr) > 0 && symArr[0].nestLevel == nestLevel {
		return error.NameRedeclared(id)
	}

	// prepend data
	ctx.symbols[id] = append([]SymbolInfo{newInfo}, ctx.symbols[id]...)
	return nil
}

// LookupSymbol - find the corresponded value from ID,
// if nothing found, return error
func (ctx *SymbolTable) LookupSymbol(id string) (ZnValue, *error.Error) {
	symArr, ok := ctx.symbols[id]
	if !ok {
		return nil, error.NameNotDefined(id)
	}

	// find the nearest level of value
	if symArr == nil || len(symArr) == 0 {
		return nil, error.NameNotDefined(id)
	}
	return symArr[0].value, nil
}

// SetSymbolValue - after variable is defined, set the value
func (ctx *Context) SetSymbolValue(id string, obj ZnValue) *error.Error {
	symArr, ok := ctx.symbols[id]
	if !ok {
		return error.NameNotDefined(id)
	}

	if symArr != nil && len(symArr) > 0 {
		symArr[0].value = obj
		if symArr[0].isConstant {
			return error.AssignToConstant()
		}
		return nil
	}

	return error.NameNotDefined(id)
}

func printSymbols(ctx *Context) string {
	strs := []string{}
	for k, symArr := range ctx.symbols {
		if symArr != nil {
			for _, symItem := range symArr {
				symStr := "ε"
				if symItem.value != nil {
					symStr = symItem.value.String()
				}
				strs = append(strs, fmt.Sprintf("‹%s, %d› => %s", k, symItem.nestLevel, symStr))
			}
		}

	}

	return strings.Join(strs, "\n")
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
