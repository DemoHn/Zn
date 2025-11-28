package exec

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/io"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/syntax/zh"
)

// Interpreter - MAIN CODE EXECUTION INSTANCE -
// ONE INTERPRETER -> ONE VM
type Interpreter struct {
	// version - compiler version
	version string

	// moduleCodeFinder - given a module name, the finder function aims to find it's corresponding source code for further execution - whatever from filesystem, DB, network, etc.
	// by default, the value is nil, that means the finder could not found any module code at all!
	moduleCodeFinder r.ModuleCodeFinder

	// externalLibs - loadable external libraries for 导入 statement
	// by default, ALL StandardLibs are included
	externalLibs []*r.Library

	// mainServer - [optional] the main server instance of this interpreter
	// we can build the additional server to serve incoming HTTP requests and
	// send response back.
	// Example 1: a HTTP web server, execute modules on the fly, then
	// make a response from the executed result.
	//
	// Example 2: a code playground - execute the script provided from incoming HTTP
	// requests, then make a response from the executed result
	mainServer ZnServer
}

type ZnServer interface {
	SetHandler(handler http.Handler)
	Start(connUrl string) error
}

func NewInterpreter(version string) *Interpreter {
	return &Interpreter{version: version}
}

// GetVersion - get current compiler's version
func (z *Interpreter) GetVersion() string {
	return z.version
}

// /// set functions /////
func (z *Interpreter) SetMainServer(server ZnServer, handler http.Handler) *Interpreter {
	z.mainServer = server
	z.mainServer.SetHandler(handler)
	return z
}

// TODO: better way to load external libs!
func (z *Interpreter) SetExternalLibs(libs []*r.Library) *Interpreter {
	z.externalLibs = libs
	return z
}

///// load functions //////

func (z *Interpreter) LoadScript(source []rune) *Interpreter {
	// set moduleCodeFinder
	z.moduleCodeFinder = func(isMain bool, info r.LibNameInfo) ([]rune, error) {
		// suppose the sourceCode is the mainModule ONLY
		if isMain {
			return source, nil
		} else {
			if info.LibType == r.LIB_TYPE_STD {
				return []rune{}, nil // return empty source for STD modules
			}
			// other cases, return error directly!
			return nil, zerr.NewErrorSLOT("在脚本模式下，不支持导入其他模块！")
		}
	}

	return z
}

func (z *Interpreter) LoadFile(file string) *Interpreter {
	// set moduleCodeFinder
	z.moduleCodeFinder = func(isMain bool, info r.LibNameInfo) ([]rune, error) {
		// get dir & fileName -
		// e.g. when exec "/home/user/xxxx/module/a.zn":
		//  - dir=/home/user/xxxx/module
		//  - file=a.zn
		rootDir := filepath.Dir(file)
		_, fileName := filepath.Split(file)

		var moduleFullPath string

		if isMain {
			// denote the loaded file as *mainModule*
			moduleFullPath = path.Join(rootDir, fileName)
		} else {
			switch info.LibType {
			case r.LIB_TYPE_STD:
				return []rune{}, nil // return empty source for STD modules
			case r.LIB_TYPE_VENDOR:
			case r.LIB_TYPE_CUSTOM:
				dirs := info.LibPath
				// add .zn for last item
				dirs[len(dirs)-1] = dirs[len(dirs)-1] + ".zn"

				moduleFullPath = filepath.Join(rootDir, filepath.Join(dirs...))
			}
		}
		if _, err := os.Stat(moduleFullPath); os.IsNotExist(err) {
			return nil, zerr.ModuleNotFound(info.OriginalName)
		}

		// read source code from the parsed modulePath
		in, err := io.NewFileStream(moduleFullPath)
		if err != nil {
			return nil, err
		}

		return in.ReadAll()
	}

	return z
}

func (z *Interpreter) Execute(varInputs r.ElementMap) (r.Element, error) {
	// #1. get the main source
	if z.moduleCodeFinder == nil {
		return nil, fmt.Errorf("code script/file not loaded")
	}

	finder := z.moduleCodeFinder
	// #2. load main module
	source, err := finder(true, r.LibNameInfo{
		OriginalName: "",
		LibType:      r.LIB_TYPE_CUSTOM,
		LibPath:      []string{},
	})
	if err != nil {
		return nil, err
	}

	// #3. compile the program -
	// currently from source code to AST, in the future, we will support compiling to bytecode
	parser := syntax.NewParser(source, zh.NewParserZH())
	program, err := parser.Compile()
	if err != nil {
		return nil, WrapSyntaxError(parser, MODULE_NAME_MAIN, err)
	}

	vm := r.InitVM(GlobalValues)
	vm.SetModuleCodeFinder(finder)
	vm.LoadExternalLibs(z.externalLibs)
	// #4. eval program
	rtnValue, err := EvalMainModule(vm, program, varInputs)
	if err != nil {
		return nil, WrapRuntimeError(vm, err)
	}

	// #5. get return value
	return rtnValue, nil
}

func (z *Interpreter) ExecuteVarInputText(exprStr string) (r.ElementMap, error) {
	return ExecVarInputText(exprStr)
}

// /// server functions /////
// Listen - if mainServer is set, start the server from defined `connection URL` (e.g. tcp://127.0.0.1:3862)
func (z *Interpreter) Listen(connUrl string) error {
	if z.mainServer == nil {
		return fmt.Errorf("main server not set")
	}
	return z.mainServer.Start(connUrl)
}
