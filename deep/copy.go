package deep

import (
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
)

// Copy returns a deep copy of a value, which shares no map or slice with the original.
// A value that it cannot copy, such as a pointer or a struct, is an error.
func Copy[T any](original T) (T, error) {

	const location = "deep.Copy"

	var zero T

	// Deep-copy the value
	duplicate, err := copyValue(original)

	if err != nil {
		return zero, derp.Wrap(err, location, "Duplicating value")
	}

	// A nil copy (a JSON `null`) matches no type assertion, so it is returned as the zero value
	if duplicate == nil {
		return zero, nil
	}

	// Assert the duplicate is the correct type
	typed, isCorrectType := duplicate.(T)

	if !isCorrectType {
		return zero, derp.Internal(location, "Duplicated value is wrong type", typeName(duplicate))
	}

	// Two of a kind
	return typed, nil
}

// copyValue returns a deep copy of a value of any type, as a value of that same type
func copyValue(value any) (any, error) {

	const location = "deep.copyValue"

	// Simple types are copied by value
	switch typed := value.(type) {

	case nil:
		return nil, nil

	case bool:
		return typed, nil

	case float32:
		return typed, nil

	case float64:
		return typed, nil

	case int:
		return typed, nil

	case int8:
		return typed, nil

	case int16:
		return typed, nil

	case int32:
		return typed, nil

	case int64:
		return typed, nil

	case string:
		return typed, nil

	case time.Time:
		return typed, nil
	}

	// Maps of simple types have nothing deeper to copy
	switch typed := value.(type) {

	case map[string]bool:
		return maps.Clone(typed), nil

	case map[string]int:
		return maps.Clone(typed), nil

	case map[string]int8:
		return maps.Clone(typed), nil

	case map[string]int16:
		return maps.Clone(typed), nil

	case map[string]int32:
		return maps.Clone(typed), nil

	case map[string]int64:
		return maps.Clone(typed), nil

	case map[string]float32:
		return maps.Clone(typed), nil

	case map[string]float64:
		return maps.Clone(typed), nil

	case map[string]string:
		return maps.Clone(typed), nil

	case mapof.Bool:
		return maps.Clone(typed), nil

	case mapof.Float:
		return maps.Clone(typed), nil

	case mapof.Int:
		return maps.Clone(typed), nil

	case mapof.Int64:
		return maps.Clone(typed), nil

	case mapof.String:
		return maps.Clone(typed), nil
	}

	// Maps of other values
	switch typed := value.(type) {

	case map[string]any:
		return copyMap(typed)

	case map[string]mapof.Any:
		return copyMap(typed)

	case mapof.Any:
		return copyMap(typed)
	}

	// Slices of simple types have nothing deeper to copy
	switch typed := value.(type) {

	case []bool:
		return slices.Clone(typed), nil

	case []byte:
		return slices.Clone(typed), nil

	case []float32:
		return slices.Clone(typed), nil

	case []float64:
		return slices.Clone(typed), nil

	case []int:
		return slices.Clone(typed), nil

	case []int8:
		return slices.Clone(typed), nil

	case []int16:
		return slices.Clone(typed), nil

	case []int32:
		return slices.Clone(typed), nil

	case []int64:
		return slices.Clone(typed), nil

	case []string:
		return slices.Clone(typed), nil

	case sliceof.Float:
		return slices.Clone(typed), nil

	case sliceof.Int:
		return slices.Clone(typed), nil

	case sliceof.String:
		return slices.Clone(typed), nil
	}

	// Slices of other values
	switch typed := value.(type) {

	case []any:
		return copySlice(typed)

	case []map[string]any:
		return copySlice(typed)

	case []map[string]string:
		return copySlice(typed)

	case []mapof.Any:
		return copySlice(typed)

	case sliceof.Any:
		return copySlice(typed)

	case sliceof.MapOfAny:
		return copySlice(typed)

	case sliceof.MapOfString:
		return copySlice(typed)
	}

	// Values that know how to copy themselves
	if typed, isCopier := value.(Copier); isCopier {

		result, err := typed.Copy()

		if err != nil {
			return nil, derp.Wrap(err, location, "Copying value", typeName(value))
		}

		return result, nil
	}

	// Every other map, slice, and array is copied by reflection
	return copyReflected(value)
}

// copyMap deep-copies a map, and returns the copy as the map's own type
func copyMap[M ~map[K]V, K comparable, V any](original M) (any, error) {

	const location = "deep.copyMap"

	// A nil map stays nil, as it does in maps.Clone
	if original == nil {
		return original, nil
	}

	result := make(M, len(original))

	// Iterate over each item in the map
	for key, value := range original {

		// Deep-copy the individual value
		duplicate, err := Copy(value)

		if err != nil {
			return nil, derp.Wrap(err, location, "Duplicating map value")
		}

		// Add it to the map
		result[key] = duplicate
	}

	// Success!
	return result, nil
}

// copySlice deep-copies a slice, and returns the copy as the slice's own type
func copySlice[S ~[]T, T any](original S) (any, error) {

	const location = "deep.copySlice"

	// A nil slice stays nil, as it does in slices.Clone
	if original == nil {
		return original, nil
	}

	result := make(S, 0, len(original))

	// Iterate over each item in the slice
	for _, value := range original {

		// Deep-copy the individual value
		duplicate, err := Copy(value)

		if err != nil {
			return nil, derp.Wrap(err, location, "Duplicating slice item")
		}

		// Append the duplicated, typed value to the result
		result = append(result, duplicate)
	}

	// Send in the clones
	return result, nil
}

// typeName names a value's Go type for an error detail
func typeName(value any) string {

	// Errors name the type and never the value, which can hold private data
	return fmt.Sprintf("%T", value)
}
