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
}

type ElementMap = map[string]Element

func InitVM(globals map[string]Element) *VM {
	return &VM{
		globals:    globals,
		valueStack: map[int]Scope{},
		callStack:  []CallFrame{},
		csCount:    0,
		csModuleID: -1, // 0 for main module
		modules:    []Module{},
	}
}

// AllocateModule - create empty module information
func (vm *VM) AllocateModule(name string, program *syntax.Program) *Module {
	module := Module{
		fullName:     name,
		program:      program,
		exportValues: map[string]Element{},
	}

	vm.modules = append(vm.modules, module)
	vm.csModuleID = len(vm.modules) - 1

	return &module
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
	if vm.csModuleID >= 0 && vm.csModuleID < len(vm.modules) {
		return &vm.modules[vm.csModuleID]
	}
	return nil
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
		return elem, vm.getModuleByID(vm.csModuleID), nil
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
	return elem, vm.getModuleByID(extModuleID), nil
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

func (vm *VM) getModuleByID(moduleID int) *Module {
	if moduleID >= 0 && moduleID < len(vm.modules) {
		return &vm.modules[moduleID]
	}
	return nil
}

func (vm *VM) initValueStack(moduleID int) {
	if _, ok := vm.valueStack[moduleID]; !ok {
		vm.valueStack[moduleID] = NewScope()
	}
}
