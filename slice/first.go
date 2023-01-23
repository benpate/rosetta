package slice

func First[T any](original []T, keep func(T) bool) T {
	for _, value := range original {
		if keep(value) {
			return value
		}
	}

	var empty T
	return empty
}
