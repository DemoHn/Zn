package stdlib

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

type funcExecutor = func(*r.Context, []r.Value) (r.Value, error)

var stdlibMap = map[string]*r.Module{}

func RegisterModule(name string, module *r.Module) {
	stdlibMap[name] = module
}

func FindModule(name string) (*r.Module, error) {
	if module, ok := stdlibMap[name]; ok {
		return module, nil
	}
	return nil, zerr.ModuleNotFound(name)
}

// RegisterFunction - add function into module
func RegisterFunctionForModule(m *r.Module, name string, fn funcExecutor) {
	m.AddExportValue(name, value.NewFunction(fn))
}

// RegisterClass - add class info module
func RegisterClassForModule(m *r.Module, name string, ref *value.ClassRef) {
	m.AddExportValue(name, ref)
}
