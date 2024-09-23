package maps

import (
	"cmp"
	"slices"
)

// Keys returns all keys of a map
func Keys[T cmp.Ordered, U any](value map[T]U) []T {
	keys := make([]T, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}

	return keys
}

// KeysSorted returns all keys of a map, sorted
func KeysSorted[T cmp.Ordered, U any](value map[T]U) []T {
	keys := Keys(value)
	slices.Sort(keys)

	return keys
}
