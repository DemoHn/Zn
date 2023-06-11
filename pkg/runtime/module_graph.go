package runtime

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
)

type ModuleGraph struct {
	rootDir   string
	depGraph  map[string][]string
	moduleMap map[string]*Module
}

func NewModuleGraph(rootDir string) *ModuleGraph {
	return &ModuleGraph{
		rootDir:   rootDir,
		depGraph:  map[string][]string{},
		moduleMap: map[string]*Module{},
	}
}

func (g *ModuleGraph) GetModulePath(name string) (string, error) {
	rootDir := g.rootDir
	dirs := strings.Split(name, "-")
	// add .zn for last item
	dirs[len(dirs)-1] = dirs[len(dirs)-1] + ".zn"

	path := filepath.Join(rootDir, filepath.Join(dirs...))
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return "", zerr.ModuleNotFound(name)
	}

	return path, nil
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
