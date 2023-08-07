package runtime

type ModuleGraph struct {
	depGraph     map[string][]string
	requireCache map[string]*Module
	// There's ONLY one anonymousModule (aka. main module in playground & REPL executor) allowed in one context
	anonymousModule *Module
}

type depNode struct {
	Name         string
	Dependencies []*depNode
	Visited      bool
	InStack      bool
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

func (g *ModuleGraph) CheckCircularDepedency(source string, depModule string) bool {
	// add dep record first
	g.addDepRecord(source, depModule)

	// detectCircularDependencyBFS
	// same as "find the loop in a directed graph"
	// ref: Kahn's algorithm in https://en.wikipedia.org/wiki/Topological_sorting

	// init inDegreeMap - stat how many incoming nodes in one module
	inDegreeMap := map[string]int{}
	for src, sDeps := range g.depGraph {
		inDegreeMap[src] = 0
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

	// One by one dequeue vertices from queue and enqueue
	// adjacents if indegree of adjacent becomes 0
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Iterate through all its neighbouring nodes of
		// dequeued node u and decrease their in-degree by 1
		for _, node := range g.depGraph[current] {
			inDegreeMap[node] -= 1

			if inDegreeMap[node] == 0 {
				queue = append(queue, node)
				cnt += 1
			}
		}
	}

	return cnt != len(g.depGraph)
}

func (g *ModuleGraph) addDepRecord(source string, depModule string) {
	depList, ok := g.depGraph[source]
	// add `source` key to depGraph
	if !ok {
		g.depGraph[source] = []string{}
	}
	// add `depModule` key to depGraph
	if _, ok := g.depGraph[depModule]; !ok {
		g.depGraph[depModule] = []string{}
	}
	g.depGraph[source] = append(depList, depModule)
}
