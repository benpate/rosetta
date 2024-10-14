package convert

import (
	"reflect"
)

// IsSlice returns TRUE if the value is a slice or array (Uses Reflection)
func IsSlice(value any) bool {

	// Get the obvious checks out of the way.
	switch value.(type) {

	case []any,
		[]bool,
		[]int,
		[]int64,
		[]float64,
		[]string,
		[]map[string]any:
		return true
	}

	// Otherwise, use reflection to see what's inside there...
	valueOf := reflect.ValueOf(value)

	switch valueOf.Kind() {

	// Dereference pointers (if necessary)
	case reflect.Pointer:
		return IsSlice(valueOf.Elem().Interface())

	// Arrays and slices are both valid
	case reflect.Array, reflect.Slice:
		return true
	}

	// Otherwise, nah.
	return false
}

// SliceLength returns the length of any slice
func SliceLength(value any) int {

	if value == nil {
		return 0
	}

	// Simple calculations for the common/knonw types
	switch typed := value.(type) {
	case []any:
		return len(typed)
	case []float64:
		return len(typed)
	case []int:
		return len(typed)
	case []int64:
		return len(typed)
	case []string:
		return len(typed)
	case []map[string]any:
		return len(typed)
	}

	// Reflection for unknown types
	valueOf := reflect.ValueOf(value)

	switch valueOf.Kind() {
	case reflect.Pointer:
		return SliceLength(valueOf.Elem().Interface())
	case reflect.Array, reflect.Slice:
		return valueOf.Len()
	}

	return 0
}
