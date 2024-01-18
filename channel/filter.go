package channel

type Predicate[T any] func(T) bool

// Filter returns a channel that contains only the items that pass the predicate function.
func Filter[T any](input <-chan T, predicate Predicate[T]) <-chan T {

	// Create a buffered channel for results to prevent blocking.
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
