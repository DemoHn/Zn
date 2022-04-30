package runtime

type Context struct {
	rootModule *Module
	scopeCursor *Scope
}

// NewContext -
func NewContext() *Context {
	return &Context{}
}