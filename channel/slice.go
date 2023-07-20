package channel

func Slice[T any](channel <-chan T) []T {

	result := make([]T, 0)

	for item := range channel {
		result = append(result, item)
	}

	return result
}
