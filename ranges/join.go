package ranges

import "iter"

// Join combines multiple iterators into a single iterator that yields all items from each input iterator in sequence.
func Join[T any](iterators ...iter.Seq[T]) iter.Seq[T] {

	return func(yield func(T) bool) {
		for _, iterator := range iterators {
			for item := range iterator {
				if !yield(item) {
					return
				}
			}
		}
	}
}
