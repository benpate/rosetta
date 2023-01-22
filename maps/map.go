package maps

import (
	"net/url"

	"github.com/benpate/rosetta/convert"
)

// Map implements some quality of life extensions to a standard map[string]any
type Map map[string]any

// New returns a fully initialized Map object.
func New() Map {
	return make(Map)
}

func FromURLValues(value url.Values) Map {

	result := New()

	for key := range value {
		if len(value[key]) == 1 {
			result[key] = value.Get(key)
		} else {
			result[key] = value[key]
		}
	}

	return result
}

// AsMapOfInterface returns the underlying map datastructure
func (m Map) AsMapOfInterface() map[string]any {
	return map[string]any(m)
}

// GetKeys returns all keys of the underlying map
func (m Map) GetKeys() []string {
	result := make([]string, len(m))

	index := 0
	for key := range m {
		result[index] = key
		index = index + 1
	}

	return result
}

/******************************************
 * Schema Getter Interfaces
 ******************************************/

func (m Map) GetBool(name string) (bool, bool) {
	return convert.BoolOk(m[name], false)
}

func (m Map) GetInt(name string) (int, bool) {
	return convert.IntOk(m[name], 0)
}

func (m Map) GetInt64(name string) (int64, bool) {
	return convert.Int64Ok(m[name], 0)
}

// GetInterface returns a named option without any conversion.  You get what you get.
func (m Map) GetInterface(name string) (any, bool) {
	result, ok := m[name]
	return result, ok
}

func (m Map) GetFloat(name string) (float64, bool) {
	return convert.FloatOk(m[name], 0)
}

func (m Map) GetString(name string) (string, bool) {
	return convert.StringOk(m[name], "")
}

/******************************************
 * Other Getter Interfaces
 ******************************************/

// GetSliceOfString returns a named option as a slice of strings
func (m Map) GetSliceOfString(name string) []string {
	return convert.SliceOfString(m[name])
}

// GetSliceOfInt returns a named option as a slice of int values
func (m Map) GetSliceOfInt(name string) []int {
	return convert.SliceOfInt(m[name])
}

// GetSliceOfFloat returns a named option as a slice of float64 values
func (m Map) GetSliceOfFloat(name string) []float64 {
	return convert.SliceOfFloat(m[name])
}

// GetSliceOfMap returns a named option as a slice of maps.Map objects.
func (m Map) GetSliceOfMap(name string) []Map {
	value := convert.SliceOfMap(m[name])
	result := make([]Map, len(value))

	for index := range value {
		result[index] = Map(value[index])
	}

	return result
}

// GetMap returns a named option as a maps.Map
func (m Map) GetMap(name string) Map {

	if value, ok := m[name].(Map); ok {
		return value
	}

	if value, ok := m[name].(map[string]any); ok {
		return Map(value)
	}

	return Map{}
}

/******************************************
 * Schema Setter Interfaces
 ******************************************/

// SetBool adds a boolean value into the map
func (m *Map) SetBool(name string, value bool) bool {
	(*m)[name] = value
	return true
}

// SetInt adds an int value into the map
func (m *Map) SetInt(name string, value int) bool {
	(*m)[name] = value
	return true
}

// SetInt64 adds an int64 value into the map
func (m *Map) SetInt64(name string, value int64) bool {
	(*m)[name] = value
	return true
}

// SetFloat adds an int value into the map
func (m *Map) SetFloat(name string, value float64) bool {
	(*m)[name] = value
	return true
}

// SetString adds an int value into the map
func (m *Map) SetString(name string, value string) bool {
	(*m)[name] = value
	return true
}

/******************************************
 * Schema Remover Interfaces
 ******************************************/

func (m *Map) Remove(name string) bool {
	delete(*m, name)
	return true
}

/******************************************
 * Tree Traversal
 ******************************************/

func (m *Map) GetObjectOK(name string) (any, bool) {

	if _, ok := (*m)[name]; !ok {
		(*m)[name] = make(Map)
	}

	return (*m)[name], true
}
