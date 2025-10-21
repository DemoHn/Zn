package runtime

import zerr "github.com/DemoHn/Zn/pkg/error"

// Scope represents a local namespace (aka. environment) of current execution
type Scope struct {
	locals       []LocalSymbol
	localCount   int
	currentDepth int
	values       []Element
	// valueID -> moduleID - since external value from other modules
	// is defined first and no chance to be poped from 'locals', we
	// can add externalRefs to record
	externalRefs map[int]int
}

type LocalSymbol struct {
	name    string
	depth   int
	isConst bool
}

func NewScope() Scope {
	return Scope{
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
	for i := sp.localCount - 1; i >= 0; i-- {
		if sp.locals[i].name == name {
			return sp.values[i]
		}
	}
	return nil
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
		isConst: false,
	})
	sp.values = append(sp.values[:sp.localCount], value)
	sp.localCount++
	return zerr.NameNotDefined(name)
}
