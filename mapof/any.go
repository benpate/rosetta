package mapof

import (
	"reflect"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/maps"
	"github.com/benpate/rosetta/pointer"
	"github.com/benpate/rosetta/schema"
)

type Any map[string]any

func NewAny() Any {
	return make(Any)
}

/******************************************
 * Map Manipulations
 ******************************************/

func (x Any) Keys() []string {
	return maps.KeysSorted(x)
}

func (x Any) Equal(value map[string]any) bool {
	return reflect.DeepEqual(x, Any(value))
}

func (x Any) NotEqual(value map[string]any) bool {
	return !reflect.DeepEqual(x, Any(value))
}

func (x Any) IsEmpty() bool {
	return len(x) == 0
}

func (x Any) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter Interfaces
 ******************************************/

func (x Any) GetAny(key string) any {
	result, _ := x.GetAnyOK(key)
	return result
}

func (x Any) GetAnyOK(key string) (any, bool) {
	if value, ok := x[key]; ok {
		return value, true
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

func (x Any) GetTime(key string) time.Time {
	result, _ := x.GetTimeOK(key)
	return result
}

func (x Any) GetTimeOK(key string) (time.Time, bool) {
	if value, ok := x[key]; ok {
		return convert.TimeOk(value, time.Time{})
	}
	return time.Time{}, false
}

/****************************************
 * Setter Interfaces
 ****************************************/

func (x *Any) SetAny(key string, value any) bool {
	x.makeNotNil()
	if compare.IsZero(value) {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

func (x *Any) SetBool(key string, value bool) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

func (x *Any) SetFloat(key string, value float64) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

func (x *Any) SetInt(key string, value int) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

func (x *Any) SetInt64(key string, value int64) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

func (x *Any) SetString(key string, value string) bool {
	x.makeNotNil()
	if value == "" {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

func (x *Any) SetValue(value any) error {

	if mapOfAny, ok := MapOfAny(value); ok {
		*x = mapOfAny
		return nil
	}

	return derp.InternalError("mapof.Any.SetValue", "Cannot convert value to mapof.Any", value)
}

// Append adds a new value to the provided key. If a value already exists for this key
// then it will be forced into a slice of values.
func (x *Any) Append(key string, value any) {
	x.makeNotNil()

	if original, ok := (*x)[key]; ok {

		if list, ok := original.([]any); ok {
			(*x)[key] = append(list, value)
			return
		}

		(*x)[key] = []any{original, value}
		return
	}

	(*x)[key] = []any{value}
}

func (x *Any) makeNotNil() {
	if *x == nil {
		*x = make(Any)
	}
}

/****************************************
 * Tree Traversal
 ****************************************/

func (x Any) GetPointer(key string) (any, bool) {
	result, ok := x[key]
	return result, ok
}

func (x *Any) SetObject(element schema.Element, path list.List, value any) error {

	if path.IsEmpty() {
		return derp.InternalError("mapof.Any.SetObject", "Cannot set values on empty path")
	}

	x.makeNotNil()

	head, tail := path.Split()

	if tail.IsEmpty() {
		(*x)[head] = value
		return nil
	}

	// Fall through means we need to make a child map and set the remaining value in it.
	subElement, ok := element.GetElement(head)

	if !ok {
		return derp.InternalError("mapof.Any.SetObject", "Unknown property", head)
	}

	var tempValue any

	// If we already have a value in this spot, then use it
	if subValue, ok := (*x)[head]; ok {
		tempValue = subValue
	} else {
		// Otherwise, initialize a new mapof.Any
		tempValue = make(Any)
	}

	if err := schema.SetElement(pointer.To(tempValue), subElement, tail, value); err != nil {
		return derp.Wrap(err, "mapof.Any.SetObject", "Error setting value", path)
	}

	// Reapply the updated value to the map
	(*x)[head] = tempValue

	return nil
}

func (x *Any) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

/******************************************
 * Other Getter Interfaces
 ******************************************/

func (m Any) IsZeroValue(name string) bool {
	return compare.IsZero(m[name])
}

// GetSliceofAny returns a named option as a slice of any values
func (m Any) GetSliceOfAny(name string) []any {
	return convert.SliceOfAny(m[name])
}

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

	value := m[name]

	switch typed := value.(type) {

	case Any:
		return []Any{typed}

	case *Any:
		return []Any{*typed}

	case []Any:
		return typed

	case *[]Any:
		return *typed
	}

	// Re-cast the value in a slice of Any objects
	mapOfAny := convert.SliceOfMap(value)
	result := make([]Any, len(mapOfAny))

	for index := range mapOfAny {
		result[index] = Any(mapOfAny[index])
	}

	return result
}

func (m Any) GetSliceOfPlainMap(name string) []map[string]any {

	switch typed := m[name].(type) {

	case []Any:
		result := make([]map[string]any, len(typed))
		for index, value := range typed {
			result[index] = map[string]any(value)
		}
		return result

	case []map[string]any:
		return typed

	default:
		return convert.SliceOfMap(m[name])
	}
}

func (m Any) MapOfAny() map[string]any {
	return m
}

func (m Any) MapOfString() map[string]string {
	return convert.MapOfString(m.MapOfAny())
}
