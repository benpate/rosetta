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
		[]map[string]any,
		SliceOfStringer:

		return true
	}

	// Otherwise, use reflection to see what's inside there...
	switch valueOf := reflect.ValueOf(value); valueOf.Kind() {

	// Dereference pointers (if necessary). A nil pointer holds nothing to look inside.
	case reflect.Pointer:

		if valueOf.IsNil() {
			return false
		}

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
	switch valueOf := reflect.ValueOf(value); valueOf.Kind() {

	// RULE: A nil pointer has no length, and cannot be dereferenced
	case reflect.Pointer:

		if valueOf.IsNil() {
			return 0
		}

		return SliceLength(valueOf.Elem().Interface())

	case reflect.Array, reflect.Slice:
		return valueOf.Len()
	}

	return 0
}

// readArrayGetter reads every item out of an ArrayGetter into a []any.
// It returns FALSE if the getter does not hold up its end of the interface.
func readArrayGetter(getter ArrayGetter) ([]any, bool) {

	length := getter.Length()

	// RULE: A negative length is a broken implementation, not an empty array
	if length < 0 {
		return make([]any, 0), false
	}

	result := make([]any, 0, length)

	for index := range length {

		item, ok := getter.GetIndex(index)

		if !ok {
			return make([]any, 0), false
		}

		result = append(result, item)
	}

	return result, true
}
