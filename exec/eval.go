package exec

import (
	"github.com/DemoHn/Zn/error"
	"github.com/DemoHn/Zn/syntax"
)

// eval.go evaluates program from generated AST tree with specific scopes
// common signature of eval functions:
//
// evalXXXXStmt(ctx *Context, scope Scope, node Node) *error.Error
//
// or
//
// evalXXXXExpr(ctx *Context, scope Scope, node Node) (ZnValue, *error.Error)
//
// NOTICE:
// `evalXXXXStmt` will change the value of its corresponding scope; However, `evalXXXXExpr` will export
// a ZnValue object and mostly won't change scopes (but search a variable from scope is frequently used)

// Scope - tmp Scope solution TODO: will move in the future!
type Scope interface {
	// GetValue - get variable name from current scope
	GetValue(ctx *Context, name string) (ZnValue, *error.Error)
	// GetValue - set variable value from current scope
	SetValue(ctx *Context, name string, value ZnValue) *error.Error
}

// TODO: find a better way to handle this
func duplicateValue(in ZnValue) ZnValue {
	return in
}

func evalVarDeclareStmt(ctx *Context, scope Scope, node *syntax.VarDeclareStmt) *error.Error {
	for _, vpair := range node.AssignPair {
		obj, err := evalExpression(ctx, scope, vpair.AssignExpr)
		if err != nil {
			return err
		}
		for _, v := range vpair.Variables {
			vtag := v.GetLiteral()
			finalObj := duplicateValue(obj)
			if scope.SetValue(ctx, vtag, finalObj); err != nil {
				return err
			}
		}
	}
	return nil
}

func evalWhileLoopStmt(ctx *Context, scope Scope, node *syntax.WhileLoopStmt) *error.Error {
	// TODO
}

func evalExpression(ctx *Context, scope Scope, node syntax.Expression) (ZnValue, *error.Error) {
	return nil, nil
}

// eval var assign
func evalVarAssignExpr(ctx *Context, scope Scope, expr *syntax.VarAssignExpr) (ZnValue, *error.Error) {
	// Right Side
	val, err := evalExpression(ctx, scope, expr.AssignExpr)
	if err != nil {
		return nil, err
	}

	// Left Side
	switch v := expr.TargetVar.(type) {
	case *syntax.ID:
		// set ID
		vtag := v.GetLiteral()
		err2 := ctx.SetData(vtag, val)
		return val, err2
	case *syntax.ArrayListIndexExpr:
		iv, err := getArrayListIV(ctx, scope, v)
		if err != nil {
			return nil, err
		}
		return iv.Reduce(val, true)
	default:
		return nil, error.InvalidCaseType()
	}
}

func getArrayListIV(ctx *Context, scope Scope, expr *syntax.ArrayListIndexExpr) (ZnIV, *error.Error) {
	// val # index  --> 【1，２，３】#2
	val, err := evalExpression(ctx, scope, expr.Root)
	if err != nil {
		return nil, err
	}
	idx, err := evalExpression(ctx, scope, expr.Index)
	if err != nil {
		return nil, err
	}
	switch v := val.(type) {
	case *ZnArray:
		vr, ok := idx.(*ZnDecimal)
		if !ok {
			return nil, error.InvalidExprType("integer")
		}
		return &ZnArrayIV{v, vr}, nil
	case *ZnHashMap:
		var s *ZnString
		switch x := idx.(type) {
		case *ZnDecimal:
			// transform decimal value to string
			// x.exp < 0 express that its a decimal value with point mark, not an integer
			if x.exp < 0 {
				return nil, error.InvalidExprType("integer", "string")
			}
			s = NewZnString(x.String())
		case *ZnString:
			s = x
		default:
			return nil, error.InvalidExprType("integer", "string")
		}
		return &ZnHashMapIV{v, s}, nil
	default:
		return nil, error.InvalidExprType("array", "hashmap")
	}
}
