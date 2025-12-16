package iterator

// Map converts an iterator into a slice of items.
// Deprecated: freeing up this namespace to use for new Go 1.23 range functions
func Map[In any, Out any](it Iterator, fn func(In) Out) []Out {

	var zeroValue In
	var currentValue In

	result := make([]Out, 0, it.Count())

	for it.Next(&currentValue) {
		result = append(result, fn(currentValue))
		currentValue = zeroValue
	}

	return result
}

// Slice converts an Iterator into a slice of items.
// You must include a constructor function that generates fully initialized values of the type you want to return.
// Deprecated: freeing up this namespace to use for new Go 1.23 range functions
func Slice[T any](iterator Iterator, constructor func() T) []T {

	result := make([]T, 0, iterator.Count())

	value := constructor() // nolint:scopeguard

	for iterator.Next(&value) {
		result = append(result, value)
		value = constructor()
	}

	return result
}

// Channel converts an Iterator into a channel of items.
// You must include a constructor function that generates fully initialized values of the type you want to return.
// Deprecated: freeing up this namespace to use for new Go 1.23 range functions
func Channel[T any](iterator Iterator, constructor func() T) chan T {

	result := make(chan T, 1) // Length of 1 to prevent blocking on the first item.

	go func() {

		defer close(result)

		value := constructor() // nolint:scopeguard

		for iterator.Next(&value) {
			result <- value
			value = constructor()
		}
	}()

	return result
}

// Channel converts an Iterator into a channel of items.
// You must include a constructor function that generates fully initialized values of the type you want to return.
// Deprecated: freeing up this namespace to use for new Go 1.23 range functions
func ChannelWithCancel[T any](iterator Iterator, constructor func() T, cancel <-chan bool) chan T {

	result := make(chan T, 1) // Length of 1 to prevent blocking on the first item.

	go func() {

		defer close(result)

		value := constructor() // nolint:scopeguard

		for iterator.Next(&value) {
			select {
			case <-cancel:
				return

			default:
				result <- value
			}

			value = constructor()
		}
	}()

	return result
}
