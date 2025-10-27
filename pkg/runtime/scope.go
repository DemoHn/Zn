package runtime

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
)

// Scope represents a local namespace (aka. environment) of current execution
type Scope struct {
	locals       []LocalSymbol
	localCount   int
	currentDepth int
	values       []Element
	// symbolID -> moduleID - since external value from other modules
	// is defined first and no chance to be poped from 'locals', we
	// can add externalRefs to record
	externalRefs map[int]int
}

type LocalSymbol struct {
	name    string
	depth   int
	isConst bool
}

func NewScope() *Scope {
	return &Scope{
		locals:       []LocalSymbol{},
		localCount:   0,
		currentDepth: 0,
		values:       []Element{},
		externalRefs: map[int]int{},
	}
}

func (sp *Scope) BeginScope() {
	sp.currentDepth++
}

func (sp *Scope) EndScope() {
	sp.currentDepth--

	// pop all deeper values
	for sp.localCount > 0 && sp.locals[sp.localCount-1].depth > sp.currentDepth {
		sp.localCount--
	}
}

func (sp *Scope) GetValue(name string) Element {
	symbolID := sp.getSymbolID(name)
	if symbolID >= 0 && symbolID < len(sp.values) {
		return sp.values[symbolID]
	}
	return nil
}

// when value is external, return the moduleID;if not found in external module, return -1
func (sp *Scope) GetValueWithModuleID(name string) (Element, int) {
	symbolID := sp.getSymbolID(name)
	if symbolID >= 0 && symbolID < len(sp.values) {
		extModuleID, ok := sp.externalRefs[symbolID]
		if ok {
			return sp.values[symbolID], extModuleID
		} else {
			return sp.values[symbolID], -1
		}
	}
	return nil, -1
}

// SetValue - set from existing symbol
func (sp *Scope) SetValue(name string, value Element) error {
	for i := sp.localCount - 1; i >= 0; i-- {
		if sp.locals[i].name == name {
			if sp.locals[i].isConst {
				// error: cannot change const value
				return zerr.AssignToConstant()
			}
			sp.values[i] = value
			return nil
		}
	}
	return zerr.NameNotDefined(name)
}

// DeclareValue
func (sp *Scope) DeclareValue(name string, value Element) error {
	return sp.declareValue(name, value, false)
}

func (sp *Scope) DeclareConstValue(name string, value Element) error {
	return sp.declareValue(name, value, true)
}

func (sp *Scope) DeclareExternalValue(name string, value Element, moduleID int) error {
	err := sp.declareValue(name, value, true)
	if err != nil {
		return err
	}
	// add external ref
	sp.externalRefs[sp.localCount-1] = moduleID
	return nil
}

// getSymbolID - get the latest symbolID that matches the name
// when not found, return -1
func (sp *Scope) getSymbolID(name string) int {
	for i := sp.localCount - 1; i >= 0; i-- {
		if sp.locals[i].name == name {
			return i
		}
	}
	return -1
}

// declareValue - add new symbol to scope
func (sp *Scope) declareValue(name string, value Element, isConst bool) error {
	for i := sp.localCount - 1; i >= 0; i-- {
		if sp.locals[i].depth < sp.currentDepth {
			break
		}
		if sp.locals[i].name == name {
			// redeclaration in the same depth leval is not allowed
			/*e.g.:
			{
				令 a = 1  // OK
				令 a = 2  // ERROR
			}
			*/
			if sp.locals[i].depth == sp.currentDepth {
				return zerr.NameRedeclared(name)
			}
		}
	}

	// add new symbol
	sp.locals = append(sp.locals[:sp.localCount], LocalSymbol{
		name:    name,
		depth:   sp.currentDepth,
		isConst: isConst,
	})
	sp.values = append(sp.values[:sp.localCount], value)
	sp.localCount++

	return nil
}
