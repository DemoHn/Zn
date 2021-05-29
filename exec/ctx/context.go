package ctx

import (
	"github.com/DemoHn/Zn/debug"
)

// Context is a global variable that stores current execution
// states, global configurations
type Context struct {
	// globals - stores all global variables
	globals map[string]Value
	// import - imported value from stdlib or elsewhere
	imports map[string]Value
	// fileInfo -
	fileInfo *FileInfo
	// a seperate map to store inner debug data
	// usage: call （__probe：「tagName」，variable）
	// it will record all logs (including variable value, curernt scope, etc.)
	// the value is deep-copied so don't worry - the value logged won't be changed
	_probe *debug.Probe
	// Scope -
	scope *Scope
}

// NewContext - create new Zn Context. Notice through the life-cycle
// of one code execution, there's only one running context to store all states.
func NewContext(globalsMap map[string]Value) *Context {
	return &Context{
		globals: globalsMap,
		imports: map[string]Value{},
		_probe:  debug.NewProbe(),
		scope:   NewScope(),
	}
}

// GetScope -
func (ctx *Context) GetScope() *Scope {
	return ctx.scope
}

// SetScope -
func (ctx *Context) SetScope(sp *Scope) {
	ctx.scope = sp
}

// SetFileInfo
func (ctx *Context) SetFileInfo(fileInfo *FileInfo) {
	ctx.fileInfo = fileInfo
}

// GetFileInfo -
func (ctx *Context) GetFileInfo() *FileInfo {
	return ctx.fileInfo
}
