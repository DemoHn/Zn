package exec

// ZnObject - the global Zn Object interface
type ZnObject interface {
	IsNull() bool
}

//// primitive types

// ZnString - Zn string concrete type
type ZnString struct {
	ZnObject
	Value    string
	nullFlag bool
}

// ZnInteger - Zn integer type (a sub-type of Number?)
type ZnInteger struct {
	ZnObject
	Value    int64
	nullFlag bool
}

// ZnFloat - Zn float type (a sub-type of Number?)
type ZnFloat struct {
	ZnObject
	Value    float64
	nullFlag bool
}

// ZnArray - Zn array type
type ZnArray struct {
	ZnObject
	Items    []ZnObject
	nullFlag bool
}
