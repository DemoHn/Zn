package exec

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/syntax/zh"
)

type FileExecutor struct {
	rootDir    string
	rootModule string
	context    *r.Context
}

func NewFileExecutor(rootDir string, rootModule string) *FileExecutor {
	// context is initialized after InitRootModule() executed
	return &FileExecutor{
		rootDir:    rootDir,
		rootModule: rootModule,
		context:    nil,
	}
}

func (fl *FileExecutor) RunModule() (r.Value, error) {
	// #1. init rootModule & context first
	if err := fl.initRootModule(); err != nil {
		return nil, err
	}

	// #2. parse program
	module := fl.context.GetCurrentModule()
	lexer := module.GetLexer()
	p := syntax.NewParser(lexer, zh.NewParserZH())

	program, err := p.Parse()
	if err != nil {
		return nil, WrapSyntaxError(lexer, module, err)
	}

	// #3. eval program
	if err := evalProgram(fl.context, program); err != nil {
		return nil, WrapRuntimeError(fl.context, err)
	}

	returnValue := fl.context.GetCurrentScope().GetReturnValue()
	return returnValue, nil
}

func (fl *FileExecutor) HasPrinted() bool {
	if fl.context != nil {
		return fl.context.GetHasPrinted()
	}
	return false
}

func (fl *FileExecutor) initRootModule() error {
	// #1. find filepath of current module
	moduleFile, err := fl.getModulePath(fl.rootModule)
	if err != nil {
		return err
	}

	// #2. read source code from file
	in, err := io.NewFileStream(moduleFile)
	if err != nil {
		return err
	}

	source, err := in.ReadAll()
	if err != nil {
		return err
	}

	// #3. create module & init context
	lexer := syntax.NewLexer(source)
	fl.context = r.NewContext(GlobalValues, r.NewModule(fl.rootModule, lexer))
	// set source code finder
	fl.context.SetModuleCodeFinder(fl.buildModuleCodeFinder())

	return nil
}

func (fl *FileExecutor) getModulePath(name string) (string, error) {
	rootDir := fl.rootDir
	dirs := strings.Split(name, "-")
	// add .zn for last item
	dirs[len(dirs)-1] = dirs[len(dirs)-1] + ".zn"

	path := filepath.Join(rootDir, filepath.Join(dirs...))
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return "", zerr.ModuleNotFound(name)
	}

	return path, nil
}

//// define moduleCodeFinder for finding other modules
func (fl *FileExecutor) buildModuleCodeFinder() r.ModuleCodeFinder {
	return func(s string) ([]rune, error) {
		// #1. find filepath of current module
		moduleFile, err := fl.getModulePath(s)
		if err != nil {
			return []rune{}, err
		}

		// #2. read source code from file
		in, err := io.NewFileStream(moduleFile)
		if err != nil {
			return []rune{}, err
		}

		return in.ReadAll()
	}
}
