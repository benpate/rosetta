package deep

import (
	"slices"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
)

// Remove deletes the named properties from every map within a value, changing it in place.
// A value that it cannot search, such as a function or a struct that holds a map, is an error.
func Remove(value any, propertyNames ...string) error {

	const location = "deep.Remove"

	// RULE: The error is wrapped here and nowhere deeper, so that its length does not
	// grow with the depth of the value
	if err := removeValue(value, propertyNames, 0); err != nil {
		return derp.Wrap(err, location, "Unable to remove properties")
	}

	return nil
}

// removeValue deletes the named properties from a value at one level of nesting
func removeValue(value any, propertyNames []string, depth int) error {

	const location = "deep.removeValue"

	// RULE: A walk stops at the maximum depth, which also stops a value that contains itself
	if depth > maxDepth {
		return derp.Internal(location, "Max Depth Exceeded", "Check this data structure for cycles")
	}

	// Simple types pass through (no properties to remove)
	switch value.(type) {

	case nil:
		return nil

	case bool:
		return nil

	case float32:
		return nil

	case float64:
		return nil

	case int:
		return nil

	case int8:
		return nil

	case int16:
		return nil

	case int32:
		return nil

	case int64:
		return nil

	case string:
		return nil

	case time.Time:
		return nil
	}

	// Maps get cleansed of propertyNames
	switch typed := value.(type) {

	case map[string]any:
		return recursiveRemoveFromMap(typed, propertyNames, depth)

	case map[string]mapof.Any:
		return recursiveRemoveFromMap(typed, propertyNames, depth)

	case map[string]bool:
		removeFromMap(typed, propertyNames)
		return nil

	case map[string]int:
		removeFromMap(typed, propertyNames)
		return nil

	case map[string]int8:
		removeFromMap(typed, propertyNames)
		return nil

	case map[string]int16:
		removeFromMap(typed, propertyNames)
		return nil

	case map[string]int32:
		removeFromMap(typed, propertyNames)
		return nil

	case map[string]int64:
		removeFromMap(typed, propertyNames)
		return nil

	case map[string]float32:
		removeFromMap(typed, propertyNames)
		return nil

	case map[string]float64:
		removeFromMap(typed, propertyNames)
		return nil

	case map[string]string:
		removeFromMap(typed, propertyNames)
		return nil

	case mapof.Any:
		return recursiveRemoveFromMap(typed, propertyNames, depth)

	case mapof.Float:
		removeFromMap(typed, propertyNames)
		return nil

	case mapof.Int:
		removeFromMap(typed, propertyNames)
		return nil

	case mapof.Int64:
		removeFromMap(typed, propertyNames)
		return nil

	case mapof.Bool:
		removeFromMap(typed, propertyNames)
		return nil

	case mapof.String:
		removeFromMap(typed, propertyNames)
		return nil
	}

	// Slices get each item scanned
	switch typed := value.(type) {

	case []any:
		return removeFromSlice(typed, propertyNames, depth)

	case []bool:
		return nil

	case []float32:
		return nil

	case []float64:
		return nil

	case []int:
		return nil

	case []int8:
		return nil

	case []int16:
		return nil

	case []int32:
		return nil

	case []int64:
		return nil

	case []byte:
		return nil

	case []string:
		return nil

	case sliceof.Any:
		return removeFromSlice(typed, propertyNames, depth)

	case sliceof.Float:
		return nil

	case sliceof.Int:
		return nil

	case sliceof.String:
		return nil
	}

	// Slices of Maps
	switch typed := value.(type) {

	case []map[string]any:
		return removeFromSlice(typed, propertyNames, depth)

	case []map[string]string:
		return removeFromSlice(typed, propertyNames, depth)

	case []mapof.Any:
		return removeFromSlice(typed, propertyNames, depth)

	case sliceof.MapOfAny:
		return removeFromSlice(typed, propertyNames, depth)

	case sliceof.MapOfString:
		return removeFromSlice(typed, propertyNames, depth)

	case *sliceof.Any:

		// A nil pointer leads to nothing
		if typed == nil {
			return nil
		}

		return removeFromSlice(*typed, propertyNames, depth)
	}

	// Values that know how to remove their own properties
	if typed, isRemover := value.(Remover); isRemover {
		return typed.Remove(propertyNames...)
	}

	// Every other map, slice, and array is searched by reflection
	return removeReflected(value, propertyNames, depth)
}

// recursiveRemoveFromMap deletes the named properties from a map, then from each value it holds
func recursiveRemoveFromMap[V any](
	value map[string]V, propertyNames []string, depth int,
) error {

	for key, property := range value {

		// Deleting a key while ranging over its map is safe in Go
		if slices.Contains(propertyNames, key) {
			delete(value, key)
			continue
		}

		if err := removeValue(property, propertyNames, depth+1); err != nil {
			return err
		}
	}

	// Nothing left to lose
	return nil
}

// removeFromMap deletes the named properties from a map whose values hold no other maps
func removeFromMap[V any](value map[string]V, propertyNames []string) {

	for _, key := range propertyNames {
		delete(value, key)
	}
}

// removeFromSlice deletes the named properties from every item in a slice
func removeFromSlice[T any](value []T, propertyNames []string, depth int) error {

	for index := range value {
		if err := removeValue(value[index], propertyNames, depth+1); err != nil {
			return err
		}
	}

	// Silence is golden
	return nil
}
