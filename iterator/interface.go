package iterator

// Iterator interface allows callers to iterator over a large number of items in an array/slice
// Deprecated: freeing up this namespace to use for new Go 1.23 range functions
type Iterator interface {

	// Next populates the provided value with the next item, and returns FALSE when the iterator is exhausted
	Next(any) bool

	// Count returns the total number of items available to this iterator
	Count() int
}
