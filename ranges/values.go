package ranges

import "iter"

// Values returns a new iterator that yields the provided values
func Values[T any](values ...T) iter.Seq[T] {

	return func(yield func(T) bool) {
		for _, value := range values {
			if !yield(value) {
				return
			}
		}
	}
}
