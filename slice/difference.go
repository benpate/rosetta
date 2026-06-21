package slice

// Difference returns a new slice of the elements in a that are not present in b.
func Difference[T comparable](a []T, b []T) []T {

	// Create a slice to store the difference
	result := make([]T, 0, len(a))

	// Iterate over a and add elements not in b to the diff slice
	for _, item := range a {
		if !Contains(b, item) {
			result = append(result, item)
		}
	}

	return result
}
