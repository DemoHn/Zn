package exec

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/syntax/zh"
)

type FileExecutor struct {
	rootDir  string
	mainFile string
	context  *r.Context
}

func NewFileExecutor(file string) *FileExecutor {
	// get dir & fileName -
	// e.g. when exec "/home/user/xxxx/module/a.zn":
	//  - dir=/home/user/xxxx/module
	//  - file=a.zn
	rootDir := filepath.Dir(file)
	_, fileName := filepath.Split(file)

	// context is initialized after InitRootModule() executed
	return &FileExecutor{
		rootDir:  rootDir,
		mainFile: fileName,
		context:  nil,
	}
}

func (fl *FileExecutor) RunCode(varInput map[string]r.Element) (r.Element, error) {
	// #1. init rootModule & context first
	if err := fl.initRootModule(); err != nil {
		return nil, err
	}

	// #2. parse program
	module := fl.context.GetCurrentModule()
	fl.context.SetVarInputs(varInput)
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

// initRootModule - and setup the context where rootModule = $this one
func (fl *FileExecutor) initRootModule() error {
	// #1. read source code from file
	in, err := io.NewFileStream(path.Join(fl.rootDir, fl.mainFile))
	if err != nil {
		return err
	}

	source, err := in.ReadAll()
	if err != nil {
		return err
	}

	// #2. create module & init context
	lexer := syntax.NewLexer(source)
	fl.context = r.NewContext(globalValues, r.NewMainModule(lexer))
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
