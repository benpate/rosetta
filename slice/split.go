package slice

// Split returns the first element of the slice and a slice of the remaining elements.
// An empty slice yields the zero value and the original (empty) slice.
func Split[T any](slice []T) (T, []T) {

	switch len(slice) {

	case 0:
		var empty T
		return empty, slice

	case 1:
		return slice[0], []T{}

	default:
		return slice[0], slice[1:]
	}
}
