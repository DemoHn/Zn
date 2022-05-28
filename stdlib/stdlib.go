package stdlib

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

var stdlibMap map[string]*r.Module

func RegisterModule(name string, module *r.Module) {
	stdlibMap[name] = module
}

func FindModule(name string) (*r.Module, error) {
	if module, ok := stdlibMap[name]; ok {
		return module, nil
	}
	return nil, zerr.ModuleNotFound(name)
}
