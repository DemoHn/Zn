package exec

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/syntax"
)

type funcExecutor func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error)

type paramHandler func(ctx *Context, scope *FuncScope, params []ZnValue) *error.Error

// ClosureRef - aka. Closure Exection Reference
// This structure wraps the execution logic inside the closure
// statically
type ClosureRef struct {
	Name         string
	ParamHandler paramHandler // bind & validate params before actual execution
	Executor     funcExecutor // actual execution logic
}

// BuildClosureRefFromNode -
func BuildClosureRefFromNode(name string, paramTags []*syntax.ParamItem, stmtBlock *syntax.BlockStmt) *ClosureRef {

	var executor = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
		// iterate block round I - function hoisting
		for _, stmtI := range stmtBlock.Children {
			if v, ok := stmtI.(*syntax.FunctionDeclareStmt); ok {
				fn := BuildZnFunctionFromNode(v)
				if err := bindValue(ctx, scope, v.FuncName.GetLiteral(), fn); err != nil {
					return nil, err
				}
			}
		}
		// iterate block round II
		for _, stmtII := range stmtBlock.Children {
			if _, ok := stmtII.(*syntax.FunctionDeclareStmt); !ok {
				if err := evalStatement(ctx, scope, stmtII); err != nil {
					// if recv breaks
					if err.GetCode() == error.ReturnBreakSignal {
						if extra, ok := err.GetExtra().(ZnValue); ok {
							return extra, nil
						}
					}
					return nil, err
				}
			}
		}
		return scope.GetReturnValue(), nil
	}

	var paramHandler = func(ctx *Context, scope *FuncScope, params []ZnValue) *error.Error {
		// check param length
		if len(params) != len(paramTags) {
			return error.MismatchParamLengthError(len(paramTags), len(params))
		}

		// bind params (as variable) to function scope
		for idx, paramVal := range params {
			param := paramTags[idx]
			// if param is NOT a reference type, then we need additionally
			// copy its value
			if !param.RefMark {
				paramVal = duplicateValue(paramVal)
			}
			paramName := param.ID.GetLiteral()
			if err := bindValue(ctx, scope, paramName, paramVal); err != nil {
				return err
			}
		}
		return nil
	}

	return &ClosureRef{
		Name:         name,
		ParamHandler: paramHandler,
		Executor:     executor,
	}
}

// NewClosureRef - define native function
func NewClosureRef(name string, executor funcExecutor) *ClosureRef {
	return &ClosureRef{
		Name:         name,
		ParamHandler: nil,
		Executor:     executor,
	}
}

// Exec - exec function
func (cr *ClosureRef) Exec(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
	// handle params
	if cr.ParamHandler != nil {
		if err := cr.ParamHandler(ctx, scope, params); err != nil {
			return nil, err
		}
	}
	// do execution
	return cr.Executor(ctx, scope, params)
}

// ClassRef -
type ClassRef struct {
	Name        string
	Constructor funcExecutor           // a function to initialize all properties
	GetterList  map[string]*ClosureRef // stores defined getters inside the class
	MethodList  map[string]*ClosureRef // stores defined methods inside the class
}

// NewClassRef - create new empty ClassRef
func NewClassRef(name string) *ClassRef {
	return &ClassRef{
		Name:        name,
		Constructor: nil,
		GetterList:  map[string]*ClosureRef{},
		MethodList:  map[string]*ClosureRef{},
	}
}

// BuildClassRefFromNode -
func BuildClassRefFromNode(name string, classNode *syntax.ClassDeclareStmt) *ClassRef {
	ref := &ClassRef{
		Name:        name,
		Constructor: nil,
		GetterList:  map[string]*ClosureRef{},
		MethodList:  map[string]*ClosureRef{},
	}

	// define default constrcutor
	var constructor = func(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
		obj := NewZnObject(ref)
		// init prop list
		for _, propPair := range classNode.PropertyList {
			propID := propPair.PropertyID.GetLiteral()
			expr, err := evalExpression(ctx, scope, propPair.InitValue)
			if err != nil {
				return nil, err
			}
			obj.PropList[propID] = expr
		}
		// constructor: set some properties' value
		if len(params) != len(classNode.ConstructorIDList) {
			return nil, error.MismatchParamLengthError(len(params), len(classNode.ConstructorIDList))
		}
		for idx, objParamVal := range params {
			param := classNode.ConstructorIDList[idx]
			// if param is NOT a reference, then we need to copy its value
			if !param.RefMark {
				objParamVal = duplicateValue(objParamVal)
			}
			paramName := param.ID.GetLiteral()
			obj.PropList[paramName] = objParamVal
		}

		return obj, nil
	}
	// set constructor
	ref.Constructor = constructor

	// add getters
	for _, gNode := range classNode.GetterList {
		getterTag := gNode.GetterName.GetLiteral()
		ref.GetterList[getterTag] = BuildClosureRefFromNode(getterTag, []*syntax.ParamItem{}, gNode.ExecBlock)
	}

	// add methods
	for _, mNode := range classNode.MethodList {
		mTag := mNode.FuncName.GetLiteral()
		ref.MethodList[mTag] = BuildClosureRefFromNode(mTag, mNode.ParamList, mNode.ExecBlock)
	}

	return ref
}

// Construct - yield new instance of this class
func (cr *ClassRef) Construct(ctx *Context, scope *FuncScope, params []ZnValue) (ZnValue, *error.Error) {
	return cr.Constructor(ctx, scope, params)
}

//// helpers
func bindClassGetters(ref *ClassRef, getterMap map[string]funcExecutor) {
	for key, executor := range getterMap {
		ref.GetterList[key] = NewClosureRef(key, executor)
	}
}
