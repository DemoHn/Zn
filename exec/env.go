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
