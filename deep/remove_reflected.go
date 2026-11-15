package deep

import (
	"reflect"
	"slices"

	"github.com/benpate/derp"
)

// removeReflected deletes the named properties from a value by its kind, for the types that
// Remove does not list. It returns an error for a kind that it cannot search.
func removeReflected(value any, propertyNames ...string) error {

	const location = "deep.removeReflected"

	target := reflect.ValueOf(value)
	kind := target.Kind()

	// Named simple types (such as a database's date and ID types) have no properties
	if isSimpleKind(kind) {
		return nil
	}

	switch kind {

	case reflect.Map:
		return removeReflectedMap(target, propertyNames...)

	case reflect.Slice, reflect.Array:
		return removeReflectedList(target, propertyNames...)
	}

	// RULE: A pointer, struct, function, or channel is refused, because a property hidden
	// inside one would be left in place without any warning
	message := "Cannot remove properties from unsupported data type"
	return derp.BadRequest(location, message, typeName(value))
}

// removeReflectedMap deletes the named properties from a map of any type, then from each
// value it holds
func removeReflectedMap(target reflect.Value, propertyNames ...string) error {

	const location = "deep.removeReflectedMap"

	iterator := target.MapRange()

	for iterator.Next() {

		key := iterator.Key()

		// Deleting a key while ranging over its map is safe, as it is in a `for range`
		if isNamedKey(key, propertyNames) {
			target.SetMapIndex(key, reflect.Value{})
			continue
		}

		if err := Remove(iterator.Value().Interface(), propertyNames...); err != nil {
			return derp.Wrap(err, location, "Removing from map value")
		}
	}

	return nil
}

// removeReflectedList deletes the named properties from every item in a slice or an array
func removeReflectedList(target reflect.Value, propertyNames ...string) error {

	const location = "deep.removeReflectedList"

	// Simple items have no properties
	if isSimpleKind(target.Type().Elem().Kind()) {
		return nil
	}

	for index := range target.Len() {
		if err := Remove(target.Index(index).Interface(), propertyNames...); err != nil {
			return derp.Wrap(err, location, "Removing from list item")
		}
	}

	return nil
}

// isNamedKey returns TRUE if a map key is a string that matches one of the property names
func isNamedKey(key reflect.Value, propertyNames []string) bool {

	// A key in a map[any]any is read as the value it holds
	if key.Kind() == reflect.Interface {
		key = key.Elem()
	}

	if key.Kind() != reflect.String {
		return false
	}

	return slices.Contains(propertyNames, key.String())
}
