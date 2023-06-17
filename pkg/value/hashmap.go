package value

import (
	r "github.com/DemoHn/Zn/pkg/runtime"
)

// HashMap - represents for 列表类
type HashMap struct {
	value    map[string]r.Element
	keyOrder []string
	*r.ElementModel
}

// KVPair - key-value pair, used for ZnHashMap
type KVPair struct {
	Key   string
	Value r.Element
}

// NewHashMap -
func NewHashMap(kvPairs []KVPair) *HashMap {
	hm := &HashMap{
		value:        map[string]r.Element{},
		keyOrder:     []string{},
		ElementModel: r.NewElementModel(),
	}

	for _, kvPair := range kvPairs {
		if _, ok := hm.value[kvPair.Key]; !ok {
			// append distinct value
			hm.keyOrder = append(hm.keyOrder, kvPair.Key)
		}
		hm.value[kvPair.Key] = kvPair.Value
	}

	//// register getters & methods
	hm.RegisterGetter("数目", hm.hmGetLength)
	hm.RegisterGetter("长度", hm.hmGetLength)
	hm.RegisterGetter("所有索引", hm.hmGetAllIndexes)
	hm.RegisterGetter("所有值", hm.hmGetAllValues)

	hm.RegisterMethod("读取", hm.hmExecGet)
	hm.RegisterMethod("写入", hm.hmExecSet)
	hm.RegisterMethod("移除", hm.hmExecDelete)
	return hm
}

// GetKeyOrder -
func (hm *HashMap) GetKeyOrder() []string {
	return hm.keyOrder
}

// GetValue -
func (hm *HashMap) GetValue() map[string]r.Element {
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

//// getters, setters and methods
// getters
func (hm *HashMap) hmGetLength(c *r.Context) (r.Element, error) {
	return NewNumber(float64(len(hm.value))), nil
}

func (hm *HashMap) hmGetAllIndexes(c *r.Context) (r.Element, error) {
	var strs []r.Element
	for _, keyName := range hm.keyOrder {
		strs = append(strs, NewString(keyName))
	}
	return NewArray(strs), nil
}

func (hm *HashMap) hmGetAllValues(c *r.Context) (r.Element, error) {
	var vals []r.Element
	for _, keyName := range hm.keyOrder {
		vals = append(vals, hm.value[keyName])
	}
	return NewArray(vals), nil
}

// methods
func (hm *HashMap) hmExecGet(c *r.Context, values []r.Element) (r.Element, error) {
	if err := ValidateLeastParams(values, "string+"); err != nil {
		return nil, err
	}
	var result r.Element = hm

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

func (hm *HashMap) hmExecSet(c *r.Context, values []r.Element) (r.Element, error) {
	if err := ValidateExactParams(values, "string", "any"); err != nil {
		return nil, err
	}
	// key name
	keyName := values[0].(*String).value
	hm.AppendKVPair(KVPair{keyName, values[1]})
	return values[1], nil
}

func (hm *HashMap) hmExecDelete(c *r.Context, values []r.Element) (r.Element, error) {
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
