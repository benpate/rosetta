package slice

import "math/rand/v2"

func Shuffle[T any](x []T) []T {
	rand.Shuffle(len(x), func(i, j int) { // NOSONAR: does not need to be cyptographically secure.
		x[i], x[j] = x[j], x[i]
	})
	return x
}
