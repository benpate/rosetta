package mapof

import (
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/schema"
)

type Any map[string]any

func NewAny() Any {
	return make(Any)
}

/****************************************
 * Path Getters
 ****************************************/

func (x Any) GetAny(key string) any {
	result, _ := x.GetAnyOK(key)
	return result
}

func (x Any) GetAnyOK(key string) (any, bool) {
	if value, ok := x[key]; ok {
		return convert.Interface(value), true
	}
	return nil, false

}

func (x Any) GetBool(key string) bool {
	result, _ := x.GetBoolOK(key)
	return result
}

func (x Any) GetBoolOK(key string) (bool, bool) {
	if value, ok := x[key]; ok {
		return convert.BoolOk(value, false)
	}
	return false, false
}

func (x Any) GetFloat(key string) float64 {
	result, _ := x.GetFloatOK(key)
	return result
}

func (x Any) GetFloatOK(key string) (float64, bool) {
	if value, ok := x[key]; ok {
		return convert.FloatOk(value, 0)
	}
	return 0, false
}

func (x Any) GetInt(key string) int {
	result, _ := x.GetIntOK(key)
	return result
}

func (x Any) GetIntOK(key string) (int, bool) {
	if value, ok := x[key]; ok {
		return convert.IntOk(value, 0)
	}
	return 0, false
}

func (x Any) GetInt64(key string) int64 {
	result, _ := x.GetInt64OK(key)
	return result
}

func (x Any) GetInt64OK(key string) (int64, bool) {
	if value, ok := x[key]; ok {
		return convert.Int64Ok(value, 0)
	}
	return 0, false
}

func (x Any) GetString(key string) string {
	result, _ := x.GetStringOK(key)
	return result
}

func (x Any) GetStringOK(key string) (string, bool) {
	if value, ok := x[key]; ok {
		return convert.StringOk(value, "")
	}
	return "", false
}

/****************************************
 * Path Setters
 ****************************************/

func (x *Any) SetAny(key string, value any) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Any) SetBool(key string, value bool) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Any) SetFloat(key string, value float64) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Any) SetInt(key string, value int) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Any) SetInt64(key string, value int64) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Any) SetString(key string, value string) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Any) makeNotNil() {
	if *x == nil {
		*x = make(Any)
	}
}

/****************************************
 * Tree Traversal
 ****************************************/

func (x *Any) GetObject(key string) (any, bool) {
	result, ok := (*x)[key]
	return result, ok
}

func (x *Any) SetObject(element schema.Element, path list.List, value any) error {

	return derp.NewInternalError("mapof.Object.SetObject", "Cannot set object values in mapo.Any")

	/*

		if path.IsEmpty() {
			return derp.NewInternalError("mapof.Object.SetObject", "Cannot set values on empty path")
		}

		x.makeNotNil()

		head, tail := path.Split()

		if tail.IsEmpty() {
			(*x)[head] = value
			return nil
		}

		subElement, ok := element.GetElement(head)

		if !ok {
			return derp.NewInternalError("mapof.Object.SetObject", "Unknown property", head)
		}

		tempValue := (*x)[head]

		if err := schema.SetElement(&tempValue, subElement, tail, value); err != nil {
			return derp.Wrap(err, "mapof.Object.SetObject", "Error setting value", path)
		}

		// Reapply the updated value to the map
		(*x)[head] = tempValue

		return nil
	*/
}

func (x *Any) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

/****************************************
 * Other Getters
 ****************************************/

// GetKeys returns all keys of the underlying map
func (m Any) Keys() []string {
	result := make([]string, len(m))

	index := 0
	for key := range m {
		result[index] = key
		index = index + 1
	}

	return result
}

/******************************************
 * Other Getter Interfaces
 ******************************************/

// GetSliceOfString returns a named option as a slice of strings
func (m Any) GetSliceOfString(name string) []string {
	return convert.SliceOfString(m[name])
}

// GetSliceOfInt returns a named option as a slice of int values
func (m Any) GetSliceOfInt(name string) []int {
	return convert.SliceOfInt(m[name])
}

// GetSliceOfFloat returns a named option as a slice of float64 values
func (m Any) GetSliceOfFloat(name string) []float64 {
	return convert.SliceOfFloat(m[name])
}

// GetMap returns a named option as a mapof.Any
func (m Any) GetMap(name string) Any {

	if value, ok := m[name].(Any); ok {
		return value
	}

	if value, ok := m[name].(map[string]any); ok {
		return Any(value)
	}

	return NewAny()
}

// GetSliceOfMap returns a named option as a slice of mapof.Any objects.
func (m Any) GetSliceOfMap(name string) []Any {
	value := convert.SliceOfMap(m[name])
	result := make([]Any, len(value))

	for index := range value {
		result[index] = Any(value[index])
	}

	return result
}
