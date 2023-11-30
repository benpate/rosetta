package slice

// Unique returns a new slice with all duplicate values removed.
func Unique[T comparable](original []T) []T {

	result := make([]T, 0, len(original))

	for _, item := range original {

		found := false

		for _, existing := range result {
			if item == existing {
				found = true
				break
			}
		}

		if !found {
			result = append(result, item)
		}
	}

	return result
}
