package zinc

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	zerr "github.com/DemoHn/Zn/pkg/error"
	"github.com/DemoHn/Zn/pkg/exec"
	"github.com/DemoHn/Zn/pkg/io"
	runtime "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/server"
	"github.com/DemoHn/Zn/pkg/syntax"
	"github.com/DemoHn/Zn/pkg/syntax/zh"
	"github.com/DemoHn/Zn/pkg/value"
)

type Element = runtime.Element
type Context = runtime.Context
type ModuleCodeFinder = runtime.ModuleCodeFinder

type ZnNumber = value.Number
type ZnString = value.String
type ZnBool = value.Bool
type ZnArray = value.Array
type ZnHashMap = value.HashMap
type ZnObject = value.Object
type ZnNull = value.Null

var NewZnNumber = value.NewNumber
var NewZnString = value.NewString
var NewZnBool = value.NewBool
var NewZnArray = value.NewArray
var NewZnHashMap = value.NewHashMap
var NewZnObject = value.NewObject
var NewZnNull = value.NewNull

const ZINC_VERSION = "rev07"

type ZnCompiler struct {
	// version - compiler version
	version string

	// moduleCodeFinder - given a module name, the finder function aims to find it's corresponding source code for further execution - whatever from filesystem, DB, network, etc.
	// by default, the value is nil, that means the finder could not found any module code at all!
	moduleCodeFinder ModuleCodeFinder
}

// NewCompiler - new ZnCompiler object
func NewCompiler() *ZnCompiler {
	return &ZnCompiler{
		version: ZINC_VERSION,
	}
}

// GetVersion - get current compiler's version
func (z *ZnCompiler) GetVersion() string {
	return z.version
}

func (z *ZnCompiler) LoadScript(source []rune) *ZnCompiler {
	// set moduleCodeFinder
	z.moduleCodeFinder = func(isMainModule bool, moduleName string) ([]rune, error) {
		// suppose the sourceCode is the mainModule ONLY
		if isMainModule {
			return source, nil
		} else {
			// thus there's no code for other modules!
			return []rune{}, nil
		}
	}

	return z
}

func (z *ZnCompiler) LoadFile(file string) *ZnCompiler {
	// set moduleCodeFinder
	z.moduleCodeFinder = func(isMainModule bool, moduleName string) ([]rune, error) {
		// get dir & fileName -
		// e.g. when exec "/home/user/xxxx/module/a.zn":
		//  - dir=/home/user/xxxx/module
		//  - file=a.zn
		rootDir := filepath.Dir(file)
		_, fileName := filepath.Split(file)

		var moduleFullPath string

		if isMainModule {
			// denote the loaded file as *mainModule*
			moduleFullPath = path.Join(rootDir, fileName)
		} else {
			dirs := strings.Split(moduleName, "-")
			// add .zn for last item
			dirs[len(dirs)-1] = dirs[len(dirs)-1] + ".zn"

			moduleFullPath = filepath.Join(rootDir, filepath.Join(dirs...))
			if _, err := os.Stat(moduleFullPath); errors.Is(err, os.ErrNotExist) {
				return nil, zerr.ModuleNotFound(moduleName)
			}
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

func (z *ZnCompiler) Execute(varInput map[string]Element) (Element, error) {
	// #1. get the main source
	if z.moduleCodeFinder == nil {
		return nil, fmt.Errorf("code script/file not loaded")
	}

	finder := z.moduleCodeFinder
	// #2. load main module
	source, err := finder(true, "")
	if err != nil {
		return nil, err
	}

	// #3. parse the program
	parser := syntax.NewParser(source, zh.NewParserZH())
	program, err := parser.Parse()
	if err != nil {
		return nil, exec.WrapSyntaxError(parser, "", err)
	}

	// create context
	runContext := runtime.NewContext(exec.GlobalValues, runtime.NewMainModule(program.Lines))
	runContext.SetModuleCodeFinder(z.moduleCodeFinder)
	runContext.SetVarInputs(varInput)

	// #4. eval program
	if err := exec.EvaluateProgram(runContext, program); err != nil {
		return nil, exec.WrapRuntimeError(runContext, err)
	}

	// #5. get return value
	return runContext.GetCurrentScope().GetReturnValue(), nil
}

type ZnServer struct {
	reqHandler     http.HandlerFunc
	pmServerConfig server.ZnPMServerConfig
}

func NewServer() *ZnServer {
	var defaultPMServerConfig = server.ZnPMServerConfig{
		InitProcs: 20,
		MaxProcs:  100,
		Timeout:   60,
	}

	return &ZnServer{
		reqHandler:     nil,
		pmServerConfig: defaultPMServerConfig,
	}
}

/// SetHandler() / SetPlaygroundHandler() / SetHTTPHandler()

func (s *ZnServer) SetHandler(handler http.HandlerFunc) *ZnServer {
	s.reqHandler = handler
	return s
}

func (s *ZnServer) SetPlaygroundHandler() *ZnServer {
	compiler := NewCompiler()
	s.reqHandler = func(w http.ResponseWriter, r *http.Request) {
		source, varInput, err := server.ReadRequestForPlayground(r)
		if err != nil {
			server.WriteResponseForPlayground(w, nil, err)
		} else {
			rtnValue, err := compiler.LoadScript(source).Execute(varInput)
			server.WriteResponseForPlayground(w, rtnValue, err)
		}
	}

	return s
}

func (s *ZnServer) SetHTTPHandler(entryFile string) *ZnServer {
	compiler := NewCompiler()
	s.reqHandler = func(w http.ResponseWriter, r *http.Request) {
		reqObj, err := server.ConstructHTTPRequestObject(r)
		if err != nil {
			server.SendHTTPResponse(nil, err, w)
			return
		}

		varInput := map[string]runtime.Element{
			"当前请求": reqObj,
		}
		// execute code
		rtnValue, err := compiler.LoadFile(entryFile).Execute(varInput)
		server.SendHTTPResponse(rtnValue, err, w)
	}
	return s
}

func (s *ZnServer) SetPMServerConfig(cfg server.ZnPMServerConfig) *ZnServer {
	s.pmServerConfig = cfg

	return s
}

/// Launch - by default launch PMServer
func (s *ZnServer) Launch(connUrl string) error {
	return s.LaunchPMServer(connUrl)
}

// PMServer works like php-fpm: at the very beginnning, we prefork some childWorker *processes* standby; when entering a new request, one of the childWorkers will undertake the request and return responses.
//
// The benefit of using PMServer is *safe*: when the execution of request
// is too long/dead loop, this childWorker will be killed automatically to avoid using resources too much.
func (s *ZnServer) LaunchPMServer(connUrl string) error {
	pmServer := server.NewZnPMServer(s.reqHandler)
	return pmServer.Start(connUrl, s.pmServerConfig)
}

func (s *ZnServer) LaunchThreadServer(connUrl string) error {
	threadServer := server.NewZnThreadServer(s.reqHandler)
	return threadServer.Start(connUrl)
}
