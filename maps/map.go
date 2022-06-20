package maps

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/path"
)

// Map implements some quality of life extensions to a standard map[string]interface{}
type Map map[string]interface{}

// NewMap returns a fully initialized Map object.
func NewMap() Map {
	return Map(map[string]interface{}{})
}

// AsMapOfInterface returns the underlying map datastructure
func (m Map) AsMapOfInterface() map[string]interface{} {
	return map[string]interface{}(m)
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
func (m Map) GetInterface(name string) interface{} {
	return m[name]
}

// GetString returns a named option as a string type.
func (m Map) GetString(name string) string {
	return convert.String(m[name])
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

// GetSliceOfMap returns a named option as a slice of map.Map objects.
func (m Map) GetSliceOfMap(name string) []Map {
	value := convert.SliceOfMap(m[name])
	result := make([]Map, len(value))

	for index := range value {
		result[index] = Map(value[index])
	}

	return result
}

// GetMap returns a named option as a map.Map
func (m Map) GetMap(name string) Map {

	if value, ok := m[name].(Map); ok {
		return value
	}

	if value, ok := m[name].(map[string]interface{}); ok {
		return Map(value)
	}

	return Map{}
}

/****************************
 * Setters
 ****************************/

// SetBool adds a boolean value into the map
func (m Map) SetBool(name string, value bool) {
	m[name] = value
}

// SetInt adds an int value into the map
func (m Map) SetInt(name string, value int) {
	m[name] = value
}

// SetInt64 adds an int64 value into the map
func (m Map) SetInt64(name string, value int64) {
	m[name] = value
}

// SetFloat adds an int value into the map
func (m Map) SetFloat(name string, value float64) {
	m[name] = value
}

// SetString adds an int value into the map
func (m Map) SetString(name string, value string) {
	m[name] = value
}

/****************************
 * Path interfaces
 ****************************/

// GetPath implements the path.Getter interface
func (m Map) GetPath(name string) (interface{}, bool) {

	head, tail := list.Split(name, ".")

	if tail == "" {
		result, ok := m[head]
		return result, ok
	}

	return path.GetOK(m[head], tail)
}

// SetPath implements the path.Setter interface
func (m Map) SetPath(name string, value interface{}) error {

	head, tail := list.Split(name, ".")

	if tail == "" {
		m[head] = value
		return nil
	}

	return path.Set(m[head], tail, value)
}

// DeletePath implements the path.Deleter interface
func (m Map) DeletePath(name string) error {

	head, tail := list.Split(name, ".")

	if tail == "" {
		delete(m, head)
		return nil
	}

	temp := m[head]

	if err := path.Delete(temp, tail); err != nil {
		return derp.Wrap(err, "map.Map.DeletePath", "Error deleting from child element")
	}

	m[head] = temp

	return nil
}
