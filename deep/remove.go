package deep

import (
	"slices"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
)

// Remove deletes the named properties from every map within a value, changing it in place.
// A value that it cannot search, such as a pointer or a struct, is an error.
func Remove(value any, propertyNames ...string) error {

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
		return recursiveRemoveFromMap(typed, propertyNames...)

	case map[string]mapof.Any:
		return recursiveRemoveFromMap(typed, propertyNames...)

	case map[string]bool:
		removeFromMap(typed, propertyNames...)
		return nil

	case map[string]int:
		removeFromMap(typed, propertyNames...)
		return nil

	case map[string]int8:
		removeFromMap(typed, propertyNames...)
		return nil

	case map[string]int16:
		removeFromMap(typed, propertyNames...)
		return nil

	case map[string]int32:
		removeFromMap(typed, propertyNames...)
		return nil

	case map[string]int64:
		removeFromMap(typed, propertyNames...)
		return nil

	case map[string]float32:
		removeFromMap(typed, propertyNames...)
		return nil

	case map[string]float64:
		removeFromMap(typed, propertyNames...)
		return nil

	case map[string]string:
		removeFromMap(typed, propertyNames...)
		return nil

	case mapof.Any:
		return recursiveRemoveFromMap(typed, propertyNames...)

	case mapof.Float:
		removeFromMap(typed, propertyNames...)
		return nil

	case mapof.Int:
		removeFromMap(typed, propertyNames...)
		return nil

	case mapof.Int64:
		removeFromMap(typed, propertyNames...)
		return nil

	case mapof.Bool:
		removeFromMap(typed, propertyNames...)
		return nil

	case mapof.String:
		removeFromMap(typed, propertyNames...)
		return nil
	}

	// Slices get each item scanned
	switch typed := value.(type) {

	case []any:
		return removeFromSlice(typed, propertyNames...)

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

	case []string:
		return nil

	case sliceof.Any:
		return removeFromSlice(typed, propertyNames...)

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
		return removeFromSlice(typed, propertyNames...)

	case []map[string]string:
		return removeFromSlice(typed, propertyNames...)

	case []mapof.Any:
		return removeFromSlice(typed, propertyNames...)

	case sliceof.MapOfAny:
		return removeFromSlice(typed, propertyNames...)

	case sliceof.MapOfString:
		return removeFromSlice(typed, propertyNames...)
	}

	// Values that know how to remove their own properties
	if typed, isRemover := value.(Remover); isRemover {
		return typed.Remove(propertyNames...)
	}

	// Every other map, slice, and array is searched by reflection
	return removeReflected(value, propertyNames...)
}

// recursiveRemoveFromMap deletes the named properties from a map, then from each value it holds
func recursiveRemoveFromMap[V any](value map[string]V, propertyNames ...string) error {

	const location = "deep.recursiveRemoveFromMap"

	for key, property := range value {

		// Deleting a key while ranging over its map is safe in Go
		if slices.Contains(propertyNames, key) {
			delete(value, key)
			continue
		}

		if err := Remove(property, propertyNames...); err != nil {
			return derp.Wrap(err, location, "Removing from map value")
		}
	}

	// Nothing left to lose
	return nil
}

// removeFromMap deletes the named properties from a map whose values hold no other maps
func removeFromMap[V any](value map[string]V, propertyNames ...string) {

	for _, key := range propertyNames {
		delete(value, key)
	}
}

// removeFromSlice deletes the named properties from every item in a slice
func removeFromSlice[T any](value []T, propertyNames ...string) error {

	const location = "deep.removeFromSlice"

	for index := range value {
		if err := Remove(value[index], propertyNames...); err != nil {
			return derp.Wrap(err, location, "Removing from slice item")
		}
	}

	// Silence is golden
	return nil
}
