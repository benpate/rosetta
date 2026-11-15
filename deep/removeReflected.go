package deep

import (
	"reflect"

	"github.com/benpate/derp"
)

// removeReflected deletes the matching properties from a value by its kind, for the types that
// Remove does not list. It returns an error for a kind that it cannot search.
func removeReflected(value any, match Matcher, depth int) error {

	const location = "deep.removeReflected"

	target := reflect.ValueOf(value)

	// Types that hold no map, slice, or pointer (such as a database's date and ID types)
	// have no properties
	if isSimpleType(target.Type()) {
		return nil
	}

	switch target.Kind() {

	case reflect.Map:
		return removeReflectedMap(target, match, depth)

	case reflect.Slice, reflect.Array:
		return removeReflectedList(target, match, depth)
	}

	// RULE: Every other pointer, struct, function, or channel is refused, because a
	// property hidden inside one would be left in place without any warning
	message := "Cannot remove properties from unsupported data type"
	return derp.Internal(location, message, typeName(value))
}

// removeReflectedMap deletes the matching properties from a map of any type, then from each
// value it holds
func removeReflectedMap(
	target reflect.Value, match Matcher, depth int,
) error {

	for iterator := target.MapRange(); iterator.Next(); {

		// Deleting a key while ranging over its map is safe, as it is in a `for range`
		if key := iterator.Key(); isMatchingKey(key, match) {
			target.SetMapIndex(key, reflect.Value{})
			continue
		}

		child := iterator.Value().Interface()

		if err := removeValue(child, match, depth+1); err != nil {
			return err
		}
	}

	return nil
}

// removeReflectedList deletes the matching properties from every item in a slice or an array
func removeReflectedList(
	target reflect.Value, match Matcher, depth int,
) error {

	// Simple items have no properties
	if isSimpleType(target.Type().Elem()) {
		return nil
	}

	for index := range target.Len() {
		child := target.Index(index).Interface()

		if err := removeValue(child, match, depth+1); err != nil {
			return err
		}
	}

	return nil
}

// isMatchingKey returns TRUE if a map key is a string that the Matcher selects
func isMatchingKey(key reflect.Value, match Matcher) bool {

	// A key in a map[any]any is read as the value it holds
	if key.Kind() == reflect.Interface {
		key = key.Elem()
	}

	if key.Kind() != reflect.String {
		return false
	}

	return match(key.String())
}
