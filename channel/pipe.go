package channel

// Pipe reads the contents of one channel directly into another channel.
func Pipe[T any](input <-chan T, output chan<- T) {

	go func() {

		defer close(output)

		for item := range input {
			output <- item
		}
	}()
}

// PipeWithCancel reads the contents of one channel directly into another channel.
// If the "done" channel is closed, the pipe will stop.
func PipeWithCancel[T any](input <-chan T, output chan<- T, done <-chan struct{}) {

	go func() {

		defer close(output)

		for {
			select {
			case item, ok := <-input:
				if !ok {
					return
				}
				output <- item
			case <-done:
				return
			}
		}
	}()
}
