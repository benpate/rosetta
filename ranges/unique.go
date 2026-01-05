package ranges

import "iter"

// Unique returns an iterator that yields only unique items from the provided iterator.
func Unique[T comparable](fn iter.Seq[T]) iter.Seq[T] {

	return func(yield func(T) bool) {
		seen := make(map[T]struct{})

		// Loop over the original iterator
		for item := range fn {

			// If the item has already been seen, then move along
			if _, exists := seen[item]; exists {
				continue
			}

			// Add the item to the list of "seen" items
			seen[item] = struct{}{}

			// Yield the unique item, and quit if requested
			if !yield(item) {
				return
			}
		}
	}
}
