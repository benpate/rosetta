package mapof

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/loose"
	"github.com/benpate/rosetta/maps"
	"github.com/benpate/rosetta/schema"
)

// LooseTemplate is a map of string keys to arbitrary values, whose strings may be templates that
// Evaluate renders. Read it only through its methods: a stored template is a loose.Template.
type LooseTemplate map[string]any

// NewLooseTemplate returns a new, initialized LooseTemplate map.
func NewLooseTemplate() LooseTemplate {
	return make(LooseTemplate)
}

// ParseLooseTemplate copies a map into a LooseTemplate, compiling every string that holds
// a template.
func ParseLooseTemplate(value map[string]any) LooseTemplate {

	result := make(LooseTemplate, len(value))

	for key, item := range value {
		result[key] = wrapTemplate(item)
	}

	return result
}

// wrapTemplate returns a string that compiles as a template as a loose.Template, and every
// other value, including every other string, unchanged.
func wrapTemplate(value any) any {

	source, isString := value.(string)

	if !isString {
		return value
	}

	if compiled := loose.NewTemplate(source); compiled.IsTemplate() {
		return compiled
	}

	return source
}

// UnmarshalJSON decodes a JSON object, compiling every string that holds a template.
func (x *LooseTemplate) UnmarshalJSON(data []byte) error {

	const location = "mapof.LooseTemplate.UnmarshalJSON"

	var value map[string]any

	if err := json.Unmarshal(data, &value); err != nil {
		return derp.Wrap(err, location, "Invalid JSON")
	}

	*x = ParseLooseTemplate(value)

	// All systems go.
	return nil
}

// Evaluate renders the value at the key against the provided data if it is a template, and
// otherwise returns it as a string.
func (x LooseTemplate) Evaluate(key string, data any) (string, error) {

	const location = "mapof.LooseTemplate.Evaluate"

	switch value := x[key].(type) {

	case loose.Template:
		result, err := value.Execute(data)

		if err != nil {
			return "", derp.Wrap(err, location, "Unable to evaluate template", key)
		}

		return result, nil

	// A map literal written in Go is never parsed, so its templates are compiled here, per call
	case string:
		result, err := loose.NewTemplate(value).Execute(data)

		if err != nil {
			return "", derp.Wrap(err, location, "Unable to evaluate template", key)
		}

		return result, nil

	default:
		return convert.String(value), nil
	}
}

// unwrap returns the value at the key, with a loose.Template replaced by its source string.
func (x LooseTemplate) unwrap(key string) (any, bool) {

	value, ok := x[key]

	if wrapped, isTemplate := value.(loose.Template); isTemplate {
		return wrapped.String(), ok
	}

	return value, ok
}

// Clone returns a deep copy of this map, so that a caller can change it, at any depth, without
// changing the original. Its loose.Template values are shared, because they never change.
func (x LooseTemplate) Clone() LooseTemplate {

	if x == nil {
		return nil
	}

	result := make(LooseTemplate, len(x))

	for key, value := range x {
		result[key] = deepCopy(value)
	}

	return result
}

// deepCopy returns a copy of every map and slice within a value, keeping their types.
// Pointers, structs, and arrays are returned as they are.
func deepCopy(value any) any {

	original := reflect.ValueOf(value)

	switch original.Kind() {

	case reflect.Map:
		return deepCopyMap(original)

	case reflect.Slice:
		return deepCopySlice(original)
	}

	// Everything else is copied by value, or shared
	return value
}

// deepCopyMap returns a copy of a map, with every value deep-copied, keeping a nil map nil
func deepCopyMap(original reflect.Value) any {

	if original.IsNil() {
		return original.Interface()
	}

	elementType := original.Type().Elem()
	result := reflect.MakeMapWithSize(original.Type(), original.Len())
	iterator := original.MapRange()

	for iterator.Next() {
		result.SetMapIndex(iterator.Key(), deepCopyItem(iterator.Value(), elementType))
	}

	return result.Interface()
}

