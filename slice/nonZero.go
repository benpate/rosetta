package slice

import "github.com/benpate/rosetta/convert"

// NonZero filters out all zero values from a slice
func NonZero[T comparable](original []T) []T {
	return Filter(original, func(item T) bool {
		return !convert.IsZeroValue(item)
	})
}
