package debug

import (
	"fmt"
	"time"
)

// Probe - debug probe
// log variable into when （__probe：xx，xx） is called
type Probe struct {
	info map[string][]ProbeLog
}

// ProbeLog -
type ProbeLog struct {
	probeTime time.Time
	// original value - DON'T USE ZnValue here to avoid circular dependency!
	value  interface{}
	valStr string
}

// NewProbe -
func NewProbe() *Probe {
	return &Probe{
		info: map[string][]ProbeLog{},
	}
}

// AddLog - add probe data to log
func (pb *Probe) AddLog(tag string, value interface{}) {
	fmt.Println(value)
}
