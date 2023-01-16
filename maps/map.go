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

// GetInterface returns a named option without any conversion.  You get what you get.
func (m Map) GetInterface(name string) any {
	return m[name]
}

// GetString returns a named option as a string type.
func (m Map) GetString(name string) string {
	return convert.String(m[name])
}

// GetBytes returns a named option as a slice of bytes.
func (m Map) GetBytes(name string) []byte {
	return convert.Bytes(m[name])
}

// GetBool returns a named option as a bool type.
func (m Map) GetBool(name string) bool {
	return convert.Bool(m[name])
}

// GetInt returns a named option as an int type.
func (m Map) GetInt(name string) int {
	return convert.Int(m[name])
}

// GetInt64 returns a named option as an int64 type.
func (m Map) GetInt64(name string) int64 {
	return convert.Int64(m[name])
}

// GetFloat returns a named option as a float type.
func (m Map) GetFloat(name string) float64 {
	return convert.Float(m[name])
}

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

/****************************
 * Setters
 ****************************/

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

/****************************
 * Deleter
 ****************************/

func (m *Map) Delete(name string) bool {
	delete(*m, name)
	return true
}

/****************************
 * Tree Traversal
 ****************************/

func (m *Map) GetChild(name string) (any, bool) {

	if _, ok := (*m)[name]; !ok {
		(*m)[name] = make(Map)
	}

	return (*m)[name], true
}
