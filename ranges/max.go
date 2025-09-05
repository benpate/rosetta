package ranges

import "iter"

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
