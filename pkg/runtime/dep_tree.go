package runtime

import (
	"errors"
	zerr "github.com/DemoHn/Zn/pkg/error"
	"os"
	"path/filepath"
	"strings"
)

// DependencyTree - manage all dependencies
type DependencyTree struct {
	rootDir string
	moduleCacheMap map[string]*Module
}

func NewDependencyTree() *DependencyTree {
	return &DependencyTree{
		rootDir:        "",
		moduleCacheMap: map[string]*Module{},
	}
}

func (dp *DependencyTree) SetRootDir(rootDir string) {
	dp.rootDir = rootDir
}

func (dp *DependencyTree) GetModulePath(name string) (string, error) {
	rootDir := dp.rootDir
	dirs := strings.Split(name, "-")
	// add .zn for last item
	dirs[len(dirs)-1] = dirs[len(dirs)-1] + ".zn"

	path := filepath.Join(rootDir, filepath.Join(dirs...))
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return "", zerr.ModuleNotFound(name)
	}

	return path, nil
}

func (dp *DependencyTree) BuildModuleCache(module *Module) {
	dp.moduleCacheMap[module.name] = module
}

func (dp *DependencyTree) GetModuleCache(name string) *Module {
	if m, ok := dp.moduleCacheMap[name]; ok {
		return m
	}
	return nil
}
