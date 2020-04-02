package exec

// Env marks the property of current execution block and stores some necessary value.
// Values may differ from various block types.
type Env interface {
	GetID() int
}

// EnvBlockType - enum block types
type EnvBlockType int

var (
	programIDcount int
	loopIDcount    int
	funcIDcount    int
)

// declare blockTypes
const (
	BlockProg EnvBlockType = 1
	BlockFunc EnvBlockType = 2
	BlockLoop EnvBlockType = 3
)

// ProgramEnv - env data for Program block
type ProgramEnv struct {
	id                int
	forwardDeclareMap map[string]ZnValue
}

// NewProgramEnv -
func NewProgramEnv() *ProgramEnv {
	return &ProgramEnv{0, map[string]ZnValue{}}
}

// GetID -
func (pe *ProgramEnv) GetID() int {
	return pe.id
}

// LoopEnv - env data under loop statement (每当)
type LoopEnv struct {
	id int
}

// GetID -
func (le *LoopEnv) GetID() int {
	return le.id
}

// FuncEnv - env data of a function (如何)
type FuncEnv struct {
	id                int
	forwardDeclareMap map[string]ZnValue
}

// GetID -
func (le *FuncEnv) GetID() int {
	return le.id
}

// DupProgramEnv -
func DupProgramEnv(env Env) *ProgramEnv {
	programIDcount++
	return &ProgramEnv{
		id:                programIDcount,
		forwardDeclareMap: map[string]ZnValue{},
	}
}

// DupLoopEnv -
func DupLoopEnv(env Env) *LoopEnv {
	loopIDcount++
	return &LoopEnv{
		id: loopIDcount,
	}
}

// DupFuncEnv -
func DupFuncEnv(env Env) *FuncEnv {
	funcIDcount++
	return &FuncEnv{
		id:                funcIDcount,
		forwardDeclareMap: map[string]ZnValue{},
	}
}

/**
// env.go manages Env struct that stores some essential info of the current nearest scope
// (block)

// EnvBlockType - enum block types
type EnvBlockType int

// Env - environment
type Env struct {
	*lex.LineStack
	idBlockFUNC  int
	idBlockWHILE int
	blockType    EnvBlockType
}

// declare blockTypes
const (
	BlockPROG  EnvBlockType = 1
	BlockFUNC  EnvBlockType = 2
	BlockWHILE EnvBlockType = 3
)

// NewEnv - new Env struct
func NewEnv(ls *lex.LineStack) *Env {
	return &Env{ls, 0, 0, BlockPROG}
}

// Dup - duplicate an new env
func (e *Env) Dup(blockType EnvBlockType) *Env {
	newEnv := new(Env)
	// copy items one by one
	newEnv.LineStack = e.LineStack
	newEnv.idBlockFUNC = e.idBlockFUNC
	newEnv.idBlockWHILE = e.idBlockWHILE
	newEnv.blockType = blockType

	// add ID
	switch blockType {
	case BlockFUNC:
		newEnv.idBlockFUNC++
	case BlockWHILE:
		newEnv.idBlockWHILE++
	}
	return newEnv
}

// GetBlockType -
func (e *Env) GetBlockType() EnvBlockType {
	return e.blockType
}
*/
