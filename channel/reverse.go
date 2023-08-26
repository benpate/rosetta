package channel

// Reverse reads an entire channel into a stack, then writes it back out in reverse order
// Since it keeps the entire channel in memory, it should not be used for unbounded data sets
func Reverse[T any](input <-chan T) <-chan T {

	result := make(chan T)

	go func() {

		defer close(result)

		// This is a stack that we'll read in reverse order
		items := make([]T, 0)

		// Read all items from the input channel
		for item := range input {

			// Push the item onto the stack
			items = append(items, item)
		}

		// Write all items back to the output channel
		for i := len(items) - 1; i >= 0; i-- {

			// Send the item to the caller
			result <- items[i]
		}
	}()

	// Return the output channel to the caller
	return result
}
