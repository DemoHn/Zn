//// evalFunctionCall() is a core procedure, but the logic is very complicated. Thus we put the logic into a separate file.

package exec

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/value"
)

// compileFunction - create a Function object (with default param handler logic)
// from Zn code (*syntax.BlockStmt). It's the constructor of 如何XX or (anoymous function in the future)
func compileFunction(vm *r.VM, node *syntax.FunctionDeclareStmt) *value.Function {
	var mainLogicHandler = func(params []r.Element) (r.Element, error) {
		// 2. do eval exec block
		return evalExecBlock(vm, node.ExecBlock, params)
	}

	return value.NewFunction(mainLogicHandler)
}

// （显示：A、B、C），得到D
func evalFunctionCall(vm *r.VM, expr *syntax.FuncCallExpr) (r.Element, error) {
	// match & get funcName
	funcName, err := MatchIDName(expr.FuncName)
	if err != nil {
		return nil, err
	}

	// exec params
	params, err := exprsToValues(vm, expr.Params)
	if err != nil {
		return nil, err
	}

	resultVal, err := execDirectFunction(vm, funcName, params)
	if err != nil {
		return nil, err
	}

	// if exec function call succeed, then the non-nil `resultVal` will be exported.
	// However, if `得到 [someVar]` semi statement is defined, we will bind the `resultVal` to `someVar` first before ending the procedure.
	if expr.YieldResult != nil {
		// add yield result
		ytag, err := MatchIDName(expr.YieldResult)
		if err != nil {
			return nil, err
		}
		// bind yield result
		if err := vm.DeclareConstElement(ytag, resultVal); err != nil {
			return nil, err
		}
	}

	// return result
	return resultVal, nil
}

func execMethodFunction(vm *r.VM, root r.Element, funcName *r.IDName, params []r.Element) (r.Element, error) {
	switch robj := root.(type) {
	case *value.Object:
		_, refModule, err := vm.FindElementWithModule(r.NewIDName(robj.GetObjectName()))
		if err != nil {
			return nil, err
		}
		fnCallFrame := r.NewFunctionCallFrame(refModule, root)
		vm.PushCallFrame(fnCallFrame)
	default:
		// for other types, we suppose it is from native code -
		// usually for internal types like Number, String, Boolean, etc.
		fnCallFrame := r.NewFunctionCallFrame(r.NativeCodeModule, root)
		vm.PushCallFrame(fnCallFrame)
	}

	elem, err := root.ExecMethod(funcName.GetLiteral(), params)
	if err == nil {
		vm.PopCallFrame()
	}

	return elem, err
}

// direct function: defined as standalone function instead of the method of
// a model
func execDirectFunction(vm *r.VM, funcName *r.IDName, params []r.Element) (r.Element, error) {
	elem, module, err := vm.FindElementWithModule(funcName)
	if err != nil {
		return nil, err
	}
	// pushCallFrame
	fnCallFrame := r.NewFunctionCallFrame(module, nil)
	vm.PushCallFrame(fnCallFrame)

	// assert value is function type
	fn, ok := elem.(*value.Function)
	if !ok {
		return nil, zerr.InvalidFuncVariable(funcName.GetLiteral())
	}

	if elem, err := fn.Exec(nil, params); err != nil {
		return nil, err
	} else {
		vm.PopCallFrame()
		return elem, nil
	}
}
