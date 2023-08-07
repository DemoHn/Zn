package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
)

// Context is a global variable that stores current execution
// states, global configurations
type Context struct {
	// globals - stores all global variables
	globals map[string]Element
	// hasPrinted - if stdout has been used to output message before program end, set `hasPrinted` -> true; so that after message is done
	hasPrinted bool

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

//// setters
func (ctx *Context) SetModuleCodeFinder(finder ModuleCodeFinder) {
	ctx.moduleCodeFinder = finder
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
	ctx.GetCurrentModule().SetCurrentLine(line)
}

//// enter & exist modules
func (ctx *Context) EnterModule(module *Module) {
	// add require cache
	ctx.moduleGraph.AddRequireCache(module)
	ctx.moduleStack = append(ctx.moduleStack, module)
}

func (ctx *Context) ExitModule() {
	sLen := len(ctx.moduleStack)
	if sLen > 0 {
		ctx.moduleStack = ctx.moduleStack[:sLen-1]
	}
}

func (ctx *Context) FindModuleCache(name string) *Module {
	return ctx.moduleGraph.FindRequireCache(name)
}

func (ctx *Context) CheckDepedency(depModule string) error {
	sourceModule := ctx.GetCurrentModule().GetName()
	// if target = source module
	if depModule == sourceModule {
		return zerr.ImportSameModule(depModule)
	}
	// if same module import twice
	if deps, ok := ctx.moduleGraph.depGraph[sourceModule]; ok {
		for _, dep := range deps {
			if dep == depModule {
				return zerr.DuplicateModule(depModule)
			}
		}
	}

	// add dep graph

	if ctx.moduleGraph.CheckCircularDepedency(sourceModule) {

	}
	return nil
}

//// scope symbols getters / setters

// FindSymbol - find symbol in the context from current scope
// up to its root scope
func (ctx *Context) FindSymbol(name string) (Element, error) {
	// find on globals first
	if symVal, inGlobals := ctx.globals[name]; inGlobals {
		return symVal, nil
	}

	return ctx.GetCurrentModule().FindScopeValue(name)
}

// SetSymbol -
func (ctx *Context) SetSymbol(name string, value Element) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	// ...then in symbols
	return ctx.GetCurrentModule().SetScopeValue(name, value)
}

// BindSymbol - bind non-const value with re-declaration check on same scope
func (ctx *Context) BindSymbol(name string, value Element) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}

	return ctx.GetCurrentModule().BindSymbol(name, value, false, true)
}

// BindSymbolDecl - bind value for declaration statement - that variables could be re-bind.
func (ctx *Context) BindSymbolDecl(name string, value Element, isConst bool) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}

	return ctx.GetCurrentModule().BindSymbol(name, value, isConst, false)
}

// BindScopeSymbolDecl - bind value for declaration statement - that variables could be re-bind.
func (ctx *Context) BindScopeSymbolDecl(scope *Scope, name string, value Element) error {
	if _, inGlobals := ctx.globals[name]; inGlobals {
		return zerr.NameRedeclared(name)
	}
	if scope != nil {
		scope.SetSymbolValue(name, value, false)
	}
	return nil
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
