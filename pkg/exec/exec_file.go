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
	// #1. parse rootModule
	program, err := fl.parseMainModule()
	if err != nil {
		return nil, err
	}

	// #2. create main module context
	module := r.NewMainModule(program.Lines)

	runContext := r.NewContext(globalValues, module)
	runContext.SetModuleCodeFinder(fl.buildModuleCodeFinder())
	runContext.SetVarInputs(varInput)

	fl.context = runContext

	// #3. eval program
	if err := evalProgram(fl.context, program); err != nil {
		return nil, WrapRuntimeError(fl.context, err)
	}

	returnValue := fl.context.GetCurrentScope().GetReturnValue()
	return returnValue, nil
}

// initRootModule - and setup the context where rootModule = $this one
func (fl *FileExecutor) parseMainModule() (*syntax.Program, error) {
	// #1. read source code from file
	in, err := io.NewFileStream(path.Join(fl.rootDir, fl.mainFile))
	if err != nil {
		return nil, err
	}

	source, err := in.ReadAll()
	if err != nil {
		return nil, err
	}

	// #2. parse source
	parser := syntax.NewParser(source, zh.NewParserZH())
	program, err := parser.Parse()
	if err != nil {
		return nil, WrapSyntaxError(parser, "", err)
	}

	return program, err
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
	return func(isMainModule bool, s string) ([]rune, error) {
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
