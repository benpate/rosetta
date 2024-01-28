package channel

import "github.com/davecgh/go-spew/spew"

// Beep is a simple channel that dumps all of its input to the console.
// It can be used to debug a channel's contents without disrupting the flow of data.
func Beep[T any](in <-chan T) <-chan T {

	result := make(chan T)

	go func() {
		defer close(result)
		for value := range in {
			spew.Dump(value)
			result <- value
		}
	}()

	return result
}
