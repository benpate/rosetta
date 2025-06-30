package iterator

// Iterator interface allows callers to iterator over a large number of items in an array/slice
// Deprecated: freeing up this namespace to use for new Go 1.23 range functions
type Iterator interface {
	Next(any) bool
	Count() int
}
