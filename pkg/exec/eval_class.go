package exec

import (
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/value"
)

// compileClass -
func compileClass(upperCtx *r.Context, classID *r.IDName, classNode *syntax.ClassDeclareStmt) (*value.ClassModel, error) {
	className := classID.GetLiteral()
	ref := value.NewClassModel(className, upperCtx.GetCurrentModule())

	// init prop list and its default value
	for _, propPair := range classNode.PropertyList {
		propID := propPair.PropertyID.GetLiteral()
		element, err := evalExpression(upperCtx, propPair.InitValue)
		if err != nil {
			return nil, err
		}

		ref.DefineProperty(propID, element)
	}

	// add getters
	for _, gNode := range classNode.GetterList {
		getterTag := gNode.Name.GetLiteral()
		ref.DefineCompProperty(getterTag, compileFunction(upperCtx, gNode.ParamList, gNode.ExecBlock, gNode.CatchBlocks))
	}

	// add methods
	for _, mNode := range classNode.MethodList {
		mTag := mNode.Name.GetLiteral()
		ref.DefineMethod(mTag, compileFunction(upperCtx, mNode.ParamList, mNode.ExecBlock, mNode.CatchBlocks))
	}

	return ref, nil
}
