package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/syntax"
)

type VM struct {
	globals map[string]Element
	// module-level & scope-level local value stacks
	// including global & local variables
	// KEY: moduleID
	// VALUE: the ScopeStack of corresponding module
	valueStack map[int]Scope

	// callStack - store all call frames
	callStack []CallFrame
	// csCount - length of callStack
	csCount int
	// csModuleID - index of current module (the moduleID at the top of callStack)
	csModuleID int

	// modules - allocates modules by ID & stores export values
	modules []Module
	// moduleGraph - record a module dependency graph to detect circular dependency
	moduleGraph *ModuleGraph

	// moduleCodeFinder - HOWTO get the source code of a module
	moduleCodeFinder ModuleCodeFinder
}

type ElementMap = map[string]Element

func InitVM(globals map[string]Element) *VM {
	return &VM{
		globals:          globals,
		valueStack:       map[int]Scope{},
		callStack:        []CallFrame{},
		csCount:          0,
		csModuleID:       -1, // 0 for main module
		modules:          []Module{},
		moduleCodeFinder: nil,
	}
}

func (vm *VM) GetModuleCodeFinder() ModuleCodeFinder {
	return vm.moduleCodeFinder
}

func (vm *VM) SetModuleCodeFinder(moduleCodeFinder ModuleCodeFinder) {
	vm.moduleCodeFinder = moduleCodeFinder
}

// AllocateModule - create empty module information
func (vm *VM) AllocateModule(name string, program *syntax.Program) *Module {
	extModuleID := vm.moduleGraph.AddModule(vm.csModuleID, name, program)
	vm.csModuleID = extModuleID

	return vm.moduleGraph.GetModuleByID(extModuleID)
}

func (vm *VM) FindModuleByName(name string) *Module {
	moduleID, exists := vm.moduleGraph.GetIDFromName(name)
	if !exists {
		return nil
	}
	return vm.moduleGraph.GetModuleByID(moduleID)
}

func (vm *VM) CheckDepedency(name string) error {
	moduleID, exists := vm.moduleGraph.GetIDFromName(name)
	if exists {
		// check circular dependency
		if vm.moduleGraph.CheckCircularDepedency(vm.csModuleID, moduleID) {
			return zerr.ModuleCircularDependency()
		}
	}
	// no existing module found - no dependency problem will be found
	return nil
}

// PushCallFrame - push a call frame onto the call stack
// and update the current call stack cursor accordingly.
func (vm *VM) PushCallFrame(callFrame *CallFrame) {
	vm.callStack = append(vm.callStack, *callFrame)
	vm.csCount += 1
	vm.csModuleID = callFrame.moduleID
	vm.initValueStack(vm.csModuleID)
}

func (vm *VM) PopCallFrame() *CallFrame {
	if vm.csCount < 0 {
		return nil
	}
	vm.csCount -= 1

	currentCF := &vm.callStack[vm.csCount-1]
	vm.csModuleID = currentCF.moduleID
	return currentCF
}

func (vm *VM) GetCurrentModule() *Module {
	return vm.moduleGraph.GetModuleByID(vm.csModuleID)
}

func (vm *VM) GetCurrentCallFrame() *CallFrame {
	return vm.getCurrentCallFrame()
}

func (vm *VM) GetThisValue() Element {
	callFrame := vm.getCurrentCallFrame()
	if callFrame != nil {
		return callFrame.thisValue
	}
	return nil
}

func (vm *VM) BeginScope() {
	scope := vm.getCurrentScope()
	if scope != nil {
		scope.BeginScope()
	}
}

// EndScope - end current scope
func (vm *VM) EndScope() {
	scope := vm.getCurrentScope()
	if scope != nil {
		scope.EndScope()
	}
}

func (vm *VM) GetCurrentModuleID() int {
	return vm.csModuleID
}

// SetCurrentLine
func (vm *VM) SetCurrentLine(line int) {
	frame := vm.getCurrentCallFrame()
	if frame != nil {
		frame.SetCurrentLine(line)
	}
}

func (vm *VM) FindElement(name *IDName) (Element, error) {
	nameStr := name.GetLiteral()
	// look for global values first
	if elem, ok := vm.globals[nameStr]; ok {
		return elem, nil
	}
	// then look for local values
	elem := vm.getCurrentScope().GetValue(nameStr)
	if elem == nil {
		return nil, zerr.NameNotDefined(nameStr)
	}
	return elem, nil
}

func (vm *VM) FindElementWithModule(name *IDName) (Element, *Module, error) {
	nameStr := name.GetLiteral()
	// look for global values first
	if elem, ok := vm.globals[nameStr]; ok {
		return elem, vm.moduleGraph.GetModuleByID(vm.csModuleID), nil
	}
	// then look for local values
	elem, moduleID := vm.getCurrentScope().GetValueWithModuleID(nameStr)
	if elem == nil {
		return nil, nil, zerr.NameNotDefined(nameStr)
	}

	extModuleID := vm.csModuleID
	if moduleID >= 0 {
		extModuleID = moduleID
	}
	return elem, vm.moduleGraph.GetModuleByID(extModuleID), nil
}

// DeclareElement
func (vm *VM) DeclareElement(name *IDName, elem Element) error {
	scope := vm.getCurrentScope()
	nameStr := name.GetLiteral()
	if scope == nil {
		return zerr.NameNotDefined(nameStr)
	}
	if _, inGlobals := vm.globals[nameStr]; inGlobals {
		return zerr.NameRedeclared(nameStr)
	}
	return scope.DeclareValue(name.GetLiteral(), elem)
}

// DeclareConstElement -
func (vm *VM) DeclareConstElement(name *IDName, elem Element) error {
	scope := vm.getCurrentScope()
	nameStr := name.GetLiteral()
	if scope == nil {
		return zerr.NameNotDefined(nameStr)
	}
	if _, inGlobals := vm.globals[nameStr]; inGlobals {
		return zerr.NameRedeclared(nameStr)
	}
	return scope.DeclareConstValue(name.GetLiteral(), elem)
}

func (vm *VM) SetElement(name *IDName, elem Element) error {
	scope := vm.getCurrentScope()
	nameStr := name.GetLiteral()
	if scope == nil {
		return zerr.NameNotDefined(nameStr)
	}
	return scope.SetValue(nameStr, elem)
}

// // internal functions
func (vm *VM) getCurrentCallFrame() *CallFrame {
	return &vm.callStack[vm.csCount-1]
}

func (vm *VM) getCurrentScope() *Scope {
	if scope, ok := vm.valueStack[vm.csModuleID]; ok {
		return &scope
	}
	return nil
}

func (vm *VM) initValueStack(moduleID int) {
	if _, ok := vm.valueStack[moduleID]; !ok {
		vm.valueStack[moduleID] = NewScope()
	}
}
