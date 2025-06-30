package slice

import "iter"

// Range returns an iterator that yields each value in a slice.
func Range[V any](value []V) iter.Seq2[int, V] {

	return func(yield func(int, V) bool) {
		for index, item := range value {
			if !yield(index, item) {
				return
			}
		}
	}
}
