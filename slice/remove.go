package slice

func RemoveAt[T comparable](slice []T, index int) []T {

	// Bounds check
	if index < 0 || index >= len(slice) {
		return slice
	}

	// Remove the index from the slice
	return append(slice[:index], slice[index+1:]...)
}
