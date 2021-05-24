package val

import "github.com/DemoHn/Zn/error"

// HashMap - represents for 列表类
type HashMap struct {
	value    map[string]ctx.Value
	keyOrder []string
}

// KVPair - key-value pair, used for ZnHashMap
type KVPair struct {
	Key   string
	ctx.Value ctx.Value
}

// NewHashMap -
func NewHashMap(kvPairs []KVPair) *HashMap {
	hm := &HashMap{
		value:    map[string]ctx.Value{},
		keyOrder: []string{},
	}

	for _, kvPair := range kvPairs {
		if _, ok := hm.value[kvPair.Key]; !ok {
			// append distinct value
			hm.keyOrder = append(hm.keyOrder, kvPair.Key)
		}
		hm.value[kvPair.Key] = kvPair.ctx.Value
	}

	return hm
}

// GetProperty -
func (hm *HashMap) GetProperty(ctx *ctx.Context, name string) (ctx.Value, *error.Error) {
	switch name {
	case "数目", "长度":
		return NewDecimalFromInt(len(hm.value), 0), nil
	case "所有索引":
		strs := []ctx.Value{}
		for _, keyName := range hm.keyOrder {
			strs = append(strs, NewString(keyName))
		}
		return NewArray(strs), nil
	case "所有值":
		vals := []ctx.Value{}
		for _, keyName := range hm.keyOrder {
			vals = append(vals, hm.value[keyName])
		}
		return NewArray(vals), nil
	}
	return nil, error.PropertyNotFound(name)
}

// SetProperty -
func (hm *HashMap) SetProperty(ctx *ctx.Context, name string, value ctx.Value) *error.Error {
	return error.PropertyNotFound(name)
}

// ExecMethod -
func (hm *HashMap) ExecMethod(ctx *ctx.Context, name string, values []ctx.Value) (ctx.Value, *error.Error) {
	switch name {
	case "读取":
		if err := validateExactParams(values, "string"); err != nil {
			return nil, err
		}
		// key name
		keyName := values[0].(*String).value
		val, ok := hm.value[keyName]
		if ok {
			return val, nil
		}
		return NewNull(), nil
	case "写入":
		if err := validateExactParams(values, "string", "any"); err != nil {
			return nil, err
		}
		// key name
		keyName := values[0].(*String).value
		val, ok := hm.value[keyName]
		if ok {
			hm.value[keyName] = values[1] // update value
			return val, nil
		}
		// insert new key
		hm.value[keyName] = values[1]
		hm.keyOrder = append(hm.keyOrder, keyName)
		return val, nil
	case "移除":
		if err := validateExactParams(values, "string"); err != nil {
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
	return nil, error.MethodNotFound(name)
}
