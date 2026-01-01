package ranges

import "iter"

// Filter returns a new iterator that yields only the items from the input iterator
// for which the predicate function returns true.
func Filter[T any](iterator iter.Seq[T], predicate func(T) bool) iter.Seq[T] {

	return func(yield func(T) bool) {
		for item := range iterator {
			if predicate(item) {
				if !yield(item) {
					return
				}
			}
		}
	}
}
