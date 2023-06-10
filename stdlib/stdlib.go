package stdlib

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

type funcExecutor = func(*r.Context, []r.Value) (r.Value, error)

var stdlibMap = map[string]*r.ModuleOLD{}

func RegisterModule(name string, module *r.ModuleOLD) {
	stdlibMap[name] = module
}

func FindModule(name string) (*r.ModuleOLD, error) {
	if module, ok := stdlibMap[name]; ok {
		return module, nil
	}
	return nil, zerr.ModuleNotFound(name)
}

// RegisterFunction - add function into module
func RegisterFunctionForModule(m *r.ModuleOLD, name string, fn funcExecutor) {
	m.RegisterValue(name, value.NewFunction(name, fn))
}

// RegisterClass - add class info module
func RegisterClassForModule(m *r.ModuleOLD, name string, ref *value.ClassRef) {
	m.RegisterValue(name, ref)
}
