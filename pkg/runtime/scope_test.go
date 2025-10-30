package runtime

import (
	"strconv"
	"strings"
	"testing"
)

type MockValue struct {
	value string
}

// impl Element (in value.go) interface for MockValue
func (m MockValue) GetProperty(name string) (Element, error) {
	return nil, nil
}

func (m MockValue) SetProperty(name string, value Element) error {
	return nil
}

func (m MockValue) ExecMethod(name string, values []Element) (Element, error) {
	return nil, nil
}

////// helper functions //////
/**
scope snapshot format example:
locals = T1,1,true;T2,1,false
localCount = 2
currentDepth = 2
values = AAA,BBB
*/
func assertSnapshot(t *testing.T, sp *Scope, scopeSnapshot string) {
	// 1. parse and assert snapshot string and compare with Scope
	lines := strings.Split(strings.ReplaceAll(scopeSnapshot, "\r\n", "\n"), "\n")
	snapshotItemMap := make(map[string]string)
	for _, line := range lines {
		parts := strings.Split(line, "=")
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		snapshotItemMap[key] = value
	}
	// assert locals
	locals := strings.Split(snapshotItemMap["locals"], ";")
	// assert locals length
	if len(locals) != sp.localCount {
		t.Errorf("localCount = %d, want %d", sp.localCount, len(locals))
	}
	if len(locals) != len(sp.locals) {
		t.Errorf("locals length = %d, want %d", len(sp.locals), len(locals))
	}

	for idx, local := range locals {
		parts := strings.Split(local, ",")
		if len(parts) != 3 {
			t.Errorf("local = %s, want 3 parts", local)
		}
		spLocal := sp.locals[idx]

		if parts[0] != spLocal.name {
			t.Errorf("idx=[%d] local.name = %s, want %s", idx, spLocal.name, parts[0])
		}
		if parts[1] != strconv.Itoa(spLocal.depth) {
			t.Errorf("idx=[%d] local.depth = %d, want %s", idx, spLocal.depth, parts[1])
		}
		if parts[2] != strconv.FormatBool(spLocal.isConst) {
			t.Errorf("idx=[%d] local.isConst = %t, want %s", idx, spLocal.isConst, parts[2])
		}
	}
	// assert localCount
	if snapshotItemMap["localCount"] != strconv.Itoa(sp.localCount) {
		t.Errorf("localCount = %s, want %s", snapshotItemMap["localCount"], strconv.Itoa(sp.localCount))
	}
	// assert currentDepth
	if snapshotItemMap["currentDepth"] != strconv.Itoa(sp.currentDepth) {
		t.Errorf("currentDepth = %s, want %s", snapshotItemMap["currentDepth"], strconv.Itoa(sp.currentDepth))
	}
	// assert values
	values := strings.Split(snapshotItemMap["values"], ",")
	if len(values) != len(sp.values) {
		t.Errorf("values length = %d, want %d", len(sp.values), len(values))
	}
	for idx, value := range values {
		if value != sp.values[idx].(MockValue).value {
			t.Errorf("idx=[%d] value = %s, want %s", idx, sp.values[idx].(MockValue).value, value)
		}
	}
}

// test beginScope() -
func TestScope_BeginScopeAndAddValue(t *testing.T) {
	// action
	initScope := NewScope()
	initScope.BeginScope()
	initScope.DeclareValue("T1", MockValue{"AAA"})
	initScope.DeclareValue("T2", MockValue{"BBB"})

	// snapshot
	scopeSnapshot := `
locals = T1,1,false;T2,1,false
localCount = 2
currentDepth = 1
values = AAA,BBB
`

	assertSnapshot(t, initScope, scopeSnapshot)
}

// test declareValue() -
func TestScope_DeclareValueWithNestedScope(t *testing.T) {
	initScope := NewScope()

	// depth=1
	initScope.BeginScope()
	initScope.DeclareValue("T1", MockValue{"aa"})
	initScope.DeclareValue("T2", MockValue{"bb"})
	// depth=2
	initScope.BeginScope()
	initScope.DeclareValue("T3", MockValue{"cc"})
	// depth=3
	initScope.BeginScope()
	initScope.DeclareValue("T4", MockValue{"dd"})
	initScope.EndScope()
	// depth=2
	initScope.DeclareValue("T5", MockValue{"ee"})
	initScope.EndScope()
	// depth=1
	initScope.DeclareValue("T6", MockValue{"ff"})

	// snapshot
	scopeSnapshot := `
locals = T1,1,false;T2,1,false;T6,1,false
localCount = 3
currentDepth = 1
values = aa,bb,ff
`
	assertSnapshot(t, initScope, scopeSnapshot)
}

func TestScope_SetScopeValue(t *testing.T) {
	initScope := NewScope()

	// Step1: declare value
	// depth=1
	initScope.BeginScope()
	initScope.DeclareValue("T1", MockValue{"aa"})
	initScope.DeclareValue("T2", MockValue{"bb"})

	// assert snapshot first
	scopeSnapshot := `
locals = T1,1,false;T2,1,false
localCount = 2
currentDepth = 1
values = aa,bb
`
	assertSnapshot(t, initScope, scopeSnapshot)

	// Step2: set value to cc
	initScope.SetValue("T1", MockValue{"cc"})

	// snapshot
	scopeSnapshot = `
locals = T1,1,false;T2,1,false
localCount = 2
currentDepth = 1
values = cc,bb
`
	// Step3: add new scope
	initScope.BeginScope()
	initScope.DeclareValue("T3", MockValue{"dd"})
	// another T2, same name but different depth
	initScope.DeclareValue("T2", MockValue{"ee"})

	// only the latter T2's value will be updated
	initScope.SetValue("T2", MockValue{"ff"})

	scopeSnapshot = `
locals = T1,1,false;T2,1,false;T3,2,false;T2,2,false
localCount = 4
currentDepth = 2
values = cc,bb,dd,ff
`
	assertSnapshot(t, initScope, scopeSnapshot)
}

func TestScope_SetConstValue_SHDFAIL(t *testing.T) {
	initScope := NewScope()
	initScope.BeginScope()
	initScope.DeclareValue("T1", MockValue{"aa"})
	initScope.DeclareConstValue("T2", MockValue{"bb"})
	// initScope.SetValue("T1", MockValue{"bb"})

	scopeSnapshot := `
locals = T1,1,false;T2,1,true
localCount = 2
currentDepth = 1
values = aa,bb`
	assertSnapshot(t, initScope, scopeSnapshot)

	// Step2: try to set a const value - SHDFAIL
	err := initScope.SetValue("T2", MockValue{"cc"})
	if err == nil {
		t.Errorf("SetValue should fail, but no error returned")
	}
}

func TestScope_RedeclareValue_SHDFAIL(t *testing.T) {
	initScope := NewScope()
	initScope.BeginScope()
	initScope.DeclareValue("T1", MockValue{"aa"})
	initScope.BeginScope()
	// another scope
	initScope.DeclareValue("T1", MockValue{"bb"})

	scopeSnapshot := `
locals = T1,1,false;T1,2,false
localCount = 2
currentDepth = 2
values = aa,bb
`
	assertSnapshot(t, initScope, scopeSnapshot)

	// Step2: try to redeclare a value - SHDFAIL
	err := initScope.DeclareValue("T1", MockValue{"cc"})
	if err == nil {
		t.Errorf("DeclareValue should fail, but no error returned")
	}
}
