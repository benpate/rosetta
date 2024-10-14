package convert

import "reflect"

// IsMap returns TRUE if the value is a map (Uses Reflection)
func IsMap(value any) bool {

	// Get the obvious checks out of the way.
	switch value.(type) {
	case map[string]any, map[string]string, map[string][]string:
		return true

	case MapOfAnyGetter:
		return true
	}

	// Otherwise, use reflection to see what's inside there...
	valueOf := reflect.ValueOf(value)

	switch valueOf.Kind() {

	// Dereference pointers (if necessary)
	case reflect.Pointer:
		return IsSlice(valueOf.Elem().Interface())

	// Arrays and slices are both valid
	case reflect.Map:
		return true
	}

	// Otherwise, nah.
	return false
}
