package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Map map[string]interface{}

func NewMap() *Map {
	return &Map{}
}
func NewMaps(body interface{}) *Map {
	e := &Map{}
	if body == nil {
		return e
	}
	switch body.(type) {
	case string:
		json.Unmarshal([]byte(body.(string)), e)
		break
	case []byte:
		json.Unmarshal(body.([]byte), e)
		break
	case map[string]interface{}:
		for k, v := range body.(map[string]interface{}) {
			e.Set(k, v)
		}
		break
	case *map[string]interface{}:
		for k, v := range *(body.(*map[string]interface{})) {
			e.Set(k, v)
		}
		break
	default:
		bts, err := json.Marshal(body)
		if err == nil {
			json.Unmarshal(bts, e)
		}
		break
	}
	return e
}

func (e *Map) Get(key string) (interface{}, bool) {
	v, ok := (*e)[key]
	return v, ok
}
func (e *Map) Set(key string, val interface{}) {
	(*e)[key] = val
}
func (e *Map) Map() map[string]interface{} {
	return *e
}
func (e *Map) ToBytes() []byte {
	bts, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bts
}
func (e *Map) ToString() string {
	bts, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(bts)
}
func (e *Map) GetString(key string) string {
	v, ok := e.Get(key)
	if !ok {
		return ""
	}
	switch v.(type) {
	case float32:
		return fmt.Sprintf("%d", int64(v.(float32)))
	case float64:
		return fmt.Sprintf("%d", int64(v.(float64)))
	}
	return fmt.Sprintf("%v", v)
}
func (e *Map) GetInt(key string) (int64, error) {
	v, ok := e.Get(key)
	if !ok {
		return 0, errors.New("not found key")
	}
	switch v.(type) {
	case int:
		return v.(int64), nil
	case string:
		return strconv.ParseInt(v.(string), 10, 64)
	case int64:
		return v.(int64), nil
	case float32:
		return int64(v.(float32)), nil
	case float64:
		return int64(v.(float64)), nil
	}
	return 0, errors.New("not found")
}
func (e *Map) GetFloat(key string) (float64, error) {
	v, ok := e.Get(key)
	if !ok {
		return 0, errors.New("not found key")
	}
	switch v.(type) {
	case int:
		return float64(v.(int)), nil
	case string:
		return strconv.ParseFloat(v.(string), 64)
	case int64:
		return float64(v.(int64)), nil
	case float32:
		return float64(v.(float32)), nil
	case float64:
		return v.(float64), nil
	}
	return 0, errors.New("not found")
}
func (e *Map) GetBool(key string) bool {
	v, ok := e.Get(key)
	if ok {
		switch v.(type) {
		case bool:
			return v.(bool)
		}
	}
	return false
}
