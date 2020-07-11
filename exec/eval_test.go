package exec

import "testing"

type znValueMap map[string]ZnValue
type programOKSuite struct {
	name           string
	program        string
	symbols        znValueMap
	expReturnValue ZnValue
	// root: <symbolMap>
	// root.iterate[1]: <symbolMap>
	// etc...
	expSymbols map[string]znValueMap
}

func TestA(t *testing.T) {

}

func assertProgrm(suite programOKSuite) {
	// new scope
	//ctx := NewContext()
}
