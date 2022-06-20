package slice

func Filter[T any](original []T, keep func(T) bool) []T {

	result := make([]T, len(original))

	index := 0

	for _, value := range original {
		if keep(value) {
			result[index] = value
			index = index + 1
		}
	}

	result = result[:index]
	return result
}
