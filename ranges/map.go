package ranges

import "iter"

// Map transforms a rangeFunc of type T into a rangeFunc of type U using the provided transform function.
func Map[T any, U any](iterator iter.Seq[T], transform func(T) U) iter.Seq[U] {

	return func(yield func(U) bool) {
		for item := range iterator {
			if mappedItem := transform(item); !yield(mappedItem) {
				return
			}
		}
	}
}
