package slice

func Difference[T comparable](a []T, b []T) []T {

	// Create a slice to store the difference
	result := make([]T, len(a))

	// Iterate over a and add elements not in b to the diff slice
	for _, item := range a {
		if !Contains(b, item) {
			result = append(result, item)
		}
	}

	return result
}
