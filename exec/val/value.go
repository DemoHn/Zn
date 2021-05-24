package val

import (
	"regexp"

	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/syntax"
)

type funcExecutor func(*Context, []Value) (Value, *error.Error)

// Value is the base unit to present a value (aka. variable) - including number, string, array, function, object...
// All kinds of values in Zn language SHOULD implement this interface.
//
// Basically there're 3 methods:
//
// 1. GetProperty - fetch the value from property list of a specific name
// 2. SetProperty - set the value of some property
// 3. ExecMethod - execute one method from method list
type Value interface {
	GetProperty(*Context, string) (Value, *error.Error)
	SetProperty(*Context, string, Value) *error.Error
	ExecMethod(*Context, string, []Value) (Value, *error.Error)
}

// ClosureRef - aka. Closure Exection Reference
// It's the structure of a closure which wraps execution logic.
// The executor could be either a bunch of code or some native code.
type ClosureRef struct {
	ParamHandler funcExecutor
	Executor     funcExecutor // closure execution logic
}

// BuildClosureFromNode - create a closure (with default param handler logic)
// from Zn code (*syntax.BlockStmt). It's the constructor of 如何XX or (anoymous function in the future)
func BuildClosureFromNode(paramTags []*syntax.ParamItem, stmtBlock *syntax.BlockStmt) ClosureRef {
	var executor = func(ctx *Context, params []Value) (Value, *error.Error) {
		// iterate block round I - function hoisting
		// NOTE: function hoisting means bind function definitions at the beginning
		// of execution so that even if "function execution" statement is before
		// "function definition" statement.
		for _, stmtI := range stmtBlock.Children {
			if v, ok := stmtI.(*syntax.FunctionDeclareStmt); ok {
				fn := BuildFunctionFromNode(v)
				if err := bindValue(ctx, v.FuncName.GetLiteral(), fn); err != nil {
					return nil, err
				}
			}
		}
		// iterate block round II - execution of rest code blocks
		for _, stmtII := range stmtBlock.Children {
			if _, ok := stmtII.(*syntax.FunctionDeclareStmt); !ok {
				if err := evalStatement(ctx, stmtII); err != nil {
					// if recv breaks
					if err.GetCode() == error.ReturnBreakSignal {
						if extra, ok := err.GetExtra().(Value); ok {
							return extra, nil
						}
					}
					return nil, err
				}
			}
		}
		return ctx.scope.returnValue, nil
	}

	var paramHandler = func(ctx *Context, params []Value) (Value, *error.Error) {
		// check param length
		if len(params) != len(paramTags) {
			return nil, error.MismatchParamLengthError(len(paramTags), len(params))
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
			if err := bindValue(ctx, paramName, paramVal); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	return ClosureRef{
		ParamHandler: paramHandler,
		Executor:     executor,
	}
}

// NewClosure - wraps a closure from native code (Golang code)
func NewClosure(paramHandler funcExecutor, executor funcExecutor) ClosureRef {
	return ClosureRef{
		ParamHandler: paramHandler,
		Executor:     executor,
	}
}

// Exec - execute a closure - accepts input params, execute from closure exeuctor and
// yields final result
func (cs *ClosureRef) Exec(ctx *Context, params []Value) (Value, *error.Error) {
	if cs.ParamHandler != nil {
		if _, err := cs.ParamHandler(ctx, params); err != nil {
			return nil, err
		}
	}
	if cs.Executor == nil {
		return nil, error.NewErrorSLOT("执行逻辑不能为空")
	}
	// do execution
	return cs.Executor(ctx, params)
}

// ClassRef - aka. Class Definition Reference
// It defines the structure of a class, including compPropList, methodList and propList.
// All instances created from this class MUST inherits from those configurations.
type ClassRef struct {
	// Name - class name
	Name string
	// Constructor defines default logic (mostly for initialization) when a new instance
	// is created by "x 成为 C：P，Q，R"
	Constructor funcExecutor
	// PropList defines all property name of a class, each item COULD NOT BE neither append nor removed
	PropList []string
	// CompPropList - CompProp stands for "Computed Property", which means the value is get or set
	// from a pre-defined function. Computed property offers more extensions for manipulations
	// of properties.
	CompPropList map[string]ClosureRef
	// MethodList - stores all available methods defintion of class
	MethodList map[string]ClosureRef
}

// BuildClassFromNode -
func BuildClassFromNode(name string, classNode *syntax.ClassDeclareStmt) ClassRef {
	ref := ClassRef{
		Name:         name,
		Constructor:  nil,
		PropList:     []string{},
		CompPropList: map[string]ClosureRef{},
		MethodList:   map[string]ClosureRef{},
	}

	// define default constrcutor
	var constructor = func(ctx *Context, params []Value) (Value, *error.Error) {
		obj := NewObject(ref)
		// init prop list
		for _, propPair := range classNode.PropertyList {
			propID := propPair.PropertyID.GetLiteral()
			expr, err := evalExpression(ctx, propPair.InitValue)
			if err != nil {
				return nil, err
			}
			obj.propList[propID] = expr
			ref.PropList = append(ref.PropList, propID)
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
			obj.propList[paramName] = objParamVal
		}

		return obj, nil
	}
	// set constructor
	ref.Constructor = constructor

	// add getters
	for _, gNode := range classNode.GetterList {
		getterTag := gNode.GetterName.GetLiteral()
		ref.CompPropList[getterTag] = BuildClosureFromNode([]*syntax.ParamItem{}, gNode.ExecBlock)
	}

	// add methods
	for _, mNode := range classNode.MethodList {
		mTag := mNode.FuncName.GetLiteral()
		ref.MethodList[mTag] = BuildClosureFromNode(mNode.ParamList, mNode.ExecBlock)
	}

	return ref
}

// NewClassRef - create new empty ClassRef
func NewClassRef(name string) ClassRef {
	return ClassRef{
		Name:         name,
		Constructor:  nil,
		PropList:     []string{},
		CompPropList: map[string]ClosureRef{},
		MethodList:   map[string]ClosureRef{},
	}
}

// Construct - yield new instance of this class
func (cr *ClassRef) Construct(ctx *Context, params []Value) (Value, *error.Error) {
	return cr.Constructor(ctx, params)
}

//// param validators

// validateExactParams is a helper function that asserts input params type where each parameter
// should exactly match the list of typeStr.
// valid typeStr are one of the followings:
//
// 1. decimal  --> *Decimal
// 2. string   --> *String
// 3. array    --> *Array
// 4. hashmap  --> *HashMap
// 5. bool     --> *Bool
// 6. object   --> *Object
// 7. function --> *Function
// 8. any      --> any value type
func validateExactParams(values []Value, typeStr ...string) *error.Error {
	if len(values) != len(typeStr) {
		return error.ExactParamsError(len(typeStr))
	}
	for idx, v := range values {
		if err := validateOneParam(v, typeStr[idx]); err != nil {
			return err
		}
	}
	return nil
}

// validateLeastParams is a helper function similar to validateExactParams(),
// however, it handles the situation that target params are variadic, like
// [int, int, string...], the length of target params varies from 0~N.
//
// To validate variadic params, we introduce wildcard (* or +) at the end of typeStr to
// mark variadic part, like "string+", "bool*"
//
// e.g.:
// ["decimal", "string+"] means the FIRST param is a decimal, and the FOLLOWING params
// are all strings (must have ONE string param)
//
// ["decimal", "bool", "string*"] means the FIRST param is a decimal, the SECOND param is a bool, and the FOLLOWING params
// are all strings (allow 0 string params)
func validateLeastParams(values []Value, typeStr ...string) *error.Error {
	for idx, t := range typeStr {
		// find if there's wildcard
		re := regexp.MustCompile(`(\w+)(\*|\+)?`)
		matches := re.FindStringSubmatch(t)
		// match: [_, name, wildcard]
		switch matches[2] {
		// matches 0 or more params
		case "*", "+":
			if matches[2] == "+" && idx > len(values) {
				return error.NewErrorSLOT("通配符需要至少一个参数")
			}
			for i := idx; i < len(values); i++ {
				if err := validateOneParam(values[i], matches[1]); err != nil {
					return err
				}
			}
			break // "break" HERE is to break the outer for-loop!
		default:
			if err := validateOneParam(values[idx], t); err != nil {
				return err
			}
		}
	}
	return nil
}

// validateAllParams doesn't limit the length of input values; instead, it requires all the parameters
// to have same value type denoted by `typeStr`
// e.g. validateAllParams([]Value{“1”, “2”, “3”}, "string")
func validateAllParams(values []Value, typeStr string) *error.Error {
	for _, v := range values {
		if err := validateOneParam(v, typeStr); err != nil {
			return err
		}
	}
	return nil
}

func validateOneParam(v Value, typeStr string) *error.Error {
	valid := true
	switch typeStr {
	case "decimal":
		if _, ok := v.(*Decimal); !ok {
			valid = false
		}
	case "string":
		if _, ok := v.(*String); !ok {
			valid = false
		}
	case "array":
		if _, ok := v.(*Array); !ok {
			valid = false
		}
	case "hashmap":
		if _, ok := v.(*HashMap); !ok {
			valid = false
		}
	case "bool":
		if _, ok := v.(*Bool); !ok {
			valid = false
		}
	case "object":
		if _, ok := v.(*Object); !ok {
			valid = false
		}
	case "function":
		if _, ok := v.(*Function); !ok {
			valid = false
		}
	case "any":
		valid = true
	}
	if valid == false {
		return error.InvalidParamType(typeStr)
	}
	return nil
}
