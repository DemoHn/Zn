package runtime

type ModuleGraph struct {
	depGraph     map[string][]string
	requireCache map[string]*Module
	// There's ONLY one anonymousModule (aka. main module in playground & REPL executor) allowed in one context
	anonymousModule *Module
}

func NewModuleGraph() *ModuleGraph {
	return &ModuleGraph{
		depGraph:     map[string][]string{},
		requireCache: map[string]*Module{},
	}
}

func (g *ModuleGraph) AddRequireCache(module *Module) {
	// requireCaches are exclusive - that is, if one Module object has occupied one requireCache name,other object could not modify/replace the cache.
	if module.IsAnonymous() {
		if g.anonymousModule == nil {
			g.anonymousModule = module
		}
	} else {
		if _, ok := g.requireCache[module.name]; !ok {

			g.requireCache[module.name] = module
		}
	}
}

func (g *ModuleGraph) FindRequireCache(name string) *Module {
	if m, ok := g.requireCache[name]; ok {
		return m
	}
	return nil
}

func (g *ModuleGraph) AddModuleDepRecord(source string, depModule string) {
	depList, ok := g.depGraph[source]
	if !ok {
		g.depGraph[source] = []string{}
	}
	g.depGraph[source] = append(depList, depModule)
}
