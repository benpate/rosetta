package channel

// Predicate is a function that reports whether a value should be included by Filter.
type Predicate[T any] func(T) bool

// Filter returns a channel that contains only the items that pass the predicate function.
func Filter[T any](input <-chan T, predicate Predicate[T]) <-chan T {

	// Buffer one result, so the goroutine can run one item ahead of the reader.
	result := make(chan T, 1)

	go func() {
		defer close(result)

		for item := range input {
			if predicate(item) {
				result <- item
			}
		}
	}()

	return result
}
