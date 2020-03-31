package exec

import (
	"github.com/DemoHn/Zn/lex"
)

// Env marks the property of current execution block and stores some necessary value.
// Values may differ from various block types.
type Env interface {
	GetLineStack() *lex.LineStack
	GetID() int
}

// ProgramEnv - env data for Program block
type ProgramEnv struct {
	lineStack         *lex.LineStack
	id                int
	forwardDeclareMap map[string]ZnValue
}

// NewProgramEnv -
func NewProgramEnv(ls *lex.LineStack) *ProgramEnv {
	return &ProgramEnv{ls, 0, map[string]ZnValue{}}
}

// GetLineStack -
func (pe *ProgramEnv) GetLineStack() *lex.LineStack {
	return pe.lineStack
}

// GetID -
func (pe *ProgramEnv) GetID() int {
	return pe.id
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
