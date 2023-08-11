package runtime

type ModuleGraph struct {
	// external modules (usually written in zinc)
	externalDepGraph map[string][]string
	// internal modules (usually standard libraray, plugin)
	internalDepGraph map[string][]string

	requireCache map[string]*Module
	// There's ONLY one anonymousModule (aka. main module in playground & REPL executor) allowed in one context
	anonymousModule *Module
}

func NewModuleGraph() *ModuleGraph {
	return &ModuleGraph{
		externalDepGraph: map[string][]string{},
		internalDepGraph: map[string][]string{},
		requireCache:     map[string]*Module{},
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

func (g *ModuleGraph) CheckCircularDepedency(source string, depModule string, internal bool) bool {
	// add dep record first
	g.addDepRecord(source, depModule, internal)

	// for internal dep modules, no need to check circular dependency
	if internal {
		return false
	}
	return g.checkCircularDepedencyBFS()
}

func (g *ModuleGraph) addDepRecord(srcModule string, depModule string, internal bool) {
	if internal {
		depList, ok := g.internalDepGraph[srcModule]
		// add `source` key to depGraph
		if !ok {
			g.internalDepGraph[srcModule] = []string{}
		}
		// add `depModule` key to depGraph
		if _, ok := g.internalDepGraph[depModule]; !ok {
			g.internalDepGraph[depModule] = []string{}
		}
		g.internalDepGraph[srcModule] = append(depList, depModule)
	} else {
		depList, ok := g.externalDepGraph[srcModule]
		// add `source` key to depGraph
		if !ok {
			g.externalDepGraph[srcModule] = []string{}
		}
		// add `depModule` key to depGraph
		if _, ok := g.externalDepGraph[depModule]; !ok {
			g.externalDepGraph[depModule] = []string{}
		}
		g.externalDepGraph[srcModule] = append(depList, depModule)
	}
}

func (g *ModuleGraph) checkCircularDepedencyBFS() bool {
	// detectCircularDependencyBFS
	// same as "find the loop in a directed graph"
	// ref: Kahn's algorithm in https://en.wikipedia.org/wiki/Topological_sorting

	// init inDegreeMap - stat how many incoming nodes in one module
	inDegreeMap := map[string]int{}
	for src, sDeps := range g.externalDepGraph {
		if _, ok := inDegreeMap[src]; !ok {
			inDegreeMap[src] = 0
		}

		for _, dep := range sDeps {
			if _, ok := inDegreeMap[dep]; !ok {
				inDegreeMap[dep] = 0
			}
			inDegreeMap[dep] += 1
		}
	}

	// create an queue and enqueue all vertices with indegree 0
	queue := []string{}
	for module, n := range inDegreeMap {
		if n == 0 {
			queue = append(queue, module)
		}
	}

	// Initialize count of visited vertices
	cnt := 1
	nodesNum := len(g.externalDepGraph)

	// One by one dequeue vertices from queue and enqueue
	// adjacents if indegree of adjacent becomes 0
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Iterate through all its neighbouring nodes of
		// dequeued node u and decrease their in-degree by 1
		for _, node := range g.externalDepGraph[current] {
			inDegreeMap[node] -= 1

			if inDegreeMap[node] == 0 {
				queue = append(queue, node)
				cnt += 1
			}
		}
	}

	return cnt != nodesNum
}
