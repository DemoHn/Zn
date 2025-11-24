package value

import (
	"fmt"

	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type Function struct {
	name         string
	logicHandler r.FuncExecutor
}

func NewFunction(executor r.FuncExecutor) *Function {
	return &Function{
		name:         "",
		logicHandler: executor,
	}
}

func (fn *Function) String() string {
	if fn.name == "" {
		return "‹某方法›"
	}
	return fmt.Sprintf("‹方法·%s›", fn.name)
}

func (fn *Function) SetName(name string) *Function {
	fn.name = name
	return fn
}

// Exec - execute the Function Object - accepts input params, execute from closure executor and
// yields final result
func (fn *Function) Exec(thisValue r.Element, params []r.Element) (r.Element, error) {
	fnLogicHandler := fn.logicHandler
	result, err := fnLogicHandler(thisValue, params)
	// convert error to exception
	if err != nil {
		switch err.(type) {
		case *zerr.SyntaxError:
		case *zerr.SemanticError:
		case *zerr.IOError:
		case *zerr.Signal:
			// return the original error AS IS
			return nil, err
		case *Exception:
			return nil, err
		case *zerr.RuntimeError:
			return nil, NewException(err.Error())
		default:
			// for other types of error (native errors), wrap the error as an Exception
			return nil, NewException(err.Error())
		}
	}
	return result, nil
}

// impl Value interface
// GetProperty -
func (fn *Function) GetProperty(name string) (r.Element, error) {
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (fn *Function) SetProperty(name string, value r.Element) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (fn *Function) ExecMethod(name string, values []r.Element) (r.Element, error) {
	return nil, zerr.MethodNotFound(name)
}

func (fn *Function) Exportable() bool {
	return true
}
