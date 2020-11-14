package exec

import "github.com/DemoHn/Zn/error"

// HashMap -
type HashMap struct {
	value    map[string]Value
	keyOrder []string
}

// KVPair - key-value pair, used for ZnHashMap
type KVPair struct {
	Key   string
	Value Value
}

// NewZnHashMap -
func NewZnHashMap(kvPairs []KVPair) *HashMap {
	hm := &HashMap{
		value:    map[string]Value{},
		keyOrder: []string{},
	}

	for _, kvPair := range kvPairs {
		hm.value[kvPair.Key] = kvPair.Value
		hm.keyOrder = append(hm.keyOrder, kvPair.Key)
	}

	return hm
}

// GetProperty -
func (hm *HashMap) GetProperty(ctx *Context, name string) (Value, *error.Error) {
	switch name {
	case "长度":
		return NewDecimalFromInt(len(hm.value), 0), nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (hm *HashMap) SetProperty(ctx *Context, name string, value Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (hm *HashMap) ExecMethod(ctx *Context, name string, values []Value) (Value, *error.Error) {
	return nil, error.MethodNotFound(name)
}
