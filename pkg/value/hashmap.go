package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)
// TODO - add methods
// HashMap - represents for 列表类
type HashMap struct {
	value    map[string]r.Value
	keyOrder []string
}

// KVPair - key-value pair, used for ZnHashMap
type KVPair struct {
	Key   string
	Value r.Value
}

// NewHashMap -
func NewHashMap(kvPairs []KVPair) *HashMap {
	hm := &HashMap{
		value:    map[string]r.Value{},
		keyOrder: []string{},
	}

	for _, kvPair := range kvPairs {
		if _, ok := hm.value[kvPair.Key]; !ok {
			// append distinct value
			hm.keyOrder = append(hm.keyOrder, kvPair.Key)
		}
		hm.value[kvPair.Key] = kvPair.Value
	}

	return hm
}

// GetKeyOrder -
func (hm *HashMap) GetKeyOrder() []string {
	return hm.keyOrder
}

// GetValue -
func (hm *HashMap) GetValue() map[string]r.Value {
	return hm.value
}

// AppendKVPair -
func (hm *HashMap) AppendKVPair(pair KVPair) {
	key := pair.Key
	value := pair.Value
	_, ok := hm.value[key]
	if ok {
		hm.value[key] = value
		return
	}
	// insert new key
	hm.value[key] = value
	hm.keyOrder = append(hm.keyOrder, key)
}

// GetProperty -
func (hm *HashMap) GetProperty(c *r.Context, name string) (r.Value, error) {
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (hm *HashMap) SetProperty(c *r.Context, name string, value r.Value) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (hm *HashMap) ExecMethod(c *r.Context, name string, values []r.Value) (r.Value, error) {
	return nil, zerr.MethodNotFound(name)
}
