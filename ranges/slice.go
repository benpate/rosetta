package ranges

import "iter"

// Slice converts an iter.Seq to a slice of T
func Slice[T any](rangeFunc iter.Seq[T]) []T {

	result := make([]T, 0)

	for value := range rangeFunc {
		result = append(result, value)
	}

	return result
}
