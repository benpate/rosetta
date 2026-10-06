package mapof

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/compare"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/loose"
	"github.com/benpate/rosetta/maps"
)

// Template is a map of string keys to arbitrary values, whose strings may be templates that
// its getters render against data. Read it only through its methods.
type Template map[string]any

// NewTemplate returns a new, initialized Template map.
func NewTemplate() Template {
	return make(Template)
}

// ParseTemplate copies a map into a Template, compiling every string that holds
// a template.
func ParseTemplate(value map[string]any) Template {

	result := make(Template, len(value))

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
func (x *Template) UnmarshalJSON(data []byte) error {

	const location = "mapof.Template.UnmarshalJSON"

	var value map[string]any

	if err := json.Unmarshal(data, &value); err != nil {
		return derp.Wrap(err, location, "Invalid JSON")
	}

	*x = ParseTemplate(value)

	// All systems go.
	return nil
}

// Evaluate renders the value at the key against the provided data if it is a template, and
// otherwise returns it as a string. Unlike the getters, it reports a failed render.
func (x Template) Evaluate(key string, data any) (string, error) {

	const location = "mapof.Template.Evaluate"

	switch value := x[key].(type) {

	case loose.Template:
		result, err := value.Execute(data)

		if err != nil {
			return "", derp.Wrap(err, location, "Unable to evaluate template", key)
		}

		return result, nil

	// A map literal written in Go is never parsed, so its templates are compiled here, once, into a
	// cache beside the map. Writing them back into the map would race with its other readers.
	case string:
		result, err := loose.CachedTemplate(value).Execute(data)

		if err != nil {
			return "", derp.Wrap(err, location, "Unable to evaluate template", key)
		}

		return result, nil

	default:
		return convert.String(value), nil
	}
}

// unwrap returns the value at the key, with a loose.Template replaced by its source string.
func (x Template) unwrap(key string) (any, bool) {

	value, ok := x[key]

	if wrapped, isTemplate := value.(loose.Template); isTemplate {
		return wrapped.String(), ok
	}

	return value, ok
}

// Clone returns a deep copy of this map, so that a caller can change it, at any depth, without
// changing the original. Its loose.Template values are shared, because they never change.
func (x Template) Clone() Template {

	if x == nil {
		return nil
	}

	result := make(Template, len(x))

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
func (x Template) Length() int {
	return len(x)
}

// Keys returns the map's keys in sorted order.
func (x Template) Keys() []string {
	return maps.KeysSorted(x)
}

// IsMap returns TRUE, declaring this type a map for schema traversal. Implements schema.MapTyper.
func (x Template) IsMap() bool {
	return true
}

// Equal returns TRUE if this map's original values deeply equal the provided map.
func (x Template) Equal(value map[string]any) bool {
	return reflect.DeepEqual(x.MapOfAny(), value)
}

// NotEqual returns TRUE if this map's original values do not deeply equal the provided map.
func (x Template) NotEqual(value map[string]any) bool {
	return !x.Equal(value)
}

// IsEmpty returns TRUE if the map contains no elements.
func (x Template) IsEmpty() bool {
	return len(x) == 0
}

// NotEmpty returns TRUE if the map contains one or more elements.
func (x Template) NotEmpty() bool {
	return len(x) > 0
}

/******************************************
 * Getter Interfaces
 ******************************************/

// GetAny returns the value for the key, with a template rendered against the data, or nil if the
// key is not present.
func (x Template) GetAny(key string, data any) any {
	result, _ := x.GetAnyOK(key, data)
	return result
}

// GetAnyOK returns the value for the key, with a template rendered against the data, and TRUE if
// it is present, or (nil, false) if not.
func (x Template) GetAnyOK(key string, data any) (any, bool) {
	return x.evaluateAny(key, data)
}

// GetBool returns the value for the key, rendered against the data and coerced to bool, or FALSE if
// absent or not coercible.
func (x Template) GetBool(key string, data any) bool {
	result, _ := x.GetBoolOK(key, data)
	return result
}

// GetBoolOK returns the value for the key, rendered against the data and coerced to bool, with TRUE
// if it was present and coercible.
func (x Template) GetBoolOK(key string, data any) (value bool, ok bool) {
	if result, present := x.evaluateScalar(key, data); present {
		return convert.BoolOk(result, false)
	}
	return false, false
}

// GetFloat returns the value for the key, rendered against the data and coerced to float64, or 0 if
// absent or not coercible.
func (x Template) GetFloat(key string, data any) float64 {
	result, _ := x.GetFloatOK(key, data)
	return result
}

// GetFloatOK returns the value for the key, rendered against the data and coerced to float64, with
// TRUE if it was present and coercible.
func (x Template) GetFloatOK(key string, data any) (value float64, ok bool) {
	if result, present := x.evaluateScalar(key, data); present {
		return convert.FloatOk(result, 0)
	}
	return 0, false
}

// GetInt returns the value for the key, rendered against the data and coerced to int, or 0 if
// absent or not coercible.
func (x Template) GetInt(key string, data any) int {
	result, _ := x.GetIntOK(key, data)
	return result
}

// GetIntOK returns the value for the key, rendered against the data and coerced to int, with TRUE
// if it was present and coercible.
func (x Template) GetIntOK(key string, data any) (value int, ok bool) {
	if result, present := x.evaluateScalar(key, data); present {
		return convert.IntOk(result, 0)
	}
	return 0, false
}

// GetInt64 returns the value for the key, rendered against the data and coerced to int64, or 0 if
// absent or not coercible.
func (x Template) GetInt64(key string, data any) int64 {
	result, _ := x.GetInt64OK(key, data)
	return result
}

// GetInt64OK returns the value for the key, rendered against the data and coerced to int64, with
// TRUE if it was present and coercible.
func (x Template) GetInt64OK(key string, data any) (value int64, ok bool) {
	if result, present := x.evaluateScalar(key, data); present {
		return convert.Int64Ok(result, 0)
	}
	return 0, false
}

// GetString returns the value for the key, rendered against the data and coerced to string, or ""
// if absent or not coercible.
func (x Template) GetString(key string, data any) string {
	result, _ := x.GetStringOK(key, data)
	return result
}

// GetStringOK returns the value for the key, rendered against the data and coerced to string, with
// TRUE if it was present and coercible.
func (x Template) GetStringOK(key string, data any) (value string, ok bool) {
	if result, present := x.evaluateAny(key, data); present {
		return convert.StringOk(result, "")
	}
	return "", false
}

// GetTime returns the value for the key, rendered against the data and coerced to time.Time, or the
// zero time if absent or not coercible.
func (x Template) GetTime(key string, data any) time.Time {
	result, _ := x.GetTimeOK(key, data)
	return result
}

// GetTimeOK returns the value for the key, rendered against the data and coerced to time.Time, with
// TRUE if it was present and coercible.
func (x Template) GetTimeOK(key string, data any) (value time.Time, ok bool) {
	if result, present := x.evaluateScalar(key, data); present {
		return convert.TimeOk(result, time.Time{})
	}
	return time.Time{}, false
}

// evaluateScalar returns the value for the key as evaluateAny does, trimming the whitespace from a
// rendered template so that it converts to a bool, number, or time.
func (x Template) evaluateScalar(key string, data any) (any, bool) {

	result, present := x.evaluateAny(key, data)

	if text, isText := result.(string); isText && (data != nil) {
		return strings.TrimSpace(text), present
	}

	return result, present
}

// evaluateAny returns the value for the key, with every string rendered against the data as a
// template. Nil data returns the source unchanged, and a failed render returns "".
func (x Template) evaluateAny(key string, data any) (any, bool) {

	value, present := x[key]

	if !present {
		return nil, false
	}

	var compiled loose.Template

	switch typed := value.(type) {

	case loose.Template:
		compiled = typed

	// A map literal written in Go is never parsed, so its templates come from the shared cache
	case string:
		compiled = loose.CachedTemplate(typed)

	default:
		return value, true
	}

	// Without data there is nothing to render, so return the source
	if data == nil {
		return compiled.String(), true
	}

	// Errors are swallowed on purpose: a getter has no error to return (see Evaluate)
	result, err := compiled.Execute(data)

	if err != nil {
		return "", true
	}

	return result, true
}

/****************************************
 * Setter Interfaces
 ****************************************/

// SetAny stores a non-zero value at the key (deleting it when zero), compiling a templated
// string.
func (x *Template) SetAny(key string, value any) bool {

	x.makeNotNil()

	if compare.IsZero(value) {
		delete(*x, key)
		return true
	}

	(*x)[key] = wrapTemplate(value)
	return true
}

// SetBool stores a bool value at the key.
func (x *Template) SetBool(key string, value bool) bool {
	x.makeNotNil()
	(*x)[key] = value
	return true
}

// SetFloat stores a non-zero float64 at the key (deleting the key when the value is zero).
func (x *Template) SetFloat(key string, value float64) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

// SetInt stores a non-zero int at the key (deleting the key when the value is zero).
func (x *Template) SetInt(key string, value int) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

// SetInt64 stores a non-zero int64 at the key (deleting the key when the value is zero).
func (x *Template) SetInt64(key string, value int64) bool {
	x.makeNotNil()
	if value == 0 {
		delete(*x, key)
	} else {
		(*x)[key] = value
	}
	return true
}

// SetString stores a non-empty string at the key (deleting it when empty), compiling a template.
func (x *Template) SetString(key string, value string) bool {

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
func (x *Template) SetValue(value any) error {

	const location = "mapof.Template.SetValue"

	mapOfAny, ok := MapOfAny(value)

	if !ok {
		return derp.Internal(location, "Cannot convert value to mapof.Template", value)
	}

	*x = ParseTemplate(mapOfAny)
	return nil
}

// Append adds a new value to the provided key. If a value already exists for this key
// then it will be forced into a slice of original values, which are never templates.
func (x *Template) Append(key string, value any) {

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
func (x *Template) makeNotNil() {
	if *x == nil {
		*x = make(Template)
	}
}

/****************************************
 * Tree Traversal
 ****************************************/

// GetPointer returns the original value for the key (implements the schema PointerGetter
// interface).
func (x Template) GetPointer(key string) (any, bool) {
	return x.unwrap(key)
}

// SetKey stores the value under the key, including a zero value, compiling a string that holds
// a template (implements the schema KeySetter interface).
func (x *Template) SetKey(key string, value any) error {
	x.makeNotNil()
	(*x)[key] = wrapTemplate(value)
	return nil
}

// Remove deletes the key from the map.
func (x *Template) Remove(key string) bool {
	x.makeNotNil()
	delete(*x, key)
	return true
}

/******************************************
 * Other Getter Interfaces
 ******************************************/

// IsZeroValue returns TRUE if the named property is absent or holds a zero value, without
// rendering it.
func (x Template) IsZeroValue(name string) bool {
	return compare.IsZero(x.GetAny(name, nil))
}

// GetSliceOfAny returns a named property, with a template rendered against the data, as a slice of
// any values
func (x Template) GetSliceOfAny(name string, data any) []any {
	return convert.SliceOfAny(x.GetAny(name, data))
}

// GetSliceOfString returns a named property, with a template rendered against the data, as a slice
// of strings
func (x Template) GetSliceOfString(name string, data any) []string {
	return convert.SliceOfString(x.GetAny(name, data))
}

// GetSliceOfInt returns a named property, with a template rendered against the data, as a slice of
// int values
func (x Template) GetSliceOfInt(name string, data any) []int {
	return convert.SliceOfInt(x.GetAny(name, data))
}

// GetSliceOfFloat returns a named property, with a template rendered against the data, as a slice
// of float64 values
func (x Template) GetSliceOfFloat(name string, data any) []float64 {
	return convert.SliceOfFloat(x.GetAny(name, data))
}

// GetMap returns a named property as a mapof.Any. Strings within it are not templates.
func (x Template) GetMap(name string, data any) Any {
	return x.GetMapOfAny(name, data)
}

// GetMapOfAny returns a named property as a mapof.Any. Strings within it are not templates.
func (x Template) GetMapOfAny(name string, data any) Any {
	return Any{name: x.GetAny(name, data)}.GetMapOfAny(name)
}

// GetMapOfString returns a named property as a mapof.String. Strings within it are not templates.
func (x Template) GetMapOfString(name string, data any) String {
	return Any{name: x.GetAny(name, data)}.GetMapOfString(name)
}

// GetSliceOfMap returns a named property as a slice of mapof.Any objects. Strings within them are
// not templates.
func (x Template) GetSliceOfMap(name string, data any) []Any {
	return Any{name: x.GetAny(name, data)}.GetSliceOfMap(name)
}

// GetSliceOfPlainMap returns a named property as a slice of plain map[string]any objects. Strings
// within them are not templates.
func (x Template) GetSliceOfPlainMap(name string, data any) []map[string]any {
	return Any{name: x.GetAny(name, data)}.GetSliceOfPlainMap(name)
}

// MapOfAny implements the MapOfAny interface.
// It returns a copy of this map, with every loose.Template replaced by its source string
func (x Template) MapOfAny() map[string]any {

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
func (x Template) MapOfString() map[string]string {
	return convert.MapOfString(x.MapOfAny())
}
