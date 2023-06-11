package runtime

type ModuleGraph struct {
	depGraph  map[string][]string
	moduleMap map[string]*Module
}

func NewModuleGraph() *ModuleGraph {
	return &ModuleGraph{
		depGraph:  map[string][]string{},
		moduleMap: map[string]*Module{},
	}
}

func (g *ModuleGraph) AddModule(module *Module) {
	g.moduleMap[module.name] = module
}

func (g *ModuleGraph) FindModule(name string) *Module {
	if m, ok := g.moduleMap[name]; ok {
		return m
	}
	return nil
}
