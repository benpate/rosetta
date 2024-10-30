package slice

// Find returns the first element in the slice that satisfies the provided function.
func Find[T any](slice []T, f func(T) bool) (T, bool) {

	for _, value := range slice {
		if f(value) {
			return value, true
		}
	}

	var result T
	return result, false
}
