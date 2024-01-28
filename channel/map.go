package channel

// Map applies a mapping function to every value in the input channel, passing the mapped
// values out to the output channel.
func Map[Input any, Output any](input <-chan Input, mapper func(Input) Output) <-chan Output {

	result := make(chan Output, 1)

	go func() {
		defer close(result)
		for value := range input {
			result <- mapper(value)
		}
	}()

	return result
}
