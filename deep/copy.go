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
// A value that it cannot copy, such as a function or a struct that holds a map, is an error.
func Copy[T any](original T) (T, error) {

	const location = "deep.Copy"

	// Start the walk at the top level
	result, err := copyAs(original, 0)

	// RULE: The error is wrapped here and nowhere deeper, so that its length does not
	// grow with the depth of the value
	if err != nil {
		return result, derp.Wrap(err, location, "Unable to copy value")
	}

	// Seeing double
	return result, nil
}

// CopyWithDepth continues a copy that is already under way, at the given level of nesting.
// It is for a Copier, which passes its own depth plus one for each value that it holds.
func CopyWithDepth[T any](original T, depth int) (T, error) {

	// The error is returned as it is, because Copy wraps it at the top of the walk
	return copyAs(original, depth)
}

// copyAs deep-copies a value at one level of nesting, and returns it as the original's type
func copyAs[T any](original T, depth int) (T, error) {

	const location = "deep.copyAs"

	var zero T

	// Deep-copy the value
	duplicate, err := copyValue(original, depth)

	if err != nil {
		return zero, err
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
func copyValue(value any, depth int) (any, error) {

	const location = "deep.copyValue"

	// RULE: A walk stops at the maximum depth, which also stops a value that contains itself
	if depth > maxDepth {
		return nil, derp.Internal(location, "Max Depth Exceeded", "Check value for cycles")
	}

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
		return asAny(copyMap(typed, depth))

	case map[string]mapof.Any:
		return asAny(copyMap(typed, depth))

	case mapof.Any:
		return asAny(copyMap(typed, depth))
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
		return asAny(copySlice(typed, depth))

	case []map[string]any:
		return asAny(copySlice(typed, depth))

	case []map[string]string:
		return asAny(copySlice(typed, depth))

	case []mapof.Any:
		return asAny(copySlice(typed, depth))

	case sliceof.Any:
		return asAny(copySlice(typed, depth))

	case sliceof.MapOfAny:
		return asAny(copySlice(typed, depth))

	case sliceof.MapOfString:
		return asAny(copySlice(typed, depth))

	case *sliceof.Any:
		return asAny(copySlicePointer(typed, depth))
	}

	// Values that know how to copy themselves
	if typed, isCopier := value.(Copier); isCopier {

		// RULE: The Copier is told its depth, so that a cycle which passes through it still ends
		return asAny(typed.Copy(depth))
	}

	// Every other map, slice, and array is copied by reflection
	return copyReflected(value, depth)
}

// typeName names a value's Go type for an error detail
func typeName(value any) string {

	// Errors name the type and never the value, which can hold private data
	return fmt.Sprintf("%T", value)
}

// copyMap deep-copies a map, and returns the copy as the map's own type
func copyMap[M ~map[K]V, K comparable, V any](original M, depth int) (M, error) {

	// A nil map stays nil, as it does in maps.Clone
	if original == nil {
		return original, nil
	}

	// Make a new map of the same type and size
	result := make(M, len(original))

	// Iterate over each item in the map
	for key, value := range original {

		// Deep-copy the individual value
		duplicate, err := copyAs(value, depth+1)

		if err != nil {
			return nil, err
		}

		// Add it to the map
		result[key] = duplicate
	}

	// Success!
	return result, nil
}

// asAny returns a copy as an `any`, and a nil in place of the copy when there is an error
func asAny[T any](value T, err error) (any, error) {

	// RULE: A typed nil stored in an `any` is not nil, so an error returns an untyped nil
	if err != nil {
		return nil, err
	}

	return value, nil
}

// copySlice deep-copies a slice, and returns the copy as the slice's own type
func copySlice[S ~[]T, T any](original S, depth int) (S, error) {

	// A nil slice stays nil, as it does in slices.Clone
	if original == nil {
		return original, nil
	}

	// Make a new slice of the same type, with room for every item
	result := make(S, 0, len(original))

	// Iterate over each item in the slice
	for _, value := range original {

		// Deep-copy the individual value
		duplicate, err := copyAs(value, depth+1)

		if err != nil {
			return nil, err
		}

		// Append the duplicated, typed value to the result
		result = append(result, duplicate)
	}

	// Send in the clones
	return result, nil
}

// copySlicePointer deep-copies the slice that a pointer leads to, and returns a pointer to
// the copy. A nil pointer stays nil.
func copySlicePointer[S ~[]T, T any](original *S, depth int) (*S, error) {

	// A nil pointer leads to nothing
	if original == nil {
		return nil, nil
	}

	// Copy the slice itself, at the pointer's own depth
	result, err := copySlice(*original, depth)

	if err != nil {
		return nil, err
	}

	// Point to the copy
	return &result, nil
}