// deepCopySlice returns a copy of a slice, with every item deep-copied, keeping a nil slice nil
func deepCopySlice(original reflect.Value) any {

	if original.IsNil() {
		return original.Interface()
	}

	elementType := original.Type().Elem()
	result := reflect.MakeSlice(original.Type(), original.Len(), original.Len())

	for index := range original.Len() {
		result.Index(index).Set(deepCopyItem(original.Index(index), elementType))
	}

	return result.Interface()
}

// deepCopyItem copies one map entry or slice item, keeping a nil item as the zero value of its
// container's element type
func deepCopyItem(value reflect.Value, elementType reflect.Type) reflect.Value {

	// A zero reflect.Value passed to SetMapIndex deletes the key, so nil must be typed
	copied := deepCopy(value.Interface())

	if copied == nil {
		return reflect.Zero(elementType)
	}

	return reflect.ValueOf(copied)
}

/******************************************
 * Map Manipulations
 ******************************************/

// Length returns the number of elements in the map
func (x LooseTemplate) Length() int {
	return len(x)
}

// Keys returns the map's keys in sorted order.
func (x LooseTemplate) Keys() []string {
	return maps.KeysSorted(x)
}

// IsMap returns TRUE, declaring this type a map for schema traversal. Implements schema.MapTyper.
func (x LooseTemplate) IsMap() bool {
	return true
}

// Equal returns TRUE if this map's original values deeply equal the provided map.
func (x LooseTemplate) Equal(value map[string]any) bool {
	return reflect.DeepEqual(x.MapOfAny(), value)
}

// NotEqual returns TRUE if this map's original values do not deeply equal the provided map.
func (x LooseTemplate) NotEqual(value map[string]any) bool {
	return !x.Equal(value)
}

// IsEmpty returns TRUE if the map contains no elements.
func (x LooseTemplate) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the map contains one or more elements.
func (x LooseTemplate) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter Interfaces
 ******************************************/

// GetAny returns the original value for the key, or nil if the key is not present.
func (x LooseTemplate) GetAny(key string) any {
	result, _ := x.GetAnyOK(key)
	return result
}

// GetAnyOK returns the original value for the key and TRUE if it is present, or (nil, false) if
// not.
func (x LooseTemplate) GetAnyOK(key string) (any, bool) {
	return x.unwrap(key)
}

// GetBool returns the value for the key coerced to bool, or FALSE if absent or not coercible.
func (x LooseTemplate) GetBool(key string) bool {
	result, _ := x.GetBoolOK(key)
	return result
}

// GetBoolOK returns the value for the key coerced to bool, with TRUE if it was present and
// coercible.
func (x LooseTemplate) GetBoolOK(key string) (value bool, ok bool) {
	if original, present := x.unwrap(key); present {
		return convert.BoolOk(original, false)
	}
	return false, false
}

// GetFloat returns the value for the key coerced to float64, or 0 if absent or not coercible.
func (x LooseTemplate) GetFloat(key string) float64 {
	result, _ := x.GetFloatOK(key)
	return result
}

// GetFloatOK returns the value for the key coerced to float64, with TRUE if it was present and
// coercible.
func (x LooseTemplate) GetFloatOK(key string) (value float64, ok bool) {
	if original, present := x.unwrap(key); present {
		return convert.FloatOk(original, 0)
	}
	return 0, false
}

// GetInt returns the value for the key coerced to int, or 0 if absent or not coercible.
func (x LooseTemplate) GetInt(key string) int {
	result, _ := x.GetIntOK(key)
	return result
}

// GetIntOK returns the value for the key coerced to int, with TRUE if it was present and coercible.
func (x LooseTemplate) GetIntOK(key string) (value int, ok bool) {
	if original, present := x.unwrap(key); present {
		return convert.IntOk(original, 0)
	}
	return 0, false
}

// GetInt64 returns the value for the key coerced to int64, or 0 if absent or not coercible.
func (x LooseTemplate) GetInt64(key string) int64 {
	result, _ := x.GetInt64OK(key)
	return result
}

