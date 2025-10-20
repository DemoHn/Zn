package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/syntax"
)

type ModuleGraph struct {
	graph [][2]int // array of [startModuleID, importModuleID]
}

type Module struct {
	fullName string
	// program stores sourceLines & AST - usually for error displaying
	program *syntax.Program
	// exportValues - all classes and functions are exported for external
	// imports - so here we insert all exportable values to this map after first scan
	// note: all export values are constants.
	exportValues map[string]Element
}

func (m *Module) AddExportValue(name string, value Element) error {
	if _, ok := m.exportValues[name]; ok {
		return zerr.NameRedeclared(name)
	}
	m.exportValues[name] = value
	return nil
}

func NewModuleGraph() *ModuleGraph {
	return &ModuleGraph{
		graph: [][2]int{},
	}
}

func (g *ModuleGraph) AddDependency(srcModuleID int, depModuleID int) {
	g.graph = append(g.graph, [2]int{srcModuleID, depModuleID})
}
func (g *ModuleGraph) CheckCircularDepedency(srcModuleID int, depModuleID int) bool {
	return g.checkCircularDepedencyDFS()
}

func (g *ModuleGraph) checkCircularDepedencyDFS() bool {
	// detectCircularDependencyDFS
	// same as "find the loop in a directed graph"
	visited := make([]bool, len(g.graph))
	var dfs func(int) bool
	dfs = func(node int) bool {
		if visited[node] {
			return true
		}
		visited[node] = true
		for _, neighbor := range g.getNeighbors(node) {
			if dfs(neighbor) {
				return true
			}
		}
		return false
	}

	for i, _ := range visited {
		if dfs(i) {
			return true
		}
	}
	return false
}

func (g *ModuleGraph) getNeighbors(node int) []int {
	neighbors := []int{}
	for _, edge := range g.graph {
		if edge[0] == node {
			neighbors = append(neighbors, edge[1])
		}
	}
	return neighbors
}
