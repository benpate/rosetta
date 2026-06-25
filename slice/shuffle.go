package slice

import "math/rand/v2"

// Shuffle randomizes the order of the elements in the slice in place, and returns it.
func Shuffle[T any](x []T) []T {
	rand.Shuffle(len(x), func(i, j int) { // NOSONAR: does not need to be cryptographically secure.
		x[i], x[j] = x[j], x[i]
	})
	return x
}
