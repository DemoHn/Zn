package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
)

// Context is a global variable that stores current execution
// states, global configurations
type Context struct {
	// globals - stores all global variables
	globals map[string]Element

	// varInputs - predefine input values before executing the program
	// used for input stmt: `输入XX、YY、ZZ`
	varInputs map[string]Element

	// hasPrinted - if stdout has been used to output message before program end, set `hasPrinted` -> true; so that after message is done
	hasPrinted bool

	// currentLine - current execution lineNum (index, start from 0)
	currentLine      int
	currentRefModule *Module

	// modulegraph - store module dependency & all preloaded modules
	moduleGraph *ModuleGraph

	// moduleCodeFinder - given a module name, the finder function aims to find it's corresponding source code for further execution - whatever from filesystem, DB, network, etc.
	// by default, the value is nil, that means the finder could not found any module code at all!
	moduleCodeFinder ModuleCodeFinder
	// callStack - get current call module & line for traceback
	callStack []CallInfo
	// current execution moduleStack. The last one represents for current execution module. Must be NON-EMPTY at initialization
	moduleStack []*Module
}

type CallInfo struct {
	*Module
	LastLineIdx int
}

/* ModuleCodeFinder - input module name, output its source code or return error. The example finder shows how to find source code from module name, where each module corresponds to a "<moduleName>.zn" text file on disk.

```go
func finder (name string) ([]rune, error) {
	path := fmt.Sprintf("./%s.zn", name)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("source code of module '%s' not found", name)
	}

	in, err := io.NewFileStream(path)
	if err != nil {
		return err
	}

	return in.ReadAll()
}
```
*/
type ModuleCodeFinder func(string) ([]rune, error)

// NewContext - create new Zn Context. Notice through the life-cycle
// of one code execution, there's only one running context to store all states.
// NOTE: initModule DO NOT accept nil value at initialization!!
func NewContext(globalsMap map[string]Element, initModule *Module) *Context {
	// init module dep graph
	graph := NewModuleGraph()
	graph.AddRequireCache(initModule)

	return &Context{
		globals:          globalsMap,
		hasPrinted:       false,
		moduleGraph:      graph,
		moduleCodeFinder: nil,
		callStack:        []CallInfo{},
		moduleStack:      []*Module{initModule},
		currentRefModule: initModule,
	}
}

//// getters

func (ctx *Context) GetCurrentModule() *Module {
	sLen := len(ctx.moduleStack)
	if sLen > 0 {
		return ctx.moduleStack[sLen-1]
	}
	return nil
}

func (ctx *Context) GetCallStack() []CallInfo {
	return ctx.callStack
}

func (ctx *Context) GetCurrentScope() *Scope {
	module := ctx.GetCurrentModule()
	if module != nil {
		return module.GetCurrentScope()
	}
	return nil
}

func (ctx *Context) GetHasPrinted() bool {
	return ctx.hasPrinted
}

func (ctx *Context) GetModuleCodeFinder() ModuleCodeFinder {
	return ctx.moduleCodeFinder
}

func (ctx *Context) GetVarInputs() map[string]Element {
	return ctx.varInputs
}

//// setters
func (ctx *Context) SetModuleCodeFinder(finder ModuleCodeFinder) {
	ctx.moduleCodeFinder = finder
}

func (ctx *Context) SetVarInputs(varInputs map[string]Element) {
	ctx.varInputs = varInputs
}

//// scope operation
func (ctx *Context) FindParentScope() *Scope {
	module := ctx.GetCurrentModule()
	if module != nil {
		sLen := len(module.scopeStack)

		if sLen > 1 {
			return module.scopeStack[sLen-2]
		}
	}
	return nil
}

// PushScope - create new scope with same module from parent scope
func (ctx *Context) PushScope() *Scope {
	sp := ctx.GetCurrentScope()
	if sp == nil {
		return nil
	}

	return ctx.GetCurrentModule().PushScope()
}

func (ctx *Context) PopScope() {
	ctx.GetCurrentModule().PopScope()
}

// SetCurrentLine - set lineIdx to current running scope
func (ctx *Context) SetCurrentLine(line int) {
	ctx.currentLine = line
}

func (ctx *Context) GetCurrentLine() int {
	return ctx.currentLine
}

func (ctx *Context) SetCurrentRefModule(refModule *Module) {
	ctx.currentRefModule = refModule
}

func (ctx *Context) GetCurrentRefModule() *Module {
	return ctx.currentRefModule
}

func (ctx *Context) PushCallStack() {
	ctx.callStack = append(ctx.callStack, CallInfo{
		Module:      ctx.currentRefModule,
		LastLineIdx: ctx.currentLine,
	})
}

func (ctx *Context) PopCallStack() {
	stackLen := len(ctx.callStack)
	if stackLen == 0 {
		return
	}
	ctx.callStack = ctx.callStack[:stackLen-1]
}

