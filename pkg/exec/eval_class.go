package exec

import (
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/value"
)

// compileClass -
func compileClass(vm *r.VM, classID *r.IDName, classNode *syntax.ClassDeclareStmt) (*value.ClassModel, error) {
	className := classID.GetLiteral()
	ref := value.NewClassModel(className)

	// init prop list and its default value
	for _, propPair := range classNode.PropertyList {
		propID := propPair.PropertyID.GetLiteral()
		element, err := evalExpression(vm, propPair.InitValue)
		if err != nil {
			return nil, err
		}

		ref.DefineProperty(propID, element)
	}

	// add getters
	for _, gNode := range classNode.GetterList {
		getterTag := gNode.Name.GetLiteral()
		ref.DefineCompProperty(getterTag, compileFunction(vm, gNode))
	}

	// add methods
	for _, mNode := range classNode.MethodList {
		mTag := mNode.Name.GetLiteral()
		ref.DefineMethod(mTag, compileFunction(vm, mNode))
	}

	return ref, nil
}
