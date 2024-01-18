package channel

// Slice returns a slice of all items in a channel.
func Slice[T any](channel <-chan T) []T {

	result := make([]T, 0)

	for item := range channel {
		result = append(result, item)
	}

	return result
}

// FromSlice posts every item from a slice to a channel, and then closes the channel.
func FromSlice[T any](slice []T) <-chan T {

	result := make(chan T)

	go func() {
		defer close(result)

		for _, item := range slice {
			result <- item
		}
	}()

	return result
}
