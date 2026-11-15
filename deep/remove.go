package deep

import (
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
)

// Remove deletes the named properties from every map within a value, changing it in place.
// A value that it cannot search, such as a function or a struct that holds a map, is an error.
func Remove(value any, propertyNames ...string) error {

	const location = "deep.Remove"

	// The names are not copied, because nothing keeps the Matcher once this call returns
	match := matchNames(propertyNames)

	// RULE: The error is wrapped here and nowhere deeper, so that its length does not
	// grow with the depth of the value
	if err := removeValue(value, match, 0); err != nil {
		return derp.Wrap(err, location, "Unable to remove properties")
	}

	// Gone, baby, gone
	return nil
}

// RemoveFunc deletes every property that the Matcher selects from every map within a value,
// changing it in place. It returns the same errors as Remove, and one for a nil Matcher.
func RemoveFunc(value any, match Matcher) error {

	const location = "deep.RemoveFunc"

	// RULE: A nil Matcher is refused, because removing nothing would look like success
	if match == nil {
		return derp.Internal(location, "Matcher must not be nil")
	}

	// RULE: The error is wrapped here and nowhere deeper, so that its length does not
	// grow with the depth of the value
	if err := removeValue(value, match, 0); err != nil {
		return derp.Wrap(err, location, "Unable to remove properties")
	}

	// Poof
	return nil
}

// RemoveWithDepth continues a removal that is already under way, at the given level of
// nesting. It is for a Remover, which passes its own depth plus one for each value it holds.
func RemoveWithDepth(value any, depth int, match Matcher) error {

	const location = "deep.RemoveWithDepth"

	// RULE: A nil Matcher is refused, because removing nothing would look like success
	if match == nil {
		return derp.Internal(location, "Matcher must not be nil")
	}

	// The error is returned as it is, because Remove wraps it at the top of the walk
	return removeValue(value, match, depth)
}

// removeValue deletes the matching properties from a value at one level of nesting
func removeValue(value any, match Matcher, depth int) error {

	const location = "deep.removeValue"

	// RULE: A walk stops at the maximum depth, which also stops a value that contains itself
	if depth > maxDepth {
		return derp.Internal(location, "Max Depth Exceeded", "Check value for cycles")
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

	// Maps get cleansed of matching properties
	switch typed := value.(type) {

	case map[string]any:
		return recursiveRemoveFromMap(typed, match, depth)

	case map[string]mapof.Any:
		return recursiveRemoveFromMap(typed, match, depth)

	case map[string]bool:
		removeFromMap(typed, match)
		return nil

	case map[string]int:
		removeFromMap(typed, match)
		return nil

	case map[string]int8:
		removeFromMap(typed, match)
		return nil

	case map[string]int16:
		removeFromMap(typed, match)
		return nil

	case map[string]int32:
		removeFromMap(typed, match)
		return nil

	case map[string]int64:
		removeFromMap(typed, match)
		return nil

	case map[string]float32:
		removeFromMap(typed, match)
		return nil

	case map[string]float64:
		removeFromMap(typed, match)
		return nil

	case map[string]string:
		removeFromMap(typed, match)
		return nil

	case mapof.Any:
		return recursiveRemoveFromMap(typed, match, depth)

	case mapof.Float:
		removeFromMap(typed, match)
		return nil

	case mapof.Int:
		removeFromMap(typed, match)
		return nil

	case mapof.Int64:
		removeFromMap(typed, match)
		return nil

	case mapof.Bool:
		removeFromMap(typed, match)
		return nil

	case mapof.String:
		removeFromMap(typed, match)
		return nil
	}

	// Slices get each item scanned
	switch typed := value.(type) {

	case []any:
		return removeFromSlice(typed, match, depth)

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
		return removeFromSlice(typed, match, depth)

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
		return removeFromSlice(typed, match, depth)

	case []map[string]string:
		return removeFromSlice(typed, match, depth)

	case []mapof.Any:
		return removeFromSlice(typed, match, depth)

	case sliceof.MapOfAny:
		return removeFromSlice(typed, match, depth)

	case sliceof.MapOfString:
		return removeFromSlice(typed, match, depth)

	case *sliceof.Any:

		// A nil pointer leads to nothing
		if typed == nil {
			return nil
		}

		return removeFromSlice(*typed, match, depth)
	}

	// Values that know how to remove their own properties
	if typed, isRemover := value.(Remover); isRemover {

		// RULE: The Remover is told its depth, so that a cycle which passes through it still ends
		return typed.Remove(depth, match)
	}

	// Every other map, slice, and array is searched by reflection
	return removeReflected(value, match, depth)
}

// recursiveRemoveFromMap deletes the matching properties from a map, then from each value
// it holds
func recursiveRemoveFromMap[V any](value map[string]V, match Matcher, depth int) error {

	for key, property := range value {

		// Deleting a key while ranging over its map is safe in Go
		if match(key) {
			delete(value, key)
			continue
		}

		// Search the values that remain, one level deeper
		if err := removeValue(property, match, depth+1); err != nil {
			return err
		}
	}

	// Nothing left to lose
	return nil
}

// removeFromMap deletes the matching properties from a map whose values hold no other maps
func removeFromMap[V any](value map[string]V, match Matcher) {

	// Deleting a key while ranging over its map is safe in Go
	for key := range value {
		if match(key) {
			delete(value, key)
		}
	}
}

// removeFromSlice deletes the matching properties from every item in a slice
func removeFromSlice[T any](value []T, match Matcher, depth int) error {

	for index := range value {
		if err := removeValue(value[index], match, depth+1); err != nil {
			return err
		}
	}

	// Silence is golden
	return nil
}
