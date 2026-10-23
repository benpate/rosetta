package iterator

// Iterator interface allows callers to iterate over a large number of items, one at a time
// Deprecated: freeing up this namespace to use for new Go 1.23 range functions
type Iterator interface {

	// Next populates the provided value with the next item, and returns FALSE when the iterator is exhausted
	Next(any) bool

	// Count returns the total number of items available to this iterator
	Count() int
}
