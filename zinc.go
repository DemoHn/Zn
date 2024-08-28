package zinc

import (
	"net/http"
	"os"

	"github.com/DemoHn/Zn/pkg/exec"
	runtime "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/server"
)

type Element = runtime.Element

// default serverHandlers - serverHandler tells zinc compiler how to
// parse & execute HTTP requests [from user to the built-in server]
//
// Normally, you could use the pre-defined default handlers to help accomplish your tasks; also, DIY your handler is feasible - just define a
// `func(w http.ResponseWriter, r *http.Request)` function (aka. type `http.HandlerFunc`) is OK.
var DefaultPlaygroundHandler = server.PlaygroundHandler
var DefaultHTTPHandlerEntryWrap = server.HTTPHandlerWithEntry
var DefaultPMServerConfig = server.ZnPMServerConfig{
	InitProcs: 20,
	MaxProcs:  100,
	Timeout:   60,
}

const ZINC_VERSION = "rev07"

type ZnCompiler struct {
	version        string
	PMServerConfig server.ZnPMServerConfig
}

// NewCompiler - new ZnCompiler object
func NewCompiler() *ZnCompiler {
	return &ZnCompiler{
		version:        ZINC_VERSION,
		PMServerConfig: DefaultPMServerConfig,
	}
}

// GetVersion - get current compiler's version
func (cp *ZnCompiler) GetVersion() string {
	return cp.version
}

// PMServer works like php-fpm: at the very beginnning, we prefork some childWorker *processes* standby; when entering a new request, one of the childWorkers will undertake the request and return responses.
//
// The benefit of using PMServer is *safe*: when the execution of request
// is too long/dead loop, this childWorker will be killed automatically to avoid using resources too much.
func (cp *ZnCompiler) StartPMServer(url string, handler http.HandlerFunc) error {
	childWorkerFlag := false
	pmServer := server.NewZnPMServer(handler)

	// check if --child-worker flag exists
	for _, arg := range os.Args {
		if arg == "--child-worker" {
			childWorkerFlag = true
			break
		}
	}

	///// run child worker if  --child-worker = true & preForkChild env is "OK"
	if childWorkerFlag && os.Getenv(server.EnvPreforkChildKey) == server.EnvPreforkChildVal {
		// start child worker to handle requests
		if err := pmServer.StartWorker(); err != nil {
			return err
		}
	} else {
		//// otherwise, just listen to the server
		if err := pmServer.StartMaster(url, cp.PMServerConfig); err != nil {
			return err
		}
	}

	return nil
}

func (cp *ZnCompiler) StartThreadServer(url string, handler http.HandlerFunc) error {
	return server.NewZnThreadServer(handler).Start(url)
}

// Run - exec a code snippet without any varInput or historial variables
func (cp *ZnCompiler) Run(code []byte) (Element, error) {
	return exec.NewPlaygroundExecutor(code).RunCode(map[string]Element{})
}

func (cp *ZnCompiler) ExecFile(file string) (Element, error) {
	return exec.NewFileExecutor(file).RunCode(map[string]Element{})
}
