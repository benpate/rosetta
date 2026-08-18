package sliceof

import (
	"strconv"
)

// sliceIndex returns the index and TRUE if it is non-negative and below every provided
// maximum, or (0, FALSE) if it falls outside those bounds.
func sliceIndex(index int, maximums ...int) (int, bool) {

	if index < 0 {
		return 0, false
	}

	for _, max := range maximums {
		if index >= max {
			return 0, false
		}
	}

	return index, true
}

// sliceStringIndex parses a string key into a bounds-checked slice index.
func sliceStringIndex(key string, maximums ...int) (int, bool) {

	index, err := strconv.Atoi(key)

	if err != nil {
		return 0, false
	}

	return sliceIndex(index, maximums...)
}

// growSlice appends zero values to the slice until the provided index is addressable.
func growSlice[T any, S ~[]T](value *S, length int) {

	for len(*value) <= length {
		var item T
		*value = append(*value, item)
	}
}
