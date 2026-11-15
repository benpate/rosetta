package deep

import (
	"slices"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
)

// Remove deletes the provided propertyName(s) from the
// data structure
func Remove(value any, propertyNames ...string) error {

	// Simple types pass through (no properties to remove)
	switch value.(type) {

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
	}

	// Maps get cleansed of propertyNames
	switch typed := value.(type) {

	case map[string]any:
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

	case sliceof.MapOfAny:
		return removeFromSlice(typed, propertyNames...)

	case sliceof.MapOfString:
		return removeFromSlice(typed, propertyNames...)
	}

	// Interfaces
	if typed, isRemover := value.(Remover); isRemover {
		return typed.Remove(propertyNames...)
	}

	// Otherwise, nothing to do here.
	return nil
}

func recursiveRemoveFromMap[V any](value map[string]V, propertyNames ...string) error {

	const location = "deep.removeFromMap"

	for key, property := range value {

		if slices.Contains(propertyNames, key) {
			delete(value, key)
			continue
		}

		if err := Remove(property, propertyNames...); err != nil {
			return derp.Wrap(err, location, "removing from ")
		}

		// Set the (updated?) property back into this map
		value[key] = property
	}

	return nil
}

func removeFromMap[K comparable, V any](value map[K]V, propertyNames ...K) {

	const location = "deep.removeFromMap"

	for _, key := range propertyNames {
		delete(value, key)
	}
}

func removeFromSlice[T any](value []T, propertyNames ...string) error {

	const location = "deep.removeFromSlice"

	for index := range value {
		if err := Remove(value[index], propertyNames...); err != nil {
			return derp.Wrap(err, location, "Removing property")
		}
	}

	return nil
}
