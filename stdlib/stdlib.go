package stdlib

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

type Library struct {
	Name         string
	ExportValues map[string]r.Element
}

var STDLIB_MAP = map[string]*Library{}

func NewLibrary(name string) *Library {
	return &Library{
		Name:         name,
		ExportValues: map[string]r.Element{},
	}
}

func (l *Library) AddExportValue(name string, value r.Element) {
	l.ExportValues[name] = value
}

func RegisterLibrary(library *Library) {
	STDLIB_MAP[library.Name] = library
}

func FindLibrary(name string) (*Library, error) {
	if library, ok := STDLIB_MAP[name]; ok {
		return library, nil
	}
	return nil, zerr.LibraryNotFound(name)
}

// RegisterFunction - add function into module
func RegisterFunctionForLibrary(lib *Library, name string, fn r.FuncExecutor) {
	lib.AddExportValue(name, value.NewFunction(fn))
}

// RegisterClass - add class info module
func RegisterClassForLibrary(lib *Library, name string, ref *value.ClassModel) {
	lib.AddExportValue(name, ref)
}
