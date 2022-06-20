package slice

func Map[T1 any, T2 any](source []T1, delta func(T1) T2) []T2 {

	result := make([]T2, len(source))

	for index := range source {
		result[index] = delta(source[index])
	}

	return result
}
