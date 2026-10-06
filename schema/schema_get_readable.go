package schema

import (
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/benpate/derp"
)

// isListPastEnd returns TRUE if the object is a list, and the key reads past its end
func isListPastEnd(object any, key string) bool {

	lister, isList := readableObject(object).(ArrayGetter)
	return isList && isPastEnd(key, lister.Length())
}

// zeroValueAt returns the zero value of the element at the end of a path: "" for a string, 0 for
// a number, FALSE for a boolean, and nil for anything else.  A path the schema does not define is
// an error.
func zeroValueAt(element Element, path string) (any, error) {

	const location = "schema.zeroValueAt"

	leaf, ok := element.GetElement(path)

	if !ok {
		return nil, derp.BadRequest(location, "Invalid property", path)
	}

	switch typed := leaf.(type) {

	case String:
		return "", nil

	case Boolean:
		return false, nil

	case Integer:
		if typed.BitSize == 64 {
			return int64(0), nil
		}
		return 0, nil

	case Number:
		return float64(0), nil
	}

	return nil, nil
}

// readableObject returns an object that the getters can reach: a list or map held by value is
// read through a pointer to a copy, and anything else is returned unchanged
func readableObject(object any) any {

	if _, ok := object.(PointerGetter); ok {
		return object
	}

	return pointerToCopy(object)
}

// scalarValue returns the value a pointer to a scalar points to, and any other value unchanged
func scalarValue(value any) any {

	reflected := reflect.ValueOf(value)

	if (reflected.Kind() != reflect.Pointer) || reflected.IsNil() {
		return value
	}

	if !isScalarKind(reflected.Elem().Kind()) {
		return value
	}

	return reflected.Elem().Interface()
}

// hasEmptySegment returns TRUE if a dotted path has an empty segment, such as "name." or ".name"
func hasEmptySegment(path string) bool {
	return (path != "") && slices.Contains(strings.Split(path, "."), "")
}

// isPastEnd returns TRUE if a key reads past the end of a list of this length: an index at or
// beyond the length (including one too large for an int), or "next"
func isPastEnd(key string, length int) bool {

	if key == "next" {
		return true
	}

	if index, err := strconv.Atoi(key); err == nil {
		return index >= length
	}

	return isDigits(key)
}

// pointerToCopy returns a pointer to a copy of a list or map held by value: its own type when
// that pointer is a PointerGetter, otherwise the rosetta type with its layout.  Any other value,
// including a struct, is returned unchanged.
func pointerToCopy(object any) any {

	if object == nil {
		return object
	}

	// RULE: Only lists and maps. A struct held by value stays unreadable, because Emissary's
	// newStartupStream relies on SetAll refusing a whole model.Content
	original := reflect.ValueOf(object)

	if kind := original.Kind(); (kind != reflect.Slice) && (kind != reflect.Map) {
		return object
	}

	copied := reflect.New(original.Type())
	copied.Elem().Set(original)

	if _, ok := copied.Interface().(PointerGetter); ok {
		return copied.Interface()
	}

	for _, containerType := range convertibleContainers {
		if original.Type().ConvertibleTo(containerType) {
			converted := reflect.New(containerType)
			converted.Elem().Set(original.Convert(containerType))
			return converted.Interface()
		}
	}

	return object
}

// isDigits returns TRUE for a non-empty string of ASCII digits
func isDigits(value string) bool {

	if value == "" {
		return false
	}

	for _, character := range value {
		if (character < '0') || (character > '9') {
			return false
		}
	}

	return true
}

// isScalarKind returns TRUE for a boolean, number, or string kind
func isScalarKind(kind reflect.Kind) bool {

	switch kind {

	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}

	return false
}
