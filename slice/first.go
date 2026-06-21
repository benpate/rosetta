package slice

// First returns the first element for which keep returns TRUE, or the zero value if none match.
func First[T any](original []T, keep func(T) bool) T {
	for _, value := range original {
		if keep(value) {
			return value
		}
	}

	var empty T
	return empty
}
