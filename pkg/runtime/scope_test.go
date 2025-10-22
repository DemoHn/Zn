package runtime

import (
	"testing"
)

type MockValue struct {
	value int
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

// test beginScope() -
func TestScope_BeginScopeAndAddValue(t *testing.T) {
	initScope := NewScope()
	initScope.BeginScope()
	initScope.DeclareValue("T1", MockValue{1.0})
	initScope.DeclareValue("T2", MockValue{2.0})

	// assert currentDepth = 1
	if initScope.currentDepth != 1 {
		t.Errorf("currentDepth = %d, want 1", initScope.currentDepth)
	}
	// assert localCount = 2
	if initScope.localCount != 2 || len(initScope.locals) != 2 {
		t.Errorf("localCount = %d, want 2", initScope.localCount)
	}
	// assert values to be Number(1.0) and Number(2.0)
	if initScope.values[0].(MockValue).value != 1 || initScope.values[1].(MockValue).value != 2 {
		t.Errorf("values = %v, want [1.0, 2.0]", initScope.values)
	}
}

// test declareValue() -
func TestScope_DeclareValue(t *testing.T) {
	initScope := NewScope()
	initScope.BeginScope()
	initScope.DeclareValue("T1", MockValue{1.0})
	initScope.DeclareValue("T2", MockValue{2.0})
	initScope.BeginScope()
	initScope.DeclareValue("T1", MockValue{3.0})

	// assert currentDepth = 2
	if initScope.currentDepth != 2 {
		t.Errorf("currentDepth = %d, want 2", initScope.currentDepth)
	}
	// assert localCount = 2
	if initScope.localCount != 3 || len(initScope.locals) != 3 {
		t.Errorf("localCount = %d, want 3", initScope.localCount)
	}
}
