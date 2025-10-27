package runtime

import (
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/syntax"
)

type ModuleGraph struct {
	modules       []*Module
	graph         [][2]int // array of [startModuleID, importModuleID]
	moduleNameMap map[string]int
}

type Module struct {
	fullName string
	// program stores sourceLines & AST - usually for error displaying
	program *syntax.Program
	// exportValues - all classes and functions are exported for external
	// imports - so here we insert all exportable values to this map after first scan
	// note: all export values are constants.
	exportValues ElementMap
}

type LibNameInfo struct {
	// libName original string
	OriginalName string
	// parsed libType
	LibType uint8
	// separate full libstring to subPath e.g.: "A-B-C" -> []string{"A", "B", "C"}
	LibPath []string
}

const (
	LIB_TYPE_STD    = 1
	LIB_TYPE_VENDOR = 2
	LIB_TYPE_CUSTOM = 3
)

type ModuleCodeFinder func(isMain bool, info LibNameInfo) ([]rune, error)

func (m *Module) GetName() string {
	return m.fullName
}

func (m *Module) AddExportValue(name string, value Element) error {
	if _, ok := m.exportValues[name]; ok {
		return zerr.NameRedeclared(name)
	}
	m.exportValues[name] = value
	return nil
}

func (m *Module) GetExportValue(name string) (Element, error) {
	if elem, ok := m.exportValues[name]; ok {
		return elem, nil
	} else {
		return nil, zerr.NameNotDefined(name)
	}
}

func (m *Module) GetAllExportValues() ElementMap {
	return m.exportValues
}

func NewSTDModule(name string) *Module {
	return &Module{
		fullName:     name,
		program:      nil, // STD module uses internal code
		exportValues: map[string]Element{},
	}
}

func NewModuleGraph() *ModuleGraph {
	return &ModuleGraph{
		modules:       []*Module{},
		graph:         [][2]int{},
		moduleNameMap: map[string]int{},
	}
}

// AddModule - create empty module information
// srcModuleID: moduleID of current module
// name: added module name
// program: parsed program
func (g *ModuleGraph) AddModule(srcModuleID int, name string, program *syntax.Program) int {
	g.modules = append(g.modules, &Module{
		fullName:     name,
		program:      program,
		exportValues: map[string]Element{},
	})

	extModuleID := len(g.modules) - 1
	g.graph = append(g.graph, [2]int{srcModuleID, extModuleID})
	g.moduleNameMap[name] = extModuleID
	return extModuleID
}

func (g *ModuleGraph) AddDependency(srcModuleID int, name string, depModuleID int) {
	g.graph = append(g.graph, [2]int{srcModuleID, depModuleID})
	g.moduleNameMap[name] = depModuleID
}

func (g *ModuleGraph) GetModuleByID(moduleID int) *Module {
	if moduleID >= 0 && moduleID < len(g.modules) {
		return g.modules[moduleID]
	}
	return nil
}

func (g *ModuleGraph) GetIDFromName(name string) (int, bool) {
	moduleID, ok := g.moduleNameMap[name]
	return moduleID, ok
}

func (g *ModuleGraph) CheckCircularDepedency(srcModuleID int, depModuleID int) bool {
	return g.checkCircularDepedencyDFS()
}

func (g *ModuleGraph) checkCircularDepedencyDFS() bool {
	// Build adjacency list
	adj := make(map[int][]int)
	for _, e := range g.graph {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		// Make sure v also exists in map (even if it has no outgoing edges)
		if _, ok := adj[v]; !ok {
			adj[v] = []int{}
		}
	}

	// 0 = unvisited, 1 = visiting, 2 = done
	color := make(map[int]int)

	var dfs func(int) bool
	dfs = func(u int) bool {
		color[u] = 1
		for _, v := range adj[u] {
			if color[v] == 1 {
				// found a back edge
				return true
			}
			if color[v] == 0 && dfs(v) {
				return true
			}
		}
		color[u] = 2
		return false
	}

	for node := range adj {
		if color[node] == 0 {
			if dfs(node) {
				return true
			}
		}
	}
	return false
}

// parseLibName - parse libName into LibNameInfo
// libName - string to be parsed, e.g. "A-B-C"
// return LibNameInfo object with originalName and libType set to 0 (LIB_TYPE_STD) and libPath set to empty slice
// LibNameInfo fields:
// originalName - original string passed to parseLibName
// libType - parsed libType (LIB_TYPE_STD, LIB_TYPE_VENDOR, LIB_TYPE_CUSTOM)
// libPath - separate full libstring to subPath, e.g. "A-B-C" -> []string{"A", "B", "C"}
func ParseLibName(libName string) LibNameInfo {
	if strings.HasPrefix(libName, "@") {
		return LibNameInfo{
			OriginalName: libName,
			LibType:      LIB_TYPE_STD,
			LibPath:      strings.Split(libName[1:], "-"),
		}
	}

	return LibNameInfo{
		OriginalName: libName,
		LibType:      LIB_TYPE_CUSTOM,
		LibPath:      strings.Split(libName, "-"),
	}
}
