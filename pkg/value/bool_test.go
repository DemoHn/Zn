package value

import (
	"testing"
)

func TestBool_String(t *testing.T) {
	// case 1
	boolTrue := NewBool(true)
	if boolTrue.String() != "真" {
		t.Errorf("bool string value not expected as 真")
	}
	boolFalse := NewBool(false)
	if boolFalse.String() != "假" {
		t.Errorf("bool string value not expected as 假")
	}
}
