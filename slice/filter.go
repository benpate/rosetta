package slice

// Filter returns a new slice containing only the elements for which keep returns TRUE.
func Filter[T any](original []T, keep func(T) bool) []T {

	result := make([]T, 0, len(original))

	for _, value := range original {
		if keep(value) {
			result = append(result, value)
		}
	}

	return result
}
