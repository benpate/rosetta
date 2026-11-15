package deep

import (
	"reflect"

	"github.com/benpate/derp"
)

// copyReflected deep-copies a value by its kind, for the types that copyValue does not list.
// It keeps the value's own type, and returns an error for a kind that it cannot copy.
func copyReflected(value any, depth int) (any, error) {

	const location = "deep.copyReflected"

	original := reflect.ValueOf(value)

	// Types that hold no map, slice, or pointer (such as a database's date and ID types)
	// were already copied when they were passed by value
	if isSimpleType(original.Type()) {
		return value, nil
	}

	switch original.Kind() {

	case reflect.Map:
		return copyReflectedMap(original, depth)

	case reflect.Slice:
		return copyReflectedSlice(original, depth)

	case reflect.Array:
		return copyReflectedArray(original, depth)
	}

	// RULE: Every other pointer, struct, function, or channel is refused, because a copy
	// that still shared it would not be a deep copy
	return nil, derp.Internal(location, "Cannot deep copy unsupported data type", typeName(value))
}

// copyReflectedMap deep-copies a map of any type, keeping a nil map nil
func copyReflectedMap(original reflect.Value, depth int) (any, error) {

	if original.IsNil() {
		return original.Interface(), nil
	}

	itemType := original.Type().Elem()
	result := reflect.MakeMapWithSize(original.Type(), original.Len())
	iterator := original.MapRange()

	// Iterate over each item in the map
	for iterator.Next() {

		duplicate, err := copyReflectedItem(iterator.Value(), itemType, depth+1)

		if err != nil {
			return nil, err
		}

		result.SetMapIndex(iterator.Key(), duplicate)
	}

	return result.Interface(), nil
}

// copyReflectedSlice deep-copies a slice of any type, keeping a nil slice nil
func copyReflectedSlice(original reflect.Value, depth int) (any, error) {

	if original.IsNil() {
		return original.Interface(), nil
	}

	itemType := original.Type().Elem()
	result := reflect.MakeSlice(original.Type(), original.Len(), original.Len())

	// Simple items have nothing deeper to copy
	if isSimpleType(itemType) {
		reflect.Copy(result, original)
		return result.Interface(), nil
	}

	// Iterate over each item in the slice
	for index := range original.Len() {

		duplicate, err := copyReflectedItem(original.Index(index), itemType, depth+1)

		if err != nil {
			return nil, err
		}

		result.Index(index).Set(duplicate)
	}

	return result.Interface(), nil
}

// copyReflectedArray deep-copies an array of any type
func copyReflectedArray(original reflect.Value, depth int) (any, error) {

	// An array of simple items never arrives here, because copyReflected returns it first
	itemType := original.Type().Elem()
	result := reflect.New(original.Type()).Elem()

	// Iterate over each item in the array
	for index := range original.Len() {

		duplicate, err := copyReflectedItem(original.Index(index), itemType, depth+1)

		if err != nil {
			return nil, err
		}

		result.Index(index).Set(duplicate)
	}

	return result.Interface(), nil
}

// copyReflectedItem deep-copies one item of a container, as a value its container can hold
func copyReflectedItem(
	item reflect.Value, itemType reflect.Type, depth int,
) (reflect.Value, error) {

	const location = "deep.copyReflectedItem"

	duplicate, err := copyValue(item.Interface(), depth)

	if err != nil {
		return reflect.Value{}, err
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

// copierType and removerType are the Copier and Remover interfaces, as reflection names them
var (
	copierType  = reflect.TypeFor[Copier]()
	removerType = reflect.TypeFor[Remover]()
)

// isSimpleType returns TRUE for a type that holds no map, slice, pointer, or interface at any
// level: a boolean, number, or string, and an array or struct made only of simple types
func isSimpleType(valueType reflect.Type) bool {

	// RULE: A type that copies or searches itself is never skipped as simple, so that its
	// own method is called wherever it is held
	if valueType.Implements(copierType) || valueType.Implements(removerType) {
		return false
	}

	switch valueType.Kind() {

	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true

	case reflect.Array:
		return isSimpleType(valueType.Elem())

	case reflect.Struct:

		for index := range valueType.NumField() {
			if !isSimpleType(valueType.Field(index).Type) {
				return false
			}
		}

		return true
	}

	return false
}
