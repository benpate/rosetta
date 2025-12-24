package ranges

import "iter"

// Limit limits the number of items returned by an iter.Seq to the specified maximum
func Limit[T any](max int, iterator iter.Seq[T]) iter.Seq[T] {

	return func(yield func(T) bool) {
		count := 0

		for item := range iterator {
			if count >= max {
				return
			}

			if !yield(item) {
				return
			}
			count++
		}
	}
}