//// enter & exist modules
func (ctx *Context) EnterModule(module *Module) {
	// add require cache
	ctx.moduleGraph.AddRequireCache(module)
	ctx.moduleStack = append(ctx.moduleStack, module)
	ctx.currentRefModule = module
}

func (ctx *Context) ExitModule() {
	sLen := len(ctx.moduleStack)
	if sLen > 0 {
		ctx.moduleStack = ctx.moduleStack[:sLen-1]
		// get last one as refModule
		ctx.currentRefModule = ctx.moduleStack[len(ctx.moduleStack)-1]
	}
}

func (ctx *Context) FindModuleCache(name string) *Module {
	return ctx.moduleGraph.FindRequireCache(name)
}

func (ctx *Context) CheckDepedency(depModule string, internal bool) error {
	sourceModule := ctx.GetCurrentModule().GetName()
	// if target = source module (external)
	// NOTICE:external modules have SAME name with internal modules is allowed
	if depModule == sourceModule && !internal {
		return zerr.ImportSameModule(depModule)
	}
	// if same module import twice ()
	depList := ctx.moduleGraph.externalDepGraph
	if internal {
		depList = ctx.moduleGraph.internalDepGraph
	}
	if deps, ok := depList[sourceModule]; ok {
		for _, dep := range deps {
			if dep == depModule {
				return zerr.DuplicateModule(depModule)
			}
		}
	}

	// add dep graph
	if ctx.moduleGraph.CheckCircularDepedency(sourceModule, depModule, internal) {
		return zerr.ModuleCircularDependency()
	}
	return nil
}

//// scope symbols getters / setters

// FindElement - find symbol value in the context from current scope
// up to its root scope
func (ctx *Context) FindElement(name *IDName) (Element, error) {
	nameStr := name.GetLiteral()
	// find on globals first
	if symVal, inGlobals := ctx.globals[nameStr]; inGlobals {
		return symVal, nil
	}

	return ctx.GetCurrentModule().FindScopeValue(nameStr)
}

func (ctx *Context) FindSymbol(name *IDName) (SymbolInfo, error) {
	// find on globals first
	nameStr := name.GetLiteral()
	if symVal, inGlobals := ctx.globals[nameStr]; inGlobals {
		return SymbolInfo{
			value:   symVal,
			isConst: true,
			module:  nil,
		}, nil
	}

	return ctx.GetCurrentModule().FindScopeSymbol(nameStr)
}

// SetSymbol -
func (ctx *Context) SetSymbol(name *IDName, value Element) error {
	nameStr := name.GetLiteral()
	if _, inGlobals := ctx.globals[nameStr]; inGlobals {
		return zerr.NameRedeclared(nameStr)
	}
	// ...then in symbols
	return ctx.GetCurrentModule().SetScopeValue(nameStr, value)
}

// BindSymbol - bind non-const value with re-declaration check on same scope
func (ctx *Context) BindSymbol(name *IDName, value Element) error {
	nameStr := name.GetLiteral()
	if _, inGlobals := ctx.globals[nameStr]; inGlobals {
		return zerr.NameRedeclared(nameStr)
	}

	return ctx.GetCurrentModule().BindSymbol(nameStr, value, false, true)
}

// BindSymbolConst - bind const value with re-declaration check on same scope
func (ctx *Context) BindSymbolConst(name *IDName, value Element) error {
	nameStr := name.GetLiteral()
	if _, inGlobals := ctx.globals[nameStr]; inGlobals {
		return zerr.NameRedeclared(nameStr)
	}

	return ctx.GetCurrentModule().BindSymbol(nameStr, value, true, true)
}

// BindSymbolDecl - bind value for declaration statement - that variables could be re-bind.
func (ctx *Context) BindSymbolDecl(name *IDName, value Element, isConst bool) error {
	nameStr := name.GetLiteral()
	if _, inGlobals := ctx.globals[nameStr]; inGlobals {
		return zerr.NameRedeclared(nameStr)
	}

	return ctx.GetCurrentModule().BindSymbol(nameStr, value, isConst, false)
}

// BindScopeSymbolDecl - bind value for declaration statement - that variables could be re-bind.
func (ctx *Context) BindScopeSymbolDecl(scope *Scope, name *IDName, value Element) error {
	nameStr := name.GetLiteral()
	if _, inGlobals := ctx.globals[nameStr]; inGlobals {
		return zerr.NameRedeclared(nameStr)
	}
	if scope != nil {
		scope.SetSymbolValue(nameStr, value, false, ctx.GetCurrentModule())
	}
	return nil
}

// BindSymbolDecl - bind value for declaration statement - that variables could be re-bind.
func (ctx *Context) BindImportSymbol(name string, value Element, refModule *Module) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}

	return ctx.GetCurrentModule().BindImportSymbol(name, value, refModule)
}

// GetThisValue -
func (ctx *Context) GetThisValue() Element {
	m := ctx.GetCurrentModule()
	return m.GetCurrentScope().thisValue
}

// MarkHasPrinted - called by `显示` function only
func (ctx *Context) MarkHasPrinted() {
	ctx.hasPrinted = true
}
