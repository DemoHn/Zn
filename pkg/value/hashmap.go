package value

import (
	zerr "github.com/DemoHn/Zn/pkg/error"
	r "github.com/DemoHn/Zn/pkg/runtime"
)

type hmGetterFunc func(*HashMap, *r.Context) (r.Value, error)
type hmMethodFunc func(*HashMap, *r.Context, []r.Value) (r.Value, error)

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
	hmGetterMap := map[string]hmGetterFunc{
		"数目": hmGetLength,
		"长度": hmGetLength,
		"所有索引": hmGetAllIndexes,
		"所有值": hmGetAllValues,
	}
	if fn, ok := hmGetterMap[name]; ok {
		return fn(hm, c)
	}
	return nil, zerr.PropertyNotFound(name)
}

// SetProperty -
func (hm *HashMap) SetProperty(c *r.Context, name string, value r.Value) error {
	return zerr.PropertyNotFound(name)
}

// ExecMethod -
func (hm *HashMap) ExecMethod(c *r.Context, name string, values []r.Value) (r.Value, error) {
	hmMethodMap := map[string]hmMethodFunc{
		"读取": hmExecGet,
		"写入": hmExecSet,
		"移除": hmExecDelete,
	}
	if fn, ok := hmMethodMap[name]; ok {
		return fn(hm, c, values)
	}
	return nil, zerr.MethodNotFound(name)
}

//// getters, setters and methods

// getters
func hmGetLength(hm *HashMap, c *r.Context) (r.Value, error) {
	return NewNumber(float64(len(hm.value))), nil
}

func hmGetAllIndexes(hm *HashMap, c *r.Context) (r.Value, error) {
	var strs []r.Value
	for _, keyName := range hm.keyOrder {
		strs = append(strs, NewString(keyName))
	}
	return NewArray(strs), nil
}

func hmGetAllValues(hm *HashMap, c *r.Context) (r.Value, error) {
	var vals []r.Value
	for _, keyName := range hm.keyOrder {
		vals = append(vals, hm.value[keyName])
	}
	return NewArray(vals), nil
}

// methods
func hmExecGet(hm *HashMap, c *r.Context, values []r.Value) (r.Value, error) {
	if err := ValidateLeastParams(values, "string+"); err != nil {
		return nil, err
	}
	var result r.Value = hm

	for _, keyName := range values {
		keyNameStr := keyName.(*String).value
		if cursorHM, ok := result.(*HashMap); ok {
			if val, ok2 := cursorHM.value[keyNameStr]; ok2 {
				result = val
			} else {
				return NewNull(), nil
			}
		} else {
			return NewNull(), nil
		}
	}
	return result, nil
}

func hmExecSet(hm *HashMap, c *r.Context, values []r.Value) (r.Value, error) {
	if err := ValidateExactParams(values, "string", "any"); err != nil {
		return nil, err
	}
	// key name
	keyName := values[0].(*String).value
	hm.AppendKVPair(KVPair{keyName, values[1]})
	return values[1], nil
}

func hmExecDelete(hm *HashMap, c *r.Context, values []r.Value) (r.Value, error) {
	if err := ValidateExactParams(values, "string"); err != nil {
		return nil, err
	}
	// key name
	keyName := values[0].(*String).value
	val, ok := hm.value[keyName]
	if ok {
		delete(hm.value, keyName)
		// find & delete key from keyOrder
		for idx, vk := range hm.keyOrder {
			// if found, delete item from keyOrder to stop loop
			if vk == keyName {
				hm.keyOrder = append(hm.keyOrder[:idx], hm.keyOrder[idx+1:]...)
			}
		}
		return val, nil
	}
	// if delete key not found in hashmap, return null directly
	return NewNull(), nil
}