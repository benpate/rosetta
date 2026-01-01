package ranges

import "iter"

// Empty returns an empty RangeFunc of any type.  It yields no values.
func Empty[T any]() iter.Seq[T] {
	return func(yield func(T) bool) {}
}
