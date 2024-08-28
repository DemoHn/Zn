package zinc

import (
	"github.com/DemoHn/Zn/pkg/exec"
	runtime "github.com/DemoHn/Zn/pkg/runtime"
)

type Element = runtime.Element

const ZINC_VERSION = "rev07"

type ZnCompiler struct {
	version string
}

// NewCompiler - new ZnCompiler object
func NewCompiler() *ZnCompiler {
	return &ZnCompiler{
		version: ZINC_VERSION,
	}
}

// GetVersion - get current compiler's version
func (cp *ZnCompiler) GetVersion() string {
	return cp.version
}

// Run - exec a code snippet without any varInput or historial variables
func (cp *ZnCompiler) Run(code []byte) (Element, error) {
	return exec.NewPlaygroundExecutor(code).RunCode(map[string]Element{})
}
