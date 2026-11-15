package deep

import (
	"fmt"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
)

// Copy returns a deep copy of a value, which shares no map or slice with the original
func Copy(value any) (any, error) {

	// Simple types
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
	}

	// Maps
	switch typed := value.(type) {

	case map[string]any:
		return copyMap(typed)

	case map[string]bool:
		return copyMap(typed)

	case map[string]int:
		return copyMap(typed)

	case map[string]int8:
		return copyMap(typed)

	case map[string]int16:
		return copyMap(typed)

	case map[string]int32:
		return copyMap(typed)

	case map[string]int64:
		return copyMap(typed)

	case map[string]float32:
		return copyMap(typed)

	case map[string]float64:
		return copyMap(typed)

	case map[string]string:
		return copyMap(typed)

	case mapof.Any:
		result, err := copyMap(typed)
		return mapof.Any(result), err

	case mapof.Bool:
		result, err := copyMap(typed)
		return mapof.Bool(result), err

	case mapof.Float:
		result, err := copyMap(typed)
		return mapof.Float(result), err

	case mapof.Int:
		result, err := copyMap(typed)
		return mapof.Int(result), err

	case mapof.Int64:
		result, err := copyMap(typed)
		return mapof.Int64(result), err

	case mapof.String:
		result, err := copyMap(typed)
		return mapof.String(result), err

	}

	// Slices
	switch typed := value.(type) {

	case []any:
		return copySlice(typed)

	case []bool:
		return copySlice(typed)

	case []float32:
		return copySlice(typed)

	case []float64:
		return copySlice(typed)

	case []int:
		return copySlice(typed)

	case []int8:
		return copySlice(typed)

	case []int16:
		return copySlice(typed)

	case []int32:
		return copySlice(typed)

	case []int64:
		return copySlice(typed)

	case []string:
		return copySlice(typed)

	case sliceof.Any:
		result, err := copySlice(typed)
		return sliceof.Any(result), err

	case sliceof.Float:
		result, err := copySlice(typed)
		return sliceof.Float(result), err

	case sliceof.Int:
		result, err := copySlice(typed)
		return sliceof.Int(result), err

	case sliceof.String:
		result, err := copySlice(typed)
		return sliceof.String(result), err
	}

	// Slices of Maps
	switch typed := value.(type) {

	case []map[string]any:
		return copySlice(typed)

	case []map[string]string:
		return copySlice(typed)

	case sliceof.MapOfAny:
		result, err := copySlice(typed)
		return sliceof.MapOfAny(result), err

	case sliceof.MapOfString:
		result, err := copySlice(typed)
		return sliceof.MapOfString(result), err
	}

	// Interfaces
	if typed, isCopier := value.(Copier); isCopier {
		return typed.Copy()
	}

	// Otherwise, just no.
	typeName := fmt.Sprintf("%T", value)
	return nil, derp.BadRequest("deep.Copy", "Cannot deep copy unsupported data type", typeName)
}

// copyMap deep-copies a map
func copyMap[K comparable, V any](original map[K]V) (map[K]V, error) {

	const location = "deep.copyMap"

	result := make(map[K]V, len(original))

	// Iterated over each item in the map
	for key, value := range original {

		// Deep-copy the individual value
		duplicate, err := copyAs(value)

		if err != nil {
			return result, derp.Wrap(err, location, "Duplicating map key")
		}

		// Append to the map
		result[key] = duplicate
	}

	// Success!
	return result, nil
}

// copySlice deep-copies a slice
func copySlice[T any](original []T) ([]T, error) {

	const location = "deep.copySlice"

	result := make([]T, 0, len(original))

	for _, value := range original {

		// Deep-copy the individual value
		duplicate, err := copyAs(value)

		if err != nil {
			return result, derp.Wrap(err, location, "Duplicating slice index")
		}

		// Append the duplicated, typed value to the result
		result = append(result, duplicate)
	}

	return result, nil
}

// copyAs deep-copies a value and returns the copy as the original's own type
func copyAs[T any](original T) (T, error) {

	const location = "deep.copyAs"

	var zero T

	// Deep-copy the value
	duplicate, err := Copy(original)

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
		typeName := fmt.Sprintf("%T", duplicate)
		return zero, derp.Internal(location, "Duplicated value is wrong type", typeName)
	}

	// Two of a kind
	return typed, nil
}
