package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
)

type ModuleGraph struct {
	depGraph  map[string][]string
	moduleMap map[string]*Module
	// There's ONLY one anonymousModule (aka. main module in playground & REPL executor) allowed in one context
	anonymousModule *Module
}

func NewModuleGraph() *ModuleGraph {
	return &ModuleGraph{
		depGraph:  map[string][]string{},
		moduleMap: map[string]*Module{},
	}
}

func (g *ModuleGraph) AddModule(module *Module) error {
	if module.IsAnonymous() {
		if g.anonymousModule != nil {
			return zerr.MoreAnonymousModule()
		}
		g.anonymousModule = module
		return nil
	}
	if _, ok := g.moduleMap[module.name]; ok {
		return zerr.ModuleHasDefined(module.name)
	}
	g.moduleMap[module.name] = module
	return nil
}

func (g *ModuleGraph) FindModule(name string) *Module {
	if m, ok := g.moduleMap[name]; ok {
		return m
	}
	return nil
}
