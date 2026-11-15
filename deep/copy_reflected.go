package deep

import (
	"reflect"

	"github.com/benpate/derp"
)

// copyReflected deep-copies a value by its kind, for the types that copyValue does not list.
// It keeps the value's own type, and returns an error for a kind that it cannot copy.
func copyReflected(value any) (any, error) {

	const location = "deep.copyReflected"

	original := reflect.ValueOf(value)
	kind := original.Kind()

	// Named simple types (such as a database's date and ID types) are copied by value
	if isSimpleKind(kind) {
		return value, nil
	}

	switch kind {

	case reflect.Map:
		return copyReflectedMap(original)

	case reflect.Slice:
		return copyReflectedSlice(original)

	case reflect.Array:
		return copyReflectedArray(original)
	}

	// RULE: A pointer, struct, function, or channel is refused, because a copy that still
	// shared it would not be a deep copy
	return nil, derp.BadRequest(location, "Cannot deep copy unsupported data type", typeName(value))
}

// copyReflectedMap deep-copies a map of any type, keeping a nil map nil
func copyReflectedMap(original reflect.Value) (any, error) {

	const location = "deep.copyReflectedMap"

	if original.IsNil() {
		return original.Interface(), nil
	}

	itemType := original.Type().Elem()
	result := reflect.MakeMapWithSize(original.Type(), original.Len())
	iterator := original.MapRange()

	// Iterate over each item in the map
	for iterator.Next() {

		duplicate, err := copyReflectedItem(iterator.Value(), itemType)

		if err != nil {
			return nil, derp.Wrap(err, location, "Duplicating map value")
		}

		result.SetMapIndex(iterator.Key(), duplicate)
	}

	return result.Interface(), nil
}

// copyReflectedSlice deep-copies a slice of any type, keeping a nil slice nil
func copyReflectedSlice(original reflect.Value) (any, error) {

	const location = "deep.copyReflectedSlice"

	if original.IsNil() {
		return original.Interface(), nil
	}

	itemType := original.Type().Elem()
	result := reflect.MakeSlice(original.Type(), original.Len(), original.Len())

	// Simple items have nothing deeper to copy
	if isSimpleKind(itemType.Kind()) {
		reflect.Copy(result, original)
		return result.Interface(), nil
	}

	// Iterate over each item in the slice
	for index := range original.Len() {

		duplicate, err := copyReflectedItem(original.Index(index), itemType)

		if err != nil {
			return nil, derp.Wrap(err, location, "Duplicating slice item")
		}

		result.Index(index).Set(duplicate)
	}

	return result.Interface(), nil
}

// copyReflectedArray deep-copies an array of any type
func copyReflectedArray(original reflect.Value) (any, error) {

	const location = "deep.copyReflectedArray"

	itemType := original.Type().Elem()

	// An array of simple items was already copied when it was passed by value
	if isSimpleKind(itemType.Kind()) {
		return original.Interface(), nil
	}

	result := reflect.New(original.Type()).Elem()

	// Iterate over each item in the array
	for index := range original.Len() {

		duplicate, err := copyReflectedItem(original.Index(index), itemType)

		if err != nil {
			return nil, derp.Wrap(err, location, "Duplicating array item")
		}

		result.Index(index).Set(duplicate)
	}

	return result.Interface(), nil
}

// copyReflectedItem deep-copies one item of a container, as a value its container can hold
func copyReflectedItem(item reflect.Value, itemType reflect.Type) (reflect.Value, error) {

	const location = "deep.copyReflectedItem"

	duplicate, err := copyValue(item.Interface())

	if err != nil {
		return reflect.Value{}, derp.Wrap(err, location, "Duplicating value")
	}

	// A nil copy (a JSON `null`) is stored as the container's zero value
	if duplicate == nil {
		return reflect.Zero(itemType), nil
	}

	result := reflect.ValueOf(duplicate)

	// RULE: A Copier can return any type, and storing the wrong one would panic
	if !result.Type().AssignableTo(itemType) {
		message := "Duplicated value is wrong type"
		return reflect.Value{}, derp.Internal(location, message, typeName(duplicate))
	}

	return result, nil
}

// isSimpleKind returns TRUE for the kinds that hold no map, slice, or pointer: booleans,
// numbers, and strings
func isSimpleKind(kind reflect.Kind) bool {

	switch kind {

	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	}

	return false
}
