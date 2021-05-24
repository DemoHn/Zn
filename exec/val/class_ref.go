package val

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/exec/ctx"
	"github.com/DemoHn/Zn/syntax"
)

// BuildClassFromNode -
func BuildClassFromNode(name string, classNode *syntax.ClassDeclareStmt) ctx.ClassRef {
	ref := ctx.ClassRef{
		Name:         name,
		Constructor:  nil,
		PropList:     []string{},
		CompPropList: map[string]ctx.ClosureRef{},
		MethodList:   map[string]ctx.ClosureRef{},
	}

	// define default constrcutor
	var constructor = func(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
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

// NewClassRef - create new empty ctx.ClassRef
func NewClassRef(name string) ctx.ClassRef {
	return ctx.ClassRef{
		Name:         name,
		Constructor:  nil,
		PropList:     []string{},
		CompPropList: map[string]ClosureRef{},
		MethodList:   map[string]ClosureRef{},
	}
}

// Construct - yield new instance of this class
func (cr *ctx.ClassRef) Construct(c *ctx.Context, params []ctx.Value) (ctx.Value, *error.Error) {
	return cr.Constructor(ctx, params)
}
