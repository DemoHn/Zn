package zinc

import (
	"github.com/DemoHn/Zn/pkg/exec"
	runtime "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/server"
	"github.com/DemoHn/Zn/pkg/value"

	// stdlibs
	libFile "github.com/DemoHn/Zn/stdlib/file"
	libHttp "github.com/DemoHn/Zn/stdlib/http"
	libJson "github.com/DemoHn/Zn/stdlib/json"
)

type Element = runtime.Element
type ElementMap = runtime.ElementMap
type Library = runtime.Library

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

// servers & handlers
var NewPlaygroundHandler = server.NewZnPlaygroundHandler
var NewHttpHandler = server.NewZnHttpHandler

var NewPMServer = server.NewZnPMServer
var NewThreadServer = server.NewZnThreadServer

const ZINC_VERSION = "rev08"

var StandardLibs = []*runtime.Library{
	libHttp.Export(),
	libJson.Export(),
	libFile.Export(),
}

// ZnInterpreter - MAIN CODE EXECUTION INSTANCE -
// ONE INTERPRETER -> ONE VM
type ZnInterpreter = exec.Interpreter

// NewInterpreter - new ZnInterpreter object
func NewInterpreter() *ZnInterpreter {
	interpreter := exec.NewInterpreter(ZINC_VERSION).
		SetExternalLibs(StandardLibs)

	return interpreter
}
