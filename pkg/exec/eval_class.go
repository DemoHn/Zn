package exec

import (
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/value"
)

// compileClass -
func compileClass(upperCtx *r.Context, name string, classNode *syntax.ClassDeclareStmt) (*value.ClassModel, error) {
	ref := value.NewClassModel(name, upperCtx.GetCurrentModule())

	// set default constructor
	ref.Constructor = func(c *r.Context, params []r.Element) (r.Element, error) {
		// create new object when exec ONLY
		obj := value.NewObject(ref)

		return obj, nil
	}

	// init prop list and its default value
	for _, propPair := range classNode.PropertyList {
		propID := propPair.PropertyID.GetLiteral()
		element, err := evalExpression(upperCtx, propPair.InitValue)
		if err != nil {
			return nil, err
		}

		ref.PropList[propID] = element
	}

	// add getters
	for _, gNode := range classNode.GetterList {
		getterTag := gNode.GetterName.GetLiteral()
		ref.CompPropList[getterTag] = compileFunction(upperCtx, []*syntax.ParamItem{}, gNode.ExecBlock)
	}

	// add methods
	for _, mNode := range classNode.MethodList {
		mTag := mNode.FuncName.GetLiteral()
		ref.MethodList[mTag] = compileFunction(upperCtx, mNode.ParamList, mNode.ExecBlock)
	}

	return ref, nil
}

func bindClassConstructor(model *value.ClassModel, fn *value.Function) {
	model.Constructor = func(c *r.Context, params []r.Element) (r.Element, error) {
		obj := value.NewObject(model)

		// exec constructor logic (last value is useless)
		if _, err := fn.Exec(c, obj, params); err != nil {
			return nil, err
		}
		return obj, nil
	}
}
