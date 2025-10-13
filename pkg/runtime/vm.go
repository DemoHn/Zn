package runtime

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
	modules *ModuleGraph
}

type ElementMap = map[string]Element

func InitVM(globals map[string]Element) *VM {
	return &VM{
		globals:    globals,
		valueStack: map[int]Scope{},
		callStack:  []CallFrame{},
		csCount:    0,
		csModuleID: -1,
		modules:    nil,
	}
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

func (vm *VM) GetCurrentCallFrame() *CallFrame {
	return vm.getCurrentCallFrame()
}

func (vm *VM) GetThisValue() Element {
	return vm.getCurrentCallFrame().thisValue
}

// // internal functions
func (vm *VM) getCurrentCallFrame() *CallFrame {
	return &vm.callStack[vm.csCount-1]
}

func (vm *VM) initValueStack(moduleID int) {
	if _, ok := vm.valueStack[moduleID]; !ok {
		vm.valueStack[moduleID] = NewScope()
	}
}