// GetInt64OK returns the value for the key coerced to int64, with TRUE if it was present and
// coercible.
func (x LooseTemplate) GetInt64OK(key string) (value int64, ok bool) {
	if original, present := x.unwrap(key); present {
		return convert.Int64Ok(original, 0)
	}
	return 0, false
}

// GetString returns the value for the key coerced to string, without rendering a template, or ""
// if absent or not coercible.
func (x LooseTemplate) GetString(key string) string {
	result, _ := x.GetStringOK(key)
	return result
}

// GetStringOK returns the value for the key coerced to string, without rendering a template, with
// TRUE if it was present and coercible.
func (x LooseTemplate) GetStringOK(key string) (value string, ok bool) {
	if original, present := x.unwrap(key); present {
		return convert.StringOk(original, "")
	}
	return "", false
}

// GetTime returns the value for the key coerced to time.Time, or the zero time if absent or not
// coercible.
func (x LooseTemplate) GetTime(key string) time.Time {
	result, _ := x.GetTimeOK(key)
	return result
}

// GetTimeOK returns the value for the key coerced to time.Time, with TRUE if it was present and
// coercible.
func (x LooseTemplate) GetTimeOK(key string) (value time.Time, ok bool) {
	if original, present := x.unwrap(key); present {
		return convert.TimeOk(original, time.Time{})
	}
	return time.Time{}, false
}

/****************************************
 * Setter Interfaces
 ****************************************/

// SetAny stores a non-zero value at the key (deleting it when zero), compiling a templated
// string.
func (x *LooseTemplate) SetAny(key string, value any) bool {

	x.makeNotNil()

	if compare.IsZero(value) {
		delete(*x, key)
		return true
	}

	(*x)[key] = wrapTemplate(value)
	return true
}

// SetBool stores a bool value at the key.
func (x *LooseTemplate) SetBool(key string, value bool) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

// SetFloat stores a non-zero float64 at the key (deleting the key when the value is zero).
func (x *LooseTemplate) SetFloat(key string, value float64) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

// SetInt stores a non-zero int at the key (deleting the key when the value is zero).
func (x *LooseTemplate) SetInt(key string, value int) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

// SetInt64 stores a non-zero int64 at the key (deleting the key when the value is zero).
func (x *LooseTemplate) SetInt64(key string, value int64) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

// SetString stores a non-empty string at the key (deleting it when empty), compiling a template.
func (x *LooseTemplate) SetString(key string, value string) bool {

	x.makeNotNil()

	if value == "" {
		delete(*x, key)
		return true
	}

	(*x)[key] = wrapTemplate(value)
	return true
}

// SetValue replaces the entire map with the provided value, if it can be converted to a map,
// compiling every string that holds a template.
func (x *LooseTemplate) SetValue(value any) error {

	const location = "mapof.LooseTemplate.SetValue"

	mapOfAny, ok := MapOfAny(value)

	if !ok {
		return derp.Internal(location, "Cannot convert value to mapof.LooseTemplate", value)
	}

	*x = ParseLooseTemplate(mapOfAny)
	return nil
}

// Append adds a new value to the provided key. If a value already exists for this key
// then it will be forced into a slice of original values, which are never templates.
func (x *LooseTemplate) Append(key string, value any) {

	x.makeNotNil()

	if original, ok := x.unwrap(key); ok {

		if items, ok := original.([]any); ok {
			(*x)[key] = append(items, value)
			return
		}

		(*x)[key] = []any{original, value}
		return
	}

	(*x)[key] = []any{value}
}

// makeNotNil allocates the backing map if the receiver currently points to a nil map.
func (x *LooseTemplate) makeNotNil() {
	if *x == nil {
		*x = make(LooseTemplate)
	}
}

/****************************************
 * Tree Traversal
 ****************************************/

// GetPointer returns the original value for the key (implements the schema PointerGetter
// interface).
func (x LooseTemplate) GetPointer(key string) (any, bool) {
	return x.unwrap(key)
}

