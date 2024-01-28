package channel

// Limit returns a channel that will receive at most the specified number of items from the input channel.
// When the maximum is reached, Limit will close the "done" channel, to communicate to other goroutines
// that they can stop sending items.
func Limit[T any](maximum int, input <-chan T, done chan<- struct{}) <-chan T {

	result := make(chan T)

	go func() {

		defer close(result)
		defer close(done)

		if maximum <= 0 {
			return
		}

		for item := range input {

			result <- item

			maximum--

			if maximum <= 0 {
				return
			}
		}
	}()

	return result
}