// SetObject descends the path (creating child maps as needed) and sets the value, compiling a
// string that holds a template (implements the schema ObjectSetter interface).
func (x *LooseTemplate) SetObject(element schema.Element, path list.List, value any) error {

	const location = "mapof.LooseTemplate.SetObject"

	if path.IsEmpty() {
		return derp.Internal(location, "Cannot set values on empty path")
	}

	x.makeNotNil()

	head, tail := path.Split()

	if tail.IsEmpty() {
		(*x)[head] = wrapTemplate(value)
		return nil
	}

	// Fall through means we need to make a child map and set the remaining value in it.
	subElement, ok := element.GetElement(head)

	if !ok {
		return derp.Internal(location, "Invalid property", head)
	}

	// Child maps are plain Any maps: only this map's own strings are templates
	var subMap Any

	if existing, ok := (*x)[head].(Any); ok {
		subMap = existing
	} else {
		subMap = make(Any)
	}

	if err := schema.SetProperty(subElement, &subMap, tail.String(), value); err != nil {
		return derp.Wrap(err, location, "Setting value", path)
	}

	// Reapply the (mutated) child map back into this map.
	(*x)[head] = subMap

	return nil
}

// Remove deletes the key from the map.
func (x *LooseTemplate) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

/******************************************
 * Other Getter Interfaces
 ******************************************/

// IsZeroValue returns TRUE if the named property is absent or holds a zero value.
func (x LooseTemplate) IsZeroValue(name string) bool {
	return compare.IsZero(x.GetAny(name))
}

// GetSliceOfAny returns a named property as a slice of any values
func (x LooseTemplate) GetSliceOfAny(name string) []any {
	return convert.SliceOfAny(x.GetAny(name))
}

// GetSliceOfString returns a named property as a slice of strings
func (x LooseTemplate) GetSliceOfString(name string) []string {
	return convert.SliceOfString(x.GetAny(name))
}

// GetSliceOfInt returns a named property as a slice of int values
func (x LooseTemplate) GetSliceOfInt(name string) []int {
	return convert.SliceOfInt(x.GetAny(name))
}

// GetSliceOfFloat returns a named property as a slice of float64 values
func (x LooseTemplate) GetSliceOfFloat(name string) []float64 {
	return convert.SliceOfFloat(x.GetAny(name))
}

// GetMap returns a named property as a mapof.Any
func (x LooseTemplate) GetMap(name string) Any {
	return x.GetMapOfAny(name)
}

// GetMapOfAny returns a named property as a mapof.Any
func (x LooseTemplate) GetMapOfAny(name string) Any {
	return Any{name: x.GetAny(name)}.GetMapOfAny(name)
}

// GetMapOfString returns a named property as a mapof.String
func (x LooseTemplate) GetMapOfString(name string) String {
	return Any{name: x.GetAny(name)}.GetMapOfString(name)
}

// GetSliceOfMap returns a named property as a slice of mapof.Any objects.
func (x LooseTemplate) GetSliceOfMap(name string) []Any {
	return Any{name: x.GetAny(name)}.GetSliceOfMap(name)
}

// GetSliceOfPlainMap returns a named property as a slice of plain map[string]any objects.
func (x LooseTemplate) GetSliceOfPlainMap(name string) []map[string]any {
	return Any{name: x.GetAny(name)}.GetSliceOfPlainMap(name)
}

// MapOfAny implements the MapOfAny interface.
// It returns a copy of this map, with every loose.Template replaced by its source string
func (x LooseTemplate) MapOfAny() map[string]any {

	if x == nil {
		return nil
	}

	result := make(map[string]any, len(x))

	for key := range x {
		result[key], _ = x.unwrap(key)
	}

	return result
}

// MapOfString implements the MapOfString interface.
// It returns this map's original values as a map[string]string
func (x LooseTemplate) MapOfString() map[string]string {
	return convert.MapOfString(x.MapOfAny())
}
